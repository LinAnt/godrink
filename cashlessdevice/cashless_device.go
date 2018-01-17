package cashlessdevice

import (
	"encoding/binary"

	serial "go.bug.st/serial.v1"
)

const (
	cmdResponseLength = 6

	cmdBeginSession     byte = 0x03
	cmdCancelSession    byte = 0x04
	cmdVendDenied       byte = 0x06
	cmdVenApproved      byte = 0x05
	cmdRevalueApproved  byte = 0x0D
	cmdRevalueDenied    byte = 0x0E
	cmdGetStatus        byte = 0xFB
	cmdInterfaceReset   byte = 0xFE
	cmdSetCountryCode   byte = 0xFE
	cmdSetScalingFactor byte = 0xFE
	cmdSetOptions       byte = 0xFE
	cmdSetDecimal       byte = 0xFE
	cmdSetAUX           byte = 0xFE
	cmdSetRTC           byte = 0xFE
	cmdGetRTC           byte = 0xFE
)

var (
	cmdAck  = []byte{0xFC, 0xFC, 0xFC, 0xFC, 0xFC, 0xFC}
	cmdNack = []byte{0xFD, 0xFD, 0xFD, 0xFD, 0xFD, 0xFD}
)

// CashlessDevice represents a MDB compatible serial device
type CashlessDevice struct {
	serial.Port
	serialMode *serial.Mode
	AckCh      chan bool
}

type Msg struct{}

// SendCancelSession command sends cancel session to the cashless device.
func (cd *CashlessDevice) SendCancelSession() (err error) {
	return cd.sendCommand([]byte{cmdCancelSession, cmdCancelSession})
}

func (cd *CashlessDevice) initResponseListener() error {
	go func() {

		tmp := make([]byte, 50)
		event := Inp
		for {

			_, err := fd.Read(tmp)
			if err != nil {
				close(ret)
				break
			}

			if err := binary.Read(bytes.NewBuffer(tmp), binary.LittleEndian, &event); err != nil {
				panic(err)
			}

			ret <- event

		}
	}()
	}

	return nil
}
func (cd *CashlessDevice) sendCommand(cmd []byte) (err error) {
	err = cd.ResetOutputBuffer()
	if err != nil {
		return err
	}
	_, err = cd.Write(cmd)
	return err
}

func calculateCrc(c []byte) byte {
	sum := uint32(0x00)
	for _, b := range c {
		sum += uint32(b)
	}

	chk := make([]byte, 4)
	binary.BigEndian.PutUint32(chk, sum)
	return chk[3]
}
func validateCrc(c []byte) bool {
	// check edgecases
	if len(c) < 2 {
		return false
	} else if len(c) == 2 {
		return c[0] == c[1]
	}

	sum := uint32(0)
	crcByte := c[len(c)-1]

	for _, b := range c[:len(c)-2] {
		sum += uint32(b)
	}

	chk := make([]byte, 4)
	binary.BigEndian.PutUint32(chk, sum)
	return crcByte == chk[3]
}

func responseHandler(ackChannel chan bool, msgChannel chan Msg){
	firstByte := make([]byte, 1)
	for{
		
		select {}
	}
}