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

// Initialise the bussiness logic or any handshake with back end server
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

	fmt.Println(mapUsers)

}

func getAllUsers() (users []User, err error) {
	users = []User{}

	backendURL := "http://127.0.0.1:5000/api/v1/allusers/?key=2c78afaf-97da-4816-bbee-9ad239abb296"

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

	backendURL := "http://127.0.0.1:5000/api/v1/addnewuser/?key=2c78afaf-97da-4816-bbee-9ad239abb296"

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

	backendURL := "http://127.0.0.1:5000/api/v1/edituser/?key=2c78afaf-97da-4816-bbee-9ad239abb296"

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

/*
// Request for new session ID for new session
func userRequestSessionID() {

}

// get the user username, password and user type
func userGetUser() {

}

// returns boolean true of false if the user is a valid
*/
