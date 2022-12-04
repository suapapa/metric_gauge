package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func getToken(user, pass string) (string, error) {
	data := map[string]interface{}{
		"username": user,
		"password": pass,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", vraptorAPIURL+"/login", body)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}

	if token, ok := result["access_token"]; ok {
		return token.(string), nil
	}

	return "", fmt.Errorf("failed to get token: %v", result)
}

// ---

type VraptorResource struct {
	CPU struct {
		CPUCount int     `json:"cpuCount"`
		CPUUsage float64 `json:"cpuUsage"`
	} `json:"cpu"`
	Memory struct {
		MemoryTotal int64   `json:"memoryTotal"`
		MemoryUsed  int     `json:"memoryUsed"`
		MemoryUsage float64 `json:"memoryUsage"`
	} `json:"memory"`
	Disk struct {
		DiskTotal int64   `json:"diskTotal"`
		DiskUsed  int64   `json:"diskUsed"`
		DiskUsage float64 `json:"diskUsage"`
	} `json:"disk"`
	DiskIo struct {
		DiskIoRead  int `json:"diskIoRead"`
		DiskIoWrite int `json:"diskIoWrite"`
	} `json:"diskIo"`
	NetworkIo struct {
		NetworkIoRead  int `json:"networkIoRead"`
		NetworkIoWrite int `json:"networkIoWrite"`
	} `json:"networkIo"`
}

func (r *VraptorResource) Get() error {
	req, err := http.NewRequest("GET", vraptorAPIURL+"/resource", nil)
	if err != nil {
		return errors.Wrap(err, "failed to get resource")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+vraptorToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to get resource")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return errors.Wrap(err, "failed to get resource")
	}
	return nil
}

// ---

type VraptorTemperature struct {
	Fahrenheit bool    `json:"fahrenheit"`
	Value      float64 `json:"value"`
}

func (t *VraptorTemperature) Get() error {
	req, err := http.NewRequest("GET", vraptorAPIURL+"/temperature", nil)
	if err != nil {
		return errors.Wrap(err, "failed to get resource")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+vraptorToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to get resource")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(t)
	if err != nil {
		return errors.Wrap(err, "failed to get resource")
	}
	return nil
}
