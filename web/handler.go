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
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

const (
	managementTemplatePath  = "/templates/management.html"
	connectingTemplatePath  = "/templates/connecting.html"
	operationalTemplatePath = "/templates/operational.html"
	refreshingTemplatePath  = "/templates/refreshing.html"
)

var resourcesPath = filepath.Join(os.Getenv("SNAP"), "static")

// Data interface representing any data included in a template
type Data interface{}

// SsidsData dynamic data to fulfill the SSIDs page template
type SsidsData struct {
	Ssids []string
}

// ConnectingData dynamic data to fulfill the connect result page template
type ConnectingData struct {
	Ssid string
}

type noData struct {
}

func execTemplate(w http.ResponseWriter, templatePath string, data Data) {
	templateAbsPath := filepath.Join(resourcesPath, templatePath)
	t, err := template.ParseFiles(templateAbsPath)
	if err != nil {
		fmt.Printf("Error loading the template at %v : %v\n", templatePath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		fmt.Printf("Error executing the template at %v : %v\n", templatePath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func managementHandler(w http.ResponseWriter, r *http.Request, ssids []string) {
	data := SsidsData{Ssids: ssids}
	execTemplate(w, managementTemplatePath, data)
}

func operationalHandler(w http.ResponseWriter, r *http.Request) {
	execTemplate(w, operationalTemplatePath, noData{})
}
