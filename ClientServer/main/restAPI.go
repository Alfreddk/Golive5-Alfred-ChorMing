package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func getAllItems() (items []Item, err error) {

	items = []Item{}

	backendURL := "http://127.0.0.1:5000/api/v1/allitems/?key=2c78afaf-97da-4816-bbee-9ad239abb296"

	resp, err := http.Get(backendURL)
	if err != nil {
		return items, fmt.Errorf("Error: POST request - %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		err := json.Unmarshal(respData, &items)
		if err != nil {
			return items, fmt.Errorf("Error: JSON unmarshaling session - %v", err)
		}

		return items, nil
	}

	return items, errors.New("Error: resp.StatusCode is not 200")
}

func addNewItem(item Item) error {

	backendURL := "http://127.0.0.1:5000/api/v1/addnewitem/?key=2c78afaf-97da-4816-bbee-9ad239abb296"

	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("Error: JSON marshaling - %v", err)
	}

	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error: POST request - %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("Error: resp.StatusCode is not 200")
}

func editItem(item Item) error { // alfred 23.06.2022: not tested...

	backendURL := "http://127.0.0.1:5000/api/v1/edititem/?key=2c78afaf-97da-4816-bbee-9ad239abb296"

	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("Error: JSON marshaling - %v", err)
	}

	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error: POST request - %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("Error: resp.StatusCode is not 200")
}
