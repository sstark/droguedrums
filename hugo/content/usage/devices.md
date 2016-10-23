---
date: 2016-10-21T10:35:20+02:00
menu:
    main:
        parent: Usage
title: Controlling devices
weight: 100
---

droguedrums has been tested with a Vermona DRM1-MKIII drum synthesizer and
[fluidsynth](http://www.fluidsynth.org). In theory everything that can receive
MIDI events can be controlled, but the main focus of the program is rhythms.

Currently no midi note off events are ever sent. With some sounds that can lead
to notes playing forever. With percussive sounds this should not be a problem
usually.

Program change messages are not implemented, so you have to map notes to
different banks or midi channels.

