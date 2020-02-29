package network

import (
	"github.com/godbus/dbus"
	"github.com/pkg/errors"
)

// Nm network manager
type Nm interface {
	Ssids() ([]string, error)
	Connected() (bool, error)
}

type nm struct {
	dbusBase
}

// NewNm returns a new Network Manager object
func NewNm() (Nm, error) {
	return newNm()
}

// Ssids returns a string list of the available WiFi ssids to connect to
func (m *nm) Ssids() ([]string, error) {
	var ssids []string

	devs, err := m.wifiDevices()
	if err != nil {
		return ssids, err
	}

	if len(devs) == 0 {
		return ssids, nil
	}

	// take the first wifi device (any is valid) and get the access points
	aps, err := devs[0].accessPoints()
	if err != nil {
		return ssids, err
	}

	for _, ap := range aps {
		ssid, err := ap.ssid()
		if err != nil {
			return ssids, err
		}
		ssids = append(ssids, ssid)
	}

	return ssids, nil
}

// Connected returns true if any WiFi device is connected to a network
func (m *nm) Connected() (bool, error) {
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

func newNm() (*nm, error) {
	c, err := dbus.SystemBus()
	if err != nil {
		return nil, errors.Errorf("Error getting system dbus reference: %v", err)
	}

	db := newDbusBase(c, "/org/freedesktop/NetworkManager")
	return &nm{db}, nil
}

func (m *nm) newDev(path string) *dev {
	return &dev{
		newDbusBase(m.c, path),
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
