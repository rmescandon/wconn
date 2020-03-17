package network

const (
	// Interface
	networkManagerSettingsInterface = networkManagerInterface + ".Settings"

	// Methods
	networkManagerSettingsListConnections = networkManagerSettingsInterface + ".ListConnections"
)

type settings struct {
	dbusBase
}

func (s *settings) listConnections() ([]*conn, error) {
	var paths []string
	err := s.o.Call(networkManagerSettingsListConnections, 0).Store(&paths)
	if err != nil {
		return nil, err
	}

	var cs []*conn
	for _, p := range paths {
		cs = append(cs, s.newConn(p))
	}
	return cs, nil
}

func (s *settings) wifiConns() ([]*conn, error) {
	cs, err := s.listConnections()
	if err != nil {
		return nil, err
	}

	var wCs []*conn
	for _, c := range cs {
		b, err := c.isWifi()
		if err != nil {
			return nil, err
		}
		if b {
			wCs = append(wCs, c)
		}
	}
	return wCs, nil
}
