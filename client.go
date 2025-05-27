package vemetric

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var version = "v0.0.5" // This will be overridden by build flags

type Client struct {
	token  string
	host   string
	hc     *http.Client
	ctx    context.Context
}

// Opts lets callers override host or timeout.
type Opts struct {
	Token   string
	// Host is optional. If not provided, defaults to "https://hub.vemetric.com"
	Host    string
	// Default timeout is 3 seconds.
	Timeout time.Duration
	Context context.Context
}

var (
	ErrBadStatus = errors.New("vemetric: non-2xx status code")
)

// New returns a ready client.
func New(o *Opts) (*Client, error) {
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

	ctx := context.Background()
	if o.Context != nil {
		ctx = o.Context
	}

	return &Client{
		token: o.Token,
		host:  host,
		hc: &http.Client{
			Timeout: timeout,
		},
		ctx: ctx,
	}, nil
}

func (c *Client) TrackEvent(opts *TrackEventOpts) error {
	if opts == nil || opts.EventName == "" {
		return errors.New("vemetric: event name required")
	}

	return c.post("/e", opts)
}

func (c *Client) UpdateUser(opts *UpdateUserOpts) error {
	if opts == nil || opts.UserIdentifier == "" {
		return errors.New("vemetric: user identifier required")
	}

	return c.post("/u", opts)
}

func (c *Client) post(path string, body any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s", c.host, path)
	httpReq, err := http.NewRequestWithContext(c.ctx, http.MethodPost, url, bytes.NewReader(b))
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