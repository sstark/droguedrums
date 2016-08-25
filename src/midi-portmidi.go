package main

import (
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
	defaultOut := portmidi.DefaultOutputDeviceID()
	midiOut, err = portmidi.NewOutputStream(defaultOut, 1024, 0)
	checkErr(err)
}

func closeMidi() {
	midiOut.Close()
	portmidi.Terminate()
}
