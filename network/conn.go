package network

import "github.com/pkg/errors"

const (
	// Methods
	networkManagerSettingsConnectionGetSettings = networkManagerSettingsConnectionInterface + ".GetSettings"
	networkManagerSettingsConnectionDelete      = networkManagerSettingsConnectionInterface + ".Delete"
)

type conn struct {
	dbusBase
}

func (c *conn) ssid() (string, error) {
	s, err := c.wifiSettings()
	if err != nil {
		return "", err
	}
	val, ok := s["ssid"]
	if !ok {
		return "", errors.New("Could not find 'ssid' WiFi setting")
	}

	return string(val.([]byte)), nil
}

func (c *conn) wifiSettings() (map[string]interface{}, error) {
	settings, err := c.settings()
	if err != nil {
		return nil, err
	}

	s, ok := settings["802-11-wireless"]
	if !ok {
		return nil, errors.New("Could not find WiFi settings")
	}

	return s.(map[string]interface{}), nil
}

func (c *conn) isWifi() (bool, error) {
	settings, err := c.settings()
	if err != nil {
		return false, err
	}

	_, ok := settings["802-11-wireless"]
	return ok, nil
}

func (c *conn) settings() (map[string]interface{}, error) {
	var settings map[string]interface{}
	err := c.o.Call(networkManagerSettingsConnectionGetSettings, 0).Store(&settings)
	return settings, err
}

func (c *conn) delete() error {
	return c.o.Call(networkManagerSettingsConnectionDelete, 0).Err
}
