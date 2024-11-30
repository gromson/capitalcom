package capitalcom

import (
	"context"
	"net/url"
	"strings"
)

type sentiment struct {
	*Client
}

type ClientSentiment struct {
	MarketID                string  `json:"marketId"`
	LongPositionPercentage  float64 `json:"longPositionPercentage"`
	ShortPositionPercentage float64 `json:"shortPositionPercentage"`
}

type clientSentimentsResponsePayload struct {
	ClientSentiments []ClientSentiment `json:"clientSentiments"`
}

func (s *sentiment) List(ctx context.Context, marketIDs []string) ([]ClientSentiment, error) {
	headers := s.tokens.headers()

	query := url.Values{}
	query.Set("marketIds", url.QueryEscape(strings.Join(marketIDs, ",")))

	queryString := query.Encode()
	if queryString != "" {
		queryString = "?" + queryString
	}

	res, err := get[clientSentimentsResponsePayload](ctx, s.Client, "/clientsentiment"+queryString, headers)
	if err != nil {
		return nil, err
	}

	s.tokens.updateTokens(res.httpResponse)

	return res.payload.ClientSentiments, nil
}

func (s *sentiment) Get(ctx context.Context, marketID string) (*ClientSentiment, error) {
	headers := s.tokens.headers()

	res, err := get[ClientSentiment](ctx, s.Client, "/clientsentiment/"+url.PathEscape(marketID), headers)
	if err != nil {
		return nil, err
	}

	s.tokens.updateTokens(res.httpResponse)

	return res.payload, nil
}
