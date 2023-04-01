package gps

import (
	"bufio"
	"github.com/jacobsa/go-serial/serial"
	"github.com/pkg/errors"
	"io"
	"log"
)

type SerialGPS struct {
	GPSCore
	port       string
	baudRate   uint
	dataBits   uint
	stopBits   uint
	serialPort io.ReadWriteCloser
	scanner    *bufio.Scanner
}

func NewSerialGPSReader(port string, baudRate uint, dataBits uint, stopBits uint) (gps *SerialGPS, err error) {
	log.Printf("New SerialGPS Reader creating port=%s, baudRate = %d, dataBits = %d, stopBits = %d",
		port, baudRate, dataBits, stopBits)
	serialGPS := &SerialGPS{
		port:     port,
		baudRate: baudRate,
		dataBits: dataBits,
		stopBits: stopBits,
	}
	serialGPS.GPSCore.reader = serialGPS

	options := serial.OpenOptions{
		PortName:        port,
		BaudRate:        baudRate,
		DataBits:        dataBits,
		StopBits:        stopBits,
		MinimumReadSize: 4,
	}

	serialPort, err := serial.Open(options)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open serial port %s", port)
	}
	serialGPS.serialPort = serialPort
	reader := bufio.NewReader(serialPort)
	scanner := bufio.NewScanner(reader)
	serialGPS.scanner = scanner
	return serialGPS, nil
}

func (g *SerialGPS) Readline() (string, error) {
	g.scanner.Scan()
	return g.scanner.Text(), nil
}
