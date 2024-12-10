package capitalcom

import (
	"log/slog"
	"net/http"
	"time"
)

const (
	HostDemo = "https://demo-api-capital.backend-capital.com"
	HostLive = "https://api-capital.backend-capital.com"

	APIPathV1 = "/api/v1"

	HeaderAPIKey           = "X-CAP-API-KEY" //nolint:gosec
	HeaderKeySecurityToken = "X-SECURITY-TOKEN"
	HeaderTokenCST         = "CST"
)

const (
	dateFormat = "2006-01-02T15:04:05"
)

type tokens struct {
	securityToken string
	cst           string
}

func (t *tokens) headers() http.Header {
	headers := make(http.Header)
	headers.Set(HeaderKeySecurityToken, t.securityToken) //nolint:canonicalheader
	headers.Set(HeaderTokenCST, t.cst)                   //nolint:canonicalheader

	return headers
}

func (t *tokens) updateTokens(res *http.Response) {
	t.securityToken = res.Header.Get(HeaderKeySecurityToken) //nolint:canonicalheader
	t.cst = res.Header.Get(HeaderTokenCST)                   //nolint:canonicalheader
}

// Client Capital.com API client.
type Client struct {
	apiKey     string
	httpClient *http.Client
	host       string
	apiPath    string
	logger     *slog.Logger

	tokens *tokens
}

// ClientOption is a functional for setting the config option for the client.
type ClientOption func(*Client)

// WithHTTPClient sets the HTTP client for the client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithHost sets the host for the client.
func WithHost(host string) ClientOption {
	return func(c *Client) {
		c.host = host
	}
}

// WithHostProd sets the host for the client to production.
func WithHostProd() ClientOption {
	return WithHost(HostLive)
}

// WithAPIPath sets the API path for the client.
func WithAPIPath(apiPath string) ClientOption {
	return func(c *Client) {
		c.apiPath = apiPath
	}
}

// WithLogger sets the logger for the client.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// NewClient creates a new Capital.com API client.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:     apiKey,
		httpClient: createDefaultHTTPClient(),
		host:       HostDemo,
		apiPath:    APIPathV1,
		logger:     slog.Default(),
		tokens:     &tokens{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Session access to the session.
func (c *Client) Session() *session {
	return &session{Client: c}
}

// Account access to the account resource.
func (c *Client) Account() *account {
	return &account{Client: c}
}

// Trading access to the trading resource.
func (c *Client) Trading() *trading {
	return &trading{Client: c}
}

// Positions access to the positions resource.
func (c *Client) Positions() *positions {
	return &positions{Client: c}
}

// Orders access to the orders resource.
func (c *Client) Orders() *orders {
	return &orders{Client: c}
}

// Markets access to the markets resource.
func (c *Client) Markets() *markets {
	return &markets{Client: c}
}

// Prices access to the prices resource.
func (c *Client) Prices() *prices {
	return &prices{Client: c}
}

// Sentiment access to the sentiment resource.
func (c *Client) Sentiment() *sentiment {
	return &sentiment{Client: c}
}

// Watchlists access to the watchlists resource.
func (c *Client) Watchlists() *watchlists {
	return &watchlists{Client: c}
}

func (c *Client) path(resourcePath string) string {
	return c.host + c.apiPath + resourcePath
}

func createDefaultHTTPClient() *http.Client {
	const defaultClientTimeout = 10 * time.Second

	return &http.Client{
		Timeout: defaultClientTimeout,
	}
}
