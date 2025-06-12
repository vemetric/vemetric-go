![Vemetric Go SDK](https://github.com/user-attachments/assets/2111bb87-9a4b-4325-8793-df45dede1e6e)

# The Vemetric SDK for Go

Learn more about the Vemetric Go SDK in the [official docs](https://vemetric.com/docs/sdks/go).

You can also checkout the package on the [Go Package Registry](https://pkg.go.dev/github.com/Vemetric/vemetric-go).

[![Go Reference](https://pkg.go.dev/badge/github.com/Vemetric/vemetric-go.svg)](https://pkg.go.dev/github.com/Vemetric/vemetric-go)

## Installation

```bash
go get github.com/Vemetric/vemetric-go
```

## Usage

```go
package main

import (
	"context"
	"log"

	"github.com/Vemetric/vemetric-go"
)

func main() {
	client, err := vemetric.New(&vemetric.Opts{
		Token: "YOUR_PROJECT_TOKEN",
	})
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
client, err := vemetric.New(&vemetric.Opts{
	Token:   "YOUR_PROJECT_TOKEN", // Required
	Host:    "https://hub.vemetric.com", // Optional, defaults to https://hub.vemetric.com
	Timeout: 3 * time.Second, // Optional, defaults to 3 seconds
	Async: false, // Optional, defaults to false, configures if the requests should be fired asynchronously
	AsyncBufferedChannelSize: 10, // Optional, defaults to 10
})
```
