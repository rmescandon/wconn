package network

const (
	// Interfaces
	accessPointIface = managerIface + ".AccessPoint"

	// Properties
	accessPointSsid = accessPointIface + ".Ssid"
)

type ap struct {
	dbusBase
}

func (a *ap) ssid() (string, error) {
	return a.propAsStr(accessPointSsid)
}
