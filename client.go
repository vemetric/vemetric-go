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

var version = "0.0.9" // This gets changed automatically during the release process

type request struct {
	ctx  context.Context
	path string
	body any
}

type client struct {
	token string
	host  string
	hc    *http.Client
	async bool
	q     chan request
	wg    sync.WaitGroup
}

// Configuration options for the Vemetric Client.
type Opts struct {
	// Required. This is the token of your project. You can find it in the Settings page.
	Token string
	// Host is optional. If not provided, defaults to "https://hub.vemetric.com"
	Host string
	// Default timeout is 3 seconds.
	Timeout time.Duration
	// Default is false.
	Async bool
	// Default is 10. Async buffered channel size
	AsyncBufferedChannelSize int
}

var (
	ErrBadStatus = errors.New("vemetric: non-2xx status code")
)

// New returns a new Vemetric client.
func New(o *Opts) (*client, error) {
	if o == nil || o.Token == "" {
		return nil, errors.New("vemetric: Token required")
	}

	host := "https://hub.vemetric.com"
	if o.Host != "" {
		host = o.Host
	}

	timeout := 3 * time.Second
	if o.Timeout > 0 {
		timeout = o.Timeout
	}

	asyncBufferedChannelSize := o.AsyncBufferedChannelSize

	if o.Async && asyncBufferedChannelSize == 0 {
		asyncBufferedChannelSize = 10
	}

	c := &client{
		token: o.Token,
		host:  host,
		hc: &http.Client{
			Timeout: timeout,
		},
		async: o.Async,
		q:     make(chan request, asyncBufferedChannelSize),
		wg:    sync.WaitGroup{},
	}

	c.wg.Add(1)

	go c.worker()

	return c, nil
}

// Tracks a custom event for the user with the given identifier.
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

// Updates the data of the user with the given identifier.
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
	defer res.Body.Close()

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
