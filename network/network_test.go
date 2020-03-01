package network_test

import (
	"fmt"
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

type NetworkSuite struct{}

var _ = check.Suite(&NetworkSuite{})

func (s *NetworkSuite) TestGetAvailableSsids(c *check.C) {
	// Prepare the mock of dbusObject
	mockCtrl := gomock.NewController(testCtx)
	defer mockCtrl.Finish()

	mockBusObject := mocks.NewMockBusObject(mockCtrl)
	network.MockBusObject(mockBusObject)

	// Define the dbusObject expectations
	mockBusObject.EXPECT().Call("org.freedesktop.NetworkManager.GetAllDevices", gomock.Any()).Return(
		&dbus.Call{
			Body: []interface{}{[]string{"WifiDevice1"}},
		})

	mockBusObject.EXPECT().GetProperty("org.freedesktop.NetworkManager.Device.DeviceType").Return(dbus.MakeVariant(2), nil)

	nm, err := network.NewNm()
	c.Assert(err, check.IsNil)

	ssids, err := nm.Ssids()
	c.Assert(err, check.IsNil)
	for _, s := range ssids {
		fmt.Println(s)
	}
}

// func (s *NetworkSuite) TestGetAvailableSsids(c *check.C) {
// 	nm, err := network.NewNm()
// 	c.Assert(err, check.IsNil)

// 	ssids, err := nm.Ssids()
// 	c.Assert(err, check.IsNil)
// 	for _, s := range ssids {
// 		fmt.Println(s)
// 	}
// }
