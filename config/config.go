package config

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v3"
)

type MouseConfig struct {
    ButtonMap map[uint16]byte
    ScrollSpeed int
}

var buttonCodes = map[uint16]byte {
    272: 1,     // left
    273: 2,     // right
    274: 4,     // middle
    275: 8,     // back
    276: 16,    // forward
    277: 32,    // 10
    278: 64,    // 11
    279: 128,   // 12
}

var buttonNames = map[string]uint16 {
    "left": 272,        // left
    "right": 273,       // right
    "middle": 274,      // middle
    "back": 275,        // back
    "forward": 276,     // forward
    "button10": 277,    // 10
    "button11": 278,    // 11
    "button12": 279,    // 12
}

func Parse(configPath string) (MouseConfig, error) {
    configYAML, err := ioutil.ReadFile(configPath)
    if err != nil {
        return MouseConfig{}, err
    }

    configData := make(map[interface{}]interface{})
    yaml.Unmarshal(configYAML, &configData)

    // initialize mouseConfig
    mouseConfig := MouseConfig{}
    // add default ButtonMap to mouseConfig
    mouseConfig.ButtonMap = buttonCodes

    for key, value := range configData {
        // modify ButtonMap
        if key == "buttons" {
            for input, output := range value.(map[string]interface{}) {
                mouseConfig.ButtonMap[buttonNames[input]] = buttonCodes[buttonNames[output.(string)]]
            }
        // modify ScrollSpeed
        } else if key == "scrollSpeed" {
            mouseConfig.ScrollSpeed = value.(int)
        }
    }
    return mouseConfig, nil
}
