import time
import usb_cdc
import json
from adafruit_macropad import MacroPad
import io

macropad = MacroPad()
usb_cdc.data.timeout = 0.1


class SetLight:
  def __init__(self, macropad):
    self._macropad = macropad
  def apply(self, command):
      # TODO Animations, fade in/out
    self._macropad.pixels.brightness = 0.5
    if command == "on":
        self._macropad.pixels.fill((0, 255, 0))
    else:
        self._macropad.pixels.fill((255, 0, 0))

setlight = SetLight(macropad)


while True:
    key_event = macropad.keys.events.get()
    if key_event and key_event.pressed:
        print("hello")
        usb_cdc.data.write("{}".format(key_event.key_number).encode())
        usb_cdc.data.flush()
    macropad.encoder_switch_debounced.update()
    if macropad.encoder_switch_debounced.pressed:
        print("Pressed")
        usb_cdc.data.write("press!".encode())
        usb_cdc.data.flush()
    line = usb_cdc.data.readline()
    if len(line) > 0:
      print(line.decode())
      setlight.apply(line.decode())
    time.sleep(0.4)


