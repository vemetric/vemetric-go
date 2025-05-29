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

var version = "0.0.7" // This gets changed automatically during the release process

type client struct {
	token  string
	host   string
	hc     *http.Client
	ctx    context.Context
}

// Configuration options for the Vemetric Client.
type Opts struct {
	// Required. This is the token of your project. You can find it in the Settings page.
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

	ctx := context.Background()
	if o.Context != nil {
		ctx = o.Context
	}

	return &client{
		token: o.Token,
		host:  host,
		hc: &http.Client{
			Timeout: timeout,
		},
		ctx: ctx,
	}, nil
}

// Tracks a custom event for the user with the given identifier.
func (c *client) TrackEvent(opts *TrackEventOpts) error {
	if opts == nil || opts.EventName == "" {
		return errors.New("vemetric: event name required")
	}

	return c.post("/e", opts)
}

// Updates the data of the user with the given identifier.
func (c *client) UpdateUser(opts *UpdateUserOpts) error {
	if opts == nil || opts.UserIdentifier == "" {
		return errors.New("vemetric: user identifier required")
	}

	return c.post("/u", opts)
}

func (c *client) post(path string, body any) error {
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