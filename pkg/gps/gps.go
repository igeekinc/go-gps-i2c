package gps

import (
	"github.com/adrianmo/go-nmea"
	"github.com/pkg/errors"
	"log"
)

// GPS is the end-user GPS structure to use.
type GPS interface {
	NextSentence() (nmea.Sentence, error)
	NextFix() (nmea.GGA, error)
	Readline() (string, error)
}

// GPSDevice is an interface to a GPS device.  This should be implemented for different types of devices and
// gives the ability to read/write NMEA strings with the device
// A specific implementation of a GPS device should implement GPSDevice and then composite GPSCore into its
// struct so that it provides a complete GPS intergace
type GPSDevice interface {
	Readline() (string, error)
}

// GPSCore is high level methods for GPS devices that uses a GPSDevice for low level access
type GPSCore struct {
	reader GPSDevice
}

func (g *GPSCore) NextSentence() (nmea.Sentence, error) {
	for {
		line, err := g.reader.Readline()
		if err != nil {
			return nil, errors.Wrap(err, "error in ReadLine")
		}
		sentence, err := nmea.Parse(line)
		if err == nil {
			return sentence, nil
		}
		log.Printf("Error parsing '%s', %v", line, err)
	}
}

func (g *GPSCore) NextFix() (nmea.GGA, error) {
	for {
		sentence, err := g.NextSentence()
		if err != nil {
			return nmea.GGA{}, errors.Wrap(err, "error in NextSentence")
		}
		gga, ok := sentence.(nmea.GGA)
		if ok {
			return gga, nil
		}
	}
}
