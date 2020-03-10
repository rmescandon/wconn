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

func (dbb *dbusBase) propAsBool(prop string) (bool, error) {
	v, err := dbb.prop(prop)
	if err != nil {
		return false, err
	}

	return v.Value().(bool), nil
}

func (dbb *dbusBase) propAsStrArray(prop string) ([]string, error) {
	vals, err := dbb.prop(prop)
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, v := range vals.Value().([]dbus.ObjectPath) {
		ret = append(ret, string(v))
	}
	return ret, nil
}

func (dbb *dbusBase) listen() <-chan *dbus.Signal {
	signal := make(chan *dbus.Signal, 10)
	dbb.c.Signal(signal)
	return signal
}
