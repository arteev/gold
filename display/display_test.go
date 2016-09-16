package display

import (
	"errors"
	"fmt"
	"testing"

	"github.com/arteev/gold/driver"
)

type mockDriver struct {
	GetDisplayFn      func() (driver.Display, error)
	GetDisplayInvoked bool
}

func TestDSPFields(t *testing.T) {
	var cfg = make(map[string]interface{})
	drv := &mockDriver{}
	cfg["port"] = "/dev/tty"
	dsp := &DSP{
		config: cfg,
		driver: drv,
	}
	got := dsp.Config()
	if port, ok := got["port"]; !ok || port != cfg["port"] {
		t.Errorf("Excepted: %q,got:%q\n", cfg, got)
	}

	if got := dsp.Driver(); got != drv {
		t.Errorf("Excepted: %v,got: %v\n", drv, got)
	}
}

func (d *mockDriver) GetDisplay(protocol driver.Protocol, config map[string]interface{}) (driver.Display, error) {
	d.GetDisplayInvoked = true
	return d.GetDisplayFn()
}

func TestRegister(t *testing.T) {
	name := "mockDriver"
	drv := &mockDriver{}
	Register(name, drv)
	if _, ok := drivers[name]; !ok {
		t.Fatalf("Driver not registred: %q", name)
	}

	if got := Drivers(); len(got) != 1 || got[0] != name {
		t.Errorf("Excepted Drivers()=%q, got:%q\n", name, got)
	}

	chknil := func() (result error) {
		defer func() {
			if r := recover(); r != nil {
				result = errors.New(r.(string))
			}
		}()
		Register(name, nil)
		return
	}
	if err := chknil(); err == nil || err.Error() != "dsp: Register driver is nil" {
		t.Errorf("Excepted: panic dsp: Register driver is nil,got %q\n", err)
	}

	chktwice := func() (result error) {
		defer func() {
			if r := recover(); r != nil {
				result = errors.New(r.(string))
			}
		}()
		Register(name, drv)
		return
	}
	if err := chktwice(); err == nil || err.Error() != "dsp: Register called twice for driver "+name {
		t.Errorf("Excepted: dsp: Register called twice for driver %v, got %q\n", name, err)
	}

	unregisterAllDrivers()
	if got := len(drivers); got != 0 {
		t.Errorf("Excepted: count of drivers 0, got %q", got)
	}

}

func TestOpen(t *testing.T) {
	mock := &mockDriver{}
	nameMock := "mockDriver"
	mock.GetDisplayFn = func() (driver.Display, error) {
		return nil, nil
	}
	Register(nameMock, mock)
	dsp, err := Open(nameMock, nil)
	if err != nil {
		t.Fatal(err)
	}
	dsp.GetDisplay(nil)

	if !mock.GetDisplayInvoked {
		t.Fatal("Excepted Open() to be invoked")
	}

	fake := "fake"
	if _, err := Open(fake, nil); err == nil || err.Error() != fmt.Sprintf("dsp: unknown driver %q", fake) {
		t.Errorf("Excepted: dsp: unknown driver %q", fake)
	}
}
