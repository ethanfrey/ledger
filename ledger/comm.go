package ledger

import (
	"errors"
	"fmt"

	"github.com/flynn/hid"
)

const (
	LEDGER = 0x2c97
	NANO   = 1
)

type Ledger struct {
	device Device
}

func NewLedger(dev Device) *Ledger {
	return &Ledger{
		device: dev,
	}
}

func FindLedger() (*Ledger, error) {
	devs, err := hid.Devices()
	if err != nil {
		return nil, err
	}
	for _, d := range devs {
		// TODO: ProductId filter
		if d.VendorID == LEDGER {
			ledger, err := d.Open()
			if err != nil {
				return nil, err
			}
			return NewLedger(ledger), nil
		}
	}
	return nil, errors.New("no ledger connected")
}

// A Device provides access to a HID device.
type Device interface {
	// Close closes the device and associated resources.
	Close()

	// Write writes an output report to device. The first byte must be the
	// report number to write, zero if the device does not use numbered reports.
	Write([]byte) error

	// ReadCh returns a channel that will be sent input reports from the device.
	// If the device uses numbered reports, the first byte will be the report
	// number.
	ReadCh() <-chan []byte

	// ReadError returns the read error, if any after the channel returned from
	// ReadCh has been closed.
	ReadError() error
}

func (l *Ledger) Exchange(command []byte, timeout int) ([]byte, error) {
	adpu := WrapCommandAPDU(0x0101, command, 64, false)

	// write all the packets
	err := l.device.Write(adpu[:64])
	if err != nil {
		return nil, err
	}
	for len(adpu) > 64 {
		adpu = adpu[64:]
		err = l.device.Write(adpu[:64])
		if err != nil {
			return nil, err
		}
	}

	input := l.device.ReadCh()
	for msg := range input {
		fmt.Println("** message **")
		fmt.Printf("%X\n\n", msg)
	}
	return nil, nil

	// // dataLength = 0
	// // dataStart = 2

	// // TODO: add timeout
	// // result = self.waitImpl.waitFirstResponse(timeout)
	// read := l.Device.Read(65)
	// response, err := unwrapResponseAPDU(0x0101, result, 64)
	// // error on first read is a problem
	// if err != nil {
	// 	return nil, err
	// }

	// // otherwise, continue til the end...
	// // self.device.set_nonblocking(False)
	// var result []byte
	// for err != nil {
	// 	result = append(result, response...)
	// 	read = l.Device.Read(65)
	// 	response, err = unwrapResponseAPDU(0x0101, result, 64)
	// }
	// // self.device.set_nonblocking(True)

	// swOffset := len(result) - 2
	// sw := codec.Uint16(result[swOffset:])
	// if sw != 0x9000 {
	// 	return nil, fmt.Errorf("Invalid status %04x", sw)
	// }
	// return response[:swOffset], nil
}
