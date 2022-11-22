package gadget

import (
    "os"
    //"log"
    "fmt"
    "path/filepath"
    "os/exec"
)

var base_dir = "/sys/kernel/config/usb_gadget/pimouse"

var usb_string = "0x409"
var usb_config = "c.1"
var usb_device = "hid.usb0"

var strings_dir = filepath.Join("strings", usb_string)
var configs_dir = filepath.Join("configs", usb_config)
var functions_dir = filepath.Join("functions", usb_device)

func Initialize() {

    var files = [][]string {
        {"idVendor", "0xbeaf"}, // Linux Foundation
        {"idProduct", "0x0104"}, // Multifunction Composite Gadget
        {"bcdDevice", "0x0100"}, // v1.0.0
        {"bcdUSB", "0x0200"}, // USB2

        {filepath.Join(strings_dir, "serialnumber"), "fedcba9876543210"},
        {filepath.Join(strings_dir, "manufacturer"), "PiMouse"},
        {filepath.Join(strings_dir, "product"), "USB Mouse Device"},

        {filepath.Join(configs_dir, strings_dir, "configuration"), "Config 1: ECM network"},
        {filepath.Join(configs_dir, "MaxPower"), "250"},

        {filepath.Join(functions_dir, "protocol"), "1"},
        {filepath.Join(functions_dir, "subclass"), "1"},
        {filepath.Join(functions_dir, "report_length"), "8"},
        // sudo usbhid-dump -d beaf | tail -n +2 | xxd -r -p | hidrd-convert -o spec
        {filepath.Join(functions_dir, "report_desc"),
            "\x05\x01" + // Usage Page (Desktop)
            "\x09\x02" + // Usage (Mouse)
            "\xa1\x01" + // Collection (Application)
            "\x09\x01" + // Usage Page (Pointer)
            "\xa1\x00" + // Collection (Physical)
            // mouse buttons
            "\x05\x09" + // Usage Page (Button)
            "\x19\x01" + // Usage Minimum (Button 1)
            "\x29\x08" + // Usage Maximum (Button 8)
            "\x15\x00" + // Logical Minimum (0)
            "\x25\x01" + // Logical Maximum (1)
            "\x95\x08" + // Report Count (8)
            "\x75\x01" + // Report Size (1)
            "\x81\x02" + // Input (Data, Var, Abs)
            // mouse movement
            "\x05\x01" + // Usage Page (Desktop)
            "\x09\x30" + // Usage (X)
            "\x09\x31" + // Usage (Y)
            "\x15\x81" + // Logical Minimum (-127)
            "\x25\x7f" + // Logical Maximum (127)
            "\x75\x08" + // Report Size (8)
            "\x95\x02" + // Report Count (2)
            "\x81\x06" + // Input (Data, Var, Rel)
            // mouse vertical scroll resolution multipler
            "\xa1\x02" + //  Collection (Logical)
            "\x09\x48" + //  Usage (Resolution Multiplier)
            "\x15\x00" + //  Logical Minimum (0)
            "\x25\x01" + //  Logical Maximum (1)
            "\x35\x01" + //  Physical Minimum (1)
            "\x45\x04" + //  Physical Maximum (4)
            "\x75\x02" + //  Report Size (2)
            "\x95\x01" + //  Report Count (1)
            "\xa4"     + //  Push
            "\xb1\x02" + //  Feature (Data, Var, Abs)
            // mouse vertical scroll
            "\x09\x38" + // Usage (Wheel)
            "\x15\x81" + // Logical Minimum (-127)
            "\x25\x7f" + // Logical Maximum (127)
            "\x35\x00" + // Physical Minimum (0)
            "\x45\x00" + // Physical Maximum (0)
            "\x75\x08" + // Report Size (8)
            "\x81\x06" + // Input (Data, Var, Rel)
            "\xc0"     + // End Collection
            "\xc0"     + // End Collection
            "\xc0"},     // End Collection
    }

    for _, file := range files {
        write_to_file(filepath.Join(base_dir, file[0]), file[1])
    }

    link := filepath.Join(base_dir, functions_dir)
    target := filepath.Join(base_dir, filepath.Join(configs_dir, usb_device))
    os.Symlink(link, target)

    dir, err := os.Open("/sys/class/udc/")
    check_err(err)
    file, err := dir.Readdir(0)
    check_err(err)
    write_to_file(filepath.Join(base_dir, "UDC"), file[0].Name())
}

func Destroy() {
    var base_dir string = "/sys/kernel/config/usb_gadget/pimouse"
    // clear UDC file data, don't think there is a way to do this with pure go
    //https://askubuntu.com/questions/823380/cannot-delete-residual-system-files-period-even-after-changing-permissions-as-r
    cmd := exec.Command("/usr/bin/env", "bash", "-c", "echo '' > " + filepath.Join(base_dir, "UDC"))
    cmd.Run()

    var gadget_files =  []string {
        // remove strings from configs
        filepath.Join(base_dir, configs_dir, strings_dir),
        // remove functions from configs
        filepath.Join(base_dir, configs_dir, usb_device),
        // remove configs
        filepath.Join(base_dir, configs_dir),
        // remove functions
        filepath.Join(base_dir, functions_dir),
        // remove strings
        filepath.Join(base_dir, strings_dir),
        // remove gadget
        base_dir,
    }

    for _, file := range gadget_files {
        os.Remove(file)
    }
}

func check_err(err error) {
    if err != nil {
        //log.Fatal(err)
        fmt.Println(err)
    }
}

func write_to_file(file_path string, file_content string) {
    file_dir := filepath.Dir(file_path)
    _, err := os.Stat(file_dir)
    if os.IsNotExist(err) {
        err := os.MkdirAll(file_dir, 0644)
        check_err(err)
    }

    file, err := os.Create(file_path)
    check_err(err)
    _, err = file.WriteString(file_content)
    check_err(err)
}
