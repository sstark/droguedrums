---
date: 2016-10-21T10:35:20+02:00
menu:
    main:
        parent: Usage
title: Writing input files
weight: 10
---

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
