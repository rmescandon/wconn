package main

import (
	"time"

	"log"

	"fmt"

	"github.com/rmescandon/wconn/wconn"
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

				//TODO TRACE
				fmt.Println("Connected to external wifi")

			} else {
				// disconnected to external wifi

				//TODO TRACE
				fmt.Println("Disconnected to external wifi")
			}
		case b := <-cAp:
			if b {
				// AP up

				//TODO TRACE
				fmt.Println("Local Access Point UP")

			} else {
				// AP down

				//TODO TRACE
				fmt.Println("Local Access Point DOWN")
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
