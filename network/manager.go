package network

import (
	"github.com/godbus/dbus"
	"github.com/pkg/errors"
)

const (
	// Interfaces
	networkManagerInterface = "org.freedesktop.NetworkManager"

	// Methods
	networkManagerGetAllDevices            = networkManagerInterface + ".GetAllDevices"
	networkManagerActivateConnection       = networkManagerInterface + ".ActivateConnection"
	networkManagerAddAndActivateConnection = networkManagerInterface + ".AddAndActivateConnection"
	newworkManagerEnable                   = networkManagerInterface + ".Enable"

	// Properties
	networkManagerNetworkingEnabled = networkManagerInterface + ".NetworkingEnabled"

	// Objects
	networkManagerObject         = "/org/freedesktop/NetworkManager"
	networkManagerSettingsObject = networkManagerObject + "/Settings"
)

// ConnectionState the state of the connection taking the values defined as dbus contants
type ConnectionState uint32

// Enum of States of connection
const (
	Disconnected ConnectionState = iota
	Connected
)

// Manager network manager
type Manager interface {
	Ssids() ([]string, error)
	Connected() (bool, error)
	Connect(ssid, passphrase, security, keyMgmt string) (<-chan ConnectionState, error)
	Disconnect() (<-chan ConnectionState, error)
	PruneConnections() error
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

func (m *manager) Connect(ssid, passphrase, security, keyMgmt string) (<-chan ConnectionState, error) {
	if err := m.enable(true); err != nil {
		return nil, err
	}

	d, err := m.firstAvailableWifiDevice()
	if err != nil {
		return nil, err
	}

	// Check old connections
	var oldConn *settingsConn
	cs, err := d.conns()
	if err != nil {
		return nil, err
	}
	for _, c := range cs {
		historicSsid, err := c.ssid()
		if err != nil {
			return nil, err
		}
		if historicSsid == ssid {
			oldConn = c
			break
		}
	}

	a, err := d.accessPoint(ssid)
	if err != nil {
		return nil, err
	}

	ch := make(chan ConnectionState)

	if err = d.registerStateChanged(ch); err != nil {
		return nil, err
	}

	if oldConn != nil {
		err = m.reconnect(oldConn, d, a)
	} else {
		err = m.connect(passphrase, security, keyMgmt, d, a)
	}

	return ch, err
}

func (m *manager) Disconnect() (<-chan ConnectionState, error) {
	devs, err := m.connectedWifiDevices()
	if err != nil {
		return nil, err
	}

	ch := make(chan ConnectionState)

	for _, d := range devs {
		if err = d.registerStateChanged(ch); err != nil {
			return nil, err
		}

		if err = d.disconnect(); err != nil {
			return nil, err
		}
	}

	return ch, nil
}

func (m *manager) PruneConnections() error {
	s := m.newManagerSettings()
	cs, err := s.wifiConns()
	if err != nil {
		return err
	}

	devs, err := m.wifiDevices()
	if err != nil {
		return err
	}

	var connectedUUID string
	for _, d := range devs {
		if len(connectedUUID) == 0 {
			b, err := d.isConnected()
			if err != nil {
				return err
			}
			if b {
				// TODO FIXME: This active connection is a /org/freedesktop/NetworkManager/ActiveConnection/30 instead of /org/freedesktop/NetworkManager/Settings/4
				ac, err := d.activeConnection()
				if err != nil {
					return err
				}
				connectedUUID, err = ac.uuid()
				if err != nil {
					return err
				}
			}
		}

		for _, c := range cs {
			if len(connectedUUID) > 0 {
				uuid, err := c.uuid()
				if err != nil {
					return err
				}
				if uuid == connectedUUID {
					continue
				}
			}

			err = c.delete()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func newManager() (*manager, error) {
	c, err := dbus.SystemBus()
	if err != nil {
		return nil, errors.Errorf("Error getting system dbus reference: %v", err)
	}

	db := newDbusBase(c, networkManagerObject)
	return &manager{db}, nil
}

func (m *manager) newManagerSettings() *settings {
	return &settings{
		newDbusBase(m.c, networkManagerSettingsObject),
	}
}

func (m *manager) newDev(path string) *dev {
	return &dev{
		dbusBase: newDbusBase(m.c, path),
	}
}

func (m *manager) enable(flag bool) error {
	if b, err := m.enabled(); err != nil || b {
		return err
	}
	return m.o.Call(newworkManagerEnable, 0, dbus.MakeVariant(flag)).Err
}

func (m *manager) enabled() (bool, error) {
	return m.propAsBool(networkManagerNetworkingEnabled)
}

func (m *manager) devices() ([]*dev, error) {
	var devPaths []string
	err := m.o.Call(networkManagerGetAllDevices, 0).Store(&devPaths)

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

func (m *manager) connectedWifiDevices() ([]*dev, error) {
	devs, err := m.wifiDevices()
	if err != nil {
		return nil, err
	}

	var cDevs []*dev
	for _, d := range devs {
		isConnected, err := d.isConnected()
		if err != nil {
			return nil, err
		}

		if isConnected {
			cDevs = append(cDevs, d)
		}
	}
	return cDevs, nil
}

func (m *manager) firstAvailableWifiDevice() (*dev, error) {
	devs, err := m.wifiDevices()
	if err != nil {
		return nil, err
	}

	for _, d := range devs {
		b, err := d.isConnected()
		if err != nil {
			return nil, err
		}

		if !b {
			return d, nil
		}
	}
	return nil, errors.New("Could not find an available WiFi device")
}

func (m *manager) getDeviceFromSsid(ssid string) (*dev, error) {
	devs, err := m.wifiDevices()
	if err != nil {
		return nil, err
	}

	for _, d := range devs {
		_, err := d.accessPoint(ssid)
		if err != nil {
			continue
		}

		return d, nil
	}

	return nil, errors.Errorf("Could not find a device for SSID: %v", ssid)
}

func (m *manager) connect(passphrase, security, keyMgmt string, d *dev, a *ap) error {
	settings := map[string]map[string]dbus.Variant{
		"801-11-wireless": map[string]dbus.Variant{
			"security": dbus.MakeVariant(security),
		},
		"802-11-wireless-security": map[string]dbus.Variant{
			"key-mgmt": dbus.MakeVariant(keyMgmt),
			"psk":      dbus.MakeVariant(passphrase),
		},
	}

	return m.o.Call(networkManagerAddAndActivateConnection, 0, settings, d.o.Path(), a.o.Path()).Err
}

func (m *manager) reconnect(c *settingsConn, d *dev, a *ap) error {
	return m.o.Call(networkManagerActivateConnection, 0, c.o.Path(), d.o.Path(), a.o.Path()).Err
}
