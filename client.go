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
	headers.Set(HeaderKeySecurityToken, t.securityToken)
	headers.Set(HeaderTokenCST, t.cst)

	return headers
}

func (t *tokens) updateTokens(res *http.Response) {
	t.securityToken = res.Header.Get(HeaderKeySecurityToken)
	t.cst = res.Header.Get(HeaderTokenCST)
}

type Client struct {
	apiKey     string
	httpClient *http.Client
	host       string
	apiPath    string
	logger     *slog.Logger

	tokens *tokens
}

type Option func(*Client)

func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithHost(host string) Option {
	return func(c *Client) {
		c.host = host
	}
}

func WithHostProd() Option {
	return WithHost(HostLive)
}

func WithAPIPath(apiPath string) Option {
	return func(c *Client) {
		c.apiPath = apiPath
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		httpClient: createDefaultHttpClient(),
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

func (c *Client) Session() *session {
	return &session{Client: c}
}

func (c *Client) Account() *account {
	return &account{Client: c}
}

func (c *Client) Trading() *trading {
	return &trading{Client: c}
}

func (c *Client) path(resourcePath string) string {
	return c.host + c.apiPath + resourcePath
}

func createDefaultHttpClient() *http.Client {
	const defaultClientTimeout = 10 * time.Second

	return &http.Client{
		Timeout: defaultClientTimeout,
	}
}
