# PiMouse
PiMouse allows you to remap mouse buttons for your mice by relaying mouse inputs through a Raspberry Pi 4 Model B. This project was created due to the Kensington Slimblade's lack of onboard memory. I got tired of having to configure Kensington's software on every system and wanted a device that would allow my settings to persist accross systems.

PiMouse acts as a proxy between your mouse and your computer. The program takes mouse inputs sent to the Pi, modifies them, and sends them to your host.
```
Mouse -> Raspberry Pi -> Computer
```

## Features
- Mutli mouse support
- Mouse button remapping
- Scroll speed configuration
- Mouse speed configuration (primitive)

## Requirements
For now it appears only the `Raspberry Pi 4 Model B` has the requirements for this project:
- supports USB gadget mode
- has USB type A port(s) to connect mice
- USB port that connects to the host (power and data)

The Pi Zero will not work, although is supports USB gadget mode it only has one USB port that provides power and data.

If you would like to use multiple mice with PiMouse you may need to add/use a more powerful power source with your Raspberry Pi. I use and can recommend [WaveShare's UPS](https://www.waveshare.com/wiki/UPS_HAT_(B)) however you also might be able to get away with a powered USB HUB.

## Setup
1. Setup Raspberry Pi
2. Install PiMouse
3. Connect Raspberry Pi to host over the Pi's USB type C port
4. Connect mice to Raspberry Pi's USB type A ports
5. Run PiMouse

## Install
Install
```
sudo apt install golang
sudo make install
```
Build only
```
sudo apt install golang
make build
```

## Run
Run from terminal
```
sudo pimouse
```
Run as Service
```
systemctl start pimouse.service
```
Start PiMouse on Boot
```
systemctl enable pimouse.service
```

# Config
The default config file is located at `/etc/pimouse/default.yaml`. This file will be used by any mouse connected that doesn't have a custom config.

## Custom Configs
Custom config files can be create for each mouse allowing you to have seperate configurations for each device. These files are placed in the same directory as the default config `/etc/pimouse/` and must follow this naming convention:
```
/etc/pimouse/<device-id>-event-mouse.yaml
```

You can get a list of device IDs by running:
```
ls -l /dev/input/by-id/ | cut -d' ' -f9 | grep event-mouse

usb-047d_Kensington_Slimblade_Trackball-event-mouse
```

Or simply monitor the output of `pimouse` when plugging in a mouse:
```
sudo pimouse

Connected "Kensington Slimblade Trackball"
  Device /dev/input/by-id/usb-047d_Kensington_Slimblade_Trackball-event-mouse
  Config /etc/pimouse/default.yaml
```

For example the config for my Kensington Slimblade is called:
```
/etc/pimouse/usb-047d_Kensington_Slimblade_Trackball-event-mouse.yaml
```

You can confirm your mouse is using the correct config by looking at the output of `pimouse`:
```
Connected "Kensington Slimblade Trackball"
  Device /dev/input/by-id/usb-047d_Kensington_Slimblade_Trackball-event-mouse
  Config /etc/pimouse/usb-047d_Kensington_Slimblade_Trackball-event-mouse.yaml
```

## Config Format
```
buttons:
  right: back
  back: right
  middle: forward

scrollSpeed: 1
cursorSpeed: 1
```

### Buttons
Use this section to remap mouse buttons. The first button for each entry is the button you would like to remap, the second is the action you would like to remap the button to.

For example this remaps right mouse button to back:
```
buttons:
  right: back
```

The following buttons are supported:
```
left
right
middle
back
forward
button10
button11
button12
```
