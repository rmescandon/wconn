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
)

// Enum of available server options
const (
	None Server = 0 + iota
	Management
	Operational
)

// Server defines an enum of servers
type Server int

// Current active server instance. None if any is enabled at this moment
var Current = None

// StartManagementPortal starts server in management mode
func StartManagementPortal(ssids map[string]string) error {

	if Current == Management {
		return nil
	}

	if Current == Operational {
		err := stop()
		if err != nil {
			return err
		}
	}

	// change current instance asap we manage this server
	Current = Management

	err := listenAndServe(managementRouter(ssids))
	if err != nil {
		Current = None
		return err
	}

	return nil
}

// StartOperationalPortal starts server in operational mode
func StartOperationalPortal() error {
	if Current == Operational {
		return nil
	}

	if Current == Management {
		err := stop()
		if err != nil {
			return err
		}
	}

	// change current instance asap we manage this server
	Current = Operational

	err := listenAndServe(operationalRouter())
	if err != nil {
		Current = None
		return err
	}

	return nil
}

// StopManagementPortal shutdown server management mode.
func StopManagementPortal() error {

	if Current != Management {
		return fmt.Errorf("Cannot stop management portal, as current server is other")
	}

	if !running() {
		return fmt.Errorf("Cannot stop a not running server")
	}

	err := stop()
	if err != nil {
		return err
	}

	Current = None
	return nil
}

// StopOperationalPortal shutdown server operational mode. If operational server is not up, returns error
func StopOperationalPortal() error {
	if Current != Operational {
		return fmt.Errorf("Cannot stop operational portal, as current server is other")
	}

	if !running() {
		return fmt.Errorf("Cannot stop a not running server")
	}

	err := stop()
	if err != nil {
		return err
	}

	Current = None
	return nil
}
