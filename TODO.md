todos

- definition file monitoring
  - use named pipes instead?
- seqs can refer to other seqs
- looping of parts/seqs
- command line parsing
  - midi device selection
    - https://github.com/rakyll/launchpad/blob/master/launchpad.go
- other midi libs
  - jack
  - coremidi
- velocity lanes
- more tests
- combination of parts
  - merge lanes if same size
  - cut/expand if not the same size?
  - how to merge fx?
- rework yaml parsing, it can probably be made easier
- try to help with portmidi bindings issues:
  - https://github.com/rakyll/portmidi/issues/11
  - https://github.com/rakyll/portmidi/issues/27
- make it possible to set the "step" value in parts
  to triplets (or any other fraction)
- detect loops when flattening seqs
- handle case when parts have no step value (currently
  leads to div by zero)

ideas

- bpm ramps
- fancy console output
  - https://github.com/gizak/termui
  - https://github.com/gosuri/uilive
- microtiming map
  - enable per lane dragging/rushing
  - random timing "errors"
- delay effect
- sequence fx
  - song velocity curve
  - algorithmic part selection

wishful thinking

- dynamic vim mode with part pointer
