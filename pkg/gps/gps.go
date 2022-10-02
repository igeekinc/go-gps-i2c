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
	Conn conn.Conn
	buf  []byte
}

// NewGPS opens a handle to an i2c .
func NewGPS(bus i2c.BusCloser, opts *Opts) (*GPS, error) {
	device := &GPS{
		Conn: &i2c.Dev{Bus: bus, Addr: uint16(DefaultAddr)},
		buf:  make([]byte, 0, 256),
	}
	return device, nil
}

func (g *GPS) sendCMD(cmd string) error {
	return g.Conn.Tx([]byte(cmd), nil)
}

func (g *GPS) readline() (string, error) {
	retStr := ""
	keepReading := true
	for keepReading {
		for len(g.buf) == 0 {
			err := g.Conn.Tx(nil, g.buf)
			if err != nil {
				return retStr, err
			}
			if len(g.buf) == 0 {
				// Didn't get anything, sleep for a little bit
				time.Sleep(time.Millisecond * 100)
			}
		}
		for len(g.buf) > 0 && g.buf[0] != '\n' {
			retStr += string(g.buf[0])
			g.buf = g.buf[1:]
		}
		if len(g.buf) > 0 && g.buf[0] == '\n' {
			g.buf = g.buf[1:]
			keepReading = false
		}
	}
	return retStr, nil
}
