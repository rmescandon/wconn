package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rmescandon/wconn/wconn"
	"github.com/rmescandon/wconn/web"
)

const idle = 5 * time.Second

var accessPoints map[string]string

func main() {

	cWifi := make(chan bool)
	quitWifi := make(chan bool)

	cAp := make(chan bool)
	quitAp := make(chan bool)

	cSSIDs := make(chan string)
	quitSSIDs := make(chan bool)

	accessPoints = make(map[string]string)

	//TODO implement go routine to search for SSIDs

	// search for external wifi status. Poll for changes
	go func() {
		// TODO verify if initially set to false or get initial real state
		b := false
		for {
			select {
			case <-quitWifi:
				return
			default:
				new := wconn.ConnectedToWifi()
				if new != b {
					cWifi <- new
					b = new
				}
				time.Sleep(idle)
			}
		}
	}()

	// search for AP status. Poll for changes
	go func() {
		// TODO verify if initially set to false or get initial real state
		b := false
		for {
			select {
			case <-quitAp:
				return
			default:
				new := wconn.AccessPointIsUp()
				if new != b {
					cAp <- new
					b = new
				}
				time.Sleep(idle)
			}
		}
	}()

	// search for available access points
	go func() {
		for {
			select {
			case <-quitSSIDs:
				return
			default:
				aps, err := wconn.AccessPoints()
				if err != nil {
					log.Printf("Error retrieving available access points: %v\n", err)
					continue
				}
				for _, ap := range aps {
					cSSIDs <- ap
				}
				time.Sleep(idle)
			}
		}
	}()

	for {
		select {
		case b := <-cWifi:
			if b {
				fmt.Println("WCONN - Connected to external wifi")

				if err := web.StartOperationalPortal(); err != nil {
					fmt.Printf("Error starting operational server: %v\n", err)
				}

			} else {
				fmt.Println("WCONN - Disconnected from external wifi")

				if err := web.StopOperationalPortal(); err != nil {
					log.Printf("Error shutting down operational portal: %v\n", err)
				}
			}
		case b := <-cAp:
			if b {
				fmt.Println("WCONN - Local Access Point UP")

				if err := web.StartManagementPortal(accessPoints); err != nil {
					fmt.Printf("Error starting management server: %v\n", err)
				}

			} else {
				fmt.Println("WCONN - Local Access Point DOWN")

				if err := web.StopManagementPortal(); err != nil {
					log.Printf("Error shutting down management portal: %v\n", err)
				}
			}
		case ssid := <-cSSIDs:
			// TODO. Temporary set key and value to ssid. In future maybe it is needed
			// to associate ssid to device entry in dbus
			accessPoints[ssid] = ssid
		}
	}
}
