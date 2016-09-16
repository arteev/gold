package display

import (
	"sort"
	"sync"

	"fmt"

	"github.com/arteev/gold/driver"
)

/*

    Транспорт выделить tty!
        Скорость обмена

   протокол
   режимы
   ширина
   кодовая страница

   Яркость дисплея
   Отображение курсора
   Отображение символов пользователя

   Флажки
   Курсор


*/
var (
	muDrivers sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

type DSP struct {
	driver driver.Driver
	config map[string]interface{}
}

func Open(driverName string, configDriver map[string]interface{}) (*DSP, error) {
	muDrivers.RLock()
	drv, ok := drivers[driverName]
	muDrivers.RUnlock()
	if !ok {
		return nil, fmt.Errorf("dsp: unknown driver %q", driverName)
	}
	dsp := &DSP{
		driver: drv,
		config: configDriver,
	}
	return dsp, nil
}

func Register(name string, driver driver.Driver) {
	muDrivers.Lock()
	defer muDrivers.Unlock()
	if driver == nil {
		panic("dsp: Register driver is nil")
	}
	if _, exists := drivers[name]; exists {
		panic("dsp: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func unregisterAllDrivers() {
	muDrivers.Lock()
	defer muDrivers.Unlock()
	drivers = make(map[string]driver.Driver)
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	muDrivers.RLock()
	defer muDrivers.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// Driver returns the display's underlying driver.
func (d *DSP) Driver() driver.Driver {
	return d.driver
}

func (d *DSP) GetDisplay(protocol driver.Protocol) (driver.Display, error) {
	return d.driver.GetDisplay(protocol, d.config)
}

// Config returns a map of the config driver
func (d *DSP) Config() map[string]interface{} {
	result := make(map[string]interface{})
	for key, pair := range d.config {
		result[key] = pair
	}
	return result
}
