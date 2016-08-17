
package main

import (
    "github.com/rakyll/portmidi"
    "fmt"
    "time"
    "math/rand"
    "os"
    "log"
)

var logger *log.Logger

func Debugf(format string, args ...interface{}) {
    if os.Getenv("DRD_DEBUG") == "1" {
        logger.Output(2, fmt.Sprintf(format, args...))
    }
}

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
    Debugf("makeTicker(): timing: %v", timing)
    return time.NewTicker(timing)
}

func player(s *portmidi.Stream, q chan Part) {
    eventQueue := make(chan row)
    dacapo := make(chan bool)
    ticker := time.NewTicker(time.Millisecond)
    Debugf("player(): starting player loop")
    go func() { dacapo <- true }()
    for {
        select {
        case e := <-eventQueue:
            go playChord(s, e)
            <-ticker.C
        case <-dacapo:
            Debugf("player(): dacapo")
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
    logger = log.New(os.Stderr, "", log.Lshortfile)
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
    numSets := len(sets)
    numParts := len(parts)
    numSeqs := len(seqs)

    fmt.Printf("Read %d sets, %d parts, %d seqs\n", numSets, numParts, numSeqs)
    Debugf("main(): sets: %+v", sets)
    Debugf("main(): parts: %+v", parts)
    Debugf("main(): seqs: %+v", seqs)

    if (numSets < 1) {
        fmt.Println("no sets found")
        os.Exit(1)
    }
    if (numParts < 1) {
        fmt.Println("no parts found")
        os.Exit(1)
    }
    if (numSeqs < 1) {
        fmt.Println("no seqs found")
        os.Exit(1)
    }
    if _, ok := seqs["start"]; !ok {
        fmt.Println("start sequence not found")
        os.Exit(1)
    }

    trackQueue := make(chan Part)

    go player(out, trackQueue)

    for {
        for _, part := range seqs["start"] {
            fmt.Println("next:", part)
            trackQueue <- parts[part]
        }
    }
}
