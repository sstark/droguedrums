package main

type midiNote struct {
	Channel int
	Note    int
}

type noteMap map[string]midiNote

type event struct {
	Notes      row
	Velocities row
}

const (
	midiVmax int = 127
)

func playChord(e event) {
	debugf("playChord(): %v", e)
	for i, note := range e.Notes {
		sendMidiNote(5, note, e.Velocities[i])
	}
}
