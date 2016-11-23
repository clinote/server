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
	"errors"
	"net/http"

	"github.com/clinotes/server/data"
)

// APIRequestStructVerifyUser is
type APIRequestStructVerifyUser struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

// APIRouteAccountVerify is
var APIRouteAccountVerify = Route{
	"/account/verify",
	func(res http.ResponseWriter, req *http.Request) (error, interface{}) {
		// Parse JSON request
		var reqData APIRequestStructVerifyUser
		if err := checkJSONBody(req, res, &reqData); err != nil {
			return err, nil
		}

		// Get account
		account, err := data.AccountByAddress(reqData.Address)
		if err != nil {
			return errors.New("Unknown account address"), nil
		}

		// Check if account has requested token
		_, err = account.GetToken(reqData.Token, data.TokenTypeMaintenace)
		if err != nil {
			return errors.New("Unable to use provided token"), nil
		}

		// Verify account
		account, err = account.Verify()
		if err != nil {
			return errors.New("Unable to use provided token"), nil
		}

		_, err = sendTokenWithTemplate(account.Address(), reqData.Token, conf.TemplateConfirm)
		if err != nil {
			return errors.New("Unable to send verification mail"), nil
		}

		return nil, nil
	},
}
