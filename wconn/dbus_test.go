package wconn

import (
	"fmt"
	"testing"
)

func TestConnectedToWifi(t *testing.T) {
	fmt.Printf("Connected to WIFI: %v\n", ConnectedToWifi())
}

func TestAccessPoints(t *testing.T) {
	aps, err := AccessPoints()
	if err != nil {
		t.Errorf("Error getting access points:%v\n", err)
	}

	fmt.Printf("AccessPoints:%v\n", aps)
}
