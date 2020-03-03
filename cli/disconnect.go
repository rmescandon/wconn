package cli

import (
	"github.com/greenbrew/wconn/network"
)

// DisconnectCmd command to disconnect from current connected network
type DisconnectCmd struct{}

// Execute executes the disconnection command
func (cmd *DisconnectCmd) Execute(args []string) error {
	m, err := network.NewManager()
	if err != nil {
		return err
	}

	return m.Disconnect()
}
