package network

const (
	// Interfaces
	connectionActiveIface = managerIface + ".Connection.Active"

	// Properties
	connectionActiveUUID = connectionActiveIface + ".Uuid"
)

type activeConn struct {
	dbusBase
}

func (ac *activeConn) uuid() (string, error) {
	return ac.propAsStr(connectionActiveUUID)
}
