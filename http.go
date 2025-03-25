package capitalcom

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
)

type response[TResPayload any] struct {
	httpResponse *http.Response
	payload      *TResPayload
}

type errorResponsePayload struct {
	ErrorCode string `json:"errorCode"`
}

func get[TResPayload any](
	ctx context.Context,
	c *Client,
	resourcePath string,
	headers http.Header,
) (*response[TResPayload], error) {
	return doRequest[TResPayload](ctx, c, http.MethodGet, resourcePath, nil, headers)
}

func post[TResPayload any](
	ctx context.Context,
	c *Client,
	resourcePath string,
	reqPayload any,
	headers http.Header,
) (*response[TResPayload], error) {
	reqBody, err := prepareRequestBody(reqPayload)
	if err != nil {
		return nil, err
	}

	return doRequest[TResPayload](ctx, c, http.MethodPost, resourcePath, reqBody, headers)
}

func put[TResPayload any](
	ctx context.Context,
	c *Client,
	resourcePath string,
	reqPayload any,
	headers http.Header,
) (*response[TResPayload], error) {
	reqBody, err := prepareRequestBody(reqPayload)
	if err != nil {
		return nil, err
	}

	return doRequest[TResPayload](ctx, c, http.MethodPut, resourcePath, reqBody, headers)
}

func del[TResPayload any](
	ctx context.Context,
	c *Client,
	resourcePath string,
	headers http.Header,
) (*response[TResPayload], error) {
	return doRequest[TResPayload](ctx, c, http.MethodDelete, resourcePath, nil, headers)
}

func prepareRequestBody(payload any) (*bytes.Buffer, error) {
	reqBody := new(bytes.Buffer)

	if err := json.NewEncoder(reqBody).Encode(payload); err != nil {
		return nil, NewRequestPayloadEncodingError(err)
	}

	return reqBody, nil
}

func doRequest[TResPayload any](
	ctx context.Context,
	c *Client,
	method string,
	resourcePath string,
	reqBody io.Reader,
	headers http.Header,
) (*response[TResPayload], error) {
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.path(resourcePath),
		reqBody)
	if err != nil {
		return nil, NewRequestCreationError(err)
	}

	setRequestHeaders(req, headers)

	logRequest(ctx, req, c.logger)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewHTTPRequestError(err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.With("error", err).
				Error("failed to close a response body")
		}
	}()

	logResponse(ctx, res, c.logger)

	if res.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(res)
	}

	resPayload := new(TResPayload)

	if err := json.NewDecoder(res.Body).Decode(resPayload); err != nil {
		return nil, NewResponsePayloadDecodingError(err)
	}

	return &response[TResPayload]{
		httpResponse: res,
		payload:      resPayload,
	}, nil
}

func setRequestHeaders(req *http.Request, headers http.Header) {
	req.Header.Set("Content-Type", "application/json")

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

func handleErrorResponse(res *http.Response) error {
	if res.StatusCode < http.StatusBadRequest || res.StatusCode >= http.StatusInternalServerError {
		return NewAPIError(res.StatusCode, "")
	}

	errPayload := &errorResponsePayload{}

	if err := json.NewDecoder(res.Body).Decode(errPayload); err != nil {
		return NewResponsePayloadDecodingError(err)
	}

	return NewAPIError(res.StatusCode, errPayload.ErrorCode)
}

func logRequest(ctx context.Context, req *http.Request, logger *slog.Logger) {
	if logger == nil || !logger.Enabled(ctx, slog.LevelDebug) {
		return
	}

	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		logger.With("error", err).Error("failed to dump a request")

		return
	}

	logger.With("request", dump).Debug("sending a request")
}

func logResponse(ctx context.Context, res *http.Response, logger *slog.Logger) {
	if logger == nil || !logger.Enabled(ctx, slog.LevelDebug) {
		return
	}

	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		logger.With("error", err).Error("failed to dump a response")

		return
	}

	logger.With("response", dump).Debug("received a response")
}
