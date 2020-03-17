package network

const (
	// Interfaces
	networkManagerConnectionActiveInterface = networkManagerInterface + ".Connection.Active"

	// Properties
	networkManagerConnectionActiveUUID = networkManagerConnectionActiveInterface + ".Uuid"
)

type activeConn struct {
	dbusBase
}

func (ac *activeConn) uuid() (string, error) {
	return ac.propAsStr(networkManagerConnectionActiveUUID)
}
