#!/bin/zsh

inputfile=$(mktemp $TMPDIR/XXXXXXXXX)
print "\
prog 5 47
prog 6 113
prog 7 13
channels
" >$inputfile
fluidsynth -a coreaudio -m coremidi ~/Music/FluidR3\ GM2-2.SF2 -f $inputfile
rm -f $inputfile
