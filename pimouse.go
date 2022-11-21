package main

import (
    "fmt"
    "os"
    "github.com/gvalkov/golang-evdev"
    "strings"
    "time"
    "log"
    "gadget"
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
                    /*
                    // load default.conf config if mouse_device specific one does not exist
                    var mouse_config string
                    config_dir := "/etc/pimk"
                    default_mouse_config := filepath.Join(config_dir, "default.conf")
                    custom_mouse_config := filepath.Join(config_dir, filepath.Base(mouse_path) + ".conf")

                    _, err = os.Stat(custom_mouse_config)
                    if errors.Is(err, os.ErrNotExist) {
                        fmt.Println("Using default config can't find: ", custom_mouse_config)
                        mouse_config = default_mouse_config
                    } else {
                        mouse_config = custom_mouse_config
                    }
                    */

                    // track mice that are currently connected and hooked in mice_hooked slice
                    mice_hooked = append(mice_hooked, mouse_path)
                    // hook the mouse
                    //go hook_mouse(mouse_device, mouse_config, gadget_device, mouse_path, &mice_hooked)
                    go hook_mouse(mouse_device, gadget_device)
                }
            }
        }
        time.Sleep(1 * time.Second)
    }
}

//func hook_mouse(mouse_device *evdev.InputDevice, mouse_config string, gadget_device *os.File, mouse_path string, mice_hooked *[]string) error {
func hook_mouse(mouse_device *evdev.InputDevice, gadget_device *os.File) error {
    fmt.Println(mouse_device)

    // main mouse_device event loop
    mouse_device.Grab()
    for {
        // check if events can be read from mouse_device (if mouse is still connected)
        mouse_events, err := mouse_device.Read()
        if err != nil {
            // remove mouse from hooked mice if it has been disconnected
            return err
        }

        for _, mouse_event := range mouse_events {
            /*
            mouse_event.Type
            mouse_event.Code
            mouse_event.Value
            */
            //1 2 4 8 16 32 64 128
            //[4, 0, 0] middle mouse (2)
            //[1, 0, 0] left mouse (1)
            //[2, 0, 0] right mouse (3)

            fmt.Println(mouse_event.Type, mouse_event.Code, mouse_event.Value)
            if mouse_event.Code == 0 {
                // x axis
                type_bytes(gadget_device, []byte{0, uint8(mouse_event.Value), 0})
            } else if mouse_event.Code == 1 {
                // y axis
                type_bytes(gadget_device, []byte{0, 0, uint8(mouse_event.Value)})
            }
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

func type_bytes(gadget_device *os.File, gadget_bytes []byte) {
    fmt.Println("sending:", gadget_bytes)
    //gadget_bytes = make([]byte, 8)
    //_, err := gadget_device.Write(gadget_bytes)
    gadget_device.Write(gadget_bytes)
}
