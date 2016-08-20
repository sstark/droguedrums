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

func debugf(format string, args ...interface{}) {
	if os.Getenv("DRD_DEBUG") == "1" {
		logger.Output(2, fmt.Sprintf(format, args...))
	}
}

func playChord(s *portmidi.Stream, c row) {
	debugf("playChord(): %v", c)
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
	debugf("makeTicker(): timing: %v", timing)
	return time.NewTicker(timing)
}

func player(s *portmidi.Stream, q chan part) {
	eventQueue := make(chan row)
	dacapo := make(chan bool)
	ticker := time.NewTicker(time.Millisecond)
	debugf("player(): starting player loop")
	go func() { dacapo <- true }()
	for {
		select {
		case e := <-eventQueue:
			go playChord(s, e)
			<-ticker.C
		case <-dacapo:
			debugf("player(): dacapo")
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

	drums := new(drums)
	drums.loadFromFile(drumsfile)
	sets := drums.getSets()
	parts := drums.getParts(sets)
	seqs := drums.getSeqs()
	numSets := len(sets)
	numParts := len(parts)
	numSeqs := len(seqs)

	fmt.Printf("droguedrums: %d sets, %d parts, %d seqs\n", numSets, numParts, numSeqs)
	debugf("main(): sets: %+v", sets)
	debugf("main(): parts: %+v", parts)
	debugf("main(): seqs: %+v", seqs)

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

	trackQueue := make(chan part)
	go player(out, trackQueue)

	for _, part := range seqs["precount"] {
		trackQueue <- parts[part]
	}
	for {
		for _, partname := range seqs["start"] {
			debugf("main(): next: %v", partname)
			if part, ok := parts[partname]; !ok {
				logger.Printf("unknown part \"%s\"", partname)
				// avoid busy loop when all parts are unknown
				time.Sleep(time.Millisecond * 50)
			} else {
				trackQueue <- part
			}
		}
	}
}
