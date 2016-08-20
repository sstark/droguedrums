package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type part struct {
	Name  string
	Set   string
	Step  int
	Bpm   int
	Lanes matrix
}

type midiNote struct {
	Channel int
	Note    int
}

type noteMap map[string]midiNote

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
		Name  string
		Set   string
		Step  int
		Bpm   int
		Lanes []string
	}
	Seqs []struct {
		Name  string
		Parts []string
	}
}

func (d *drums) dump() {
	fmt.Println(*d)
}

func translateKit(set noteMap, s string) int {
	return set[s].Note
}

func text2matrix(set noteMap, txt []string) matrix {
	var m []row
	for _, line := range txt {
		var r []int
		lane := strings.Split(line, " ")
		for _, elem := range lane {
			r = append(r, translateKit(set, elem))
		}
		m = append(m, row(r))
	}
	return m
}

func (d *drums) getSeqs() map[string][]string {
	seqs := make(map[string][]string)
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
		parts[inp.Name] = part{
			Name:  inp.Name,
			Set:   inp.Set,
			Step:  inp.Step,
			Bpm:   inp.Bpm,
			Lanes: text2matrix(partset, inp.Lanes),
		}
		err := parts[inp.Name].Lanes.check()
		if err != nil {
			logger.Fatalf("part \"%s\" has wrong format: %v", inp.Name, err)
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
