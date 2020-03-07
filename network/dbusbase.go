package network

import (
	"github.com/godbus/dbus"
)

type dbusBase struct {
	c *dbus.Conn
	o dbus.BusObject
}

var newDbusBase = func(c *dbus.Conn, path string) dbusBase {
	return dbusBase{
		c: c,
		o: c.Object("org.freedesktop.NetworkManager", dbus.ObjectPath(path)),
	}
}

func (dbb *dbusBase) prop(prop string) (dbus.Variant, error) {
	return dbb.o.GetProperty(prop)
}

func (dbb *dbusBase) propAsStr(prop string) (string, error) {
	v, err := dbb.prop(prop)
	if err != nil {
		return "", err
	}

	switch v.Value().(type) {
	case []byte:
		return string(v.Value().([]byte)), nil
	case string:
		return v.Value().(string), nil
	default:
		return string(v.Value().(dbus.ObjectPath)), nil
	}
}
