#!/usr/bin/env python

'''
prototype
'''

import pygame.midi
import time
import sys
import random

def playChord(notes):
    for n in notes:
        v = random.randint(100,127)
        #v = 127
        player.note_on(n, velocity=v, channel=5)

def stopChord(notes):
    for n in notes:
        player.note_off(n)

def stopAll():
    for n in range(1,9):
        player.note_off(n)

def step(bpm=120):
    time.sleep(60./bpm)

class MultiLane():

    def __init__(self):
        self.lanes = []
        self.bpm = 160
        self.kit = {
            '--': 0,
            'bd': 1,
            'lt': 2,
            'ht': 3,
            'mu': 4,
            'sd': 5,
            'ho': 6,
            'hc': 8,
            'co': 9,
            'cc': 11,
            'cl': 12
        }

    def __updateChords(self):
        self.chords = map(list, zip(*self.lanes))

    def play(self, count=1):
        for c in range(count):
            for chord in self.chords:
                print(chord)
                playChord(chord)
                step(self.bpm)

    def trig(self, count=1):
        self.play(count)
        self.flushLanes()

    def addLane(self, lane):
        self.lanes.append(lane)

    def parseLane(self, s):
        newLane = []
        for t in s.split():
            newLane.append(self.kit[t])
        self.addLane(newLane)

    def testtrack1(self):
        self.parseLane('bd -- -- -- -- bd -- --')
        self.parseLane('ho hc ho hc ho hc ho ho')
        self.parseLane('-- -- sd -- -- -- sd --')
        self.__updateChords()

    def testtrack2(self):
        self.addLane([11, 0, 0, 9, 0, 0, 0, 0])
        self.addLane([0, 2, 0, 2, 0, 0, 3, 0])
        self.addLane([0, 0, 12, 0, 0, 0, 0, 0])
        self.__updateChords()

    def testtrack3(self):
        self.addLane([1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0])
        self.addLane([6, 0, 6, 6, 8, 0, 6, 0, 6, 6, 6, 0])
        self.addLane([0, 0, 12, 0, 0, 2, 5, 0, 3, 0, 0, 2])
        self.__updateChords()

    def testtrack4(self):
        self.addLane([1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0])
        self.addLane([6, 0, 6, 6, 8, 0, 6, 0, 6, 6, 6, 0])
        self.addLane([0, 12, 0, 9, 12, 2, 5, 12, 3, 9, 12, 2])
        self.__updateChords()

    def testbreak(self):
        self.parseLane('bd -- bd bd')
        self.parseLane('hc hc hc hc')
        self.parseLane('sd sd -- sd')
        self.__updateChords()
        return self

    def flushLanes(self):
        print('flush')
        self.lanes = []


pygame.midi.init()
print(pygame.midi.get_default_output_id())
print(pygame.midi.get_device_info(2))
player = pygame.midi.Output(2)
player.set_instrument(0)

m = MultiLane()

time.sleep(0.2)

m.bpm=390
m.testbreak().trig()
#m.trig()

while True:
    try:
        m.testtrack3()
        m.trig(1)
        m.testtrack4()
        m.trig(1)
    except KeyboardInterrupt:
        offAll()
        pygame.midi.quit()
        raise
