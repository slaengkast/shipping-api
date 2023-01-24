package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type client struct {
	baseUrl string
	apiUrl  string
}

func NewClient(baseUrl string) client {
	return client{baseUrl, fmt.Sprintf("%s/%s", baseUrl, "api/shipping")}
}

type bookingResponse struct {
	Id string `json:"id"`
}

func (c client) BookShipping(origin, destination string, weight float32) (string, error) {
	data := map[string]interface{}{"origin": origin, "destination": destination, "weight": weight}
	input, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(c.apiUrl, "application/json", bytes.NewBuffer(input))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response bookingResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

func (c client) GetBooking(id string) (map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", c.apiUrl, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type healthResponse struct {
	Status string `json:"status"`
}

func (c client) Health() (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", c.baseUrl, "health"))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response healthResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Status, nil
}
