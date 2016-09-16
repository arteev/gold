package serial

import (
	"github.com/arteev/gold/display"
	"github.com/arteev/gold/driver"
	"github.com/arteev/gold/serial/com"
	"github.com/tarm/serial"
)

type SerialDriver struct {
}

func (d *SerialDriver) GetDisplay(protocol driver.Protocol, config map[string]interface{}) (driver.Display, error) {
	s := com.MustSerial(protocol)
	get := func(name string, defvalue interface{}) interface{} {
		val, ok := config[name]
		if ok {
			return val
		}
		return defvalue
	}
	c := &serial.Config{}
	c.Name = get("Name", "").(string)
	c.Baud = get("Baud", 9600).(int)
	c.Size = get("Size", byte(serial.DefaultSize)).(byte)
	c.StopBits = serial.StopBits(get("StopBits", byte(serial.Stop1)).(byte))
	c.Parity = serial.Parity(get("Parity", byte(serial.ParityNone)).(byte))
	port, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	s.CreatePort(port)
	return s, nil
}

func init() {
	display.Register("serial", &SerialDriver{})
}

//TODO: GitHub
//TODO: ReadMe
//TODO: user symbols
//TODO: Examples
//TODO: Dateks
