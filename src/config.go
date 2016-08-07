
package main

import (
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "fmt"
    "strings"
)

type Part struct {
    Name string
    Set string
    Step int
    Bpm int
    Lanes matrix
}

type MidiNote struct {
    Channel int
    Note int
}

type NoteMap map[string]MidiNote

type Drums struct {
    Sets []struct {
        Name string
        Kit []struct {
            Key string
            Channel int
            Note int
        }
    }
    Parts []struct {
        Name string
        Set string
        Step int
        Bpm int
        Lanes []string
    }
    Seqs []struct {
        Name string
        Parts []string
    }
}

func (d *Drums) Dump() {
    fmt.Println(*d)
}

func translateKit(s string) int {
    switch s {
    case "--": return 0
    case "bd": return 1
    case "sd": return 5
    case "hc": return 6
    case "ho": return 8
    case "cl": return 12
    }
    return 0
}

func text2matrix(txt []string) matrix {
    var m []row
    for _, line := range txt {
        var r []int
        lane := strings.Split(line, " ")
        for _, elem := range lane {
            r = append(r, translateKit(elem))
        }
        m = append(m, row(r))
    }
    return m
}

func (d *Drums) GetParts() []Part {
    var parts []Part
    for _, inp := range d.Parts {
        parts = append(parts, Part{
            Name: inp.Name,
            Set: inp.Set,
            Step: inp.Step,
            Bpm: inp.Bpm,
            Lanes: text2matrix(inp.Lanes),
        })
    }
    return parts
}

func (d *Drums) GetSets() map[string]NoteMap {
    sets := make(map[string]NoteMap)
    for _, set := range d.Sets {
        fmt.Println(set)
        notes := make(NoteMap)
        for _, note := range set.Kit {
            notes[note.Key] = MidiNote{
                Channel: note.Channel,
                Note: note.Note,
            }
            fmt.Println(note)
        }
        sets[set.Name] = notes
    }
    return sets
}

func (d *Drums) LoadFromFile() {
    data, err := ioutil.ReadFile("drums.yml")
    if err != nil {
        panic(err)
    }
    yaml.Unmarshal(data, d)
}
