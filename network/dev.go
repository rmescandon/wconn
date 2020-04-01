package network

import (
	"github.com/godbus/dbus"
	"github.com/pkg/errors"
)

const (
	// WifiDeviceType is the flag determinating a device on dbus
	WifiDeviceType uint32 = 2

	wifiDeviceDisconnected uint32 = 30
	// WifiDeviceConnected is the flag determinating if device is connected to a network
	WifiDeviceConnected    uint32 = 100
	wifiDeviceDeactivating uint32 = 110
)

const (
	// Interfaces
	deviceIface         = managerIface + ".Device"
	deviceWirelessIface = deviceIface + ".Wireless"

	// Properties
	deviceActiveConnection = deviceIface + ".ActiveConnection"
	deviceState            = deviceIface + ".State"
	deviceDeviceType       = deviceIface + ".DeviceType"

	// Methods
	deviceWirelessGetAccessPoints = deviceWirelessIface + ".GetAccessPoints"
	deviceWirelessRequestScan     = deviceWirelessIface + ".RequestScan"
	deviceDisconnect              = deviceIface + ".Disconnect"
	deviceAvailableConnections    = deviceIface + ".AvailableConnections"

	// Members
	stateChanged = "StateChanged"

	// Signals
	deviceStateChanged = deviceIface + "." + stateChanged
)

type dev struct {
	dbusBase
	ch <-chan *dbus.Signal
}

func (d *dev) newAp(path string) *ap {
	return &ap{
		newDbusBase(d.c, path),
	}
}

func (d *dev) newActiveConn(path string) *activeConn {
	return &activeConn{
		newDbusBase(d.c, path),
	}
}

func (d *dev) isConnected() (bool, error) {
	return d.is(deviceState, WifiDeviceConnected)
}

func (d *dev) isWifi() (bool, error) {
	return d.is(deviceDeviceType, WifiDeviceType)
}

func (d *dev) accessPoints() ([]*ap, error) {
	var apPaths []string
	err := d.o.Call(deviceWirelessGetAccessPoints, 0).Store(&apPaths)
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

func (d *dev) activeConnection() (*activeConn, error) {
	str, err := d.propAsStr(deviceActiveConnection)
	if err != nil {
		return nil, err
	}
	return d.newActiveConn(str), nil
}

func (d *dev) disconnect() error {
	return d.o.Call(deviceDisconnect, 0).Err
}

func (d *dev) conns() ([]*conn, error) {
	connPaths, err := d.propAsStrArray(deviceAvailableConnections)
	if err != nil {
		return nil, err
	}

	var conns []*conn
	for _, connPath := range connPaths {
		conns = append(conns, d.newConn(connPath))
	}
	return conns, nil
}

func (d *dev) findExistingConn(ssid string) (*conn, error) {
	cs, err := d.conns()
	if err != nil {
		return nil, err
	}

	for _, c := range cs {
		s, err := c.ssid()
		if err != nil {
			return nil, err
		}
		if s == ssid {
			return c, nil
		}
	}
	return nil, nil
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

func (d *dev) subscribeStateChanged() error {
	return d.o.AddMatchSignal(
		deviceIface,
		stateChanged,
		dbus.WithMatchObjectPath(d.o.Path()),
	).Err
}

func (d *dev) unsubscribeStateChanged() error {
	return d.o.RemoveMatchSignal(
		deviceIface,
		stateChanged,
		dbus.WithMatchObjectPath(d.o.Path()),
	).Err
}

func (d *dev) registerStateChanged(ch chan ConnectionState) error {
	d.ch = d.listen()
	if err := d.subscribeStateChanged(); err != nil {
		return err
	}

	if ch == nil {
		return errors.New("State changed channel is nil")
	}

	go func() {
		// Infinite loop until having a connection state change on this device
		for {
			signal := <-d.ch
			switch signalStatus(signal) {
			case WifiDeviceConnected:
				ch <- Connected
				return
			case wifiDeviceDeactivating:
				fallthrough
			case wifiDeviceDisconnected:
				ch <- Disconnected
				return
			}
		}
	}()

	return nil
}

func (d *dev) requestScan() error {
	return d.o.Call(deviceWirelessRequestScan, 0, map[string]interface{}{}).Err
}

func connectedSignal(s *dbus.Signal) bool {
	return signalStatus(s) == WifiDeviceConnected
}

func disconnectedSignal(s *dbus.Signal) bool {
	status := signalStatus(s)
	return status == wifiDeviceDisconnected || status == wifiDeviceDeactivating
}

func signalStatus(s *dbus.Signal) uint32 {
	if s == nil {
		return uint32(0)
	}
	return s.Body[0].(uint32)
}
