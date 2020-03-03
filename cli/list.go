package cli

import (
	"fmt"

	"github.com/greenbrew/wconn/network"
)

// ListCmd command to list available networks
type ListCmd struct{}

// Execute executes the disconnection command
func (cmd *ListCmd) Execute(args []string) error {
	m, err := network.NewManager()
	if err != nil {
		return err
	}

	ssids, err := m.Ssids()
	if err != nil {
		return err
	}

	for _, s := range ssids {
		fmt.Println(s)
	}

	return nil
}
