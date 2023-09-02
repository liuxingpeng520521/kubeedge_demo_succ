package main

import (
	"os"
	"time"
	"github.com/stianeikeland/go-rpio/v4"
)

func TurnON(pinNumber int64) {
	pin := rpio.Pin(pinNumber)
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()

	pin.High()
}

func TurnOff(pinNumber int64) {
	pin := rpio.Pin(pinNumber)
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()
	pin.Low()

}

func main() {
	for {
		TurnON(18)
		time.Sleep(2 * time.Second)
		TurnOff(18)
		time.Sleep(2 * time.Second)
	}

}
