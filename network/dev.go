package network

// wifiType is the flag determinating a device on dbus
const wifiType uint32 = 2

// connected is the flag determinating if device is connected to a network
const connected uint32 = 100

type dev struct {
	dbusBase
}

func (d *dev) newAp(path string) *ap {
	return &ap{
		newDbusBase(d.c, path),
	}
}

func (d *dev) isConnected() (bool, error) {
	return d.is("org.freedesktop.NetworkManager.Device.State", connected)
}

func (d *dev) isWifi() (bool, error) {
	return d.is("org.freedesktop.NetworkManager.Device.DeviceType", wifiType)
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
