package main

import (
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

func makeTicker(bpm int, step int) *time.Ticker {
	step = step / 4
	timing := (time.Minute / time.Duration(bpm)) / time.Duration(step)
	debugf("makeTicker(): timing: %v", timing)
	return time.NewTicker(timing)
}

func player(playQ chan part) {
	eventQueue := make(chan event)
	dacapo := make(chan bool)
	ticker := time.NewTicker(time.Millisecond)
	debugf("player(): starting player loop")
	go func() { dacapo <- true }()
	for {
		select {
		case e := <-eventQueue:
			go playChord(e)
			<-ticker.C
		case <-dacapo:
			debugf("player(): dacapo")
			currentPart := <-playQ
			ticker.Stop()
			ticker = makeTicker(currentPart.Bpm, currentPart.Step)
			fmt.Printf("> %s (%d/%d)\n", currentPart.Name, currentPart.Bpm, currentPart.Step)
			go func() {
				vmap := genVelocityMap(currentPart).transpose()
				for i, c := range currentPart.Lanes.transpose() {
					eventQueue <- event{
						Notes:      c,
						Velocities: vmap[i],
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

	var drumsfile string
	if len(os.Args) > 1 {
		drumsfile = os.Args[1]
	} else {
		drumsfile = "drums.yml"
	}

	initMidi()
	defer closeMidi()

	playQ := make(chan part)
	go player(playQ)
	feeder(drumsfile, playQ)
}
