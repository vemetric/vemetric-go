package main

import (
	"context"
	"log"

	"github.com/vemetric/vemetric-go"
)

func main() {
	client, err := vemetric.New(
		"o1rySsGlUtFCyflo",
		// Host is optional. If not provided, defaults to "https://hub.vemetric.com"
		vemetric.WithHost("http://localhost:4004"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Track an event
	err = client.TrackEvent(ctx, &vemetric.TrackEventOpts{
		EventName:      "SignupCompleted",
		UserIdentifier: "dmmIrnzUzVMJD03tjCiHXTEEgX6xIPJm",
		UserDisplayName: "TestName",
		EventData: map[string]any{
			"plan": "Pro",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Update user
	err = client.UpdateUser(ctx, &vemetric.UpdateUserOpts{
		UserIdentifier: "dmmIrnzUzVMJD03tjCiHXTEEgX6xIPJm",
		UserData: vemetric.UserData{
			Set: map[string]any{"plan": "BusinessGo"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ events sent")
}
