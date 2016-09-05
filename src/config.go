package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

const (
	figurePrefix rune = '+'
	seqPrefix    rune = ':'
	patternPause rune = '.'
	patternBeat  rune = 'x'
)

type part struct {
	Name    string
	Set     noteMap
	Figures map[string]figure
	Step    int
	Bpm     string
	Fx      []map[string]string
	Lanes   []string
}

type seq []string
type seqMap map[string]seq

type drums struct {
	Sets []struct {
		Name string
		Kit  []struct {
			Key     string
			Channel int
			Note    int
		}
	}
	Figures []struct {
		Name     string
		Key      string
		Velocity int
		Pattern  string
	}
	Parts []struct {
		Name     string
		Set      string
		Step     int
		Bpm      string
		Fx       []map[string]string
		Lanes    []string
		Genlanes []map[string]map[string]string
	}
	Seqs []struct {
		Name  string
		Parts []string
	}
}

func (d *drums) dump() {
	fmt.Println(*d)
}

func translateKit(set noteMap, s string) (channel, note int) {
	channel = set[s].Channel
	note = set[s].Note
	return
}

// renders the pattern string of a figure into individual midi events,
// translating the kit keys on the way.
func translateFigure(set noteMap, figures map[string]figure, s string) midiFigure {
	var fig midiFigure
	f := figures[s]
	c, n := translateKit(set, f.key)
	for _, elem := range f.pattern {
		switch elem {
		case patternPause:
			fig = append(fig, midiEvent{})
		case patternBeat:
			fig = append(fig, midiEvent{c, n, f.velocity})
		}
	}
	return fig
}

// text2matrix takes the lanes as strings, looks up the channel
// and note number of the events and returns them as separate
// matrices.
// It also returns a map of micro pattern figures. The index of the
// map is the position in the lane, so the player can fire it
// at the right time later.
func text2matrix(set noteMap, figures map[string]figure, txt []string) (channels, notes matrix, figuremap map[int][]midiFigure) {
	figuremap = make(map[int][]midiFigure)
	for _, line := range txt {
		var chanV, noteV []int
		lane := strings.Split(strings.TrimSpace(line), " ")
		// use our own index because we need to be able to decrease it
		// if we encounter unwanted (uncounted) characters
		var i int = -1
		for _, elem := range lane {
			i += 1
			if elem == "" {
				i -= 1
				continue
			}
			first, rest := splitFirstRune(elem)
			if first == figurePrefix {
				// generate a micro pattern
				figuremap[i] = append(figuremap[i], translateFigure(set, figures, rest))
				// also generate an empty single note
				chanV = append(chanV, 0)
				noteV = append(noteV, 0)
			} else {
				// generate a normal note
				c, n := translateKit(set, elem)
				chanV = append(chanV, c)
				noteV = append(noteV, n)
			}
		}
		channels = append(channels, row(chanV))
		notes = append(notes, row(noteV))
	}
	debugf("text2matrix(): figuremap: %v", figuremap)
	return
}

func (d *drums) getSeqs() seqMap {
	seqs := make(seqMap)
	for _, s := range d.Seqs {
		var seqparts []string
		for _, partname := range s.Parts {
			seqparts = append(seqparts, partname)
		}
		seqs[s.Name] = seqparts
	}
	return seqs
}

// take on UTF-8 safe string splitting
func splitFirstRune(s string) (first rune, rest string) {
	r := []rune(s)
	first = r[0]
	rest = string(r[1:])
	return
}

// chains up all seqs by looking up references to other
// seqs or parts and returns a list of part names.

func (sm seqMap) flatten(startAt string) []string {
	var l []string
	for _, p := range sm[startAt] {
		ok := false
		first, rest := splitFirstRune(p)
		if first == seqPrefix {
			_, ok = sm[rest]
		}
		if ok {
			l = append(l, sm.flatten(rest)...)
		} else {
			l = append(l, p)
		}
	}
	return l
}

func (d *drums) getParts(sets map[string]noteMap, figures map[string]figure) map[string]part {
	parts := make(map[string]part)
	var partsetname string
	for _, inp := range d.Parts {
		if inp.Set == "" {
			partsetname = "default"
		} else {
			partsetname = inp.Set
		}
		partset, ok := sets[partsetname]
		if !ok {
			logger.Printf("unknown set \"%s\"", partsetname)
		}
		genlanes, err := renderGenlanes(inp.Genlanes)
		if err != nil {
			logger.Print(err)
		}
		lanes := inp.Lanes
		lanes = append(lanes, genlanes...)
		debugf("getParts(): %#v", lanes)
		parts[inp.Name] = part{
			Name:    inp.Name,
			Set:     partset,
			Figures: figures,
			Step:    inp.Step,
			Bpm:     inp.Bpm,
			Fx:      inp.Fx,
			Lanes:   lanes,
		}
	}
	return parts
}

func (d *drums) getSets() map[string]noteMap {
	sets := make(map[string]noteMap)
	for _, set := range d.Sets {
		debugf("getSets(): %+v", set)
		notes := make(noteMap)
		for _, note := range set.Kit {
			notes[note.Key] = midiNote{
				Channel: note.Channel,
				Note:    note.Note,
			}
			debugf("getSets(): %+v", note)
		}
		sets[set.Name] = notes
	}
	return sets
}

func (d *drums) getFigures() map[string]figure {
	figures := make(map[string]figure)
	for _, f := range d.Figures {
		debugf("getFigures(): %+v", f)
		figures[f.Name] = figure{
			key:      f.Key,
			velocity: f.Velocity,
			pattern:  f.Pattern,
		}
	}
	return figures
}

func (d *drums) loadFromFile(fn string) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	err = yaml.Unmarshal(data, d)
	if err != nil {
		logger.Fatalf("Syntax error when reading file \"%s\". Maybe it is not proper YAML format.\n%v", fn, err)
	}
}
