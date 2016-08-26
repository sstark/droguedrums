package main

import (
	"fmt"
	"github.com/rakyll/portmidi"
)

var midiOut *portmidi.Stream

func sendMidiNote(channel, note, velocity int) {
	if note != 0 {
		midiOut.WriteShort(midiNoteOn+int64(channel), int64(note), int64(velocity))
	}
}

func initMidi(chosenPort int) error {
	err := portmidi.Initialize()
	midiDevCount := portmidi.CountDevices() // returns the number of MIDI devices
	defaultOut := portmidi.DefaultOutputDeviceID()
	var midiDeviceInfo *portmidi.DeviceInfo
	fmt.Println("MIDI outputs found:")
	for i := 0; i < midiDevCount; i++ {
		midiDeviceInfo = portmidi.Info(portmidi.DeviceID(i)) // returns info about a MIDI device
		if midiDeviceInfo.IsOutputAvailable {
			fmt.Printf("%v: ", i)
			fmt.Print(midiDeviceInfo.Interface, "/", midiDeviceInfo.Name)
			if i == int(defaultOut) {
				fmt.Print(" <default>")
			}
			if i == chosenPort {
				fmt.Print(" <selected>")
			}
			fmt.Println()
		}
	}
	var outid portmidi.DeviceID
	if chosenPort != -1 {
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
