package cli

import (
	"github.com/greenbrew/wconn/network"
)

// ScanCmd command to request scanning all accesible wifi connections
type ScanCmd struct{}

// Execute executes the request scan command
func (cmd *ScanCmd) Execute(args []string) error {
	m, err := network.NewManager()
	if err != nil {
		return err
	}

	return m.RequestScan()
}
