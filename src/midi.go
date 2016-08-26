package main

const (
	midiNoteOff = 0x80
	midiNoteOn  = 0x90
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

const (
	midiVmax int = 127
)

func playChord(e event) {
	debugf("playChord(): %v", e)
	for i, note := range e.Notes {
		sendMidiNote(e.Channels[i], note, e.Velocities[i])
	}
}
