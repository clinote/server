/**
 * clinot.es server
 * Copyright (C) 2016 Sebastian Müller
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.

 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package route

import (
	"net/http"

	"github.com/clinotes/server/data"
	"github.com/keighl/postmark"
)

// APIRequestStructCreateUser is
type APIRequestStructCreateUser struct {
	Address string `json:"address"`
}

var postmarkTemplateIDauthUserCreate = 1012641

// APIRouteAccountCreate is
var APIRouteAccountCreate = Route{
	"/account/create",
	func(res http.ResponseWriter, req *http.Request) {
		var reqData APIRequestStructCreateUser
		if ensureJSONPayload(req, res, &reqData) != nil {
			return
		}

		account := data.AccountNew(reqData.Address)
		account, err := account.Store()

		// If account cannot be created, fail
		if err != nil {
			writeJSONError(res, "Unable to create account")
			return
		}

		token := data.TokenNew(account.ID(), data.TokenTypeMaintenace)
		tokenRaw := token.Raw()
		token, err = token.Store()

		// If token cannot be created, fail and remove user
		if err != nil {
			account.Remove()
			writeJSONError(res, "Unable to create account")
			return
		}

		// Send confirmation mail using Postmark
		_, err = pmark.SendTemplatedEmail(postmark.TemplatedEmail{
			TemplateId: int64(postmarkTemplateIDauthUserCreate),
			TemplateModel: map[string]interface{}{
				"token": tokenRaw,
			},
			From:    "mail@clinot.es",
			To:      account.Address(),
			ReplyTo: "\"CLINotes\" <mail@clinot.es>",
		})

		// If mail cannot be sent, fail and remove user
		if err != nil {
			account.Remove()
			writeJSONError(res, "Unable to create account")
			return
		}

		// Done!
		writeJSONResponse(res)
	},
}
