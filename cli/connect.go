package cli

import (
	"errors"
	"time"

	"github.com/greenbrew/wconn/network"
)

// ConnectCmd command to connect to a certain available SSID
type ConnectCmd struct {
	Args struct {
		Ssid       string `positional-arg-name:"ssid" description:"The identifier for the WiFi network to connect to"`
		Passphrase string `positional-arg-name:"passphrase" description:"The passphrase to access the WiFi network"`
	} `positional-args:"yes" required:"yes"`
	Security string `short:"s" long:"security" description:"The connection security system to use" default:"802-11-wireless-security"`
	KeyMgmt  string `short:"k" long:"keymgmt" description:"The passphrase key management schema to use" default:"wpa-psk"`
}

// Execute executes the connect command
func (cmd *ConnectCmd) Execute(args []string) error {
	m, err := network.NewManager()
	if err != nil {
		return err
	}

	ch, err := m.Connect(cmd.Args.Ssid, cmd.Args.Passphrase, cmd.Security, cmd.KeyMgmt)
	if err != nil {
		return err
	}

	for {
		select {
		case st := <-ch:
			if st == network.Connected {
				return nil
			}
		case <-time.After(30 * time.Second):
			return errors.New("Timeout")
		}
	}
}
