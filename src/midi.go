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

func playChord(e event) {
	debugf("playChord(): %v", e)
	for i, note := range e.Notes {
		sendMidiNote(e.Channels[i], note, e.Velocities[i])
	}
}
