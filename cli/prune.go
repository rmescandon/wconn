package cli

import (
	"github.com/greenbrew/wconn/network"
)

// PruneCmd command to prune all not connected wifi connections
type PruneCmd struct{}

// Execute executes the prune command
func (cmd *PruneCmd) Execute(args []string) error {
	m, err := network.NewManager()
	if err != nil {
		return err
	}

	return m.PruneConnections()
}
