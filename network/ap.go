package network

type ap struct {
	dbusBase
}

func (a *ap) ssid() (string, error) {
	ssid, err := a.o.GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid")
	if err != nil {
		return "", err
	}

	return string(ssid.Value().([]byte)), nil
}
