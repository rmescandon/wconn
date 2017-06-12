/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package web

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	telnet "github.com/reiver/go-telnet"
)

var address = ":8080"
var listener net.Listener
var done chan bool

type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts incoming tcp connections
func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return tc, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// ListenAndServe starts http server in port 8080. Stops any previous one if needed
func ListenAndServe(handler http.Handler) error {

	if runningOn(address) {
		// stop previous instance before starting new one
		err := Stop()
		if err != nil {
			return fmt.Errorf("Could not stop the current running instance before starting a new one %v", err)
		}
	}

	srv := &http.Server{Addr: address, Handler: handler}
	// channel needed to communicate real server shutdown, as after calling listener.Close()
	// it can take several milliseconds to really stop the listening.
	done = make(chan bool)

	var err error
	listener, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// launch goroutine to check server state changes after startup is triggered
	go func() {
		retries := 10
		idle := 10 * time.Millisecond
		for ; !runningOn(address) && retries > 0; retries-- {
			time.Sleep(idle)
			idle *= 2
		}

		if retries == 0 {
			log.Print("Server could not be started")
			return
		}
	}()

	// launching server in a goroutine for not blocking
	go func() {
		if listener != nil {
			err := srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
			if err != nil {
				log.Printf("HTTP Server closing - %v", err)
			}
			// notify server real stop
			done <- true
		}

		close(done)
	}()

	return nil
}

// Stop stops current http server
func Stop() error {

	if !runningOn(address) {
		return fmt.Errorf("Already stopped")
	}

	if listener == nil {
		return fmt.Errorf("Already closed")
	}

	err := listener.Close()
	if err != nil {
		return err
	}
	listener = nil

	// wait for server real shutdown confirmation
	<-done
	return nil
}

func runningOn(address string) bool {

	if strings.HasPrefix(address, ":") {
		address = "localhost" + address
	}
	// telnet to check server is alive
	caller := telnet.StandardCaller
	err := telnet.DialToAndCall(address, caller)
	if err != nil {
		return false
	}
	return true
}
