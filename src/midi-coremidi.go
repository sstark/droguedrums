// +build coremidi

package main

import (
	"fmt"
	"github.com/youpy/go-coremidi"
)

var midiOut coremidi.OutputPort
var midiDestination coremidi.Destination

func sendMidiNote(channel, note, velocity int) {
	if note != 0 {
		packet := coremidi.NewPacket([]byte{byte(midiNoteOn | channel), byte(note), byte(velocity)})
		_ = packet.Send(&midiOut, &midiDestination)
	}
}

func initMidi(chosenPort int) (err error) {
	allDests, err := coremidi.AllDestinations()
	fmt.Println("midi outputs found:")
	for i, d := range allDests {
		fmt.Printf("%d: \"(%s), %s\"", i, d.Manufacturer(), d.Name())
		if i == 0 {
			fmt.Print(" ", labelPortDefault)
		}
		if i == chosenPort {
			fmt.Print(" ", labelPortSelected)
		}
		fmt.Println()
	}
	if chosenPort != defaultMidiPort {
		if chosenPort >= len(allDests) {
			logger.Fatalf("selected midi port does not exist: %d\n", chosenPort)
		}
		midiDestination = allDests[chosenPort]
	} else {
		midiDestination = allDests[0]
	}
	client, err := coremidi.NewClient("droguedrums-client")
	if err != nil {
		fmt.Println(err)
		return
	}
	midiOut, err = coremidi.NewOutputPort(client, "droguedrums-port")
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func closeMidi() {
}
