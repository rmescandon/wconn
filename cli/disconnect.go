package cli

import (
	"errors"
	"time"

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

	ch, err := m.Disconnect()
	if err != nil {
		return err
	}

	for {
		select {
		case st := <-ch:
			if st == network.Disconnected {
				return nil
			}
		case <-time.After(5 * time.Second):
			return errors.New("Timeout")
		}
	}
}
