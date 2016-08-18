package main

import (
	"fmt"
	"github.com/rakyll/portmidi"
	"log"
	"math/rand"
	"os"
	"time"
)

var logger *log.Logger

func Debugf(format string, args ...interface{}) {
	if os.Getenv("DRD_DEBUG") == "1" {
		logger.Output(2, fmt.Sprintf(format, args...))
	}
}

func playChord(s *portmidi.Stream, c row) {
	Debugf("playChord(): %v", c)
	dev := 60
	for _, i := range c {
		v := (rand.Int() % dev) - (dev / 2)
		v = 127 - (dev / 2) + v
		s.WriteShort(0x95, int64(i), int64(v))
	}
}

func makeTicker(bpm int, step int) *time.Ticker {
	step = step / 4
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
			fmt.Printf("> %s (%d/%d)\n", currentPart.Name, currentPart.Bpm, currentPart.Step)
			go func() {
				for _, c := range currentPart.Lanes.transpose() {
					eventQueue <- c
				}
				dacapo <- true
			}()
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	logger = log.New(os.Stderr, "", log.Lshortfile)

	var drumsfile string
	if len(os.Args) > 1 {
		drumsfile = os.Args[1]
	} else {
		drumsfile = "drums.yml"
	}

	err := portmidi.Initialize()
	checkErr(err)
	defer portmidi.Terminate()
	defaultOut := portmidi.DefaultOutputDeviceID()
	out, err := portmidi.NewOutputStream(defaultOut, 1024, 0)
	checkErr(err)
	defer out.Close()

	drums := new(Drums)
	drums.LoadFromFile(drumsfile)
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

	if numSets < 1 {
		logger.Fatalf("no sets found")
	}
	if numParts < 1 {
		logger.Fatalf("no parts found")
	}
	if numSeqs < 1 {
		logger.Fatalf("no seqs found")
	}
	if _, ok := seqs["start"]; !ok {
		logger.Fatalf("start sequence not found")
	}

	trackQueue := make(chan Part)
	go player(out, trackQueue)

	for _, part := range seqs["precount"] {
		trackQueue <- parts[part]
	}
	for {
		for _, part := range seqs["start"] {
			Debugf("next: %v", part)
			trackQueue <- parts[part]
		}
	}
}
