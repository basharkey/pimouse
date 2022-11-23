# pimouse

# Build Only
```
make build
```

# Install
```
sudo make install
```

## Run
```
sudo pimouse
```

## Run as Service
```
systemctl start pimouse.service
```

### Start PiMouse on Boot
```
systemctl enable pimouse.service
```

# Config
The default config file is located at `/etc/pimouse/default.yaml`. This file will be used by any mouse connected that doesn't have a custom config.

## Custom Configs
Custom config files can be create for each mouse allowing you to have seperate configurations for each device. These files are placed in the same directory as the default config (`/etc/pimouse/`) and must follow this naming convention:
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
