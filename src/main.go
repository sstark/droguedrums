
package main

import (
    "github.com/rakyll/portmidi"
    "fmt"
    "time"
    "math/rand"
)

type Kit []struct {
    Key string
    Channel int
    Note int
}

type Part struct {
    Name string
    Set string
    Step int
    Bpm int
    Lanes matrix
 }

func playChord(s *portmidi.Stream, c row) {
    fmt.Println(c)
    dev := 32
    for _, i := range c {
        v := (rand.Int() % dev) - (dev/2)
        v = 127-(dev/2) + v
        s.WriteShort(0x95, int64(i), int64(v))
    }
}

func makeTicker(bpm int, step int) *time.Ticker {
    step = step/4
    timing := (time.Minute / time.Duration(bpm)) / time.Duration(step)
    fmt.Println("New timing:", timing)
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
            fmt.Println("received new part", currentPart.Name)
            go func() {
                for _, c := range currentPart.Lanes.transpose() {
                    eventQueue <- c
                }
                dacapo <- true
            }()
        }
    }
}

func main() {
    portmidi.Initialize()
    defer portmidi.Terminate()
    defaultOut := portmidi.DefaultOutputDeviceID()
    out, _ := portmidi.NewOutputStream(defaultOut, 1024, 0)
    defer out.Close()

    trackQueue := make(chan Part)

    go player(out, trackQueue)

    for {
        trackQueue <- demo1()
        trackQueue <- demo2()
        trackQueue <- demo3()
    }
}
