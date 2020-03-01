package network_test

import (
	"testing"

	"github.com/godbus/dbus"
	"github.com/greenbrew/wconn/mocks"
	"github.com/greenbrew/wconn/network"
	check "gopkg.in/check.v1"

	"github.com/golang/mock/gomock"
)

const (
	// any value other than 2
	notWifiDeviceType uint32 = 18
	// any value other than 100
	notWifiDeviceConnected uint32 = 200
)

var testCtx *testing.T

func Test(t *testing.T) {
	testCtx = t
	check.TestingT(t)
}

type NetworkSuite struct {
	mockCtrl      *gomock.Controller
	mockBusObject *mocks.MockBusObject
}

var _ = check.Suite(&NetworkSuite{})

func (s *NetworkSuite) SetUpTest(c *check.C) {
	// Prepare the mock of dbusObject
	s.mockCtrl = gomock.NewController(testCtx)
	s.mockBusObject = mocks.NewMockBusObject(s.mockCtrl)
	network.MockBusObject(s.mockBusObject)
}

func (s *NetworkSuite) TearDownTest(c *check.C) {
	s.mockCtrl.Finish()
}

func (s *NetworkSuite) TestGetAvailableSsids(c *check.C) {
	devs := []string{
		"WifiDevice1",
		"WifiDevice2",
		"OtherDevice3",
	}

	aps := []string{
		"AccessPoint1",
		"AccessPoint2",
	}

	reachableSsids := map[string]bool{
		"MyHomeWifi":     false,
		"MyNeigbourWifi": false,
	}

	nWifiDevices := 2
	nNotWifiDevices := 1

	// Get All devices once
	s.mockBusObject.EXPECT().Call("org.freedesktop.NetworkManager.GetAllDevices", gomock.Any()).Return(
		&dbus.Call{
			Body: []interface{}{devs},
		})

	// Two devices are wifi
	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.DeviceType").Return(
		dbus.MakeVariant(uint32(network.WifiDeviceType)),
		nil,
	).Times(nWifiDevices)

	// One device is not wifi
	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.DeviceType").Return(
		dbus.MakeVariant(uint32(notWifiDeviceType)),
		nil,
	).Times(nNotWifiDevices)

	// Get the access points
	s.mockBusObject.EXPECT().Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", gomock.Any()).Return(
		&dbus.Call{
			Body: []interface{}{aps},
		})

	for ssid, _ := range reachableSsids {
		s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid").Return(
			dbus.MakeVariant(ssid),
			nil,
		)
	}

	// Execute the test
	nm, err := network.NewManager()
	c.Assert(err, check.IsNil)

	ssids, err := nm.Ssids()
	c.Assert(err, check.IsNil)
	c.Assert(ssids, check.HasLen, 2)

	for _, ssid := range ssids {
		reachableSsids[ssid] = false
	}

	for _, flag := range reachableSsids {
		c.Assert(flag, check.Equals, false)
	}
}

func (s *NetworkSuite) TestIsConnected(c *check.C) {
	devs := []string{
		"WifiDevice1",
		"WifiDevice2",
		"OtherDevice3",
	}

	nWifiDevices := 2
	nNotWifiDevices := 1

	// Get All devices once
	s.mockBusObject.EXPECT().Call("org.freedesktop.NetworkManager.GetAllDevices", gomock.Any()).Return(
		&dbus.Call{
			Body: []interface{}{devs},
		})

	// Two devices are wifi
	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.DeviceType").Return(
		dbus.MakeVariant(uint32(network.WifiDeviceType)),
		nil,
	).Times(nWifiDevices)

	// One device is not wifi
	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.DeviceType").Return(
		dbus.MakeVariant(uint32(notWifiDeviceType)),
		nil,
	).Times(nNotWifiDevices)

	// Get the Connected state
	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.State").Return(
		dbus.MakeVariant(uint32(network.WifiDeviceConnected)),
		nil,
	)

	// Execute the test
	nm, err := network.NewManager()
	c.Assert(err, check.IsNil)

	b, err := nm.Connected()
	c.Assert(err, check.IsNil)
	c.Assert(b, check.Equals, true)
}

func (s *NetworkSuite) TestIsDisconnected(c *check.C) {
	devs := []string{
		"WifiDevice1",
		"WifiDevice2",
		"OtherDevice3",
	}

	nWifiDevices := 2
	nNotWifiDevices := 1

	// Get All devices once
	s.mockBusObject.EXPECT().Call("org.freedesktop.NetworkManager.GetAllDevices", gomock.Any()).Return(
		&dbus.Call{
			Body: []interface{}{devs},
		})

	// Two devices are wifi
	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.DeviceType").Return(
		dbus.MakeVariant(uint32(network.WifiDeviceType)),
		nil,
	).Times(nWifiDevices)

	// One device is not wifi
	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.DeviceType").Return(
		dbus.MakeVariant(uint32(notWifiDeviceType)),
		nil,
	).Times(nNotWifiDevices)

	// Get the Connected state
	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.State").Return(
		dbus.MakeVariant(uint32(notWifiDeviceConnected)),
		nil,
	).Times(nWifiDevices)

	// Execute the test
	nm, err := network.NewManager()
	c.Assert(err, check.IsNil)

	b, err := nm.Connected()
	c.Assert(err, check.IsNil)
	c.Assert(b, check.Equals, false)
}
