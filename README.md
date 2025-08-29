![Vemetric Go SDK](https://github.com/user-attachments/assets/dfaf4210-c4ea-433a-84ac-48d5807034cc)


# The Vemetric SDK for Go

Learn more about the Vemetric Go SDK in the [official docs](https://vemetric.com/docs/sdks/go).

You can also checkout the package on the [Go Package Registry](https://pkg.go.dev/github.com/vemetric/vemetric-go).

[![Go Reference](https://pkg.go.dev/badge/github.com/vemetric/vemetric-go.svg)](https://pkg.go.dev/github.com/vemetric/vemetric-go)

## Installation

```bash
go get github.com/vemetric/vemetric-go
```

## Usage

```go
package main

import (
	"context"
	"log"

	"github.com/vemetric/vemetric-go"
)

func main() {
	client, err := vemetric.New("YOUR_PROJECT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Track an event
	err = client.TrackEvent(ctx, &vemetric.TrackEventOpts{
		EventName: "SignupCompleted",
		UserIdentifier: "user-id",
		EventData: map[string]any{
			"key": "value",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Update the data of a user
	err = client.UpdateUser(ctx, &vemetric.UpdateUserOpts{
		UserIdentifier: "user-id",
		UserData: vemetric.UserData{
			Set: map[string]any{"key1": "value1"},
			SetOnce: map[string]any{"key2": "value2"},
			Unset: []string{"key3"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

## Configuration

The client can be configured with the following options:

```go
client, err := vemetric.New(
    "YOUR_PROJECT_TOKEN",  // Required
    // Optional, defaults to https://hub.vemetric.com
    vemetric.WithHost("https://custom.host"),
    // Optional, defaults to a http.Client with a 3 second timeout
    vemetric.WithHTTPClient(&http.Client{Timeout: 3 * time.Second}),
    // Optional, defaults to false, configures if the requests should be fired asynchronously
    vemetric.UseAsync(),
    // Optional, defaults to 10
    vemetric.WithAsyncBufferedChannelSize(10),
})
```

The client can be combined with [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp) to automatically retry
failed requests.

```go
retryClient := retryablehttp.NewClient()

// up to max 3 retries, with exponential backoff
retryClient.RetryMax = 3

// default retry policy:
// retry requests in case of network issues, SSL certificate errors, 
// 429 Too Many Requests, or any of the 500-range response errors
retryClient.CheckRetry = retryablehttp.DefaultRetryPolicy

client, err := vemetric.New(
    "YOUR_PROJECT_TOKEN",
    vemetric.WithHTTPClient(retryClient.StandardClient()),
)
```
