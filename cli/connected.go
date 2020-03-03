package cli

import (
	"fmt"

	"github.com/greenbrew/wconn/network"
)

// ConnectedCmd command identifies if connected to WiFi network
type ConnectedCmd struct{}

// Execute executes the connected command
func (cmd *ConnectedCmd) Execute(args []string) error {
	m, err := network.NewManager()
	if err != nil {
		return err
	}

	b, err := m.Connected()
	if err != nil {
		return err
	}

	fmt.Println(b)
	return nil
}
