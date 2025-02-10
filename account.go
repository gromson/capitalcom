package capitalcom

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

type account struct {
	*Client
}

type (
	accountListResponsePayload struct {
		Accounts []Account `json:"accounts"`
	}

	Account struct {
		AccountID   string  `json:"accountId"`
		AccountName string  `json:"accountName"`
		Status      string  `json:"status"`
		AccountType string  `json:"accountType"`
		Preferred   bool    `json:"preferred"`
		Balance     Balance `json:"balance"`
		Currency    string  `json:"currency"`
		Symbol      string  `json:"symbol"`
	}

	Balance struct {
		Balance    float64 `json:"balance"`
		Deposit    float64 `json:"deposit"`
		ProfitLoss float64 `json:"profitLoss"`
		Available  float64 `json:"available"`
	}
)

// List retrieves a list of accounts associated with the authenticated user.
func (a *account) List(ctx context.Context) ([]Account, error) {
	headers := a.tokens.headers()

	res, err := get[accountListResponsePayload](ctx, a.Client, "/account", headers)
	if err != nil {
		return nil, err
	}

	a.tokens.updateTokens(res.httpResponse)

	return res.payload.Accounts, nil
}

type (
	Leverage struct {
		Current   int   `json:"current"`
		Available []int `json:"available"`
	}

	Leverages struct {
		Shares           Leverage `json:"SHARES"`           //nolint:tagliatelle
		Currencies       Leverage `json:"CURRENCIES"`       //nolint:tagliatelle
		Indices          Leverage `json:"INDICES"`          //nolint:tagliatelle
		Cryptocurrencies Leverage `json:"CRYPTOCURRENCIES"` //nolint:tagliatelle
		Commodities      Leverage `json:"COMMODITIES"`      //nolint:tagliatelle
	}

	Preferences struct {
		HedgingMode bool      `json:"hedgingMode"`
		Leverages   Leverages `json:"leverages"`
	}
)

// Preferences retrieves the preferences of the authenticated user.
func (a *account) Preferences(ctx context.Context) (*Preferences, error) {
	headers := a.tokens.headers()

	res, err := get[Preferences](ctx, a.Client, "/accounts/preferences", headers)
	if err != nil {
		return nil, err
	}

	a.tokens.updateTokens(res.httpResponse)

	return res.payload, nil
}

type (
	updatePreferencesRequestPayload struct {
		Leverages   *UpdateLeverages `json:"leverages,omitempty"`
		HedgingMode bool             `json:"hedgingMode"`
	}

	UpdateLeverages struct {
		Shares           int `json:"SHARES,omitempty"`           //nolint:tagliatelle
		Currencies       int `json:"CURRENCIES,omitempty"`       //nolint:tagliatelle
		Indices          int `json:"INDICES,omitempty"`          //nolint:tagliatelle
		Cryptocurrencies int `json:"CRYPTOCURRENCIES,omitempty"` //nolint:tagliatelle
		Commodities      int `json:"COMMODITIES,omitempty"`      //nolint:tagliatelle
	}
)

// UpdatePreferences updates the preferences of the authenticated user.
func (a *account) UpdatePreferences(
	ctx context.Context,
	leverages *UpdateLeverages,
	hedgingMode bool,
) (string, error) {
	headers := a.tokens.headers()

	reqPayload := updatePreferencesRequestPayload{
		Leverages:   leverages,
		HedgingMode: hedgingMode,
	}

	res, err := put[StatusResponsePayload](ctx, a.Client, "/accounts/preferences", reqPayload, headers)
	if err != nil {
		return "", err
	}

	a.tokens.updateTokens(res.httpResponse)

	return res.payload.Status, nil
}

type (
	activitiesResponsePayload struct {
		Activities []Activity `json:"activities"`
	}

	Activity struct {
		Date    time.Time `json:"-"`
		DateUTC time.Time `json:"-"`
		Epic    string    `json:"epic"`
		DealID  string    `json:"dealId"`
		Source  string    `json:"source"`
		Type    string    `json:"type"`
		Status  string    `json:"status"`
	}

	ActivityParams struct {
		From       time.Time
		To         time.Time
		LastPeriod int
		Detailed   bool
		DealID     string
		// Filter is a FIQL string, https://open-api.capital.com/#tag/Accounts/paths/~1api~1v1~1history~1activity/get
		Filter string
	}
)

