package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	seqNameStart    string        = "start"
	seqNamePrecount string        = "precount"
	partNameStop    string        = "_stop_"
	stopMessage     string        = "Press 'Enter' to continue..."
	unknownPartWait time.Duration = time.Millisecond * 50
	defaultMidiPort int           = -1
)

var (
	logger *log.Logger
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
	var ticker varTicker
	var eventCounter int
	var timing, timingIncrement time.Duration
	midiQueue := make(chan midiEvent)
	go processMidiQ(midiQueue)
	debugf("player(): starting player loop")
	go func() { dacapo <- true }()
	for {
		select {
		case e := <-eventQueue:
			go playChord(e, midiQueue)
			// figures are played independently to keep the timing tight for
			// normal events. But they should never be longer then the duration
			// of a player tick, so we pass the current timing to playFigure().
			for _, fig := range e.Figures {
				go playFigure(fig, timing, midiQueue)
			}
			<-ticker.C
			// variable bpm processing only if needed
			if timingIncrement != 0 {
				dur := timing + timingIncrement*time.Duration(eventCounter)
				debugf("player(): setDuration counter:%v dur:%v, timingIncrement:%v",
					eventCounter, dur, timingIncrement)
				ticker.SetDuration(dur)
				eventCounter++
			}
		case <-dacapo:
			debugf("player(): dacapo")
			currentPart := <-playQ
			go func() {
				channels, notes, figures := text2matrix(currentPart.set, currentPart.figures, currentPart.lanes)
				debugf("player(): %v", channels)
				debugf("player(): %v", notes)

				timing, timingIncrement, err = makeTiming(currentPart.bpm,
					currentPart.step,
					len(notes[0]))
				if err != nil {
					logger.Fatalf("bpm value could not be read: %s, %v", currentPart.bpm, err)
				}
				eventCounter = 0
				ticker.SetDuration(timing)

				// sanity check before transposing
				// FIXME: should probably be done in the matrix lib
				err = notes.check()
				if err != nil {
					logger.Fatalf("part \"%s\" has wrong format: %v", currentPart.name, err)
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
						Figures:    figures[i],
					}
				}
				dacapo <- true
			}()
		}
	}
}

func getDrumsfile(drumsfile string) (map[string]part, seqMap) {
	drums := new(drums)
	drums.loadFromFile(drumsfile)
	sets := drums.getSets()
	figures := drums.getFigures()
	parts := drums.getParts(sets, figures)
	seqs := drums.getSeqs()
	numSets := len(sets)
	numFigures := len(figures)
	numParts := len(parts)
	numSeqs := len(seqs)
	fmt.Printf("%d sets, %d figures, %d parts, %d seqs\n", numSets, numFigures, numParts, numSeqs)
	debugf("getDrumsfile(): sets: %+v", sets)
	debugf("getDrumsfile(): figures: %+v", figures)
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
	if _, ok := seqs[seqNameStart]; !ok {
		logger.Fatalf("start sequence not found")
	}
	fmt.Println(seqs.flatten(seqNameStart))
	return parts, seqs
}

func feeder(drumsfile string, playQ chan part) {
	parts, seqs := getDrumsfile(drumsfile)
	feedParts := func(startAt string) {
		for _, partname := range seqs.flatten(startAt) {
			debugf("feeder(): next: %v", partname)
			if partname == partNameStop {
				fmt.Println(stopMessage)
				bufio.NewReader(os.Stdin).ReadBytes('\n')
				continue
			}
			if part, ok := parts[partname]; !ok {
				logger.Printf("unknown part \"%s\"", partname)
				// avoid busy loop when all parts are unknown
				time.Sleep(unknownPartWait)
			} else {
				playQ <- part
				fmt.Printf("[%s] ", part.name)
			}
		}
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGUSR1)
	debugf("installed signal handler")
	feedParts(seqNamePrecount)
	for {
		select {
		case sig := <-sigc:
			debugf("feeder(): got signal %v, re-reading drumsfile", sig)
			fmt.Println("re-reading input file")
			parts, seqs = getDrumsfile(drumsfile)
			debugf("feeder(): done re-reading drumsfile")
		default:
			fmt.Printf("> ")
			feedParts(seqNameStart)
			fmt.Println()
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
	fmt.Printf("droguedrums %s\n", version)

	var chosenPort int
	flag.IntVar(&chosenPort, "port", defaultMidiPort, "choose output port")
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
