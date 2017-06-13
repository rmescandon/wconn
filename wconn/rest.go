package wconn

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

const (
	versionURI       = "/v1"
	configurationURI = "/configuration"
)

var socketPath = os.Getenv("SNAP_COMMON") + "/sockets/control"

type serviceResponse struct {
	Result     map[string]interface{} `json:"result"`
	Status     string                 `json:"status"`
	StatusCode int                    `json:"status-code"`
	Type       string                 `json:"type"`
}

func unixDialer(_, _ string) (net.Conn, error) {
	return net.Dial("unix", socketPath)
}

func sendHTTPRequest(uri string, method string, body io.Reader) (*serviceResponse, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: unixDialer,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	realResponse := &serviceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&realResponse); err != nil {
		return nil, err
	}

	if realResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed: %s", realResponse.Result["message"])
	}

	return realResponse, nil
}

func defaultServiceURI() string {
	return fmt.Sprintf("http://unix%s", filepath.Join(versionURI, configurationURI))
}

// AccessPointIsUp checks if wifi-ap is up
func AccessPointIsUp() bool {
	response, err := sendHTTPRequest(defaultServiceURI(), "GET", nil)
	if err != nil {
		log.Printf("Error checking if AP is up: %v\n", err)
		return false
	}

	return !response.Result["disabled"].(bool)
}
