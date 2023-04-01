package gps

import (
	"github.com/pkg/errors"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/i2c"
	"time"
)

const DefaultAddr = 0x10

type Opts struct {
}

type I2CGPS struct {
	GPSCore
	Conn     conn.Conn
	lastByte byte
}

const BYTES_READY_REG uint8 = 0xFD
const READ_REG uint8 = 0xFF

// NewI2CGPS opens a handle to an i2c gps
func NewI2CGPS(bus i2c.BusCloser, opts *Opts) (*I2CGPS, error) {
	device := &I2CGPS{
		Conn: &i2c.Dev{Bus: bus, Addr: uint16(DefaultAddr)},
	}
	device.GPSCore.reader = device
	return device, nil
}

func (g *I2CGPS) sendCMD(cmd string) error {
	return g.Conn.Tx([]byte(cmd), nil)
}

func (g *I2CGPS) Readline() (string, error) {
	retStr := ""
	//var availableBuf [1]byte
	keepReading := true
	for keepReading {
		/*
			err := g.readReg(BYTES_READY_REG, availableBuf[0:1])
			if err != nil {
				return "", errors.Wrap(err, "Could not read available bytes")
			}

		*/
		//log.Printf("availableBuf[0] = %x, availableBuf[1] = %x\n", availableBuf[0], availableBuf[1])
		//available := int(availableBuf[0])*256 + int(availableBuf[1])
		available := 1
		//log.Printf("available = %d\n", available)
		goodBytes := 0
		if available > 0 {
			readBuf := make([]byte, available)
			err := g.readReg(READ_REG, readBuf)
			if err != nil {
				return "", errors.Wrap(err, "Could not read I2CGPS bytes")
			}
			for _, curByte := range readBuf {
				// The I2CGPS likes to return newlines for no data, if newline is not following CR, quash those here
				if curByte == '\n' && g.lastByte != '\r' {
					time.Sleep(500 * time.Millisecond) // GPS will only update once per second so a long sleep is good
					continue
				}
				g.lastByte = curByte
				if curByte != '\n' {
					retStr += string(curByte)
					goodBytes++
				} else {
					keepReading = false
				}
			}
		}
		/*
			if goodBytes == 0 {
				time.Sleep(10 * time.Millisecond)
			}

		*/
	}
	return retStr, nil
}

func (g *I2CGPS) readReg(reg uint8, b []byte) error {
	//log.Printf("Reading %d bytes from reg %x", len(b), reg)
	if err := g.Conn.Tx([]byte{reg}, b); err != nil {
		return errors.Wrap(err, "readReg failed")
	}
	return nil
}
