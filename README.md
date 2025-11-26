# Capital.com API Client for Go

A Go client library for the [Capital.com](https://capital.com) trading platform REST API. This library provides type-safe access to all Capital.com API endpoints, including session management, trading operations, market data, and account management.

## Features

- **Complete API Coverage**: Full support for all Capital.com REST API endpoints
- **Type-Safe**: Leverages Go generics and strong typing for compile-time safety
- **Automatic Token Management**: Handles authentication tokens transparently
- **Password Encryption**: Optional RSA encryption for secure authentication
- **Demo & Live Support**: Easily switch between demo and live environments
- **Zero External Dependencies**: Uses only standard library (except testify for tests)

## Installation

```bash
go get github.com/gromson/capitalcom
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/gromson/capitalcom"
)

func main() {
    // Create a new client (uses demo environment by default)
    client := capitalcom.NewClient(
        "your-api-key",
        "your-identifier",
        "your-password",
    )

    ctx := context.Background()

    // Create a session
    session, err := client.Session().CreateNew(ctx, false)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Logged in as: %s\n", session.ClientID)
    fmt.Printf("Current account: %s\n", session.CurrentAccountID)

    // Get account list
    accounts, err := client.Account().List(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, account := range accounts {
        fmt.Printf("Account: %s (Balance: %.2f %s)\n",
            account.AccountID,
            account.Balance.Balance,
            account.Currency)
    }
}
```

## Authentication

The client handles authentication automatically using a three-token system:
- **API Key**: Provided during client creation
- **Security Token**: Obtained from session creation
- **CST Token**: Obtained from session creation and updated with each request

### Basic Authentication

```go
client := capitalcom.NewClient(apiKey, identifier, password)

// Create session (tokens are managed automatically)
session, err := client.Session().CreateNew(ctx, false)
```

### Encrypted Password Authentication

For enhanced security, password encryption is handled automatically:

```go
client := capitalcom.NewClient(apiKey, identifier, password)

// Create session with encrypted password (second parameter = true)
// The client will automatically fetch the encryption key and encrypt the password
session, err := client.Session().CreateNew(ctx, true)
if err != nil {
    log.Fatal(err)
}
```

## Configuration

### Use Live Environment

```go
client := capitalcom.NewClient(apiKey, identifier, password,
    capitalcom.WithHostProd(), // Switch to live API
)
```

### Custom HTTP Client

```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}

client := capitalcom.NewClient(apiKey, identifier, password,
    capitalcom.WithHTTPClient(httpClient),
)
```

### Enable Logging

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

client := capitalcom.NewClient(apiKey, identifier, password,
    capitalcom.WithLogger(logger),
)
```

## Usage Examples

### Session Management

```go
// Create new session
session, err := client.Session().CreateNew(ctx, false)

// Get session details
sessionData, err := client.Session().Details(ctx)

// Switch active account
accountStatus, err := client.Session().SwitchActiveAccount(ctx, "ACCOUNT_ID")

// Logout
status, err := client.Session().LogOut(ctx)
```

### Account Operations

```go
// List all accounts
accounts, err := client.Account().List(ctx)

// Get account preferences
prefs, err := client.Account().Preferences(ctx)

// Update leverage preferences
status, err := client.Account().UpdatePreferences(ctx,
    &capitalcom.UpdateLeverages{
        Cryptocurrencies: 10,
        Shares:          5,
    },
    false, // hedgingMode
)

// Get activity history
activities, err := client.Account().ActivityHistory(ctx, capitalcom.ActivityParams{
    From:     time.Now().AddDate(0, -1, 0), // Last month
    To:       time.Now(),
    Detailed: true,
})

// Get transaction history
transactions, err := client.Account().TransactionHistory(ctx, capitalcom.TransactionParams{
    Type: capitalcom.TransactionTypeTrade,
    From: time.Now().AddDate(0, 0, -7), // Last week
})
```

### Trading - Positions

```go
// Open a new position
dealRef, err := client.Positions().Open(ctx, capitalcom.OpenPositionRequest{
    Direction: capitalcom.PositionDirectionBuy,
    Epic:      "BTCUSD",
    Size:      0.01,
    UpdatePositionRequest: capitalcom.UpdatePositionRequest{
        GuaranteedStop: false,
        StopLevel:      95000.00,  // Stop loss
        ProfitLevel:    105000.00, // Take profit
    },
})

// List all open positions
positions, err := client.Positions().List(ctx)

for _, pos := range positions {
    fmt.Printf("Position: %s %s %.4f @ %.2f (P/L: %.2f)\n",
        pos.Position.Direction,
        pos.Market.Epic,
        pos.Position.Size,
        pos.Position.Level,
        pos.Position.UPL)
}

// Get specific position
position, err := client.Positions().Get(ctx, "DEAL_ID")

// Update position (modify stop/profit levels)
dealRef, err = client.Positions().Update(ctx, "DEAL_ID", capitalcom.UpdatePositionRequest{
    StopLevel:   96000.00,
    ProfitLevel: 110000.00,
})

// Close a position
err = client.Positions().Close(ctx, "DEAL_ID")
```

### Trading - Working Orders

```go
// Create a limit order
dealRef, err := client.Orders().Create(ctx, capitalcom.CreateOrderRequest{
    Direction: capitalcom.PositionDirectionBuy,
    Epic:      "ETHUSD",
    Size:      0.1,
    Type:      capitalcom.LimitOrder,
    UpdateOrderRequest: capitalcom.UpdateOrderRequest{
        Level:          3500.00,
        GoodTillDate:   time.Now().AddDate(0, 0, 7), // Valid for 7 days
        GuaranteedStop: false,
        StopLevel:      3400.00,
        ProfitLevel:    3700.00,
    },
})

// Create a stop order
dealRef, err = client.Orders().Create(ctx, capitalcom.CreateOrderRequest{
    Direction: capitalcom.PositionDirectionSell,
    Epic:      "GBPUSD",
    Size:      1000,
    Type:      capitalcom.StopOrder,
    UpdateOrderRequest: capitalcom.UpdateOrderRequest{
        Level:       1.2500,
        StopLevel:   1.2550,
        ProfitLevel: 1.2400,
    },
})

// List all working orders
orders, err := client.Orders().List(ctx)

// Update an order
dealRef, err = client.Orders().Update(ctx, "DEAL_ID", capitalcom.UpdateOrderRequest{
    Level:       3550.00,
    StopLevel:   3450.00,
    ProfitLevel: 3750.00,
})

// Delete an order
dealRef, err = client.Orders().Delete(ctx, "DEAL_ID")
```

### Trade Confirmations

```go
// Get deal confirmation
deal, err := client.Trading().Confirm(ctx, "DEAL_REFERENCE")

fmt.Printf("Deal Status: %s\n", deal.Status)
fmt.Printf("Deal ID: %s\n", deal.DealID)
fmt.Printf("Level: %.2f\n", deal.Level)
```

### Market Data

```go
// Browse market categories
categories, err := client.Markets().Categories(ctx)

for _, cat := range categories {
    fmt.Printf("Category: %s (ID: %s)\n", cat.Name, cat.ID)
}

// Browse subcategories within a category
subcategories, err := client.Markets().Subcategories(ctx, "195969", 100)

// Search markets by epic or search term
markets, err := client.Markets().Details(ctx, capitalcom.DetailsParams{
    SearchTerm: "bitcoin",
})

// Get multiple markets by EPIC codes
markets, err = client.Markets().Details(ctx, capitalcom.DetailsParams{
    Epics: []string{"BTCUSD", "ETHUSD"},
})

for _, market := range markets {
    fmt.Printf("Market: %s\n", market.InstrumentName)
    fmt.Printf("Bid: %.2f, Offer: %.2f\n", market.Bid, market.Offer)
}

// Get detailed information for a single market
detail, err := client.Markets().Detail(ctx, "BTCUSD")

fmt.Printf("Market: %s\n", detail.Instrument.Name)
fmt.Printf("Bid: %.2f, Offer: %.2f\n",
    detail.Snapshot.Bid,
    detail.Snapshot.Offer)
fmt.Printf("Min Trade Size: %.4f\n",
    detail.DealingRules.MinDealSize.Value)
```

### Price History

```go
// Get historical prices
prices, err := client.Prices().History(ctx, "BTCUSD", capitalcom.PricesParams{
    Resolution: capitalcom.ResolutionHour,
    Max:        100,
})

for _, price := range prices.Prices {
    fmt.Printf("%s: O=%.2f H=%.2f L=%.2f C=%.2f\n",
        price.SnapshotTime.Format(time.RFC3339),
        price.OpenPrice.Bid,
        price.HighPrice.Bid,
        price.LowPrice.Bid,
        price.ClosePrice.Bid)
}

// Get prices for date range
prices, err = client.Prices().History(ctx, "ETHUSD", capitalcom.PricesParams{
    Resolution: capitalcom.ResolutionDay,
    From:       time.Now().AddDate(0, -1, 0), // Last month
    To:         time.Now(),
})
```

### Client Sentiment

```go
// Get sentiment for specific market
sentiment, err := client.Sentiment().Get(ctx, "BTCUSD")

fmt.Printf("Long: %.1f%%, Short: %.1f%%\n",
    sentiment.LongPositionPercentage,
    sentiment.ShortPositionPercentage)

// Get sentiment for multiple markets
sentiments, err := client.Sentiment().List(ctx, []string{"BTCUSD", "ETHUSD", "GBPUSD"})
```

### Watchlists

```go
// Create a watchlist
resp, err := client.Watchlists().Create(ctx, capitalcom.CreateWatchlistRequest{
    Name:  "Crypto Portfolio",
    Epics: []string{"BTCUSD", "ETHUSD"},
})
watchlistID := resp.WatchlistID

// Add more markets to watchlist
status, err := client.Watchlists().AddMarket(ctx, watchlistID, "SOLUSD")

// Get markets in watchlist
markets, err := client.Watchlists().Get(ctx, watchlistID)

for _, market := range markets {
    fmt.Printf("Market: %s\n", market.Epic)
}

// Remove market from watchlist
status, err = client.Watchlists().RemoveMarket(ctx, watchlistID, "ETHUSD")

// Delete watchlist
status, err = client.Watchlists().Delete(ctx, watchlistID)
```

### Utility Functions

```go
// Check server connection
status, err := client.Ping(ctx)
fmt.Printf("API Status: %s\n", status.Status)

// Get server time
serverTime, err := client.Time(ctx)
fmt.Printf("Server Time: %s\n", serverTime.Format(time.RFC3339))
```

## Error Handling

The library provides specific error types for different failure scenarios:

```go
import "errors"

position, err := client.Positions().Get(ctx, "DEAL_ID")
if err != nil {
    // Check for API errors (4xx/5xx responses)
    var apiErr capitalcom.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error: %s (Status: %d, Code: %s)\n",
            apiErr.Error(),
            apiErr.StatusCode(),
            apiErr.ErrorCode())
        return
    }

    // Check for encoding errors
    var encErr capitalcom.RequestPayloadEncodingError
    if errors.As(err, &encErr) {
        fmt.Printf("Failed to encode request: %v\n", encErr)
        return
    }

    // Check for decoding errors
    var decErr capitalcom.ResponsePayloadDecodingError
    if errors.As(err, &decErr) {
        fmt.Printf("Failed to decode response: %v\n", decErr)
        return
    }

    // Check for HTTP request errors
    var httpErr capitalcom.HTTPRequestError
    if errors.As(err, &httpErr) {
        fmt.Printf("HTTP request failed: %v\n", httpErr)
        return
    }

    // Generic error
    fmt.Printf("Error: %v\n", err)
}
```

## Disclaimer

This is an unofficial library and is not affiliated with or endorsed by Capital.com. Use at your own risk. Trading involves substantial risk of loss and is not suitable for all investors.

## Resources

- [Capital.com API Documentation](https://capital.com/api-development-guide)
- [Capital.com Website](https://capital.com)
