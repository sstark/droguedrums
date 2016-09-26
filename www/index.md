
![droguetrees](droguetrees.jpg)

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

For the easiest possible start just download the source and one of the
precompiled binaries below and run one of the contained example files:

```
$ ./droguedrums testfiles/beat1.yml
```

Or continue reading.

## Download

### Source archive

- **Version 1.1 (2016-09-07)**: [droguedrums-1.1.zip](droguedrums-1.1.zip)
    - New in 1.1: support for coremidi on OSX, variable bpm in parts, figures (micro patterns)
- **Version 1.0 (2016-08-28)**: [droguedrums-1.0.zip](droguedrums-1.0.zip)

Download this in any case, it has all the examples.

See or get the source code repository here: [Github project](https://github.com/sstark/droguedrums)

### Precompiled

For convenience you can download precompiled versions of the binary:

- compiled on Ubuntu 16.04: [droguedrums-1.1-linux.zip](droguedrums-1.1-linux.zip) (using portmidi)
- compiled on OSX 10.9.5: [droguedrums-1.1-darwin.zip](droguedrums-1.1-darwin.zip) (using coremidi)

Just unzip into the unpacked source archive downloaded before.

## Requirements for running

Requirements for running the binary:

  - portmidi library
    - OSX:      `brew install portmidi` (not needed if you compiled with coremidi)
    - Ubuntu:   `apt-get install libportmidi0`

## Requirements for building from source

1. libportmidi headers
    - OSX:      `brew install portmidi` (not needed if you compile with coremidi)
    - Ubuntu:   `apt-get install libportmidi-dev`

2. The go programming language: <https://golang.org/>

3. (if not using coremidi, mandatory on Linux) libportmidi bindings for go
   (<https://github.com/rakyll/portmidi>): `go get github.com/rakyll/portmidi`

    > Unfortunately the current versions of the Ubuntu (v200) and Debian (v184)
    > libportmidi packages do not link properly with the portmidi go bindings,
    > which could make this step fail.  See below for a workaround.

4. (OSX only, optional) coremidi bindings for go
   (<https://github.com/youpy/go-coremidi>): `go get
   github.com/youpy/go-coremidi`

5. run `make`

6. (OSX only, optional) run `make MIDILIB=coremidi` to use coremidi directly,
   which is recommended on OSX.

7. (optional) run `make install` or copy the binary to a convenient place yourself.

# Usage

Start the program like this:

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
found so you can choose a specific one using the `-port` parameter:

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

While it is running you can tell the program to re-read its input file by sending it SIGUSR1, e. g. like this:

```
$ killall -USR1 droguedrums
```

It will play the new file after the current sequence has finished.

To stop the program press `Ctrl-C`.

# Writing input files

The YAML file format seems like a good compromise between easy machine
readability and human editing. For this reason, but also for ease of
development (making a program read yaml is simple), YAML was chosen as the
input file format for drougedrums.

The following sections will give an overview of what things you can put into a
droguedrums input file. See the full example files provided with the download
for some possibilities.

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

*Parts* are the basic elements for building patterns. They provide lanes, the
main building blocks. Each part has its own bpm value, step size and drumset.
Patterns are defined in an arbitrary number of lanes, using the keys from
the drum set defined earlier.

_Bpm_ defines how many 4th note are in one minute for this pattern. You can
also write a bpm value like "120-150", which will make the bpm rate go from 120
to 150 gradually over the duration of the part.

_Step_ defines the duration of events in a lane. This could be 8 for saying the
lane contains 8th notes or 16th, 32th and so on. Technically you could also
give values like 3 or 7.

The length of the first lane defines the length of the pattern. If you have
many lanes, all other lanes must be the same length or shorter than the first
one.

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

You can refer to other seqs by using their name as a part name prefixed by a ":".
  
#### Example

```
seqs:
    - name: precount
      parts:
          - countin
    - name: end
          - end1
          - end2
    - name: start
      parts:
          - part1
          - part2
          - part3
          - :end
```

## Advanced concepts

### Part FX

Parts can have an fx section that will modify the events in the part. To use
multiple effects you write them as a list like this:

```
fx: [rampv: 23-85, randv: 22]
```

#### randv

Usage: `fx: [randv: <randomness>]`

Random velocity. Will make the velocity values of the events vary randomly by
_randomness_ steps up or down. The velocity range is limited by midi and goes
from 0 to 127. Will be applied _after_ rampv.

#### rampv

Usage: `fx: [rampv: <start>-<end>]`

Ramp velocity. Will make the velocity go from _start_ to _end_ over the whole
part. Useful for crescendo or diminuendo, but also for limiting velocity.


### Figures

A _figure_ is a short pattern that can be played like a normal note. Just
instead it will play the defined pattern. Figures let you add micro beats,
dead notes, grace notes, triplet feels or subliminal timings in general.

The name of a figure will be used to refer to it in the lanes of a part,
prefixed by a "+". The key defines which key from the set it will play.
The pattern is a sequence of "." or "x" signs, . denoting a pause and x
denoting a trigger of the chosen key. All events in the pattern (beats and
pauses) are timed such that the first one happens at the beat and the remaining
ones are equally spaced until the next beat. This means a figure always has the
same length in time, but since you can place as many micro events in it as you
like, their spacing will change accordingly to fit into that single beat.

For instance, to get a very short grace note you could define a figure pattern
like ".....x" and place it just one beat before the actual beat. The "x" will
play very shortly before.

Examples:

```
figures:
    # some triplet patterns
    - {name: f1, key: fl, velocity:  65, pattern: .xx}
    - {name: f2, key: ce, velocity: 100, pattern: ..x}

parts:
    - name: part1
      bpm: 120
      step: 8
      lanes:
        - +f2 ce +f2 ce
```

See testfiles/beat8.yml for a more elaborate example.


### Genlanes

Genlanes are lanes with algorithmically generated events. Parts can have a
_genlanes_ section in addition to normal lanes, or even genlanes only.

#### equid

Usage: `- equid {note: <key>, length: 16, dist: 1, start: 1}`

Generates equidistant events starting at _start_, with distance _dist_.

Example:

```yaml
- equid {note: hc, length 8, dist: 1, start: 1}
- equid {note: sd, length 8, dist: 2, start: 4}
```

will generate two lanes as if manually written like this:

```yaml
- hc hc hc hc hc hc hc hc
- -- -- -- sd -- sd -- sd
```

#### sinez

Usage: `- sinez {note: <key>, length: 32, period: 1.0, xshift: 0.0, yshift: 0.0}`

Projects a sine wave onto a lane of _length_ events. Events are generated at
zero crossings of the sine through the lane.

With the default _period_ of 1.0, the period of the sine is exactly the length
of the lane. A smaller period value than 1 will shrink the wave, a higher value
will stretch it. You could think of the period value as wave length in relation
to lane length. Period values close to or over 2 will give you no events. With
_xshift_ you can move the sine by time, while a shift value of 1 or -1 means
shift by one event. With _yshift_ you can move the sine up or down by up to -1
or 1. Higher or lower values do not make sense, as the sine will not have any
zero crossings with the lane then.

The sinez function can be seen as a not-quite-random generation of events. Just
play around to find combinations of values that you find musical. For smaller
values of length or higher values of period the result of sinez can be
disappointing, since little or no events are generated.

Also, aliasing effects make the result a bit unpredictable, but not less
interesting.

Example:

```
- sinez: {note: ag, length: 32, period: 0.1, yshift: 0.9, xshift: 4}
```

will generate a lane like this:
```
- -- -- ag ag -- -- -- -- -- -- -- -- ag ag -- ag ag -- ag ag -- -- -- -- -- -- -- -- ag ag -- ag
```

# Listen

Demo patterns produced with droguedrums. The demo yaml files you can find in the
testfiles/ directory of the distribution archive.

[beat2.yml](beat2.yml) <audio src="beat2.mp3" controls></audio> (Sound: Vermona DRM1-MKIII)

[beat3.yml](beat3.yml) <audio src="beat3.mp3" controls></audio> (Sound: Vermona DRM1-MKIII)

[beat4.yml](beat4.yml) <audio src="beat4.mp3" controls></audio> (Sound: Vermona DRM1-MKIII)

[beat5.yml](beat5.yml) <audio src="beat5.mp3" controls></audio> (Sound: Vermona DRM1-MKIII)

[beat6.yml](beat6.yml) <audio src="beat6.mp3" controls></audio> (Sound: fluidsynth/FluidR3 GM2-2 soundfont)

[beat8.yml](beat8.yml) <audio src="beat8.mp3" controls></audio> (Sound: fluidsynth/FluidR3 GM2-2 soundfont)

# Controlling devices

droguedrums has been tested with a Vermona DRM1-MKIII drum synthesizer and
[fluidsynth](http://www.fluidsynth.org). In theory everything that can receive
MIDI events can be controlled, but the main focus of the program is rhythms.

Currently no midi note off events are ever sent. With some sounds that can lead
to notes playing forever. With percussive sounds this should not be a problem
usually.

Program change messages are not implemented, so you have to map notes to
different banks or midi channels.

# Workaround for libportmidi on Debian/Ubuntu

On Ubuntu (and probably on Debian) you will see the following error when
running "go get" for the portmidi bindings:

```
> go get github.com/rakyll/portmidi
# github.com/rakyll/portmidi
/usr/bin/ld: $WORK/github.com/rakyll/portmidi/_obj/portmidi.cgo2.o: undefined reference to symbol 'Pt_Start'
//usr/lib/libporttime.so.0: error adding symbols: DSO missing from command line
collect2: error: ld returned 1 exit status
```

This is because Ubuntu 16.04 does not provide the current (6 years old!)
version of libportmidi (217). A [bug
report](https://bugs.launchpad.net/ubuntu/+source/portmidi/+bug/1616384) has
been filed already.

You can work around that by editing the source code of the go bindings after
the error message above. Load the following two files in your editor:

  - $GOPATH/src/github.com/rakyll/portmidi/portmidi.go
  - $GOPATH/src/github.com/rakyll/portmidi/stream.go

and change the line that reads

```
  // #cgo LDFLAGS: -lportmidi
```

to this:

```
  // #cgo LDFLAGS: -lportmidi -lporttime
```

After this little change you can run `go get github.com/rakyll/portmidi` again
and it will use your changed copy of the source code to compile, now
sucessfully, the package.

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
