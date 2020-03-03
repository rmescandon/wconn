package cli

// Command the struct holding the cli commands
var Command Cmd

//Cmd struct with the commands to expose through cli
type Cmd struct {
	Connect    ConnectCmd    `command:"connect" description:"Connects to a WiFi network"`
	Disconnect DisconnectCmd `command:"disconnect" description:"Disconnects from current WiFi connected network"`
	List       ListCmd       `command:"list" alias:"ls" description:"Lists available WiFi networks"`
	Connected  ConnectedCmd  `command:"connected" description:"Returns true if connected to a WiFi network"`
}
