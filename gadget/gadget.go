package gadget

import (
    "os"
    "path/filepath"
    "os/exec"
)

var dirBase = "/sys/kernel/config/usb_gadget/pimouse"

var usbString = "0x409"
var usbConfig = "c.1"
var usbDevice = "hid.usb0"

var dirStrings = filepath.Join("strings", usbString)
var dirConfigs = filepath.Join("configs", usbConfig)
var dirFunctions = filepath.Join("functions", usbDevice)

func Initialize() error {

    var files = [][]string {
        {"idVendor", "0xbeaf"}, // Linux Foundation
        {"idProduct", "0x0104"}, // Multifunction Composite Gadget
        {"bcdDevice", "0x0100"}, // v1.0.0
        {"bcdUSB", "0x0200"}, // USB2

        {filepath.Join(dirStrings, "serialnumber"), "fedcba9876543210"},
        {filepath.Join(dirStrings, "manufacturer"), "PiMouse"},
        {filepath.Join(dirStrings, "product"), "USB Mouse Device"},

        {filepath.Join(dirConfigs, dirStrings, "configuration"), "Config 1: ECM network"},
        {filepath.Join(dirConfigs, "MaxPower"), "250"},

        {filepath.Join(dirFunctions, "protocol"), "1"},
        {filepath.Join(dirFunctions, "subclass"), "1"},
        {filepath.Join(dirFunctions, "report_length"), "8"},
        // sudo usbhid-dump -d beaf | tail -n +2 | xxd -r -p | hidrd-convert -o spec
        {filepath.Join(dirFunctions, "report_desc"),
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
        writeFile(filepath.Join(dirBase, file[0]), file[1])
    }

    link := filepath.Join(dirBase, dirFunctions)
    target := filepath.Join(dirBase, filepath.Join(dirConfigs, usbDevice))
    os.Symlink(link, target)

    dir, err := os.Open("/sys/class/udc/")
    if err != nil {
        return err
    }

    file, err := dir.Readdir(0)
    if err != nil {
        return err
    }

    err = writeFile(filepath.Join(dirBase, "UDC"), file[0].Name())
    if err != nil {
        return err
    }
    return nil
}

func Destroy() {
    var dirBase string = "/sys/kernel/config/usb_gadget/pimouse"
    // clear UDC file data, don't think there is a way to do this with pure go
    //https://askubuntu.com/questions/823380/cannot-delete-residual-system-files-period-even-after-changing-permissions-as-r
    cmd := exec.Command("/usr/bin/env", "bash", "-c", "echo '' > " + filepath.Join(dirBase, "UDC"))
    cmd.Run()

    var gadgetFiles =  []string {
        // remove strings from configs
        filepath.Join(dirBase, dirConfigs, dirStrings),
        // remove functions from configs
        filepath.Join(dirBase, dirConfigs, usbDevice),
        // remove configs
        filepath.Join(dirBase, dirConfigs),
        // remove functions
        filepath.Join(dirBase, dirFunctions),
        // remove strings
        filepath.Join(dirBase, dirStrings),
        // remove gadget
        dirBase,
    }

    for _, file := range gadgetFiles {
        os.Remove(file)
    }
}

func writeFile(filePath string, fileContent string) error {
    fileDir := filepath.Dir(filePath)
    _, err := os.Stat(fileDir)
    if os.IsNotExist(err) {
        err := os.MkdirAll(fileDir, 0644)
        if err != nil {
            return err
        }
    }

    file, err := os.Create(filePath)
    if err != nil {
        return err
    }

    _, err = file.WriteString(fileContent)
    if err != nil {
        return err
    }
    return nil
}
