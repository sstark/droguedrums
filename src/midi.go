package main

const (
	midiNoteOff int = 0x80
	midiNoteOn  int = 0x90
	midiVmax    int = 127
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
}

type midiEvent struct {
	channel  int
	note     int
	velocity int
}

func playChord(e event, q chan midiEvent) {
	debugf("playChord(): %v", e)
	for i, note := range e.Notes {
		q <- midiEvent{e.Channels[i], note, e.Velocities[i]}
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
