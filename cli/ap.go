package cli

import (
	"errors"
	"time"

	"github.com/greenbrew/wconn/network"
)

// ApCmd command to start a hotspot
type ApCmd struct {
	Args struct {
		Ssid string `positional-arg-name:"ssid" description:"The identifier for the WiFi network to connect to"`
		Psk  string `positional-arg-name:"pre-shared-key" description:"The pre-shared key to access the WiFi network"`
		Addr string `positional-arg-name:"addr" description:"The address of the hotspot"`
		Pref uint32 `positional-arg-name:"pref" description:"Prefix for the hotspot"`
	} `positional-args:"yes" required:"yes"`
	Security string `short:"s" long:"security" description:"The connection security system to use" default:"802-11-wireless-security"`
	KeyMgmt  string `short:"k" long:"keymgmt" description:"The passphrase key management schema to use" default:"wpa-psk"`
}

// Execute executes the hotspot launch command
func (cmd *ApCmd) Execute(args []string) error {
	m, err := network.NewManager()
	if err != nil {
		return err
	}

	ch, err := m.StartHotspot(cmd.Args.Ssid, cmd.Args.Psk, cmd.Security, cmd.KeyMgmt, cmd.Args.Addr, cmd.Args.Pref)
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
