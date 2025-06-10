package vemetric

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    *Opts
		wantErr bool
	}{
		{
			name: "valid options",
			opts: &Opts{
				Token: "test-token",
			},
			wantErr: false,
		},
		{
			name:    "nil options",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "empty token",
			opts: &Opts{
				Token: "",
			},
			wantErr: true,
		},
		{
			name: "custom host",
			opts: &Opts{
				Token: "test-token",
				Host:  "https://custom.host",
			},
			wantErr: false,
		},
		{
			name: "custom timeout",
			opts: &Opts{
				Token:   "test-token",
				Timeout: 5 * time.Second,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("New() returned nil client when no error expected")
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

	client, err := New(&Opts{
		Token: "test-token",
		Host:  server.URL,
	})
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
				EventData: map[string]any{
					"key": "value",
				},
			},
			wantErr: false,
		},
		{
			name: "nil options",
			opts: nil,
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

	client, err := New(&Opts{
		Token: "test-token",
		Host:  server.URL,
	})
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

	client, err := New(&Opts{
		Token:   "test-token",
		Host:    server.URL,
	})
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
