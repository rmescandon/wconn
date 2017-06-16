package wconn

import (
	"log"

	"github.com/godbus/dbus"
)

const wifiType uint32 = 2

// AccessPoint holds ssid and dbus path for a network Ap
type AccessPoint struct {
	SSID string
	Path string
}

// ConnectedToWifi returns true if detected any wifi device is connected to an external network
func ConnectedToWifi() bool {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("Error getting system dbus reference: %v\n", err)
		return false
	}

	// get a reference to all devices
	obj := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	var devices []string
	err = obj.Call("org.freedesktop.NetworkManager.GetAllDevices", 0).Store(&devices)
	if err != nil {
		log.Printf("Error getting all devices: %v\n", err)
	}

	// filter only wifi devices and return true if any is connected
	for _, d := range devices {
		devicePath := dbus.ObjectPath(d)
		obj := conn.Object("org.freedesktop.NetworkManager", devicePath)

		deviceType, err := obj.GetProperty("org.freedesktop.NetworkManager.Device.DeviceType")
		if err != nil {
			log.Printf("Error getting wifi devices: %v\n", err)
			continue
		}

		if deviceType.Value() == nil {
			break
		}
		if deviceType.Value() != wifiType {
			continue
		}

		state, err := obj.GetProperty("org.freedesktop.NetworkManager.Device.State")
		if err != nil {
			log.Printf("Error getting device state: %v\n", err)
			continue
		}

		if dbus.Variant.Value(state) == uint32(100) {
			return true
		}
	}
	return false
}

// AccessPoints returns available access points to connect to
func AccessPoints() ([]string, error) {

	var ssids []string

	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("Error getting system dbus reference: %v\n", err)
		return nil, err
	}

	// get a reference to all devices
	obj := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	var devices []string
	err = obj.Call("org.freedesktop.NetworkManager.GetAllDevices", 0).Store(&devices)
	if err != nil {
		log.Printf("Error getting all devices: %v\n", err)
		return nil, err
	}

	// filter only wifi devices and return true if any is connected
	for _, d := range devices {
		devicePath := dbus.ObjectPath(d)
		obj := conn.Object("org.freedesktop.NetworkManager", devicePath)

		deviceType, err := obj.GetProperty("org.freedesktop.NetworkManager.Device.DeviceType")
		if err != nil {
			log.Printf("Error getting wifi devices: %v\n", err)
			continue
		}

		if deviceType.Value() == nil {
			break
		}

		if deviceType.Value() != wifiType {
			continue
		}

		var accessPoints []string
		err = obj.Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", 0).Store(&accessPoints)
		if err != nil {
			log.Printf("Error getting available access points: %v\n", err)
			continue
		}

		for _, ap := range accessPoints {
			apPath := dbus.ObjectPath(ap)
			obj := conn.Object("org.freedesktop.NetworkManager", apPath)

			ssid, err := obj.GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid")
			if err != nil {
				log.Printf("Error getting ssid: %v\n", err)
				continue
			}

			ssidStr := string(ssid.Value().([]byte))
			if len(ssidStr) < 1 {
				continue
			}

			ssids = append(ssids, ssidStr)
		}
	}

	return ssids, nil
}
