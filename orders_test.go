package capitalcom_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gromson/capitalcom"
	"github.com/stretchr/testify/require"
)

func Test_CreateOrderRequestMarshalSuccessfully(t *testing.T) {
	t.Parallel()

	// Arrange
	goodTillDate := time.Date(2025, 3, 5, 12, 23, 43, 0, time.UTC)

	underTest := capitalcom.CreateOrderRequest{
		Direction: capitalcom.PositionDirectionBuy,
		Epic:      "BTCUSD",
		Size:      0.0034,
		Type:      capitalcom.StopOrder,
		UpdateOrderRequest: capitalcom.UpdateOrderRequest{
			Level:          103263.45,
			GoodTillDate:   goodTillDate,
			GuaranteedStop: true,
			TrailingStop:   false,
			StopLevel:      102456.65,
			StopDistance:   0,
			StopAmount:     0,
			ProfitLevel:    0,
			ProfitDistance: 0,
			ProfitAmount:   0,
		},
	}

	expectedJSON := `{
"direction": "BUY",
"epic": "BTCUSD",
"size": 0.0034,
"type": "STOP",
"level": 103263.45,
"goodTillDate": "2025-03-05T12:23:43",
"guaranteedStop": true,
"stopLevel": 102456.65
}`

	// Act
	actualJSON, err := json.Marshal(underTest)

	// Assert
	require.NoError(t, err)
	require.JSONEq(t, expectedJSON, string(actualJSON))
}
