
package main

import (
    "github.com/rakyll/portmidi"
    "fmt"
    "time"
    "math/rand"
    "os"
)

func playChord(s *portmidi.Stream, c row) {
    fmt.Println(c)
    dev := 60
    for _, i := range c {
        v := (rand.Int() % dev) - (dev/2)
        v = 127-(dev/2) + v
        s.WriteShort(0x95, int64(i), int64(v))
    }
}

func makeTicker(bpm int, step int) *time.Ticker {
    step = step/4
    timing := (time.Minute / time.Duration(bpm)) / time.Duration(step)
    fmt.Println("timing:", timing)
    return time.NewTicker(timing)
}

func player(s *portmidi.Stream, q chan Part) {
    eventQueue := make(chan row)
    dacapo := make(chan bool)
    ticker := time.NewTicker(time.Millisecond)
    fmt.Println("starting player loop")
    go func() { dacapo <- true }()
    for {
        select {
        case e := <-eventQueue:
            go playChord(s, e)
            <-ticker.C
        case <-dacapo:
            currentPart := <-q
            ticker.Stop()
            ticker = makeTicker(currentPart.Bpm, currentPart.Step)
            fmt.Println("part:", currentPart.Name)
            go func() {
                for _, c := range currentPart.Lanes.transpose() {
                    eventQueue <- c
                }
                dacapo <- true
            }()
        }
    }
}

func checkErr (err error) {
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func main() {
    err := portmidi.Initialize()
    checkErr(err)
    defer portmidi.Terminate()
    defaultOut := portmidi.DefaultOutputDeviceID()
    out, err := portmidi.NewOutputStream(defaultOut, 1024, 0)
    checkErr(err)
    defer out.Close()

    drums := new(Drums)
    drums.LoadFromFile()
    sets := drums.GetSets()
    parts := drums.GetParts(sets)
    seqs := drums.GetSeqs()
    fmt.Println(sets)
    fmt.Println(parts)
    fmt.Println(seqs)

    trackQueue := make(chan Part)

    go player(out, trackQueue)

    for {
        for _, part := range seqs["start"] {
            fmt.Println("next:", part)
            trackQueue <- parts[part]
        }
    }
}
