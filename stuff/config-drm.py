#!/usr/bin/env python

'''
configure a Vermona DRM1-MKIII in learn mode
'''

import pygame.midi
import time
import sys

pygame.midi.init()

print(pygame.midi.get_default_output_id())
print(pygame.midi.get_device_info(2))
player = pygame.midi.Output(2)
player.set_instrument(0)

for n in [1, 2, 3, 4, 5, 8, 11, 12]:
    print(n)
    player.note_on(n, velocity=127, channel=5)
    time.sleep(0.5)
    player.note_off(n)

pygame.midi.quit()
