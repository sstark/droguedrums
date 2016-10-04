package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

const (
	figurePrefix   rune   = '+'
	seqPrefix      rune   = ':'
	patternPause   rune   = '.'
	patternBeat    rune   = 'x'
	laneSplitStr   string = " "
	choiceSplitStr string = "|"
	defaultStep    int    = 8
	defaultSet     string = "default"
)

type part struct {
	name    string
	set     noteMap
	figures map[string]figure
	step    int
	bpm     string
	fx      []map[string]string
	lanes   []string
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
	// pick a random note from choiceSplitStr separated choices
	ss := strings.Split(s, choiceSplitStr)
	pick := ss[rand.Intn(len(ss))]
	channel = set[pick].Channel
	note = set[pick].Note
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
		lane := strings.Split(strings.TrimSpace(line), laneSplitStr)
		// use our own index because we need to be able to decrease it
		// if we encounter unwanted (uncounted) characters
		var i = -1
		for _, elem := range lane {
			i++
			if elem == "" {
				i--
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
	seen := make(map[string]bool)
	return sm.flattenRecursive(startAt, seen)
}

func (sm seqMap) flattenRecursive(startAt string, seen map[string]bool) []string {
	var l []string
	for _, p := range sm[startAt] {
		ok := false
		first, rest := splitFirstRune(p)
		if first == seqPrefix {
			_, ok = sm[rest]
		}
		if ok {
			s, ok := seen[rest]
			if ok && s {
				logger.Fatalf("flattenRecursive(): loop in sequence \"%v\". You can not make circular seqs.", rest)
			} else {
				seen[rest] = true
			}
			l = append(l, sm.flattenRecursive(rest, seen)...)
			seen[rest] = false
		} else {
			l = append(l, p)
		}
	}
	return l
}

func (d *drums) getParts(sets map[string]noteMap, figures map[string]figure) map[string]part {
	parts := make(map[string]part)
	var partsetname string
	var partstep int
	for _, inp := range d.Parts {
		if inp.Set == "" {
			partsetname = defaultSet
		} else {
			partsetname = inp.Set
		}
		partset, ok := sets[partsetname]
		if inp.Step == 0 {
			partstep = defaultStep
		} else {
			partstep = inp.Step
		}
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
			name:    inp.Name,
			set:     partset,
			figures: figures,
			step:    partstep,
			bpm:     inp.Bpm,
			fx:      inp.Fx,
			lanes:   lanes,
		}
	}
	return parts
}

func (d *drums) getSets() map[string]noteMap {
	// Make sure we have a different random sequence for each invocation
	rand.Seed(time.Now().Unix())
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
	var figurevelocity int
	figures := make(map[string]figure)
	for _, f := range d.Figures {
		debugf("getFigures(): %+v", f)
		if f.Velocity == 0 {
			figurevelocity = midiVmax
		} else {
			figurevelocity = f.Velocity
		}
		figures[f.Name] = figure{
			key:      f.Key,
			velocity: figurevelocity,
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
