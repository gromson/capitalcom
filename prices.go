package capitalcom

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

type prices struct {
	*Client
}

type Resolution string

const (
	ResolutionMinute   Resolution = "MINUTE"
	ResolutionMinute5  Resolution = "MINUTE_5"
	ResolutionMinute15 Resolution = "MINUTE_15"
	ResolutionMinute30 Resolution = "MINUTE_30"
	ResolutionHour     Resolution = "HOUR"
	ResolutionHour4    Resolution = "HOUR_4"
	ResolutionDay      Resolution = "DAY"
	ResolutionWeek     Resolution = "WEEK"
)

type PricesParams struct {
	Resolution Resolution `json:"resolution"`
	Max        int        `json:"max"`
	From       time.Time  `json:"from"`
	To         time.Time  `json:"to"`
}

func (p *PricesParams) toQueryString() string {
	values := url.Values{}

	if p.Resolution != "" {
		values.Add("resolution", string(p.Resolution))
	}

	if p.Max > 0 {
		values.Add("max", strconv.Itoa(p.Max))
	}

	if !p.From.IsZero() {
		values.Add("from", p.From.Format(dateFormat))
	}

	if !p.To.IsZero() {
		values.Add("to", p.To.Format(dateFormat))
	}

	return values.Encode()
}

type (
	Prices struct {
		Prices         []Price `json:"prices"`
		InstrumentType string  `json:"instrumentType"`
	}

	Price struct {
		SnapshotTime     time.Time
		SnapshotTimeUTC  time.Time
		OpenPrice        PriceData `json:"openPrice"`
		ClosePrice       PriceData `json:"closePrice"`
		HighPrice        PriceData `json:"highPrice"`
		LowPrice         PriceData `json:"lowPrice"`
		LastTradedVolume int       `json:"lastTradedVolume"`
	}

	PriceData struct {
		Bid float64 `json:"bid"`
		Ask float64 `json:"ask"`
	}
)

func (pr *Price) UnmarshalJSON(data []byte) error {
	type Alias Price

	aux := &struct {
		SnapshotTimeString    string `json:"snapshotTime"`
		SnapshotTimeUTCString string `json:"snapshotTimeUTC"` //nolint:tagliatelle
		*Alias
	}{
		Alias: (*Alias)(pr),
	}

	if err := json.Unmarshal(data, &aux); err != nil { //nolint:musttag
		return NewResponsePayloadDecodingError(err)
	}

	var err error

	pr.SnapshotTime, err = time.Parse(dateFormat, aux.SnapshotTimeString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	pr.SnapshotTimeUTC, err = time.Parse(dateFormat, aux.SnapshotTimeUTCString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	return nil
}

func (p *prices) History(ctx context.Context, epic string, params PricesParams) (*Prices, error) {
	headers := p.tokens.headers()

	queryString := params.toQueryString()
	if queryString != "" {
		queryString = "?" + queryString
	}

	res, err := get[Prices](ctx, p.Client, "/prices/"+url.PathEscape(epic)+queryString, headers)
	if err != nil {
		return nil, err
	}

	p.tokens.updateTokens(res.httpResponse)

	return res.payload, nil
}
