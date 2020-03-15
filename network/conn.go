package network

const (
	// Interface
	networkManagerSettingsInterface           = networkManagerInterface + ".Settings"
	networkManagerSettingsConnectionInterface = networkManagerSettingsInterface + ".Connection"

	// Methods
	networkManagerSettingsConnectionGetSettings = networkManagerSettingsConnectionInterface + ".GetSettings"
	networkManagerSettingsConnectionDelete      = networkManagerSettingsConnectionInterface + ".Delete"
)

type conn struct {
	dbusBase
}

func (c *conn) ssid() (string, error) {
	var settings map[string]interface{}
	err := c.o.Call(networkManagerSettingsConnectionGetSettings, 0).Store(&settings)
	if err != nil {
		return "", err
	}

	s, ok := settings["802-11-wireless"]
	if !ok {
		return "", nil
	}

	wifiSettings := s.(map[string]interface{})
	val, ok := wifiSettings["ssid"]
	if !ok {
		return "", nil
	}

	return string(val.([]byte)), nil
}

func (c *conn) delete() error {
	return c.o.Call(networkManagerSettingsConnectionDelete, 0).Err
}
