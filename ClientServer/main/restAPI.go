package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func getAllUsers() (users []User, err error) {
	users = []User{}

	backendURL := "http://" + backendHost + ":" + backendPort + "/api/v1/allusers/?key=" + urlKey

	resp, err := http.Get(backendURL)
	if err != nil {
		return users, fmt.Errorf("Error: POST request - %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		err := json.Unmarshal(respData, &users)
		if err != nil {
			return users, fmt.Errorf("Error: JSON unmarshaling session - %v", err)
		}

		return users, nil
	}

	return users, errors.New("Error: resp.StatusCode is not 200")
}

func addNewUser(user User) error {

	backendURL := "http://" + backendHost + ":" + backendPort + "/api/v1/addnewuser/?key=" + urlKey

	jsonData, err := json.Marshal(user)
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

func editUser(user User) error {

	backendURL := "http://" + backendHost + ":" + backendPort + "/api/v1/edituser/?key=" + urlKey

	jsonData, err := json.Marshal(user)
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

func delUser(user User) error {

	backendURL := "http://" + backendHost + ":" + backendPort + "/api/v1/deleteuser/?key=" + urlKey

	jsonData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Error: JSON marshaling - %v", err)
	}

	req, err := http.NewRequest(http.MethodDelete, backendURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error: DELETE request - %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error: DELETE request - %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("Error: resp.StatusCode is not 200")

}

func getAllItems() (items []Item, err error) {

	items = []Item{}

	backendURL := "http://" + backendHost + ":" + backendPort + "/api/v1/allitems/?key=" + urlKey

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

	backendURL := "http://" + backendHost + ":" + backendPort + "/api/v1/addnewitem/?key=" + urlKey

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

func editItem(item Item) error {

	backendURL := "http://" + backendHost + ":" + backendPort + "/api/v1/edititem/?key=" + urlKey

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
