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
