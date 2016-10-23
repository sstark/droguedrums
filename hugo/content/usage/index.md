---
date: 2016-10-22T15:39:45+02:00
menu: ""
title: How to use droguedrums
weight: 0
---

![matter](../b3.jpg)

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
