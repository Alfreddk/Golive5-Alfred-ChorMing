package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"goSch/golive5/arrds"
	"goSch/golive5/database"

	uuid "github.com/satori/go.uuid"
)

// userSignUp is a function that allows the Client-Server to send user sign up details to Backend-Server for registration.
// It takes in user details as database.User type variable and sends over the data to Backend-Server in JSON format.
// Upon successful user signup, it returns a UUID typed session ID and a status code of 200.
// Upon unsuccessful user signup, it returns a nil UUID typed session ID and a status code of 422.
func userSignUp(user User) (sessionID uuid.UUID, respStatusCode int, err error) { // Ivan to revise database.User to follow input details by client via req.FormValue.

	sessionID = uuid.UUID{}
	respStatusCode = 0

	backendURL := "http://127.0.0.1:8000/api/v1/signup/"

	jsonData, err := json.Marshal(user)
	if err != nil {
		return sessionID, respStatusCode, fmt.Errorf("Error: JSON marshaling - %v", err)
	}

	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return sessionID, respStatusCode, fmt.Errorf("Error: POST request - %v", err)
	}

	var session database.Session

	if resp.StatusCode == http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		err := json.Unmarshal(respData, &session)
		if err != nil {
			return sessionID, respStatusCode, fmt.Errorf("Error: JSON unmarshaling session - %v", err)
		}

		sessionID = session.ID
		respStatusCode = resp.StatusCode
		return sessionID, respStatusCode, nil
	}

	respStatusCode = resp.StatusCode

	return sessionID, respStatusCode, nil
}

// userLogin is a function that allows the Client-Server to send user login details to Backend-Server for authentication.
// It takes in username and password as string type variable and sends over the data to Backend-Server in JSON format.
// Upon successful user login, it returns a UUID typed session ID, user role in string, an arrds.Items type of slice of items listed by user and a status code of 200.
// Upon unsuccessful user login, it returns a nil UUID typed session ID, empty string for user role, an empty slice (arrds.Items type) and a status code other than 200.
func userLogin(username string, password string) (sessionID uuid.UUID, role string, items arrds.Items, respStatusCode int, err error) {

	sessionID = uuid.UUID{}
	role = ""
	items = arrds.Items{}
	respStatusCode = 0

	type UserPass struct {
		Username string
		Password string
	}
	userpass := UserPass{username, password}

	backendURL := "http://127.0.0.1:8000/api/v1/login/"

	jsonData, err := json.Marshal(userpass)
	if err != nil {
		return sessionID, role, items, respStatusCode, fmt.Errorf("Error: JSON marshaling - %v", err)
	}

	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return sessionID, role, items, respStatusCode, fmt.Errorf("Error: POST request - %v", err)
	}

	var session database.Session

	if resp.StatusCode == http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		err := json.Unmarshal(respData, &session)
		if err != nil {
			return sessionID, role, items, respStatusCode, fmt.Errorf("Error: JSON unmarshaling session - %v", err)
		}

		sessionID = session.ID
		role = session.Role

		err = json.Unmarshal([]byte(session.Data), &items)
		if err != nil {
			return sessionID, role, items, respStatusCode, fmt.Errorf("Error: JSON unmarshaling items - %v", err)
		}

		respStatusCode = resp.StatusCode

		return sessionID, role, items, respStatusCode, nil
	}

	return sessionID, role, items, respStatusCode, nil
}

// userLogout is a function that allows the Client-Server to send username and the updated slice of items listed by user to Backend-Server for storing.
// It takes in username as string type variable, slice of items as arrds.Item variable and sends over the data to Backend-Server in JSON format.
// Upon successful user logout, it returns a status code of 200.
func userLogout(username string, items arrds.Items) (respStatusCode int, err error) {

	respStatusCode = 0

	var userAndItems arrds.MyGiveGets
	userAndItems = arrds.MyGiveGets{username, items}

	backendURL := "http://127.0.0.1:8000/api/v1/logout/"

	jsonData, err := json.Marshal(userAndItems)
	if err != nil {
		return respStatusCode, fmt.Errorf("Error: JSON marshaling - %v", err)
	}

	fmt.Println(string(jsonData))

	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return respStatusCode, fmt.Errorf("Error: POST request - %v", err)
	}

	respStatusCode = resp.StatusCode

	return respStatusCode, nil
}

/*
// Initialise the bussiness logic or any handshake with back end server
func userInit() {

}

// Request for new session ID for new session
func userRequestSessionID() {

}

// get the user username, password and user type
func userGetUser() {

}

// returns boolean true of false if the user is a valid
*/
