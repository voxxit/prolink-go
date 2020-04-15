package main

import (
	"fmt"
	"time"

	"github.com/voxxit/prolink-go"
	"github.com/voxxit/prolink-go/mixstatus"
)

func main() {
	network, err := prolink.Connect()
	if err != nil {
		panic(err)
	}

	if err := network.AutoConfigure(3 * time.Second); err != nil {
		fmt.Println(err)
	}

	dj := network.CDJStatusMonitor()
	rb := network.RemoteDB()

	config := mixstatus.Config{
		AllowedInterruptBeats: 8,
		BeatsUntilReported:    128,
		TimeBetweenSets:       10 * time.Second,
	}

	handler := func(event mixstatus.Event, status *prolink.CDJStatus) {
		fmt.Printf("Event: %s\n", event)
		fmt.Println(status)

		if status.TrackID != 0 {
			track, err := rb.GetTrack(status.TrackKey())
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(track)
		}

		fmt.Println("---")
	}

	processor := mixstatus.NewProcessor(
		config,
		mixstatus.HandlerFunc(handler),
	)

	dj.AddStatusHandler(processor)

	<-make(chan bool)
}
