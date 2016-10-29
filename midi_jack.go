// +build linux jack
// +build !coremidi
// +build !portmidi

package main

import (
	"errors"
	"fmt"
	"github.com/xthexder/go-jack"
	_ "time"
)

const (
	jackClientName = "droguedrums"
	jackPortName = "output"
)

var midiOut *jack.Port
var jackClient jack.Client

func sendMidiNote(channel, note, velocity int) {
	 if note != 0 {
		ev := &jack.MidiData{
			Time: 0,
			Buffer: []byte{
				byte(midiNoteOn|channel),
				byte(note),
				byte(velocity),
			},
		}
		buf := &[]byte{}
		midiOut.MidiEventWrite(ev, buf)
        }
}

func initMidi(chosenPort int) error {
	var err error
	jackClient, _ := jack.ClientOpen(jackClientName, jack.NoStartServer)
	if jackClient == nil {
		fmt.Println("Could not connect to jack server")
		return err
	}
	jackClient.Activate()
	midiOut = jackClient.PortRegister(jackPortName, jack.DEFAULT_MIDI_TYPE, jack.PortIsOutput, 0)
	fmt.Println(midiOut)
	ports := jackClient.GetPorts("", "", jack.PortIsInput)
	fs := jackClient.GetPortByName("fluidsynth:midi")
	fmt.Println(fs)
	ret := jackClient.ConnectPorts(midiOut, fs)
	if ret != 0 {
		err = errors.New("could not connect ports")
	}
	fmt.Println(ports)
	//time.Sleep(time.Second * 10)
	return err
}

func closeMidi() {
	jackClient.Close()
}
