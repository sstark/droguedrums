
# About

droguedrums is a text file based drum sequencer for the command line.

It reads a definition of drum sets, parts and sequences from a text file and
plays it back as midi. No sounds are provided. Development is done with the
following workflow in mind:

  - define rhythmical patterns using a simple text file
  - feed midi data to a hardware drum synthesizer
  - process and record the audio from hardware

There is a single input file containing definitions for drum kits, parts and
sequences. The file is specified at the command line.

# Installation

## Download binary

- compiled on Ubuntu 16.04: ...
- compiled on OSX 10.9.5: ...

## Requirements for running

Requirements for running the binary:

  - portmidi library
    - OSX:      `brew install portmidi`
    - Ubuntu:   `apt-get install libportmidi0`

## Requirements for building from source

1. libportmidi headers
    - OSX:      `brew install portmidi`
    - Ubuntu:   `apt-get install libportmidi-dev`

    Unfortunately the current versions of the Ubuntu (v200) and Debian (v187)
    libportmidi packages do not link properly with the portmidi go bindings.
    See below for a workaround.
  
2. The go programming language: https://golang.org/

3. libportmidi bindings for go: (https://github.com/rakyll/portmidi) `go get
   github.com/rakyll/portmidi`

4. run `make`

5. (optional) run `make install`

# Usage

start like this:

```
$ droguedrums myfile.yml
droguedrums 0.9 (built 2016-08-26T13:16:29+0200)
midi outputs found:
0: "CoreMIDI/FluidSynth virtual port (45160)" <default>
1: "CoreMIDI/FluidSynth virtual port (45253)"
input file: myfile.yml
-- player starting --
1 sets, 6 parts, 1 seqs
> demo1 (160/8)
...
```

It will try to choose a working midi out port. But it will also list all it
found so you can choose a sepcific one using the `-port` parameter:

```
$ droguedrums -port 1 myfile.yml
droguedrums 0.9 (built 2016-08-26T13:16:29+0200)
midi outputs found:
0: "CoreMIDI/FluidSynth virtual port (45160)" <default>
1: "CoreMIDI/FluidSynth virtual port (45253)" <selected>
input file: myfile.yml
-- player starting --
1 sets, 6 parts, 1 seqs
> demo1 (160/8)
...
```

# Writing input files

The YAML file format seems like a good compromise between easy machine
readability and human editing. For this reason, but also for ease of
development (making a program read yaml is simple), YAML was chosen as the
input file format for drougedrums.

## Input file sections

### Sets

A *set* maps arbitrary keys to midi channels and note numbers. Each part has a
specific drumset associated with it. The keys are used later in the parts as a
reference to a specific note to be played.

In the examples two-character keys are used throughout, just to make the
formatting concise. You can, however use any number of characters for note
keys. Also the pause symbol "--" is completely arbitrary and can be changed to
your liking.

#### Example

```
sets:
    - name: default
      kit:
          - {key: --, channel:  5, note:  0} # pause
          - {key: bd, channel:  5, note:  1} # bass drum
          - {key: hc, channel:  5, note:  2} # high hat closed
          - {key: sd, channel:  5, note:  5} # snare drum
```

### Parts

*Parts* are the building blocks of rhythmic patterns. Each part has its own bpm
value, step size and drumset. Patterns are defined in an arbitrary number of
lanes, using the strings from the drum set defined earlier.

#### Example

```
parts:
    - name: countin
      bpm: 100
      step: 4
      lanes:
          - hc hc hc hc

    - name: part1
      bpm: 100
      step: 16
      lanes:
          - hc hc hc hc hc hc hc hc hc hc hc hc hc hc hc hc
          - bd -- -- -- -- -- bd -- -- -- bd -- -- -- bd --
          - -- sd -- -- sd sd -- -- -- sd -- -- -- sd -- sd
```

### Seqs

*Seqs* are playlists of parts. They can define whole projects or simply play
back some parts in a loop.

There are some reserved seq names:

  - *start*: This is always the first sequence played.
  - *precount*: This will be played only once at the beginning, before the
    start seq.
  
#### Example

```
seqs:
    - name: precount
      parts:
          - countin
    - name: start
      parts:
          - part1
          - part2
          - part3
          - end
```

## Advanced concepts

### Part FX

Parts can have an fx section that will modify the events in the part. To use
multiple effects you write them as a list like this:

```
fx: [rampv: 23-85, randv: 22]
```

See the full example files provided with droguedrums for some possibilities.

#### randv

Usage: `fx: [randv: <randomness>]`

Random velocity. Will make the velocity values of the events vary randomly by
_randomness_ steps up or down. The velocity range is limited by midi and goes
from 0 to 127. Will be applied _after_ rampv.

#### rampv

Usage: `fx: [rampv: <start>-<end>]`

Ramp velocity. Will make the velocity go from _start_ to _end_ over the whole
part. Useful for cresecndo or diminuendo, but also for limiting velocity.

### Genlanes

Genlanes are lanes with algorithmically generated events.

# Controlling devices

droguedrums has been tested with a Vermona DRM1-MKIII and fluidsynth. In theory
everything that can receive MIDI events can be controlled, but the main focus of
the program is rhythms.

Currently no midi note off events are ever sent. With some sounds that can lead
to notes playing forever. With percussive sounds this should not be a problem
usually.

Program change messages are not implemented, so you have to map notes to
different banks or midi channels.

# Workaround for libportmidi on Debian/Ubuntu

...

# License

Copyright © 2016 Sebastian Stark <sstark@mailbox.org>

Permission to use, copy, modify, and distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
