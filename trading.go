package capitalcom

import (
	"context"
	"encoding/json"
	"net/url"
	"time"
)

type trading struct {
	*Client
}

type (
	Deal struct {
		Date           time.Time      `json:"-"`
		Status         string         `json:"status"`
		DealStatus     string         `json:"dealStatus"`
		Epic           string         `json:"epic"`
		DealReference  string         `json:"dealReference"`
		DealID         string         `json:"dealId"`
		AffectedDeals  []AffectedDeal `json:"affectedDeals"`
		Level          float64        `json:"level"`
		Size           float64        `json:"size"`
		Direction      string         `json:"direction"`
		GuaranteedStop bool           `json:"guaranteedStop"`
		TrailingStop   bool           `json:"trailingStop"`
	}

	AffectedDeal struct {
		DealID string `json:"dealId"`
		Status string `json:"status"`
	}
)

func (d *Deal) UnmarshalJSON(data []byte) error {
	type alias Deal

	aux := struct {
		DateString string `json:"date"`
		*alias
	}{
		alias: (*alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	var err error

	d.Date, err = time.Parse(dateFormat, aux.DateString)
	if err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	return nil
}

// Confirm retrieves the confirmation details for a given deal reference.
func (t *trading) Confirm(ctx context.Context, dealReference string) (*Deal, error) {
	headers := t.tokens.headers()

	res, err := get[Deal](ctx, t.Client, "/confirms/"+url.PathEscape(dealReference), headers)
	if err != nil {
		return nil, err
	}

	t.tokens.updateTokens(res.httpResponse)

	return res.payload, nil
}
