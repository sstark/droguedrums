
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

type Player struct {
    Bpm int
    Lane []string
    Kit Kit
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

func player(s *portmidi.Stream, ticker *time.Ticker, q chan matrix) {
    eventQueue := make(chan row)
    dacapo := make(chan bool)
    fmt.Println("starting player loop")
    go func() { dacapo <- true }()
    for {
        select {
        case <-ticker.C:
            go playChord(s, <-eventQueue)
        case <-dacapo:
            currentTrack := <-q
            fmt.Println("received new track")
            go func() {
                for _, c := range currentTrack.transpose() {
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

    ticker := time.NewTicker(160 * time.Millisecond)
    trackQueue := make(chan matrix)

    go player(out, ticker, trackQueue)

    for {
        trackQueue <- demo1()
        trackQueue <- demo2()
    }
}
