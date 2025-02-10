package capitalcom

import (
	"context"
	"encoding/json"
	"net/url"
	"time"
)

type dealReferenceResponsePayload struct {
	DealReference string `json:"dealReference"`
}

type positions struct {
	*Client
}

type PositionDirection string

const (
	PositionDirectionBuy  PositionDirection = "BUY"
	PositionDirectionSell PositionDirection = "SELL"
)

type (
	positionsResponsePayload struct {
		Positions []PositionDetail `json:"positions"`
	}

	PositionDetail struct {
		Position Position `json:"position"`
		Market   Market   `json:"market"`
	}

	Position struct {
		ContractSize   int               `json:"contractSize"`
		CreatedDate    time.Time         `json:"-"`
		CreatedDateUTC time.Time         `json:"-"`
		DealID         string            `json:"dealId"`
		DealReference  string            `json:"dealReference"`
		WorkingOrderID string            `json:"workingOrderId"`
		Size           int               `json:"size"`
		Leverage       int               `json:"leverage"`
		UPL            float64           `json:"upl"`
		Direction      PositionDirection `json:"direction"`
		Level          float64           `json:"level"`
		Currency       string            `json:"currency"`
		GuaranteedStop bool              `json:"guaranteedStop"`
	}

	Market struct {
		InstrumentName           string    `json:"instrumentName"`
		Expiry                   string    `json:"expiry"`
		MarketStatus             string    `json:"marketStatus"`
		Epic                     string    `json:"epic"`
		Symbol                   string    `json:"symbol"`
		InstrumentType           string    `json:"instrumentType"`
		LotSize                  int       `json:"lotSize"`
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
		ScalingFactor            int       `json:"scalingFactor"`
		MarketModes              []string  `json:"marketModes"`
	}
)

func (m *Market) UnmarshalJSON(data []byte) error {
	type alias Market

	aux := &struct {
		UpdateTimeString    string `json:"updateTime"`
		UpdateTimeUTCString string `json:"updateTimeUTC"` //nolint:tagliatelle
		*alias
	}{
		alias: (*alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	var err error

	m.UpdateTime, err = time.Parse(dateFormat, aux.UpdateTimeString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	m.UpdateTimeUTC, err = time.Parse(dateFormat, aux.UpdateTimeUTCString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	return nil
}

func (p *Position) UnmarshalJSON(data []byte) error {
	type alias Position

	aux := &struct {
		CreatedDateString    string `json:"createdDate"`
		CreatedDateUTCString string `json:"createdDateUTC"` //nolint:tagliatelle
		*alias
	}{
		alias: (*alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	var err error

	p.CreatedDate, err = time.Parse(dateFormat, aux.CreatedDateString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	p.CreatedDateUTC, err = time.Parse(dateFormat, aux.CreatedDateUTCString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	return nil
}

// List retrieves all open positions for the authenticated user.
func (p *positions) List(ctx context.Context) ([]PositionDetail, error) {
	headers := p.tokens.headers()

	res, err := get[positionsResponsePayload](ctx, p.Client, "/positions", headers)
	if err != nil {
		return nil, err
	}

	p.tokens.updateTokens(res.httpResponse)

	return res.payload.Positions, nil
}

type (
	// OpenPositionRequest represents the payload to open a new position.
	OpenPositionRequest struct {
		Direction PositionDirection `json:"direction"`

		// Epic is an instrument epic identifier
		Epic string `json:"epic"`

		// Size is a deal size
		Size float64 `json:"size"`

		UpdatePositionRequest
	}

	// UpdatePositionRequest represents the payload to update a position.
	// todo: make validation or a builder that won't allow to create an invalid instance
	UpdatePositionRequest struct {
		// GuaranteedStop
		// - Default value: false
		// - If GuaranteedStop equals true, then set StopLevel, StopDistance or StopAmount
		// - Cannot be set if TrailingStop is true
		// - Cannot be set if hedgingMode is true
		GuaranteedStop bool `json:"guaranteedStop,omitempty"`

		// TrailingStop
		// - Default value: false
		// - If TrailingStop equals true, then set StopDistance
		// - Cannot be set if GuaranteedStop is true
		TrailingStop bool `json:"trailingStop,omitempty"`

		// StopLevel is a price level when a stop loss will be triggered
		StopLevel float64 `json:"stopLevel,omitempty"`

		// StopDistance Required parameter if TrailingStop is true
		StopDistance float64 `json:"stopDistance,omitempty"`

		// StopAmount Loss amount when a stop loss will be triggered
		StopAmount float64 `json:"stopAmount,omitempty"`

		// ProfitLevel is a price level when a take profit will be triggered
		ProfitLevel float64 `json:"profitLevel,omitempty"`

		// ProfitDistance is a distance between current and take profit triggering price
		ProfitDistance float64 `json:"profitDistance,omitempty"`

		// ProfitAmount is a profit amount when a take profit will be triggered
		ProfitAmount float64 `json:"profitAmount,omitempty"`
	}
)

// Open opens a new position with the specified parameters.
func (p *positions) Open(ctx context.Context, req OpenPositionRequest) (string, error) {
	headers := p.tokens.headers()

	res, err := post[dealReferenceResponsePayload](ctx, p.Client, "/positions", req, headers)
	if err != nil {
		return "", err
	}

	p.tokens.updateTokens(res.httpResponse)

	return res.payload.DealReference, nil
}

// Get retrieves the position details for the specified deal.
func (p *positions) Get(ctx context.Context, dealID string) (*PositionDetail, error) {
	headers := p.tokens.headers()

	res, err := get[PositionDetail](ctx, p.Client, "/positions/"+url.PathEscape(dealID), headers)
	if err != nil {
		return nil, err
	}

	p.tokens.updateTokens(res.httpResponse)

	return res.payload, nil
}

// Update updates the position for the specified deal.
func (p *positions) Update(ctx context.Context, dealID string, req UpdatePositionRequest) (string, error) {
	headers := p.tokens.headers()

	res, err := put[dealReferenceResponsePayload](ctx, p.Client, "/positions/"+url.PathEscape(dealID), req, headers)
	if err != nil {
		return "", err
	}

	return res.payload.DealReference, nil
}

// Close closes the position for the specified deal.
func (p *positions) Close(ctx context.Context, dealID string) error {
	headers := p.tokens.headers()

	res, err := del[dealReferenceResponsePayload](ctx, p.Client, "/positions/"+url.PathEscape(dealID), headers)
	if err != nil {
		return err
	}

	p.tokens.updateTokens(res.httpResponse)

	return nil
}
