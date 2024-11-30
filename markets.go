package capitalcom

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type markets struct {
	*Client
}

type (
	navigationNodesResponsePayload struct {
		Nodes []NavigationNode `json:"nodes"`
	}

	NavigationNode struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

func (m *markets) Categories(ctx context.Context) ([]NavigationNode, error) {
	headers := m.tokens.headers()

	res, err := get[navigationNodesResponsePayload](ctx, m.Client, "/marketnavigation", headers)
	if err != nil {
		return nil, err
	}

	m.tokens.updateTokens(res.httpResponse)

	return res.payload.Nodes, nil
}

func (m *markets) Subcategories(ctx context.Context, nodeID string, limit int) ([]NavigationNode, error) {
	headers := m.tokens.headers()

	query := url.Values{}

	if limit != 0 {
		query.Add("limit", strconv.Itoa(limit))
	}

	queryString := query.Encode()
	if queryString != "" {
		queryString = "?" + queryString
	}

	res, err := get[navigationNodesResponsePayload](ctx,
		m.Client, "/marketnavigation/"+url.PathEscape(nodeID)+queryString,
		headers)
	if err != nil {
		return nil, err
	}

	m.tokens.updateTokens(res.httpResponse)

	return res.payload.Nodes, nil
}

type (
	marketsResponsePayload struct {
		Markets []Market `json:"markets"`
	}

	DetailsParams struct {
		SearchTerm string
		Epics      []string
	}
)

func (p DetailsParams) toQueryString() string {
	values := url.Values{}

	if p.SearchTerm != "" {
		values.Add("searchTerm", p.SearchTerm)
	}

	if len(p.Epics) > 0 {
		values.Add("epics", url.QueryEscape(strings.Join(p.Epics, ",")))
	}

	return values.Encode()
}

func (m *markets) Details(ctx context.Context, params DetailsParams) ([]Market, error) {
	headers := m.tokens.headers()

	queryString := params.toQueryString()
	if queryString != "" {
		queryString = "?" + queryString
	}

	res, err := get[marketsResponsePayload](ctx, m.Client, "/markets"+queryString, headers)
	if err != nil {
		return nil, err
	}

	m.tokens.updateTokens(res.httpResponse)

	return res.payload.Markets, nil
}

type (
	Instrument struct {
		Epic                     string       `json:"epic"`
		Symbol                   string       `json:"symbol"`
		Expiry                   string       `json:"expiry"`
		Name                     string       `json:"name"`
		LotSize                  float64      `json:"lotSize"`
		Type                     string       `json:"type"`
		GuaranteedStopAllowed    bool         `json:"guaranteedStopAllowed"`
		StreamingPricesAvailable bool         `json:"streamingPricesAvailable"`
		Currency                 string       `json:"currency"`
		MarginFactor             float64      `json:"marginFactor"`
		MarginFactorUnit         string       `json:"marginFactorUnit"`
		OpeningHours             OpeningHours `json:"openingHours"`
		OvernightFee             OvernightFee `json:"overnightFee"`
	}

	OpeningHours struct {
		Mon  []string `json:"mon"`
		Tue  []string `json:"tue"`
		Wed  []string `json:"wed"`
		Thu  []string `json:"thu"`
		Fri  []string `json:"fri"`
		Sat  []string `json:"sat"`
		Sun  []string `json:"sun"`
		Zone string   `json:"zone"`
	}

	OvernightFee struct {
		LongRate            float64 `json:"longRate"`
		ShortRate           float64 `json:"shortRate"`
		SwapChargeTimestamp int     `json:"swapChargeTimestamp"`
		SwapChargeInterval  int     `json:"swapChargeInterval"`
	}

	DealingRules struct {
		MinStepDistance           Rule   `json:"minStepDistance"`
		MinDealSize               Rule   `json:"minDealSize"`
		MaxDealSize               Rule   `json:"maxDealSize"`
		MinSizeIncrement          Rule   `json:"minSizeIncrement"`
		MinGuaranteedStopDistance Rule   `json:"minGuaranteedStopDistance"`
		MinStopOrProfitDistance   Rule   `json:"minStopOrProfitDistance"`
		MaxStopOrProfitDistance   Rule   `json:"maxStopOrProfitDistance"`
		MarketOrderPreference     string `json:"marketOrderPreference"`
		TrailingStopsPreference   string `json:"trailingStopsPreference"`
	}

	Rule struct {
		Unit  string  `json:"unit"`
		Value float64 `json:"value"`
	}

	Snapshot struct {
		MarketStatus        string    `json:"marketStatus"`
		NetChange           float64   `json:"netChange"`
		PercentageChange    float64   `json:"percentageChange"`
		UpdateTime          time.Time `json:"updateTime"`
		DelayTime           int       `json:"delayTime"`
		Bid                 float64   `json:"bid"`
		Offer               float64   `json:"offer"`
		High                float64   `json:"high"`
		Low                 float64   `json:"low"`
		DecimalPlacesFactor int       `json:"decimalPlacesFactor"`
		ScalingFactor       int       `json:"scalingFactor"`
		MarketModes         []string  `json:"marketModes"`
	}

	MarketDetails struct {
		Instrument   Instrument   `json:"instrument"`
		DealingRules DealingRules `json:"dealingRules"`
		Snapshot     Snapshot     `json:"snapshot"`
	}
)

func (m *markets) Detail(ctx context.Context, epic string) (*MarketDetails, error) {
	headers := m.tokens.headers()

	res, err := get[MarketDetails](ctx, m.Client, "/markets/"+url.PathEscape(epic), headers)
	if err != nil {
		return nil, err
	}

	m.tokens.updateTokens(res.httpResponse)

	return res.payload, nil
}
