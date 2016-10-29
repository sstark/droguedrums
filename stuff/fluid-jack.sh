#!/bin/zsh

inputfile=$(mktemp /tmp/XXXXXXXXX)
print "\
prog 5 47
prog 6 113
prog 7 13
prog 10 11
gain 3
reverb off
channels
" >$inputfile
fluidsynth -a jack -m jack -j -o synth.verbose=yes -v "$@" /usr/share/sounds/sf2/FluidR3_GM.sf2 -f $inputfile
rm -f $inputfile
