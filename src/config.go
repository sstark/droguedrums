
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

func str2note(s string) int {
    switch s {
    case "--": return 0
    case "bd": return 1
    case "sd": return 5
    case "hc": return 6
    case "ho": return 8
    }
    return 0
}

func text2matrix(txt []string) matrix {
    var m []row
    for _, line := range txt {
        var r []int
        lane := strings.Split(line, " ")
        for _, elem := range lane {
            r = append(r, str2note(elem))
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

func (d *Drums) LoadFromFile() {
    data, err := ioutil.ReadFile("drums.yml")
    if err != nil {
        panic(err)
    }
    yaml.Unmarshal(data, d)
}
