---
date: 2016-10-21T10:35:20+02:00
menu:
    main:
        parent: Usage
title: Advanced concepts
weight: 80
---

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
