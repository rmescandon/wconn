package network

import (
	"github.com/pkg/errors"
)

const (
	// Interface
	networkManagerSettingsConnectionInterface = networkManagerSettingsInterface + ".Connection"

	// Methods
	networkManagerSettingsConnectionGetSettings = networkManagerSettingsConnectionInterface + ".GetSettings"
	networkManagerSettingsConnectionDelete      = networkManagerSettingsConnectionInterface + ".Delete"
)

type conn struct {
	dbusBase
}

func (c *conn) uuid() (string, error) {
	s, err := c.connSettings()
	if err != nil {
		return "", err
	}

	val, ok := s["uuid"]
	if !ok {
		return "", errors.New("Could not find 'uuid' setting")
	}

	return val.(string), nil
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
	return c.subsettings("802-11-wireless")
}

func (c *conn) connSettings() (map[string]interface{}, error) {
	return c.subsettings("connection")
}

func (c *conn) subsettings(key string) (map[string]interface{}, error) {
	st, err := c.getSettings()
	if err != nil {
		return nil, err
	}

	s, ok := st[key]
	if !ok {
		return nil, errors.Errorf("Could not find %s settings", key)
	}

	return s.(map[string]interface{}), nil
}

func (c *conn) isWifi() (bool, error) {
	st, err := c.getSettings()
	if err != nil {
		return false, err
	}

	_, ok := st["802-11-wireless"]
	return ok, nil
}

func (c *conn) getSettings() (map[string]interface{}, error) {
	var st map[string]interface{}
	err := c.o.Call(networkManagerSettingsConnectionGetSettings, 0).Store(&st)
	return st, err
}

func (c *conn) delete() error {
	return c.o.Call(networkManagerSettingsConnectionDelete, 0).Err
}
