package network

import (
	"github.com/godbus/dbus"
	"github.com/pkg/errors"
)

type dbusBase struct {
	c *dbus.Conn
	o dbus.BusObject
}

type nm struct {
	dbusBase
}

func newNm() (*nm, error) {
	c, err := dbus.SystemBus()
	if err != nil {
		return nil, errors.Errorf("Error getting system dbus reference: %v", err)
	}

	o := c.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")

	return &nm{dbusBase{c: c, o: o}}, nil
}

func (m *nm) newDev(path string) *dev {
	return &dev{
		dbusBase{
			c: m.c,
			o: m.c.Object("org.freedesktop.NetworkManager", dbus.ObjectPath(path)),
		},
	}
}

func (m *nm) devices() ([]*dev, error) {
	var devPaths []string
	err := m.o.Call("org.freedesktop.NetworkManager.GetAllDevices", 0).Store(&devPaths)

	var devs []*dev
	for _, devPath := range devPaths {
		devs = append(devs, m.newDev(devPath))
	}
	return devs, err
}

func (m *nm) wifiDevices() ([]*dev, error) {
	devs, err := m.devices()
	if err != nil {
		return nil, err
	}

	var wifiDevs []*dev
	for _, d := range devs {
		isWifi, err := d.isWifi()
		if err != nil {
			return nil, err
		}

		if !isWifi {
			continue
		}

		wifiDevs = append(wifiDevs, d)
	}

	return wifiDevs, nil
}

func (m *nm) connected() (bool, error) {
	devs, err := m.wifiDevices()
	if err != nil {
		return false, err
	}
	for _, d := range devs {
		isConnected, err := d.isConnected()
		if err != nil {
			return false, err
		}

		if isConnected {
			return true, nil
		}
	}
	return false, nil
}
