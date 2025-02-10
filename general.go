package capitalcom

import (
	"context"
	"encoding/json"
	"time"
)

type serverTimeResponsePayload struct {
	ServerTime time.Time `json:"-"`
}

func (p *serverTimeResponsePayload) UnmarshalJSON(data []byte) error {
	var aux struct {
		ServerTime int64 `json:"serverTime"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	p.ServerTime = time.UnixMilli(aux.ServerTime)

	return nil
}

type StatusResponsePayload struct {
	Status string `json:"status"`
}

func (c *Client) Time(ctx context.Context) (time.Time, error) {
	res, err := get[serverTimeResponsePayload](ctx, c, "/time", nil)
	if err != nil {
		return time.Time{}, err
	}

	return res.payload.ServerTime, nil
}

func (c *Client) Ping(ctx context.Context) (string, error) {
	headers := c.tokens.headers()

	res, err := get[StatusResponsePayload](ctx, c, "/ping", headers)
	if err != nil {
		return "", err
	}

	c.tokens.updateTokens(res.httpResponse)

	return res.payload.Status, nil
}
