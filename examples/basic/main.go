package main

import (
	"log"

	"github.com/Vemetric/vemetric-go"
)

func main() {
	client, err := vemetric.New(&vemetric.Opts{
		Token: "WRlW37cPSLUAbXDk76wYU",
		Host: "http://localhost:4004",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Track an event
	err = client.TrackEvent(&vemetric.TrackEventOpts{
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
	err = client.UpdateUser(&vemetric.UpdateUserOpts{
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