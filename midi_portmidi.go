// +build linux portmidi
// +build !coremidi

package main

import (
	"errors"
	"fmt"
	"github.com/rakyll/portmidi"
)

var midiOut *portmidi.Stream

func sendMidiNote(channel, note, velocity int) {
	if note != 0 {
		midiOut.WriteShort(int64(midiNoteOn|channel), int64(note), int64(velocity))
	}
}

func initMidi(chosenPort int) error {
	err := portmidi.Initialize()
	midiDevCount := portmidi.CountDevices() // returns the number of MIDI devices
	if midiDevCount == 0 {
		return errors.New("no midi device found")
	}
	defaultOut := portmidi.DefaultOutputDeviceID()
	var midiDeviceInfo *portmidi.DeviceInfo
	fmt.Println("midi outputs found:")
	for i := 0; i < midiDevCount; i++ {
		midiDeviceInfo = portmidi.Info(portmidi.DeviceID(i)) // returns info about a MIDI device
		if midiDeviceInfo.IsOutputAvailable {
			fmt.Printf("%v: ", i)
			fmt.Print("\"", midiDeviceInfo.Interface, "/", midiDeviceInfo.Name, "\"")
			if i == int(defaultOut) {
				fmt.Print(" ", labelPortDefault)
			}
			if i == chosenPort {
				fmt.Print(" ", labelPortSelected)
			}
			fmt.Println()
		}
	}
	var outid portmidi.DeviceID
	if chosenPort != defaultMidiPort {
		outid = portmidi.DeviceID(chosenPort)
	} else {
		outid = defaultOut
	}
	midiOut, err = portmidi.NewOutputStream(outid, 1024, 0)
	return err
}

func closeMidi() {
	midiOut.Close()
	portmidi.Terminate()
}
