package reader

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"reflect"
	"strings"
	"syscall"
)

const (
	inputs     = "/sys/class/input/event%d/device/uevent"
	deviceFile = "/dev/input/event%d"
	maxFiles   = 255
)

// event types
const (
	EvSYN      = 0x00
	EvKEY      = 0x01
	EvREL      = 0x02
	EvABS      = 0x03
	EvMSC      = 0x04
	EvSW       = 0x05
	EvLED      = 0x11
	EvSND      = 0x12
	EvREP      = 0x14
	EvFF       = 0x15
	EvPWR      = 0x16
	EvFFStatus = 0x17
	EvMAX      = 0x1f
)

var eventsize = int(reflect.TypeOf(InputEvent{}).Size())

// Reader represents a Vending Machine compatible reader
type Reader struct {
	dev *InputDevice
}

type InputDevice struct {
	Id   int
	Name string
}

type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

// GetReader returns a detected reader or appropriate error
func GetReader(readerName string) (*Reader, error) {
	var devices []*InputDevice

	if err := checkRoot(); err != nil {
		return nil, err
	}

	for i := 0; i < maxFiles; i++ {
		buff, err := ioutil.ReadFile(fmt.Sprintf(inputs, i))
		if err != nil {
			break
		}
		devices = append(devices, newInputDeviceReader(buff, i))
	}

	for _, d := range devices {
		if d.Name == readerName {
			return &Reader{
				dev: d,
			}, nil
		}
	}

	return nil, errors.New("no suitable reader found")
}

func checkRoot() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	if u.Uid != "0" {
		return fmt.Errorf("cannot read device files. Are you running as root?")
	}
	return nil
}

func newInputDeviceReader(buff []byte, id int) *InputDevice {
	rd := bufio.NewReader(bytes.NewReader(buff))
	rd.ReadLine()
	dev, _, _ := rd.ReadLine()
	splt := strings.Split(string(dev), "=")

	return &InputDevice{
		Id:   id,
		Name: splt[1],
	}
}

func (r *Reader) Read() (<-chan InputEvent, error) {
	ret := make(chan InputEvent, 512)

	if err := checkRoot(); err != nil {
		close(ret)
		return ret, err
	}

	fd, err := os.Open(fmt.Sprintf(deviceFile, r.dev.Id))
	if err != nil {
		close(ret)
		return ret, err
	}

	go func() {

		tmp := make([]byte, eventsize)
		event := InputEvent{}
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
	return ret, nil
}

// Keystring Returns the char of a key event
func (i *InputEvent) KeyString() string {
	return keyCodeMap[i.Code]
}
