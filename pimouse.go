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
    gadgetDevice, err := os.OpenFile(
        "/dev/hidg0",
        os.O_WRONLY,
        0000,
    )
    if err != nil {
        log.Fatal(err)
    }

    defer gadgetDevice.Close()

    var miceHooked []*evdev.InputDevice
    for {
        micePaths, err := getMicePaths("/dev/input/by-id/")
        if err != nil {
            // errors when no usb devices plugged in
            fmt.Printf("\rWaiting no devices found...")
        } else {
            next:
            for _, mousePath := range micePaths {
                // don't hook mice that are already hooked
                for _, mouseHooked := range miceHooked {
                    if mousePath == mouseHooked.Fn {
                        continue next
                    }
                }

                mouseDevice, err := evdev.Open(mousePath)
                if err != nil {
                    fmt.Println(err)
                } else {
                    // load default.conf config if mouseDevice specific one does not exist
                    var configPath string
                    configDir := "/etc/pimouse"
                    configPathDefault := filepath.Join(configDir, "default.yaml")
                    configPathCustom := filepath.Join(configDir, filepath.Base(mousePath) + ".yaml")

                    _, err = os.Stat(configPathCustom)
                    if errors.Is(err, os.ErrNotExist) {
                        configPath = configPathDefault
                    } else {
                        configPath = configPathCustom
                    }
                    config, _ := config.Parse(configPath)

                    // track mice that are currently connected and hooked in miceHooked slice
                    miceHooked = append(miceHooked, mouseDevice)

                    fmt.Printf("\nConnected \"%s\"\n", mouseDevice.Name)
                    fmt.Printf("  Device %s\n", mouseDevice.Fn)
                    fmt.Printf("  Config %s\n", configPath)
                    // hook the mouse
                    go hookMouse(mouseDevice, config, gadgetDevice, &miceHooked)
                }
            }
        }
        time.Sleep(1 * time.Second)
    }
}

func hookMouse(mouseDevice *evdev.InputDevice, config config.MouseConfig, gadgetDevice *os.File, miceHooked *[]*evdev.InputDevice) {
    // main mouseDevice event loop
    mouseDevice.Grab()
    gadgetBytes := make([]byte, 4)
    for {
        // read events from mouseDevice
        mouseEvents, err := mouseDevice.Read()

        // if events cannot be read from mouseDevice remove mouse from miceHooked and exit goroutine
        if err != nil {
            for i, mouseHooked := range *miceHooked {
                if mouseDevice == mouseHooked {
                    (*miceHooked)[i] = (*miceHooked)[len(*miceHooked)-1]
                    *miceHooked = (*miceHooked)[:len(*miceHooked)-1]
                }
            }
            fmt.Printf("\nDisconnected \"%s\"\n", mouseDevice.Name)
            return
        }

        for _, mouseEvent := range mouseEvents {
            // button event
            if mouseEvent.Type == 1 {
                if button_byte, ok := config.ButtonMap[mouseEvent.Code]; ok {
                    if mouseEvent.Value == 1 {
                        gadgetBytes[0] = button_byte
                    } else {
                        gadgetBytes[0] = byte(0)
                    }
                }
            // movement event
            } else if mouseEvent.Type == 2 {
                // x axis
                if mouseEvent.Code == 0 {
                    gadgetBytes[1] = byte(mouseEvent.Value * int32(config.CursorSpeed))
                    gadgetBytes[2] = byte(0)
                // y axis
                } else if mouseEvent.Code == 1 {
                    gadgetBytes[2] = byte(mouseEvent.Value * int32(config.CursorSpeed))
                    gadgetBytes[1] = byte(0)
                // scroll event
                } else if mouseEvent.Code == 8 {
                    gadgetBytes[3] = byte(mouseEvent.Value * int32(config.ScrollSpeed))
                }
            } else if mouseEvent.Type == 0 {
                gadgetBytes[1] = byte(0)
                gadgetBytes[2] = byte(0)
                gadgetBytes[3] = byte(0)
            }
            //fmt.Println("event:", mouseEvent.Type, mouseEvent.Code, mouseEvent.Value, "gadget:", gadgetBytes)
            gadgetDevice.Write(gadgetBytes)
        }
    }
}

func getMicePaths(basePath string) ([]string, error) {
    dir, err := os.Open(basePath)
    if err != nil {
        return nil, err
    }

    devices, err := dir.Readdir(0)
    if err != nil {
        return nil, err
    }

    var micePaths []string
    for _, device := range devices {
        if strings.Contains(device.Name(), "event-mouse") {
            micePaths = append(micePaths, basePath + device.Name())
        }
    }
    return micePaths, nil
}
