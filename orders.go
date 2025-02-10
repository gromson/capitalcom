package capitalcom

import (
	"context"
	"encoding/json"
	"net/url"
	"time"
)

type orders struct {
	*Client
}

type OrderType string

const (
	LimitOrder OrderType = "LIMIT"
	StopOrder  OrderType = "STOP"
)

type (
	workingOrdersResponsePayload struct {
		WorkingOrders []WorkingOrderDetail `json:"workingOrders"`
	}

	WorkingOrderDetail struct {
		WorkingOrderData WorkingOrderData `json:"workingOrderData"`
		MarketData       MarketData       `json:"marketData"`
	}

	WorkingOrderData struct {
		DealID          string            `json:"dealId"`
		Direction       PositionDirection `json:"direction"`
		Epic            string            `json:"epic"`
		OrderSize       float64           `json:"orderSize"`
		Leverage        float64           `json:"leverage"`
		OrderLevel      float64           `json:"orderLevel"`
		TimeInForce     string            `json:"timeInForce"`
		GoodTillDate    time.Time         `json:"-"`
		GoodTillDateUTC time.Time         `json:"-"`
		CreatedDate     time.Time         `json:"-"`
		CreatedDateUTC  time.Time         `json:"-"`
		GuaranteedStop  bool              `json:"guaranteedStop"`
		OrderType       string            `json:"orderType"`
		StopDistance    float64           `json:"stopDistance"`
		ProfitDistance  float64           `json:"profitDistance"`
		TrailingStop    bool              `json:"trailingStop"`
		CurrencyCode    string            `json:"currencyCode"`
	}

	MarketData struct {
		InstrumentName           string    `json:"instrumentName"`
		Expiry                   string    `json:"expiry"`
		MarketStatus             string    `json:"marketStatus"`
		Epic                     string    `json:"epic"`
		Symbol                   string    `json:"symbol"`
		InstrumentType           string    `json:"instrumentType"`
		LotSize                  float64   `json:"lotSize"`
		High                     float64   `json:"high"`
		Low                      float64   `json:"low"`
		PercentageChange         float64   `json:"percentageChange"`
		NetChange                float64   `json:"netChange"`
		Bid                      float64   `json:"bid"`
		Offer                    float64   `json:"offer"`
		UpdateTime               time.Time `json:"-"`
		UpdateTimeUTC            time.Time `json:"-"`
		DelayTime                int       `json:"delayTime"`
		StreamingPricesAvailable bool      `json:"streamingPricesAvailable"`
		ScalingFactor            float64   `json:"scalingFactor"`
		MarketModes              []string  `json:"marketModes"`
	}
)

