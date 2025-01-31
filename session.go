package capitalcom

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type SessionAccount struct {
	AccountType           string    `json:"accountType"`
	AccountInfo           Balance   `json:"accountInfo"`
	CurrencyIsoCode       string    `json:"currencyIsoCode"`
	CurrencySymbol        string    `json:"currencySymbol"`
	CurrentAccountID      string    `json:"currentAccountId"`
	StreamingHost         string    `json:"streamingHost"`
	Accounts              []Account `json:"accounts"`
	ClientID              string    `json:"clientId"`
	TimezoneOffset        int       `json:"timezoneOffset"`
	HasActiveDemoAccounts bool      `json:"hasActiveDemoAccounts"`
	HasActiveLiveAccounts bool      `json:"hasActiveLiveAccounts"`
	TrailingStopsEnabled  bool      `json:"trailingStopsEnabled"`
}

type EncryptionKey struct {
	EncryptionKey string    `json:"encryptionKey"`
	TimeStamp     time.Time `json:"timeStamp"`
}

func (e *EncryptionKey) UnmarshalJSON(data []byte) error {
	var aux struct {
		EncryptionKey string `json:"encryptionKey"`
		TimeStamp     int64  `json:"timeStamp"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	e.EncryptionKey = aux.EncryptionKey
	e.TimeStamp = time.UnixMilli(aux.TimeStamp)

	return nil
}

type SessionData struct {
	ClientID       string `json:"clientId"`
	AccountID      string `json:"accountId"`
	TimezoneOffset int    `json:"timezoneOffset"`
	Locale         string `json:"locale"`
	Currency       string `json:"currency"`
	StreamEndpoint string `json:"streamEndpoint"`
}

type AccountStatus struct {
	TrailingStopsEnabled  bool `json:"trailingStopsEnabled"`
	DealingEnabled        bool `json:"dealingEnabled"`
	HasActiveDemoAccounts bool `json:"hasActiveDemoAccounts"`
	HasActiveLiveAccounts bool `json:"hasActiveLiveAccounts"`
}

type createSessionPayload struct {
	Identifier        string `json:"identifier"`
	Password          string `json:"password"`
	EncryptedPassword bool   `json:"encryptedPassword,omitempty"`
}

type switchAccountPayload struct {
	AccountID string `json:"accountId"`
}

type session struct {
	*Client
}

func (s *session) CreateNew(
	ctx context.Context,
	identifier string,
	password string,
	passwordIsEncrypted bool,
) (*SessionAccount, error) {
	reqPayload := createSessionPayload{
		Identifier:        identifier,
		Password:          password,
		EncryptedPassword: passwordIsEncrypted,
	}

	headers := make(http.Header)
	headers.Set(HeaderAPIKey, s.apiKey) //nolint:canonicalheader

	res, err := post[SessionAccount](ctx, s.Client, "/session", reqPayload, headers)
	if err != nil {
		return nil, err
	}

	s.tokens.updateTokens(res.httpResponse)

	return res.payload, nil
}

func (s *session) EncryptionKey(ctx context.Context) (*EncryptionKey, error) {
	headers := make(http.Header)
	headers.Set(HeaderAPIKey, s.apiKey) //nolint:canonicalheader

	res, err := get[EncryptionKey](ctx, s.Client, "/session/encryptionKey", headers)
	if err != nil {
		return nil, err
	}

	return res.payload, nil
}

func (s *session) Details(ctx context.Context) (*SessionData, error) {
	headers := s.tokens.headers()

	res, err := get[SessionData](ctx, s.Client, "/session", headers)
	if err != nil {
		return nil, err
	}

	s.tokens.updateTokens(res.httpResponse)

	return res.payload, nil
}

func (s *session) SwitchActiveAccount(ctx context.Context, accountID string) (*AccountStatus, error) {
	header := s.tokens.headers()

	reqPayload := switchAccountPayload{AccountID: accountID}

	res, err := put[AccountStatus](ctx, s.Client, "/session", reqPayload, header)
	if err != nil {
		return nil, err
	}

	s.tokens.updateTokens(res.httpResponse)

	return res.payload, nil
}

func (s *session) LogOut(ctx context.Context) (string, error) {
	headers := s.tokens.headers()

	res, err := del[StatusResponsePayload](ctx, s.Client, "/session", headers)
	if err != nil {
		return "", err
	}

	s.tokens.updateTokens(res.httpResponse)

	return res.payload.Status, nil
}
