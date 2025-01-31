package capitalcom_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gromson/capitalcom"
	"github.com/stretchr/testify/require"
)

func TestClient_Time(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()

	location := time.FixedZone("Test/Timezone", -2*3600)
	serverTime := time.Date(2024, 11, 29, 13, 4, 6, 0, location)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload := fmt.Sprintf(`{ "serverTime": %d }`, serverTime.UnixMilli())

		_, _ = w.Write([]byte(payload))
	}))

	t.Cleanup(func() {
		srv.Close()
	})

	underTest := capitalcom.NewClient(expectedAPIKey,
		identifier,
		password,
		capitalcom.WithHTTPClient(srv.Client()),
		capitalcom.WithHost(srv.URL))

	// Act
	got, err := underTest.Time(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, serverTime, got.In(location))
}

func TestClient_Ping(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case capitalcom.APIPathV1 + "/session":
			handleSessionCreation(w, r)
		case capitalcom.APIPathV1 + "/ping":
			if !isAuthorized(w, r) {
				return
			}

			_, _ = w.Write([]byte(`{ "status": "OK" }`))
		}
	}))

	t.Cleanup(func() {
		srv.Close()
	})

	underTest := capitalcom.NewClient(expectedAPIKey,
		identifier,
		password,
		capitalcom.WithHTTPClient(srv.Client()),
		capitalcom.WithHost(srv.URL))

	// Act
	_, err := underTest.Session().CreateNew(ctx, false)
	require.NoError(t, err)

	got, err := underTest.Ping(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, "OK", got)
}
