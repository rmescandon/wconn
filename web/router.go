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
	"net/http"

	"github.com/gorilla/mux"
)

// managementHandler handles requests for web UI when AP is up
func managementRouter(ssidsMap map[string]string) *mux.Router {
	router := mux.NewRouter()

	ssids := make([]string, 0, len(ssidsMap))
	for k := range ssidsMap {
		ssids = append(ssids, k)
	}

	// only start in AP interface: 10.0.60.1
	address = "10.0.60.1:8080"

	// Pages routes
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		managementHandler(w, r, ssids)
	})
	// router.HandleFunc("/connect", ConnectHandler).Methods("POST")
	// router.HandleFunc("/hashit", HashItHandler).Methods("POST")
	// router.HandleFunc("/refresh", RefreshHandler).Methods("GET")

	// Resources path
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(resourcesPath)))
	router.PathPrefix("/static/").Handler(fs)

	return router
}

// operationalHandler handles request for web UI when connected to external WIFI
func operationalRouter() *mux.Router {
	router := mux.NewRouter()

	address = ":8080"

	router.HandleFunc("/", operationalHandler).Methods("GET")
	// router.HandleFunc("/disconnect", DisconnectHandler).Methods("GET")
	// router.HandleFunc("/hashit", HashItHandler).Methods("POST")

	// Resources path
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(resourcesPath)))
	router.PathPrefix("/static/").Handler(fs)
	return router
}
