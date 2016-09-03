package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Version   string
	BuildTime string
	logger    *log.Logger
)

func debugf(format string, args ...interface{}) {
	if os.Getenv("DRD_DEBUG") == "1" {
		logger.Output(2, fmt.Sprintf(format, args...))
	}
}

func player(playQ chan part) {
	var err error
	eventQueue := make(chan event)
	dacapo := make(chan bool)
	var ticker VarTicker
	var eventCounter int = 0
	ticker.SetDuration(time.Millisecond)
	var timing, timingIncrement time.Duration
	midiQueue := make(chan midiEvent)
	go processMidiQ(midiQueue)
	debugf("player(): starting player loop")
	go func() { dacapo <- true }()
	for {
		select {
		case e := <-eventQueue:
			<-ticker.C
			go playChord(e, midiQueue)
			// variable bpm processing only if needed
			if timingIncrement != 0 {
				dur := timing + timingIncrement*time.Duration(eventCounter)
				debugf("player(): setDuration counter:%v dur:%v, timingIncrement:%v",
					eventCounter, dur, timingIncrement)
				ticker.SetDuration(dur)
				eventCounter += 1
			}
		case <-dacapo:
			debugf("player(): dacapo")
			currentPart := <-playQ
			fmt.Printf("> %s (%s/%d)\n", currentPart.Name, currentPart.Bpm, currentPart.Step)
			go func() {
				channels, notes := text2matrix(currentPart.Set, currentPart.Lanes)
				debugf("player(): %v", channels)
				debugf("player(): %v", notes)

				timing, timingIncrement, err = makeTiming(currentPart.Bpm,
					currentPart.Step,
					len(notes[0]))
				if err != nil {
					logger.Fatalf("bpm value could not be read: %s, %v", currentPart.Bpm, err)
				}
				eventCounter = 0
				ticker.SetDuration(timing)

				// sanity check before transposing
				// FIXME: should probably be done in the matrix lib
				err = notes.check()
				if err != nil {
					logger.Fatalf("part \"%s\" has wrong format: %v", currentPart.Name, err)
				}
				debugf("player(): %+v:", currentPart)
				vmap := genVelocityMap(currentPart, notes).transpose()
				debugf("player(): vmap: %v:", vmap)
				cmap := channels.transpose()
				for i, c := range notes.transpose() {
					eventQueue <- event{
						Notes:      c,
						Velocities: vmap[i],
						Channels:   cmap[i],
					}
				}
				dacapo <- true
			}()
		}
	}
}

func getDrumsfile(drumsfile string) (map[string]part, map[string]seq) {
	drums := new(drums)
	drums.loadFromFile(drumsfile)
	sets := drums.getSets()
	parts := drums.getParts(sets)
	seqs := drums.getSeqs()
	numSets := len(sets)
	numParts := len(parts)
	numSeqs := len(seqs)
	fmt.Printf("%d sets, %d parts, %d seqs\n", numSets, numParts, numSeqs)
	debugf("getDrumsfile(): sets: %+v", sets)
	debugf("getDrumsfile(): parts: %+v", parts)
	debugf("getDrumsfile(): seqs: %+v", seqs)
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
	return parts, seqs
}

func feeder(drumsfile string, playQ chan part) {
	parts, seqs := getDrumsfile(drumsfile)
	for _, part := range seqs["precount"] {
		playQ <- parts[part]
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGUSR1)
	debugf("installed signal handler")
	for {
		select {
		case sig := <-sigc:
			debugf("feeder(): got signal %v, re-reading drumsfile", sig)
			fmt.Println("re-reading input file")
			parts, seqs = getDrumsfile(drumsfile)
			debugf("feeder(): done re-reading drumsfile", sig)
		default:
			for _, partname := range seqs["start"] {
				debugf("feeder(): next: %v", partname)
				if part, ok := parts[partname]; !ok {
					logger.Printf("unknown part \"%s\"", partname)
					// avoid busy loop when all parts are unknown
					time.Sleep(time.Millisecond * 50)
				} else {
					playQ <- part
				}
			}
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
	fmt.Printf("droguedrums %s (built %s)\n", Version, BuildTime)

	var chosenPort int
	flag.IntVar(&chosenPort, "port", -1, "choose output port")
	flag.Parse()

	var drumsfile string
	if flag.NArg() > 0 {
		drumsfile = flag.Args()[0]
	}

	err := initMidi(chosenPort)
	checkErr(err)

	defer closeMidi()

	if drumsfile == "" {
		fmt.Println("no input file")
		os.Exit(1)
	} else {
		fmt.Printf("input file: %s\n", drumsfile)
	}
	fmt.Println("-- player starting --")

	playQ := make(chan part)
	go player(playQ)
	feeder(drumsfile, playQ)
}
