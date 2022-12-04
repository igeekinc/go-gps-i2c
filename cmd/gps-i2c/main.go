package main

import (
	"fmt"
	gps2 "github.com/igeekinc/go-gps-i2c/pkg/gps"
	"log"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/pin"
	"periph.io/x/conn/v3/pin/pinreg"
	"periph.io/x/host/v3"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("1")

	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	log.Println("Found I2C bus on ftdi")

	gps, err := gps2.NewGPS(bus, nil)
	for {
		line, err := gps.Readline()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(line)
		//time.Sleep(time.Second)
	}
}

func printPin(fn string, p pin.Pin) {
	name, pos := pinreg.Position(p)
	if name != "" {
		fmt.Printf("  %-3s: %-10s found on header %s, #%d\n", fn, p, name, pos)
	} else {
		fmt.Printf("  %-3s: %-10s\n", fn, p)
	}
}
