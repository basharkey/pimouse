# TODO
## fix increased mouse lag with high polling rate mice
This is due to how long golang-evdev takes to read mouse events. To work well with 1000 Hz mice it needs to take 1 ms or less to process and send mouse events. Currently it takes 6-8 ms just to read a mouse movement event.

Perhaps I should write my own code to read events? Maybe I could make it faster?
