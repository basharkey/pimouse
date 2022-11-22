package main

import (
    "fmt"
    "os"
    "github.com/gvalkov/golang-evdev"
    "strings"
    "time"
    "log"
    "path/filepath"
    "errors"
    "gadget"
    "config"
)

func main() {
    // initialize usb gadget
    gadget.Initialize()

    // open usb gadget device for write only
    gadget_device, err := os.OpenFile(
        "/dev/hidg0",
        os.O_WRONLY,
        0000,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer gadget_device.Close()

    var mice_hooked []string
    for {
        // I don't think this will ever error for no mice being plugged in
        // errors would probably be related to permissions issues
        mice_paths, err := get_mice_paths("/dev/input/by-id/")
        if err != nil {
            log.Fatal(err)
        } else {
            next:
            for _, mouse_path := range mice_paths {
                // don't hook mice that are already hooked
                for _, mouse_hooked := range mice_hooked {
                    if mouse_path == mouse_hooked {
                        continue next
                    }
                }

                mouse_device, err := evdev.Open(mouse_path)
                if err != nil {
                    fmt.Println(err)
                } else {
                    // load default.conf config if mouse_device specific one does not exist
                    var mouse_config string
                    config_dir := "/etc/pimouse"
                    default_mouse_config := filepath.Join(config_dir, "default.yaml")
                    custom_mouse_config := filepath.Join(config_dir, filepath.Base(mouse_path) + ".yaml")

                    _, err = os.Stat(custom_mouse_config)
                    if errors.Is(err, os.ErrNotExist) {
                        fmt.Println("Using default config can't find: ", custom_mouse_config)
                        mouse_config = default_mouse_config
                    } else {
                        mouse_config = custom_mouse_config
                    }
                    mouseConfig, _ := config.Parse(mouse_config)

                    // track mice that are currently connected and hooked in mice_hooked slice
                    mice_hooked = append(mice_hooked, mouse_path)
                    // hook the mouse
                    //go hook_mouse(mouse_device, mouse_config, gadget_device, mouse_path, &mice_hooked)
                    go hook_mouse(mouse_device, mouseConfig, gadget_device)
                }
            }
        }
        time.Sleep(1 * time.Second)
    }
}

func hook_mouse(mouse_device *evdev.InputDevice, mouseConfig config.MouseConfig, gadget_device *os.File) error {
    fmt.Println(mouse_device)

    // main mouse_device event loop
    mouse_device.Grab()
    gadget_bytes := make([]byte, 4)
    for {
        // check if events can be read from mouse_device (if mouse is still connected)
        mouse_events, err := mouse_device.Read()
        if err != nil {
            // remove mouse from hooked mice if it has been disconnected
            return err
        }

        for _, mouse_event := range mouse_events {
            // button event
            if mouse_event.Type == 1 {
                if button_byte, ok := mouseConfig.ButtonMap[mouse_event.Code]; ok {
                    if mouse_event.Value == 1 {
                        gadget_bytes[0] = button_byte
                    } else {
                        gadget_bytes[0] = byte(0)
                    }
                }
            // movement event
            } else if mouse_event.Type == 2 {
                // x axis
                if mouse_event.Code == 0 {
                    gadget_bytes[1] = byte(mouse_event.Value)
                    gadget_bytes[2] = byte(0)
                // y axis
                } else if mouse_event.Code == 1 {
                    gadget_bytes[2] = byte(mouse_event.Value)
                    gadget_bytes[1] = byte(0)
                // scroll event
                } else if mouse_event.Code == 8 {
                    gadget_bytes[3] = byte(mouse_event.Value)
                }
            } else if mouse_event.Type == 0 {
                gadget_bytes[1] = byte(0)
                gadget_bytes[2] = byte(0)
                gadget_bytes[3] = byte(0)
            }
            fmt.Println("raw:", mouse_event.Type, mouse_event.Code, mouse_event.Value, "gadget:", gadget_bytes)
            gadget_device.Write(gadget_bytes)
        }
    }
}

func get_mice_paths(base_path string) ([]string, error) {
    dir, err := os.Open(base_path)
    if err != nil {
        return nil, err
    }

    devices, err := dir.Readdir(0)
    if err != nil {
        return nil, err
    }

    var mice_paths []string
    for _, device := range devices {
        if strings.Contains(device.Name(), "event-mouse") {
            mice_paths = append(mice_paths, base_path + device.Name())
        }
    }
    return mice_paths, nil
}
