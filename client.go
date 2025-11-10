package vemetric

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var version = "0.2.0" // This gets changed automatically during the release process

type request struct {
	ctx  context.Context
	path string
	body any
}

type Client interface {
	TrackEvent(ctx context.Context, opts *TrackEventOpts) error
	UpdateUser(ctx context.Context, opts *UpdateUserOpts) error
	Close()
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type client struct {
	token string
	host  string
	hc    HTTPClient
	async bool
	q     chan request
	wg    sync.WaitGroup
}

// ClientOptions holds the configuration options for the Vemetric Client.
type ClientOptions struct {
	// Host is optional. If not provided, defaults to "https://hub.vemetric.com"
	Host string
	// The HTTP client to use for requests. If not provided, it defaults to a http.Client with a 3 second timeout.
	HttpClient HTTPClient
	// Default is false.
	Async bool
	// Default is 10. Async buffered channel size
	AsyncBufferedChannelSize uint
}

type ClientOption func(*ClientOptions)

var (
	ErrBadStatus = errors.New("vemetric: non-2xx status code")
)

// WithHost sets the host for the Vemetric Client. Defaults to "https://hub.vemetric.com"
func WithHost(host string) ClientOption {
	return func(opts *ClientOptions) {
		opts.Host = host
	}
}

// WithHTTPClient sets the HTTP client to use for sending requests to Vemetric. Defaults to a http.Client with a
// 3 second timeout.
func WithHTTPClient(hc HTTPClient) ClientOption {
	return func(opts *ClientOptions) {
		opts.HttpClient = hc
	}
}

// UseAsync sets the client to use asynchronous requests. Defaults to false.
func UseAsync() ClientOption {
	return func(opts *ClientOptions) {
		opts.Async = true
	}
}

// WithAsyncBufferedChannelSize sets the size of the buffered channel for asynchronous requests. Defaults to 10.
// Only applies if UseAsync() is used.
func WithAsyncBufferedChannelSize(size uint) ClientOption {
	return func(opts *ClientOptions) {
		opts.AsyncBufferedChannelSize = size
	}
}

// New returns a new Vemetric client.
func New(token string, opts ...ClientOption) (Client, error) {
	if token == "" {
		return nil, errors.New("vemetric: Token required")
	}

	httpClient := &http.Client{Timeout: 3 * time.Second}

	// initialize options with defaults
	options := &ClientOptions{
		Host:                     "https://hub.vemetric.com",
		HttpClient:               httpClient,
		Async:                    false,
		AsyncBufferedChannelSize: 10,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Async && options.AsyncBufferedChannelSize == 0 {
		return nil, errors.New("vemetric: AsyncBufferedChannelSize must be greater than 0")
	}

	c := &client{
		token: token,
		host:  options.Host,
		hc:    options.HttpClient,
		async: options.Async,
		q:     make(chan request, options.AsyncBufferedChannelSize),
		wg:    sync.WaitGroup{},
	}

	c.wg.Add(1)

	go c.worker()

	return c, nil
}

// TrackEvent tracks a custom event for the user with the given identifier.
func (c *client) TrackEvent(ctx context.Context, opts *TrackEventOpts) error {
	if opts == nil || opts.EventName == "" {
		return errors.New("vemetric: event name required")
	}

	if !c.async {
		return c.post(ctx, "/e", opts)
	}

	c.q <- request{ctx: ctx, path: "/e", body: opts}

	return nil
}

// UpdateUser updates the data of the user with the given identifier.
func (c *client) UpdateUser(ctx context.Context, opts *UpdateUserOpts) error {
	if opts == nil || opts.UserIdentifier == "" {
		return errors.New("vemetric: user identifier required")
	}

	if !c.async {
		return c.post(ctx, "/u", opts)
	}

	c.q <- request{ctx: ctx, path: "/u", body: opts}

	return nil
}

func (c *client) Close() {
	close(c.q)

	c.wg.Wait()
}

func (c *client) post(ctx context.Context, path string, body any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s", c.host, path)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Token", c.token)
	httpReq.Header.Set("User-Agent", "vemetric-go/"+version)
	httpReq.Header.Set("V-SDK", "go")
	httpReq.Header.Set("V-SDK-Version", version)

	res, err := c.hc.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode >= 300 {
		return fmt.Errorf("%w: %s", ErrBadStatus, res.Status)
	}
	return nil
}

func (c *client) worker() {
	defer c.wg.Done()

	for request := range c.q {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		_ = c.post(ctx, request.path, request.body)

		cancel()
	}
}
