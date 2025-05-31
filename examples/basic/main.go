package main

import (
	"context"
	"log"

	"github.com/Vemetric/vemetric-go"
)

func main() {
	client, err := vemetric.New(&vemetric.Opts{
		Token: "o1rySsGlUtFCyflo",
		Host: "http://localhost:4004", // Host is optional. If not provided, defaults to "https://hub.vemetric.com"
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Track an event
	err = client.TrackEvent(ctx, &vemetric.TrackEventOpts{
		EventName: "SignupCompleted",
		UserIdentifier: "dmmIrnzUzVMJD03tjCiHXTEEgX6xIPJm",
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