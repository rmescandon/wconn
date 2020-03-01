package network

import (
	"github.com/godbus/dbus"
	"github.com/pkg/errors"
)

// Manager network manager
type Manager interface {
	Ssids() ([]string, error)
	Connected() (bool, error)
}

type manager struct {
	dbusBase
}

// NewManager returns a new Network Manager object
func NewManager() (Manager, error) {
	return newManager()
}

// Ssids returns a string list of the available WiFi ssids to connect to
func (m *manager) Ssids() ([]string, error) {
	// Idiomatic way of creating a set for not duplicating SSIDs
	ssidsSet := make(map[string]bool)
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

		ssidsSet[ssid] = true
	}

	for k := range ssidsSet {
		ssids = append(ssids, k)
	}

	return ssids, nil
}

// Connected returns true if any WiFi device is connected to a network
func (m *manager) Connected() (bool, error) {
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

func newManager() (*manager, error) {
	c, err := dbus.SystemBus()
	if err != nil {
		return nil, errors.Errorf("Error getting system dbus reference: %v", err)
	}

	db := newDbusBase(c, "/org/freedesktop/NetworkManager")
	return &manager{db}, nil
}

func (m *manager) newDev(path string) *dev {
	return &dev{
		newDbusBase(m.c, path),
	}
}

func (m *manager) devices() ([]*dev, error) {
	var devPaths []string
	err := m.o.Call("org.freedesktop.NetworkManager.GetAllDevices", 0).Store(&devPaths)

	var devs []*dev
	for _, devPath := range devPaths {
		devs = append(devs, m.newDev(devPath))
	}
	return devs, err
}

func (m *manager) wifiDevices() ([]*dev, error) {
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