func (d *MarketData) UnmarshalJSON(data []byte) error {
	type alias MarketData

	aux := &struct {
		UpdateTimeString    string `json:"updateTime"`
		UpdateTimeUTCString string `json:"updateTimeUTC"` //nolint:tagliatelle
		*alias
	}{
		alias: (*alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	var err error

	d.UpdateTime, err = time.Parse(dateFormat, aux.UpdateTimeString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	d.UpdateTimeUTC, err = time.Parse(dateFormat, aux.UpdateTimeUTCString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	return nil
}

func (a *WorkingOrderData) UnmarshalJSON(data []byte) error {
	type alias WorkingOrderData

	aux := &struct {
		GoodTillDateString    string `json:"goodTillDate"`
		GoodTillDateUTCString string `json:"goodTillDateUTC"` //nolint:tagliatelle
		CreatedDateString     string `json:"createdDate"`
		CreatedDateUTCString  string `json:"createdDateUTC"` //nolint:tagliatelle
		*alias
	}{
		alias: (*alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	var err error

	a.GoodTillDate, err = time.Parse(dateFormat, aux.GoodTillDateString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	a.GoodTillDateUTC, err = time.Parse(dateFormat, aux.GoodTillDateUTCString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	a.CreatedDate, err = time.Parse(dateFormat, aux.CreatedDateString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	a.CreatedDateUTC, err = time.Parse(dateFormat, aux.CreatedDateUTCString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	return nil
}

func (o *orders) List(ctx context.Context) ([]WorkingOrderDetail, error) {
	headers := o.tokens.headers()

	res, err := get[workingOrdersResponsePayload](ctx, o.Client, "/workingorders", headers)
	if err != nil {
		return nil, err
	}

	o.tokens.updateTokens(res.httpResponse)

	return res.payload.WorkingOrders, nil
}

type (
	// CreateOrderRequest represents the payload to create a new order.
	CreateOrderRequest struct {
		Direction PositionDirection `json:"direction"`

		// Epic is an instrument epic identifier
		Epic string `json:"epic"`

		// Size is a order size
		Size float64 `json:"size"`

		// Type is an order type
		Type OrderType `json:"type"`

		UpdateOrderRequest
	}

	// UpdateOrderRequest represents the payload to update an existing order.
	// todo: make validation or a builder that won't allow to create an invalid instance
	UpdateOrderRequest struct {
		// Level - the order price
		Level float64 `json:"level"`

		// GoodTillDate - order cancellation date in UTC time
		GoodTillDate time.Time `json:"-"`

		// GuaranteedStop must be true if a guaranteed stop is required.
		// - Default value: false
		// - If GuaranteedStop equals true, then set StopLevel, StopDistance or StopAmount
		// - Cannot be set if TrailingStop is true
		// - Cannot be set if hedgingMode is true
		GuaranteedStop bool `json:"guaranteedStop,omitempty"`

		// TrailingStop must be true if a trailing stop is required.
		// - Default value: false
		// - If TrailingStop equals true, then set StopDistance
		// - Cannot be set if GuaranteedStop is true
		TrailingStop bool `json:"trailingStop,omitempty"`

		// StopDistance is a price level when a stop loss will be triggered
		StopLevel float64 `json:"stopLevel,omitempty"`

		// StopDistance is a distance between current and stop loss triggering price.
		// Required parameter if trailingStop is true
		StopDistance float64 `json:"stopDistance,omitempty"`

		// StopAmount is a loss amount when a stop loss will be triggered
		StopAmount float64 `json:"stopAmount,omitempty"`

		// ProfitLevel is a price level when a take profit will be triggered
		ProfitLevel float64 `json:"profitLevel,omitempty"`

		// ProfitDistance is a distance between current and take profit triggering price
		ProfitDistance float64 `json:"profitDistance,omitempty"`

		// ProfitAmount is a profit amount when a take profit will be triggered
		ProfitAmount float64 `json:"profitAmount,omitempty"`
	}
)

func (ur UpdateOrderRequest) MarshalJSON() ([]byte, error) {
	type alias UpdateOrderRequest

	aux := struct {
		GoodTillDateString string `json:"goodTillDate,omitempty"`
		alias
	}{
		alias: alias(ur),
	}

	if !ur.GoodTillDate.IsZero() {
		aux.GoodTillDateString = ur.GoodTillDate.Format(dateFormat)
	}

	data, err := json.Marshal(aux)
	if err != nil {
		return nil, NewRequestPayloadEncodingError(err)
	}

	return data, nil
}

func (o *orders) Create(ctx context.Context, req CreateOrderRequest) (string, error) {
	headers := o.tokens.headers()

	res, err := post[dealReferenceResponsePayload](ctx, o.Client, "/workingorders", req, headers)
	if err != nil {
		return "", err
	}

	o.tokens.updateTokens(res.httpResponse)

	return res.payload.DealReference, nil
}

func (o *orders) Update(ctx context.Context, dealID string, req UpdateOrderRequest) (string, error) {
	headers := o.tokens.headers()

	res, err := put[dealReferenceResponsePayload](ctx, o.Client, "/workingorders/"+url.PathEscape(dealID), req, headers)
	if err != nil {
		return "", err
	}

	o.tokens.updateTokens(res.httpResponse)

	return res.payload.DealReference, nil
}

func (o *orders) Delete(ctx context.Context, dealID string) (string, error) {
	headers := o.tokens.headers()

	res, err := del[dealReferenceResponsePayload](ctx, o.Client, "/workingorders/"+url.PathEscape(dealID), headers)
	if err != nil {
		return "", err
	}

	o.tokens.updateTokens(res.httpResponse)

	return res.payload.DealReference, err
}
