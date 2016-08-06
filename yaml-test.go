
package main

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

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

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    data, err := ioutil.ReadFile("drums.yml")
    check(err)
    var d Drums
    yaml.Unmarshal(data, &d)
    fmt.Println(d.Sets)
    fmt.Println(d.Parts)
    fmt.Println(d.Seqs)
    for _, p := range d.Parts {
        if p.Name == "beat68" {
            for _, e := range p.Lanes {
                fmt.Println(e)
            }
            fmt.Println(p.Bpm)
        }
    }
}