func (a *Activity) UnmarshalJSON(data []byte) error {
	type alias Activity

	aux := &struct {
		DateString    string `json:"date"`
		DateUTCString string `json:"dateUTC"` //nolint:tagliatelle
		*alias
	}{
		alias: (*alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	var err error

	a.Date, err = time.Parse(dateFormat, aux.DateString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	a.DateUTC, err = time.Parse(dateFormat, aux.DateUTCString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	return nil
}

func (a ActivityParams) toQueryString() string {
	values := url.Values{}

	if !a.From.IsZero() {
		values.Add("from", a.From.Format(dateFormat))
	}

	if !a.To.IsZero() {
		values.Add("to", a.To.Format(dateFormat))
	}

	if a.From.IsZero() && a.To.IsZero() && a.LastPeriod > 0 {
		values.Add("lastPeriod", strconv.Itoa(a.LastPeriod))
	}

	if a.Detailed {
		values.Add("detailed", "true")
	}

	if a.DealID != "" {
		values.Add("dealId", a.DealID)
	}

	if a.Filter != "" {
		values.Add("filter", a.Filter)
	}

	return values.Encode()
}

// ActivityHistory retrieves the activity history of the authenticated user based on the provided parameters.
func (a *account) ActivityHistory(ctx context.Context, params ActivityParams) ([]Activity, error) {
	headers := a.tokens.headers()

	queryString := params.toQueryString()
	if queryString != "" {
		queryString = "?" + queryString
	}

	res, err := get[activitiesResponsePayload](ctx, a.Client, "/history/activity"+queryString, headers)
	if err != nil {
		return nil, err
	}

	a.tokens.updateTokens(res.httpResponse)

	return res.payload.Activities, nil
}

type TransactionType string

const (
	TransactionTypeInactivityFee             TransactionType = "INACTIVITY_FEE"
	TransactionTypeReserve                   TransactionType = "RESERVE"
	TransactionTypeVoid                      TransactionType = "VOID"
	TransactionTypeUnreserve                 TransactionType = "UNRESERVE"
	TransactionTypeWriteOffOrCredit          TransactionType = "WRITE_OFF_OR_CREDIT" //nolint:gosec
	TransactionTypeCreditFacility            TransactionType = "CREDIT_FACILITY"
	TransactionTypeFxCommission              TransactionType = "FX_COMMISSION"
	TransactionTypeComplaintSettlement       TransactionType = "COMPLAINT_SETTLEMENT"
	TransactionTypeDeposit                   TransactionType = "DEPOSIT"
	TransactionTypeWithdrawal                TransactionType = "WITHDRAWAL"
	TransactionTypeRefund                    TransactionType = "REFUND"
	TransactionTypeWithdrawalMoneyBack       TransactionType = "WITHDRAWAL_MONEY_BACK"
	TransactionTypeTrade                     TransactionType = "TRADE"
	TransactionTypeSwap                      TransactionType = "SWAP"
	TransactionTypeTradeCommission           TransactionType = "TRADE_COMMISSION"
	TransactionTypeTradeCommissionGsl        TransactionType = "TRADE_COMMISSION_GSL"
	TransactionTypeNegativeBalanceProtection TransactionType = "NEGATIVE_BALANCE_PROTECTION"
	TransactionTypeTradeCorrection           TransactionType = "TRADE_CORRECTION"
	TransactionTypeChargeback                TransactionType = "CHARGEBACK"
	TransactionTypeAdjustment                TransactionType = "ADJUSTMENT"
	TransactionTypeBonus                     TransactionType = "BONUS"
	TransactionTypeTransfer                  TransactionType = "TRANSFER"
	TransactionTypeCorporateAction           TransactionType = "CORPORATE_ACTION"
	TransactionTypeConversion                TransactionType = "CONVERSION"
	TransactionTypeRebate                    TransactionType = "REBATE"
	TransactionTypeTradeSlippageProtection   TransactionType = "TRADE_SLIPPAGE_PROTECTION"
)

type (
	transactionsResponsePayload struct {
		Transactions []Transaction `json:"transactions"`
	}

	TransactionParams struct {
		From       time.Time
		To         time.Time
		LastPeriod int
		Type       TransactionType
	}

	Transaction struct {
		Date            time.Time       `json:"-"`
		DateUTC         time.Time       `json:"-"`
		InstrumentName  string          `json:"instrumentName"`
		TransactionType TransactionType `json:"transactionType"`
		Note            string          `json:"note"`
		Reference       string          `json:"reference"`
		Size            string          `json:"size"`
		Currency        string          `json:"currency"`
		Status          string          `json:"status"`
	}
)

func (p *Transaction) UnmarshalJSON(data []byte) error {
	type alias Transaction

	aux := &struct {
		DateString    string `json:"date"`
		DateUTCString string `json:"dateUTC"` //nolint:tagliatelle
		*alias
	}{
		alias: (*alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	var err error

	p.Date, err = time.Parse(dateFormat, aux.DateString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	p.DateUTC, err = time.Parse(dateFormat, aux.DateUTCString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	return nil
}

func (p TransactionParams) toQueryString() string {
	values := url.Values{}

	if !p.From.IsZero() {
		values.Add("from", p.From.Format(dateFormat))
	}

	if !p.To.IsZero() {
		values.Add("to", p.To.Format(dateFormat))
	}

	if p.From.IsZero() && p.To.IsZero() && p.LastPeriod > 0 {
		values.Add("lastPeriod", strconv.Itoa(p.LastPeriod))
	}

	if p.Type != "" {
		values.Add("type", string(p.Type))
	}

	return values.Encode()
}

// TransactionHistory retrieves the transaction history of the authenticated user based on the provided parameters.
func (a *account) TransactionHistory(ctx context.Context, params TransactionParams) ([]Transaction, error) {
	headers := a.tokens.headers()

	queryString := params.toQueryString()
	if queryString != "" {
		queryString = "?" + queryString
	}

	res, err := get[transactionsResponsePayload](ctx, a.Client, "/history/transactions"+queryString, headers)
	if err != nil {
		return nil, err
	}

	a.tokens.updateTokens(res.httpResponse)

	return res.payload.Transactions, nil
}

type (
	topUpDemoAccountRequestPayload struct {
		Amount float64 `json:"amount"`
	}

	topUpDemoAccountResponsePayload struct {
		Successful string `json:"successful"`
	}
)

// TopUpDemoAccount adjusts the balance of the demo account by the specified amount.
func (a *account) TopUpDemoAccount(ctx context.Context, amount float64) (string, error) {
	headers := a.tokens.headers()

	reqPayload := topUpDemoAccountRequestPayload{
		Amount: amount,
	}

	res, err := post[topUpDemoAccountResponsePayload](ctx, a.Client, "/accounts/topUp", reqPayload, headers)
	if err != nil {
		return "", err
	}

	a.tokens.updateTokens(res.httpResponse)

	return res.payload.Successful, nil
}
