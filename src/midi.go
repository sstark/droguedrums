package main

import "time"

const (
	midiNoteOff       int    = 0x80
	midiNoteOn        int    = 0x90
	midiVmax          int    = 127
	labelPortDefault  string = "<default>"
	labelPortSelected string = "<selected>"
)

type midiNote struct {
	Channel int
	Note    int
}

type noteMap map[string]midiNote

type event struct {
	Notes      row
	Velocities row
	Channels   row
	Figures    []midiFigure
}

type midiEvent struct {
	channel  int
	note     int
	velocity int
}

type figure struct {
	key      string
	velocity int
	pattern  string
}

type midiFigure []midiEvent

func playChord(e event, q chan midiEvent) {
	debugf("playChord(): %v", e)
	for i, note := range e.Notes {
		q <- midiEvent{e.Channels[i], note, e.Velocities[i]}
	}
}

// play a list of micro pattern midi events. The events are stretched
// out equally over the duration of one major event.
func playFigure(fig midiFigure, timing time.Duration, q chan midiEvent) {
	if len(fig) < 1 {
		return
	}
	debugf("playFigure(): %v", fig)
	sleepval := timing / time.Duration(len(fig))
	for i := 0; i < len(fig); i++ {
		if fig[i].note != 0 {
			q <- fig[i]
		}
		time.Sleep(sleepval)
	}
}

func processMidiQ(q chan midiEvent) {
	debugf("processMidiQ(): started")
	var ev midiEvent
	for {
		ev = <-q
		sendMidiNote(ev.channel, ev.note, ev.velocity)
	}
}
