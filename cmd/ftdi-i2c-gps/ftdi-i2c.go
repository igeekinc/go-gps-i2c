package main

import (
	"fmt"
	"log"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/pin"
	"periph.io/x/conn/v3/pin/pinreg"
	"periph.io/x/host/v3"

	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3/ftdi"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	all := ftdi.All()
	if len(all) == 0 {
		log.Fatal("found no FTDI device on the USB bus")
	}

	// Use channel A.
	ft232h, ok := all[0].(*ftdi.FT232H)
	if !ok {
		log.Fatal("no FTDI device on the USB bus")
	}

	bus, err := ft232h.I2C(gpio.Float)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Found I2C bus on ftdi")

	if p, ok := bus.(i2c.Pins); ok {
		printPin("SCL", p.SCL())
		printPin("SDA", p.SDA())
	}

	i2cAddr := 0x76

	if dev, err = bmxx80.NewI2C(i, uint16(*i2cAddr), &opts); err != nil {
		return err
	}
	if err := bus.Close(); err != nil {
		log.Fatal(err)
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
