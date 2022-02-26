
![droguetrees](assets/droguetrees.jpg)

droguedrums is a text file based drum sequencer for the command line.

It reads a definition of drum sets, parts and sequences from a text file and
plays it back as midi. No sounds are provided. Development is done with the
following workflow in mind:

  - define rhythmical patterns using a simple text file
  - feed midi data to a hardware drum synthesizer
  - process and record the audio from hardware

There is a single input file containing definitions for drum kits, parts and
sequences. The file is specified at the command line.

![matter](assets/b2.jpg)

For the easiest possible start just download the source and one of the
precompiled binaries and run one of the contained example files:

```shell
$ ./droguedrums testfiles/beat1.yml
```

## Download

Choose a download from [Releases](releases) and unzip.

## Requirements for running

Requirements for running the binary:

  - portmidi library
    - OSX:      `brew install portmidi` (not needed if compiled with coremidi, which is the default)
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

## Workaround for libportmidi on Debian/Ubuntu

On Ubuntu (and probably on Debian) you will see the following error when
running "go get" for the portmidi bindings:

```shell
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

![matter](assets/b3.jpg)

Start the program like this:

```shell
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

```shell
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

The YAML file format seems like a good compromise between easy machine
readability and human editing. For this reason, but also for ease of
development (making a program read yaml is simple), YAML was chosen as the
input file format for drougedrums.

The following paragraphs will give an overview of what things you can put into
a droguedrums input file. See the full example files provided with the download
for some possibilities.

In the examples section of this web site you will find fully working files
linked from the individual examples.

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

~~~yaml
sets:
    - name: default
      kit:
          - {key: --, channel:  5, note:  0} # pause
          - {key: bd, channel:  5, note:  1} # bass drum
          - {key: hc, channel:  5, note:  2} # high hat closed
          - {key: sd, channel:  5, note:  5} # snare drum
~~~

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

```yaml
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

```yaml
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

### Part FX

Parts can have an fx section that will modify the events in the part. To use
multiple effects you write them as a list like this:

```yaml
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

```yaml
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

#### Alternating figures

In a figure you can specify alternative values for the "key" parameter. "|" is
used as the separation character. Each time the figure is used, another
alternative will be randomly chosen from the given ones. This can enrich the
beat a lot if used sensibly.

Example:

```yaml
figures:
    - {name: f3, key: fl|ce|cb, velocity: 80, pattern: x}
```

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

```yaml
- sinez: {note: ag, length: 32, period: 0.1, yshift: 0.9, xshift: 4}
```

will generate a lane like this:
```yaml
- -- -- ag ag -- -- -- -- -- -- -- -- ag ag -- ag ag -- ag ag -- -- -- -- -- -- -- -- ag ag -- ag
```

#### place

Usage: `- place {note: <key>, pos: p1 p2 ...}`

Sometimes all you need is place an event at the 15th beat of a bar. This can be
tedious if written out manually. With place you can do just that: It will place
<key> at the positions p1, p2, ... in the generated lane, filling the rest with
pauses. The individual positions (p1, p2, ...) need to be separated by _spaces_
and strictly increasing in the definition.

Example:

```yaml
- place: {note: fl, pos: 1 5 11}
```

will generate a lane like this:
```yaml
- fl -- -- -- fl -- -- -- -- -- fl
```

#### euclid

Usage: `- euclid {note: <key>, length: <numlen>, accents: <numacc>, rotation: <rot>}`

Uses Björklund's algorithm to generate "euclidean" rhythmical patterns of
length *<numlen>* with *<numacc>* accents. The whole pattern can be shifted by
supplying a rotation value *<rot>*.

Example:

```yaml
- euclid: {note: xy, length: 8, accents: 3}
```

will generate a lane like this (the Cuban "tresillo"):
```yaml
- xy -- -- xy -- -- xy -
```

droguedrums has been tested with a Vermona DRM1-MKIII drum synthesizer and
[fluidsynth](http://www.fluidsynth.org). In theory everything that can receive
MIDI events can be controlled, but the main focus of the program is rhythms.

Currently no midi note off events are ever sent. With some sounds that can lead
to notes playing forever. With percussive sounds this should not be a problem
usually.

Program change messages are not implemented, so you have to map notes to
different banks or midi channels.

![houses](assets/b1.jpg)

Demo patterns produced with droguedrums. The demo yaml files you can find in the
testfiles/ directory of the distribution archive.

[beat2.yml](testfiles/beat2.yml) <audio src="../beat2.mp3" controls></audio> (Sound: Vermona DRM1-MKIII)

[beat3.yml](testfiles/beat3.yml) <audio src="../beat3.mp3" controls></audio> (Sound: Vermona DRM1-MKIII)

[beat4.yml](testfiles/beat4.yml) <audio src="../beat4.mp3" controls></audio> (Sound: Vermona DRM1-MKIII)

[beat5.yml](testfiles/beat5.yml) <audio src="../beat5.mp3" controls></audio> (Sound: Vermona DRM1-MKIII)

[beat6.yml](testfiles/beat6.yml) <audio src="../beat6.mp3" controls></audio> (Sound: fluidsynth/FluidR3 GM2-2 soundfont)

[beat8.yml](testfiles/beat8.yml) <audio src="../beat8.mp3" controls></audio> (Sound: fluidsynth/FluidR3 GM2-2 soundfont)

