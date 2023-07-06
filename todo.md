# TODO
## fix increased mouse lag with high polling rate mice
This is due to how long golang-evdev takes to read mouse events. To work well with 1000 Hz mice it needs to take 1 ms or less to process and send mouse events. Currently it takes 6-8 ms just to read a mouse movement event.

Perhaps I should write my own code to read events? Maybe I could make it faster?

Write something like evhz in golang to see how fast we can read mouse events?
Can we reach 1000hz?


## Actual issue
Reading isnt the issue, pi can easily read 1000hz mice device
the problem is writing the forwarded packets to the gadget device file.

When pimouse is reading/writing it can handle 750hz
If its just reading and performing all other logic other than writing it can get 1000hz easily
