package network_test

import (
	"testing"

	"github.com/godbus/dbus"
	"github.com/greenbrew/wconn/mocks"
	"github.com/greenbrew/wconn/network"
	check "gopkg.in/check.v1"

	"github.com/golang/mock/gomock"
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

	nWifiDevices := 2
	nNotWifiDevices := 1

	// Get All devices once
	s.mockBusObject.EXPECT().Call("org.freedesktop.NetworkManager.GetAllDevices", gomock.Any()).Return(
		&dbus.Call{
			Body: []interface{}{devs},
		})

	// Two devices are wifi
	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.DeviceType").Return(
		dbus.MakeVariant(uint32(2)),
		nil,
	).Times(nWifiDevices)

	// One device is not wifi
	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.DeviceType").Return(
		dbus.MakeVariant(uint32(18)),
		nil,
	).Times(nNotWifiDevices)

	// Get the access points
	s.mockBusObject.EXPECT().Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", gomock.Any()).Return(
		&dbus.Call{
			Body: []interface{}{aps},
		})

	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid").Return(
		dbus.MakeVariant("MyHomeWifi"),
		nil,
	)

	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid").Return(
		dbus.MakeVariant("MyNeigbourWifi"),
		nil,
	)

	// Execute the test
	nm, err := network.NewNm()
	c.Assert(err, check.IsNil)

	ssids, err := nm.Ssids()
	c.Assert(err, check.IsNil)
	c.Assert(ssids, check.HasLen, 2)
	c.Assert(ssids[0], check.Equals, "MyHomeWifi")
	c.Assert(ssids[1], check.Equals, "MyNeigbourWifi")
}
