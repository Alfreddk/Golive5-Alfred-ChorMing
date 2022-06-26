package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var mapUsers = map[string]User{}

// userInit performs getAllUsers() function call to retreive all users records and maps these records to mapUsers on runtime memory.
func userInit() {
	users, err := getAllUsers()
	if err != nil {
		fmt.Println(err)
		// log error
	}

	for _, v := range users {
		mapUsers[v.Username] = v
		userLastVisit[v.Username] = v.LastLogin
	}

}

func getAllUsers() (users []User, err error) {
	users = []User{}

	//backendURL := "http://127.0.0.1:5000/api/v1/allusers/?key=2c78afaf-97da-4816-bbee-9ad239abb296"
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

	//backendURL := "http://127.0.0.1:5000/api/v1/addnewuser/?key=2c78afaf-97da-4816-bbee-9ad239abb296"
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

	//backendURL := "http://127.0.0.1:5000/api/v1/edituser/?key=2c78afaf-97da-4816-bbee-9ad239abb296"
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

	//backendURL := "http://127.0.0.1:5000/api/v1/deleteuser/?key=2c78afaf-97da-4816-bbee-9ad239abb296"
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
