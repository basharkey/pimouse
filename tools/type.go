package main

import (
    "fmt"
    "gadget"
    "os"
    "time"
    "log"
)

func main () {
    // initialize usb gadget
    gadget.Initialize()

    // open usb gadget device for write only
    gadget_device, err := os.OpenFile(
        "/dev/hidg0",
        os.O_WRONLY,
        0000,
    )
    check_err(err)
    defer gadget_device.Close()

    var val int8
    val = 0

    time.Sleep(2 * time.Second)
    type_bytes(gadget_device, []byte{255, 0, 0, byte(val)})
    type_bytes(gadget_device, []byte{0, 0, 0, 0})
}

func type_bytes(gadget_device *os.File, key_bytes []byte) {
    fmt.Println("typing:", key_bytes)
    //key_bytes = make([]byte, 8)
    _, err := gadget_device.Write(key_bytes)
	check_err(err)
}

func prepend_byte(x []byte, y byte) []byte {
    x = append(x, 0)
    copy(x[1:], x)
    x[0] = y
    return x
}

func check_err(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
