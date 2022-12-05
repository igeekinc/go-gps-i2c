package main

import (
	"fmt"
	gps2 "github.com/igeekinc/go-gps-i2c/pkg/gps"
	"log"
	"periph.io/x/conn/v3/pin"
	"periph.io/x/conn/v3/pin/pinreg"
)

func main() {
	gps, err := gps2.NewSerialGPSReader("/dev/ttyACM0", 9600, 8, 1)
	if err != nil {
		log.Fatal(err)
	}
	for {
		gga, err := gps.NextFix()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Time: %v Lat: %f Long: %f Alt: %f\n", gga.Time, gga.Latitude, gga.Longitude, gga.Altitude)
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
