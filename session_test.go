package capitalcom_test

import (
	"net/http"

	"github.com/gromson/capitalcom"
)

const (
	identifier = "client@example.com"
	password   = "password"
)

const expectedAPIKey = "apikey"

const (
	expectedKeySecurityToken = "KeySecurityToken"
	expectedCST              = "CST"
)

func handleSessionCreation(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get(capitalcom.HeaderAPIKey)

	if apiKey != expectedAPIKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	w.Header().Set(capitalcom.HeaderKeySecurityToken, expectedKeySecurityToken)
	w.Header().Set(capitalcom.HeaderTokenCST, expectedCST)

	_, _ = w.Write([]byte(`{
		"accountType": "CFD",
		"accountInfo": {
			"balance": 92.89,
			"deposit": 90.38,
			"profitLoss": 2.51,
			"available": 64.66
		},
		"currencyIsoCode": "USD",
		"currencySymbol": "$",
		"currentAccountId": "12345678901234567",
		"streamingHost": "wss://api-streaming-capital.backend-capital.com/",
		"accounts": [
			{
				"accountId": "12345678901234567",
				"accountName": "USD",
				"preferred": true,
				"accountType": "CFD",
				"currency": "USD",
				"symbol": "$",
				"balance": {
					"balance": 92.89,
					"deposit": 90.38,
					"profitLoss": 2.51,
					"available": 64.66
				}
			},
			{
				"accountId": "12345678907654321",
				"accountName": "EUR",
				"preferred": false,
				"accountType": "CFD",
				"currency": "EUR",
				"symbol": "€",
				"balance": {
					"balance": 0.0,
					"deposit": 0.0,
					"profitLoss": 0.0,
					"available": 0.0
				}
			}
		],
		"clientId": "12345678",
		"timezoneOffset": 3,
		"hasActiveDemoAccounts": true,
		"hasActiveLiveAccounts": true,
		"trailingStopsEnabled": false
	}`))
}

func isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	gotKeySecurityToken := r.Header.Get(capitalcom.HeaderKeySecurityToken)
	gotCST := r.Header.Get(capitalcom.HeaderTokenCST)

	if gotKeySecurityToken != expectedKeySecurityToken || gotCST != expectedCST {
		w.WriteHeader(http.StatusUnauthorized)

		return false
	}

	w.Header().Set(capitalcom.HeaderKeySecurityToken, expectedKeySecurityToken)
	w.Header().Set(capitalcom.HeaderTokenCST, expectedCST)

	return true
}
