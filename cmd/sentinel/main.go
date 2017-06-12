package main

import (
	"time"

	"log"

	"fmt"

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
				// connected to external wifi
				// TODO: Start operational

				//TODO TRACE
				fmt.Println("WCONN - Connected to external wifi")

				if err := web.ListenAndServe(web.OperationalHandler()); err != nil {
					fmt.Printf("Error starting operational server: %v\n", err)
				}

			} else {
				// disconnected to external wifi

				//TODO TRACE
				fmt.Println("WCONN - Disconnected from external wifi")

				//TODO IT IS NEEDED CONTROL EVERY SERVER SERPARATEDLY
				// if err := web.Stop(); err != nil {
				// 	log.Printf("Error shutting down server: %v\n", err)
				// }
			}
		case b := <-cAp:
			if b {
				// AP up

				//TODO TRACE
				fmt.Println("WCONN - Local Access Point UP")

				if err := web.ListenAndServe(web.ManagementHandler(accessPoints)); err != nil {
					fmt.Printf("Error starting management server: %v\n", err)
				}

			} else {
				// AP down

				//TODO TRACE
				fmt.Println("WCONN - Local Access Point DOWN")

				//TODO IT IS NEEDED CONTROL EVERY SERVER SERPARATEDLY
				// if err := web.Stop(); err != nil {
				// 	log.Printf("Error shutting down server: %v\n", err)
				// }
			}
		case ssid := <-cSSIDs:

			//TODO TRACE
			fmt.Printf("SSID: %v\n", ssid)

			// TODO. Temporary set key and value to ssid. In future maybe it is needed
			// to associate ssid to device entry in dbus
			accessPoints[ssid] = ssid
		}
	}
}
