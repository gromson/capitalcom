package capitalcom

import (
	"context"
	"net/url"
	"time"
)

type trading struct {
	*Client
}

type (
	Deal struct {
		Date           time.Time      `json:"date"`
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
