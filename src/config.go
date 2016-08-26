package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type part struct {
	Name  string
	Set   noteMap
	Step  int
	Bpm   int
	Fx    []map[string]string
	Lanes []string
}

type seq []string

type drums struct {
	Sets []struct {
		Name string
		Kit  []struct {
			Key     string
			Channel int
			Note    int
		}
	}
	Parts []struct {
		Name     string
		Set      string
		Step     int
		Bpm      int
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

// text2matrix takes the lanes as strings, looks up the channel
// and note number of the events and returns them as separate
// matrices.
func text2matrix(set noteMap, txt []string) (channels, notes matrix) {
	for _, line := range txt {
		var chanV, noteV []int
		lane := strings.Split(strings.TrimSpace(line), " ")
		for _, elem := range lane {
			c, n := translateKit(set, elem)
			chanV = append(chanV, c)
			noteV = append(noteV, n)
		}
		channels = append(channels, row(chanV))
		notes = append(notes, row(noteV))
	}
	return
}

func (d *drums) getSeqs() map[string]seq {
	seqs := make(map[string]seq)
	for _, s := range d.Seqs {
		var seqparts []string
		for _, partname := range s.Parts {
			seqparts = append(seqparts, partname)
		}
		seqs[s.Name] = seqparts
	}
	return seqs
}

func (d *drums) getParts(sets map[string]noteMap) map[string]part {
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
			Name:  inp.Name,
			Set:   partset,
			Step:  inp.Step,
			Bpm:   inp.Bpm,
			Fx:    inp.Fx,
			Lanes: lanes,
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
