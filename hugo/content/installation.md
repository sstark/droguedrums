---
date: 2016-10-21T09:48:15+02:00
menu: main
title: Installation
weight: 10
---

![matter](../b2.jpg)

For the easiest possible start just download the source and one of the
precompiled binaries and run one of the contained example files:

```shell
$ ./droguedrums testfiles/beat1.yml
```

## Download

Choose a download from [Releases](../releases) and unzip.

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
