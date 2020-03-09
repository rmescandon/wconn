package network

import (
	"github.com/godbus/dbus"
	"github.com/greenbrew/wconn/mocks"
)

func MockBusObject(m *mocks.MockBusObject) {
	newDbusBase = func(c *dbus.Conn, path string) dbusBase {
		return dbusBase{
			c: &dbus.Conn{},
			o: m,
		}
	}
}
