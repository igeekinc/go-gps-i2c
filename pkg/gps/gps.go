package gps

import (
	"github.com/adrianmo/go-nmea"
	"github.com/pkg/errors"
	"log"
)

type GPSReader interface {
	Readline() (string, error)
}

type GPSCore struct {
	reader GPSReader
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
