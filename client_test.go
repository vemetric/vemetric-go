package vemetric

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		token string
		opts  []ClientOption
	}

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantToken   string
		wantHost    string
		wantTimeout time.Duration
		wantAsync   bool
	}{
		{
			name: "valid default options",
			args: args{
				token: "test-token",
			},
			wantErr:     false,
			wantToken:   "test-token",
			wantHost:    "https://hub.vemetric.com",
			wantTimeout: 3 * time.Second,
			wantAsync:   false,
		},
		{
			name: "empty token",
			args: args{
				token: "",
			},
			wantErr: true,
		},
		{
			name: "custom host",
			args: args{
				token: "test-token",
				opts:  []ClientOption{WithHost("https://custom.host")},
			},
			wantErr:     false,
			wantToken:   "test-token",
			wantHost:    "https://custom.host",
			wantTimeout: 3 * time.Second,
			wantAsync:   false,
		},
		{
			name: "custom timeout",
			args: args{
				token: "test-token",
				opts:  []ClientOption{WithHTTPClient(&http.Client{Timeout: 5 * time.Second})},
			},
			wantErr:     false,
			wantToken:   "test-token",
			wantHost:    "https://hub.vemetric.com",
			wantTimeout: 5 * time.Second,
			wantAsync:   false,
		},
		{
			name: "async",
			args: args{
				token: "test-token",
				opts:  []ClientOption{UseAsync()},
			},
			wantErr:     false,
			wantToken:   "test-token",
			wantHost:    "https://hub.vemetric.com",
			wantTimeout: 3 * time.Second,
			wantAsync:   true,
		},
		{
			name: "async with zero channel size",
			args: args{
				token: "test-token",
				opts:  []ClientOption{UseAsync(), WithAsyncBufferedChannelSize(0)},
			},
			wantErr: true,
		},
		{
			name: "combine multiple options",
			args: args{
				token: "test-token",
				opts: []ClientOption{
					WithHost("https://custom.host"),
					WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
					UseAsync(),
					WithAsyncBufferedChannelSize(5),
				},
			},
			wantErr:     false,
			wantToken:   "test-token",
			wantHost:    "https://custom.host",
			wantTimeout: 5 * time.Second,
			wantAsync:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vemetricClient, err := New(tt.args.token, tt.args.opts...)
			gotErr := err != nil

			if tt.wantErr {
				if !gotErr {
					t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if gotErr {
				t.Error("Expected no error, got", err)
				return
			}

			actualClient := vemetricClient.(*client)

			if actualClient.token != tt.wantToken {
				t.Errorf("New() token = %v, want %v", actualClient.token, tt.args.token)
			}
			if actualClient.host != tt.wantHost {
				t.Errorf("New() host = %v, want %v", actualClient.host, tt.wantHost)
			}
			if actualClient.hc.(*http.Client).Timeout != tt.wantTimeout {
				t.Errorf("New() timeout = %v, want %v", actualClient.hc.(*http.Client).Timeout, tt.wantTimeout)
			}
			if actualClient.async != tt.wantAsync {
				t.Errorf("New() async = %v, want %v", actualClient.async, tt.wantAsync)
			}
		})
	}
}

func TestTrackEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Token") != "test-token" {
			t.Errorf("Expected Token header to be 'test-token', got %s", r.Header.Get("Token"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be 'application/json', got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("V-SDK") != "go" {
			t.Errorf("Expected V-SDK header to be 'go', got %s", r.Header.Get("V-SDK"))
		}

		// Verify request path
		if r.URL.Path != "/e" {
			t.Errorf("Expected request path to be '/e', got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := New("test-token", WithHost(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		opts    *TrackEventOpts
		wantErr bool
	}{
		{
			name: "valid event",
			opts: &TrackEventOpts{
				EventName:      "test-event",
				UserIdentifier: "user123",
				UserDisplayName: "John Doe",
				EventData: map[string]any{
					"key": "value",
				},
			},
			wantErr: false,
		},
		{
			name:    "nil options",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "empty event name",
			opts: &TrackEventOpts{
				EventName: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := client.TrackEvent(ctx, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrackEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Token") != "test-token" {
			t.Errorf("Expected Token header to be 'test-token', got %s", r.Header.Get("Token"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be 'application/json', got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("V-SDK") != "go" {
			t.Errorf("Expected V-SDK header to be 'go', got %s", r.Header.Get("V-SDK"))
		}

		// Verify request path
		if r.URL.Path != "/u" {
			t.Errorf("Expected request path to be '/u', got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := New("test-token", WithHost(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		opts    *UpdateUserOpts
		wantErr bool
	}{
		{
			name: "valid update",
			opts: &UpdateUserOpts{
				UserIdentifier: "user123",
				UserData: UserData{
					Set: map[string]any{
						"plan": "pro",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "nil options",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "empty user identifier",
			opts: &UpdateUserOpts{
				UserIdentifier: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := client.UpdateUser(ctx, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := New("test-token", WithHost(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Cancel the context immediately
	cancel()

	// Try to track an event
	err = client.TrackEvent(ctx, &TrackEventOpts{
		EventName: "test-event",
	})
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}
}
