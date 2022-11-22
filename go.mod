module github.com/basharkey/pimouse

go 1.15

require (
	config v1.0.0
	gadget v1.0.0
	github.com/gvalkov/golang-evdev v0.0.0-20220815104727-7e27d6ce89b6
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace gadget v1.0.0 => ./gadget

replace config v1.0.0 => ./config
