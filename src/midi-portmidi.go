package main

import (
	"fmt"
	"github.com/rakyll/portmidi"
)

var midiOut *portmidi.Stream

func sendMidiNote(channel, note, velocity int) {
	if note != 0 {
		midiOut.WriteShort(0x95, int64(note), int64(velocity))
	}
}

func initMidi() {
	err := portmidi.Initialize()
	midiDevCount := portmidi.CountDevices() // returns the number of MIDI devices
	defaultOut := portmidi.DefaultOutputDeviceID()
	var midiDeviceInfo *portmidi.DeviceInfo
	fmt.Println("MIDI outputs found:")
	for i := 0; i < midiDevCount; i++ {
		midiDeviceInfo = portmidi.Info(portmidi.DeviceID(i)) // returns info about a MIDI device
		fmt.Printf("%d: ", i)
		if midiDeviceInfo.IsOutputAvailable {
			fmt.Print(midiDeviceInfo.Interface, "/", midiDeviceInfo.Name)
		}
		if i == int(defaultOut) {
			fmt.Println(" (default)")
		} else {
			fmt.Println()
		}
	}
	midiOut, err = portmidi.NewOutputStream(defaultOut, 1024, 0)
	checkErr(err)
}

func closeMidi() {
	midiOut.Close()
	portmidi.Terminate()
}
