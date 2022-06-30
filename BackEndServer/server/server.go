//
package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"BackEndServer/logger"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var urlKey string
var hostPort string

var cfg mysql.Config // configuration for DSN

// init() initialises the server
func init() {

	// set path for the env file
	envFile := path.Join("config", ".env")

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalln("Error loading .env file: ", err)
	}

	// getting env variables SERVER_NAME, SERVER_HOST, SERVER_PORT and SERVER_URLKEY
	serverName := os.Getenv("SERVER_NAME")
	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")
	urlKey = os.Getenv("SERVER_URLKEY")

	// Create Host Port from environment variable
	hostPort = fmt.Sprintf("%s:%s", serverHost, serverPort)

	fmt.Printf("Server Name: %s\n", serverName)

	// SQL DB Data Source Name config
	cfg = mysql.Config{
		User:   os.Getenv("SQL_USER"),
		Passwd: os.Getenv("SQL_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("SQL_ADDR"),
		DBName: os.Getenv("SQL_DB"),
	}
}

// validKey validate key from the query key-value pair
func validKey(r *http.Request) bool {
	// query() get the key-value pair after URL
	v := r.URL.Query()

	if key, ok := v["key"]; ok {
		if key[0] == urlKey {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// allItems executes HTTP GET request, function calls sqlGetAllItems to retrieve all items from database and returns all items via HTTP response in JSON format.
// It also checks if a valid key has been provided.
func allItems(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		msg := "Status 404 - Invalid key."
		w.Write([]byte(msg))
		logger.Trace.Println(msg)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		return
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Trace.Panicln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer db.Close()

	buffer := sqlGetAllItems(db)

	json.NewEncoder(w).Encode(buffer)

	logger.Info.Println("All items successfully transmitted.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)
}

// addNewItem executes HTTP POST request, retrieve request body in JSON format, function calls sqlAddNewItem to add an item to database.
// It also checks if a valid key has been provided.
func addNewItem(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		msg := "Status 404 - Invalid key."
		w.Write([]byte(msg))
		logger.Trace.Println(msg)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		return
	}

	if r.Header.Get("Content-type") == "application/json" {

		var item Item

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Status 422 - Invalid JSON format"))
			logger.Trace.Println(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		} else {
			err = json.Unmarshal(reqBody, &item)
			if err != nil {
				logger.Trace.Println(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			if item.Name == "" { // To check if item sent by frontend server is empty.
				w.WriteHeader(http.StatusUnprocessableEntity)
				msg := "Error: Item is empty."
				w.Write([]byte(msg))
				logger.Trace.Println(msg)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			db, err := sql.Open("mysql", cfg.FormatDSN())
			if err != nil {
				logger.Trace.Panicln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}
			defer db.Close()

			logger.Info.Println("Item successfully retrieved.")
			logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

			sqlAddNewItem(db, item)

		}
	}
}

// editItem executes HTTP POST request, retrieve request body in JSON format, function calls sqlEditItem to update an item in database.
// It also checks if a valid key has been provided.
func editItem(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		msg := "Status 404 - Invalid key."
		w.Write([]byte(msg))
		logger.Trace.Println(msg)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		return
	}

	if r.Header.Get("Content-type") == "application/json" {

		var item Item

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Status 422 - Invalid JSON format"))
			logger.Trace.Println(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		} else {
			err = json.Unmarshal(reqBody, &item)
			if err != nil {
				logger.Trace.Println(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			if item.ID == "" { // To check if item sent by frontend server is empty.
				w.WriteHeader(http.StatusUnprocessableEntity)
				msg := "Error: Item is empty."
				w.Write([]byte(msg))
				logger.Trace.Println(msg)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			db, err := sql.Open("mysql", cfg.FormatDSN())
			if err != nil {
				logger.Trace.Panicln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}
			defer db.Close()

			logger.Info.Println("Item updates successfully retrieved.")
			logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

			sqlEditItem(db, item)
		}
	}
}

// allUsers executes HTTP GET request, function calls sqlGetAllUsers to retrieve all users from database and returns all users via HTTP response in JSON format.
// It also checks if a valid key has been provided.
func allUsers(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		msg := "Status 404 - Invalid key."
		w.Write([]byte(msg))
		logger.Trace.Println(msg)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		return
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Trace.Panicln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer db.Close()

	buffer := sqlGetAllUsers(db)

	json.NewEncoder(w).Encode(buffer)

	logger.Info.Println("All items successfully transmitted.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

}

// addNewUser executes HTTP POST request, retrieve request body in JSON format, function calls sqlAddNewUser to add an user to database.
// It also checks if a valid key has been provided.
func addNewUser(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		msg := "Status 404 - Invalid key."
		w.Write([]byte(msg))
		logger.Trace.Println(msg)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		return
	}

	if r.Header.Get("Content-type") == "application/json" {

		var user User

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Status 422 - Invalid JSON format"))
			logger.Trace.Println(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		} else {
			err = json.Unmarshal(reqBody, &user)
			if err != nil {
				logger.Trace.Println(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			if user.Name == "" { // To check if user sent by frontend server is empty.
				w.WriteHeader(http.StatusUnprocessableEntity)
				msg := "Error: User is empty."
				w.Write([]byte(msg))
				logger.Trace.Println(msg)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			db, err := sql.Open("mysql", cfg.FormatDSN())
			if err != nil {
				logger.Trace.Panicln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}
			defer db.Close()

			logger.Info.Println("User successfully retrieved.")
			logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

			sqlAddNewUser(db, user)

		}
	}
}

// editUser executes HTTP POST request, retrieve request body in JSON format, function calls sqlEditUser to update an user in database.
// It also checks if a valid key has been provided.
func editUser(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		msg := "Status 404 - Invalid key."
		w.Write([]byte(msg))
		logger.Trace.Println(msg)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		return
	}

	if r.Header.Get("Content-type") == "application/json" {

		var user User

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Status 422 - Invalid JSON format"))
			logger.Trace.Println(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		} else {
			err = json.Unmarshal(reqBody, &user)
			if err != nil {
				logger.Trace.Println(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			if user.ID == "" { // To check if user sent by frontend server is empty.
				w.WriteHeader(http.StatusUnprocessableEntity)
				msg := "Error: User is empty."
				w.Write([]byte(msg))
				logger.Trace.Println(msg)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			db, err := sql.Open("mysql", cfg.FormatDSN())
			if err != nil {
				logger.Trace.Panicln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}
			defer db.Close()

			logger.Info.Println("User updates successfully retrieved.")
			logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

			sqlEditUser(db, user)
		}
	}
}

// deleteUser executes HTTP DELETE request, retrieve request body in JSON format, function calls sqlDeleteUser to delete an user from database.
// It also checks if a valid key has been provided.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		msg := "Status 404 - Invalid key."
		w.Write([]byte(msg))
		logger.Trace.Println(msg)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		return
	}

	var user User

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Status 422 - Invalid JSON format"))
		logger.Trace.Println(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	} else {
		err = json.Unmarshal(reqBody, &user)
		if err != nil {
			logger.Trace.Println(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		}

		if user.ID == "" { // To check if user sent by frontend server is empty.
			w.WriteHeader(http.StatusUnprocessableEntity)
			msg := "Error: User is empty."
			w.Write([]byte(msg))
			logger.Trace.Println(msg)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		}

		db, err := sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			logger.Trace.Panicln(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		}
		defer db.Close()

		logger.Info.Println("User to be deleted successfully retrieved.")
		logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

		sqlDeleteUser(db, user)
	}

}

// ExecuteServer executes the backend server.
func ExecuteServer() {

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/allitems/", allItems).Methods("GET")
	router.HandleFunc("/api/v1/addnewitem/", addNewItem).Methods("POST")
	router.HandleFunc("/api/v1/edititem/", editItem).Methods("POST")
	router.HandleFunc("/api/v1/allusers/", allUsers).Methods("GET")
	router.HandleFunc("/api/v1/addnewuser/", addNewUser).Methods("POST")
	router.HandleFunc("/api/v1/edituser/", editUser).Methods("POST")
	router.HandleFunc("/api/v1/deleteuser/", deleteUser).Methods("DELETE")

	fmt.Printf("Listening on http://%s\n", hostPort)

	log.Fatal(http.ListenAndServe(hostPort, router))
}
