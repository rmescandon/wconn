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
	// Define the dbusObject expectations
	s.mockBusObject.EXPECT().Call("org.freedesktop.NetworkManager.GetAllDevices", gomock.Any()).Return(
		&dbus.Call{
			Body: []interface{}{[]string{"WifiDevice1"}},
		})

	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.DeviceType").Return(
		dbus.MakeVariant(uint32(2)),
		nil,
	)

	s.mockBusObject.EXPECT().Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", gomock.Any()).Return(
		&dbus.Call{
			Body: []interface{}{[]string{"AccessPoint1"}},
		})

	s.mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid").Return(
		dbus.MakeVariant("MyHomeWifi"),
		nil,
	)

	// Execute the test
	nm, err := network.NewNm()
	c.Assert(err, check.IsNil)

	ssids, err := nm.Ssids()
	c.Assert(err, check.IsNil)
	c.Assert(ssids, check.HasLen, 1)
	c.Assert(ssids[0], check.Equals, "MyHomeWifi")
}
