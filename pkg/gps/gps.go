package gps

import (
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/i2c"
	"time"
)

const DefaultAddr = 0x10

type Opts struct {
}

type GPS struct {
	Conn     conn.Conn
	lastByte byte
}

// NewGPS opens a handle to an i2c .
func NewGPS(bus i2c.BusCloser, opts *Opts) (*GPS, error) {
	device := &GPS{
		Conn: &i2c.Dev{Bus: bus, Addr: uint16(DefaultAddr)},
	}
	return device, nil
}

func (g *GPS) sendCMD(cmd string) error {
	return g.Conn.Tx([]byte(cmd), nil)
}

func (g *GPS) Readline() (string, error) {
	retStr := ""
	byteBuf := make([]byte, 1)
	keepReading := true
	for keepReading {
		err := g.Conn.Tx(nil, byteBuf)
		if err != nil {
			return retStr, err
		}
		curByte := byteBuf[0]
		// The GPS likes to return newlines for no data, if newline is not following CR, quash those here
		if curByte == '\n' && g.lastByte != '\r' {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		g.lastByte = curByte
		if curByte != '\n' {
			retStr += string(curByte)
		} else {
			keepReading = false
		}
	}
	return retStr, nil
}
