package network

import (
	"github.com/pkg/errors"
)

const (
	// WifiDeviceType is the flag determinating a device on dbus
	WifiDeviceType uint32 = 2

	// WifiDeviceConnected is the flag determinating if device is connected to a network
	WifiDeviceConnected uint32 = 100
)

const (
	// Interfaces
	deviceInterface         = networkManagerInterface + ".Device"
	deviceWirelessInterface = deviceInterface + ".Wireless"

	// Properties
	deviceActiveConnection = deviceInterface + ".ActiveConnection"

	// Methods
	deviceWirelessGetAccessPoints = deviceWirelessInterface + ".GetAccessPoints"
)

type dev struct {
	dbusBase
}

func (d *dev) newAp(path string) *ap {
	return &ap{
		newDbusBase(d.c, path),
	}
}

func (d *dev) isConnected() (bool, error) {
	return d.is("org.freedesktop.NetworkManager.Device.State", WifiDeviceConnected)
}

func (d *dev) isWifi() (bool, error) {
	return d.is("org.freedesktop.NetworkManager.Device.DeviceType", WifiDeviceType)
}

func (d *dev) accessPoints() ([]*ap, error) {
	var apPaths []string
	err := d.o.Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", 0).Store(&apPaths)
	if err != nil {
		return nil, err
	}

	var aps []*ap
	for _, apPath := range apPaths {
		aps = append(aps, d.newAp(apPath))
	}
	return aps, nil
}

func (d *dev) accessPoint(ssid string) (*ap, error) {
	aps, err := d.accessPoints()
	if err != nil {
		return nil, err
	}

	for _, ap := range aps {
		s, err := ap.ssid()
		if err != nil {
			return nil, err
		}

		if s == ssid {
			return ap, nil
		}
	}

	return nil, errors.Errorf("Could not find an access point for %v", ssid)
}

// func (d *dev) activeConnection() (string, error) {
// 	return d.propAsStr(deviceActiveConnection)
// }

func (d *dev) disconnect() error {
	return d.o.Call("org.freedesktop.NetworkManager.Device.Disconnect", 0).Err
}

func (d *dev) is(propertyPath string, comparationFlag uint32) (bool, error) {
	v, err := d.o.GetProperty(propertyPath)
	if err != nil {
		return false, err
	}

	if v.Value() == nil {
		return false, nil
	}

	return v.Value() == comparationFlag, nil
}
