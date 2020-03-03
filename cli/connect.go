package cli

// ConnectCmd command to connect to a certain available SSID
type ConnectCmd struct {
	Args struct {
		Ssid string `positional-arg-name:"ssid" description:"Connects to a WiFi network"`
	} `positional-args:"yes"`
}

// Execute executes the connect command
func Execute(args []string) error {

	// TODO IMPLEMENT
	return nil
}
