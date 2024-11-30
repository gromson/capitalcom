package capitalcom

import (
	"context"
	"net/url"
)

type watchlists struct {
	*Client
}

type (
	watchlistsResponsePayload struct {
		Watchlists []Watchlist `json:"watchlists"`
	}

	Watchlist struct {
		ID                     string `json:"id"`
		Name                   string `json:"name"`
		Editable               bool   `json:"editable"`
		Deleteable             bool   `json:"deleteable"`
		DefaultSystemWatchlist bool   `json:"defaultSystemWatchlist"`
	}
)

func (w *watchlists) List(ctx context.Context) ([]Watchlist, error) {
	headers := w.tokens.headers()

	res, err := get[watchlistsResponsePayload](ctx, w.Client, "/watchlists", headers)
	if err != nil {
		return nil, err
	}

	w.tokens.updateTokens(res.httpResponse)

	return res.payload.Watchlists, nil
}

type (
	CreateWatchlistRequest struct {
		Epics []string `json:"epics"`
		Name  string   `json:"name"`
	}

	WatchlistResponse struct {
		WatchlistID string `json:"watchlistId"`
		Status      string `json:"status"`
	}
)

func (w *watchlists) Create(ctx context.Context, req CreateWatchlistRequest) (*WatchlistResponse, error) {
	headers := w.tokens.headers()

	res, err := post[WatchlistResponse](ctx, w.Client, "/watchlists", req, headers)
	if err != nil {
		return nil, err
	}

	w.tokens.updateTokens(res.httpResponse)

	return res.payload, nil
}

func (w *watchlists) Get(ctx context.Context, watchlistID string) ([]Market, error) {
	headers := w.tokens.headers()

	res, err := get[marketsResponsePayload](ctx, w.Client, "/watchlists/"+url.PathEscape(watchlistID), headers)
	if err != nil {
		return nil, err
	}

	w.tokens.updateTokens(res.httpResponse)

	return res.payload.Markets, nil
}

func (w *watchlists) AddMarket(ctx context.Context, watchlistID, epic string) (string, error) {
	headers := w.tokens.headers()

	reqPayload := struct {
		Epic string `json:"epic"`
	}{
		Epic: epic,
	}

	res, err := put[StatusResponsePayload](ctx,
		w.Client,
		"/watchlists/"+url.PathEscape(watchlistID),
		reqPayload,
		headers)
	if err != nil {
		return "", err
	}

	w.tokens.updateTokens(res.httpResponse)

	return res.payload.Status, nil
}

func (w *watchlists) Delete(ctx context.Context, watchlistID string) (string, error) {
	headers := w.tokens.headers()

	res, err := delete[StatusResponsePayload](ctx,
		w.Client,
		"/watchlists/"+url.PathEscape(watchlistID),
		headers)
	if err != nil {
		return "", err
	}

	w.tokens.updateTokens(res.httpResponse)

	return res.payload.Status, nil
}

func (w *watchlists) RemoveMarket(ctx context.Context, watchlistID, epic string) (string, error) {
	headers := w.tokens.headers()

	res, err := delete[StatusResponsePayload](ctx,
		w.Client,
		"/watchlists/"+url.PathEscape(watchlistID)+"/"+url.PathEscape(epic),
		headers)
	if err != nil {
		return "", err
	}

	w.tokens.updateTokens(res.httpResponse)

	return res.payload.Status, nil
}
