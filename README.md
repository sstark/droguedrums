
droguedrums - text file based drum sequencer
============================================

Reads a definition of drumsets, parts and sequences from a text file and plays
it back as midi. No sounds are provided. Droguedrums has been developed with
the following workflow in mind:

  - define rhythmical patterns using a simple text file
  - feed midi to a hardware drum kit, e. g. a Vermona DRM1-MKIII
  - process and record the audio from hardware

A *set* maps arbitrary strings to midi channels and note numbers. Each part
has a specific drumset associated with it.

*Parts* are the building blocks of rhythmic patterns. Each part has its own bpm
value, step size and drumset. Patterns are defined in an arbitrary number of
lanes, using the strings from the drum set defined earlier.

*Seqs* are playlists of parts. They can define whole projects or simply play back
some parts in a loop.

There is a single input file holding all the above definitions. The file is
specified at the command line. During playback it can be changed and replaced
on the fly in the running instance of the program.

Currently the input file is in YAML format, mainly for ease of early
development of the program. The format is likely to change to something less
cumbersome in future.


Requirements for running the binary:

  - portmidi library
    - OSX:      `brew install portmidi`
    - Ubuntu:   `apt-get install libportmidi0`

  - Text editor of choice. Good YAML support recommended.


Build from source:

  - libportmidi headers
    - OSX:      `brew install portmidi`
    - Ubuntu:   `apt-get install libportmidi0-dev`
  
  - The go programming language: https://golang.org/

  - libportmidi bindings for go: (https://github.com/rakyll/portmidi)
    - `go get github.com/rakyll/portmidi`
