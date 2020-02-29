package network_test

import (
	"fmt"
	"testing"

	"github.com/greenbrew/wconn/network"
	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type NetworkSuite struct{}

var _ = check.Suite(&NetworkSuite{})

func (s *NetworkSuite) TestGetAvailableSsids(c *check.C) {
	nm, err := network.NewNm()
	c.Assert(err, check.IsNil)

	ssids, err := nm.Ssids()
	c.Assert(err, check.IsNil)
	for _, s := range ssids {
		fmt.Println(s)
	}
}
