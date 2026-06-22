package capitalcom_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gromson/capitalcom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	refreshedSecurityToken = "RefreshedSecurityToken"
	refreshedCST           = "RefreshedCST"
	priceHistoryJSON       = `{
		"prices": [{
			"snapshotTime": "2026-06-22T12:00:00",
			"snapshotTimeUTC": "2026-06-22T10:00:00",
			"openPrice": {"bid": 100.1, "ask": 100.2},
			"closePrice": {"bid": 101.1, "ask": 101.2},
			"highPrice": {"bid": 102.1, "ask": 102.2},
			"lowPrice": {"bid": 99.1, "ask": 99.2},
			"lastTradedVolume": 42
		}],
		"instrumentType": "CRYPTOCURRENCIES"
	}`
)

func TestPrices_History(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	requests := make([]*http.Request, 0, 3)
	httpClient := newPricesHTTPClient(&requests)

	underTest := capitalcom.NewClient(
		expectedAPIKey,
		identifier,
		password,
		capitalcom.WithHTTPClient(httpClient),
	)
	_, err := underTest.Session().CreateNew(ctx, false)
	require.NoError(t, err)

	location := time.FixedZone("UTC+2", 2*60*60)
	params := capitalcom.PricesParams{
		Resolution: capitalcom.ResolutionMinute,
		Max:        1000,
		From:       time.Date(2026, 6, 22, 12, 0, 0, 0, location),
		To:         time.Date(2026, 6, 22, 12, 15, 0, 0, location),
	}

	// Act
	actual, err := underTest.Prices().History(ctx, "BTCUSD", params)

	// Assert
	require.NoError(t, err)
	assertPriceHistory(t, actual)

	require.Len(t, requests, 2)
	assertPriceRequest(t, requests[1])

	_, err = underTest.Ping(ctx)
	require.NoError(t, err)
	require.Len(t, requests, 3)
	assert.Equal(t, refreshedSecurityToken, requests[2].Header.Get(capitalcom.HeaderKeySecurityToken))
	assert.Equal(t, refreshedCST, requests[2].Header.Get(capitalcom.HeaderTokenCST))
}

func TestPrices_HistoryReturnsTimestampDecodingError(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()

	var decodingErr capitalcom.ResponsePayloadDecodingError

	httpClient := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		switch request.URL.Path {
		case capitalcom.APIPathV1 + "/session":
			return jsonResponse(request, `{}`, tokenHeaders(expectedKeySecurityToken, expectedCST)), nil
		default:
			return jsonResponse(request, `{
				"prices": [{
					"snapshotTime": "invalid",
					"snapshotTimeUTC": "2026-06-22T10:00:00",
					"openPrice": {}, "closePrice": {}, "highPrice": {}, "lowPrice": {}
				}],
				"instrumentType": "CRYPTOCURRENCIES"
			}`, nil), nil
		}
	})}

	underTest := capitalcom.NewClient(
		expectedAPIKey,
		identifier,
		password,
		capitalcom.WithHTTPClient(httpClient),
	)
	_, err := underTest.Session().CreateNew(ctx, false)
	require.NoError(t, err)

	// Act
	_, err = underTest.Prices().History(ctx, "BTCUSD", capitalcom.PricesParams{})

	// Assert
	require.ErrorAs(t, err, &decodingErr)
}

func newPricesHTTPClient(requests *[]*http.Request) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		*requests = append(*requests, request)

		switch request.URL.Path {
		case capitalcom.APIPathV1 + "/session":
			return jsonResponse(request, `{}`, tokenHeaders(expectedKeySecurityToken, expectedCST)), nil
		case capitalcom.APIPathV1 + "/prices/BTCUSD":
			return jsonResponse(request, priceHistoryJSON, tokenHeaders(refreshedSecurityToken, refreshedCST)), nil
		case capitalcom.APIPathV1 + "/ping":
			return jsonResponse(request, `{"status":"OK"}`, nil), nil
		default:
			return jsonResponse(request, `{"errorCode":"error.not-found"}`, nil), nil
		}
	})}
}

func assertPriceHistory(t *testing.T, actual *capitalcom.Prices) {
	t.Helper()

	require.Len(t, actual.Prices, 1)
	assert.Equal(t, "CRYPTOCURRENCIES", actual.InstrumentType)
	assert.Equal(t, time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC), actual.Prices[0].SnapshotTime)
	assert.Equal(t, time.Date(2026, 6, 22, 10, 0, 0, 0, time.UTC), actual.Prices[0].SnapshotTimeUTC)
	assert.Equal(t, capitalcom.PriceData{Bid: 100.1, Ask: 100.2}, actual.Prices[0].OpenPrice)
	assert.Equal(t, capitalcom.PriceData{Bid: 101.1, Ask: 101.2}, actual.Prices[0].ClosePrice)
	assert.Equal(t, capitalcom.PriceData{Bid: 102.1, Ask: 102.2}, actual.Prices[0].HighPrice)
	assert.Equal(t, capitalcom.PriceData{Bid: 99.1, Ask: 99.2}, actual.Prices[0].LowPrice)
	assert.Equal(t, 42, actual.Prices[0].LastTradedVolume)
}

func assertPriceRequest(t *testing.T, request *http.Request) {
	t.Helper()

	assert.Equal(t, capitalcom.APIPathV1+"/prices/BTCUSD", request.URL.Path)
	assert.Equal(t, expectedKeySecurityToken, request.Header.Get(capitalcom.HeaderKeySecurityToken))
	assert.Equal(t, expectedCST, request.Header.Get(capitalcom.HeaderTokenCST))
	assert.Equal(t, "MINUTE", request.URL.Query().Get("resolution"))
	assert.Equal(t, "1000", request.URL.Query().Get("max"))
	assert.Equal(t, "2026-06-22T10:00:00", request.URL.Query().Get("from"))
	assert.Equal(t, "2026-06-22T10:15:00", request.URL.Query().Get("to"))
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func jsonResponse(request *http.Request, body string, headers http.Header) *http.Response {
	if headers == nil {
		headers = make(http.Header)
	}

	headers.Set("Content-Type", "application/json")

	return &http.Response{
		Status:        http.StatusText(http.StatusOK),
		StatusCode:    http.StatusOK,
		Header:        headers,
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Request:       request,
	}
}

func tokenHeaders(securityToken, cst string) http.Header {
	headers := make(http.Header)
	headers.Set(capitalcom.HeaderKeySecurityToken, securityToken)
	headers.Set(capitalcom.HeaderTokenCST, cst)

	return headers
}
