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

//var sqlDBConnection string
var cfg mysql.Config // configuration for DSN

// init() initialises the server
func init() {

	// set path for the env file
	//envFile := path.Join("..", "config", ".env")
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
	//fmt.Printf("User http://%s\n", hostPort)

	// SQL DB Data Source Name config
	cfg = mysql.Config{
		User:   os.Getenv("SQL_USER"),
		Passwd: os.Getenv("SQL_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("SQL_ADDR"),
		DBName: os.Getenv("SQL_DB"),
	}
}

// validate key from the query key-value pair
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

func allItems(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Trace.Fatalln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer db.Close()

	buffer := sqlGetAllItems(db)

	json.NewEncoder(w).Encode(buffer)

	logger.Info.Println("All items successfully transmitted.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)
}

func addNewItem(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-type") == "application/json" {

		var item Item

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Trace.Fatalln(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		} else {
			err = json.Unmarshal(reqBody, &item)
			if err != nil {
				logger.Trace.Fatalln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			if item.ID == "" { // To check if item sent by frontend server is empty.
				logger.Trace.Fatalln("Error: Item is empty.")
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			db, err := sql.Open("mysql", cfg.FormatDSN())
			if err != nil {
				logger.Trace.Fatalln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}
			defer db.Close()

			logger.Info.Println("Item successfully retrieved.")
			logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

			sqlAddNewItem(db, item)

		}
	}
}

func editItem(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") == "application/json" {

		var item Item

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Trace.Fatalln(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		} else {
			err = json.Unmarshal(reqBody, &item)
			if err != nil {
				logger.Trace.Fatalln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			if item.ID == "" { // To check if item sent by frontend server is empty.
				logger.Trace.Fatalln("Error: Item is empty.")
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			db, err := sql.Open("mysql", cfg.FormatDSN())
			if err != nil {
				logger.Trace.Fatalln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}
			defer db.Close()

			logger.Info.Println("Item updates successfully retrieved.")
			logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

			sqlEditItem(db, item)
		}
	}
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Trace.Fatalln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer db.Close()

	buffer := sqlGetAllUsers(db)

	json.NewEncoder(w).Encode(buffer)

	logger.Info.Println("All items successfully transmitted.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

}

func addNewUser(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-type") == "application/json" {

		var user User

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Trace.Fatalln(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		} else {
			err = json.Unmarshal(reqBody, &user)
			if err != nil {
				logger.Trace.Fatalln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			if user.ID == "" { // To check if user sent by frontend server is empty.
				logger.Trace.Fatalln("Error: Item is empty.")
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			db, err := sql.Open("mysql", cfg.FormatDSN())
			if err != nil {
				logger.Trace.Fatalln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}
			defer db.Close()

			logger.Info.Println("User successfully retrieved.")
			logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

			sqlAddNewUser(db, user)

		}
	}
}

func editUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") == "application/json" {

		var user User

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Trace.Fatalln(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		} else {
			err = json.Unmarshal(reqBody, &user)
			if err != nil {
				logger.Trace.Fatalln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			if user.ID == "" { // To check if user sent by frontend server is empty.
				logger.Trace.Fatalln("Error: Item is empty.")
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}

			db, err := sql.Open("mysql", cfg.FormatDSN())
			if err != nil {
				logger.Trace.Fatalln(err)
				logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
			}
			defer db.Close()

			logger.Info.Println("User updates successfully retrieved.")
			logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

			sqlEditUser(db, user)
		}
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	var user User

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Trace.Fatalln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	} else {
		err = json.Unmarshal(reqBody, &user)
		if err != nil {
			logger.Trace.Fatalln(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		}

		if user.ID == "" { // To check if user sent by frontend server is empty.
			logger.Trace.Fatalln("Error: Item is empty.")
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		}

		db, err := sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			logger.Trace.Fatalln(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		}
		defer db.Close()

		logger.Info.Println("User to be deleted successfully retrieved.")
		logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

		sqlDeleteUser(db, user)
	}

}

/*
// home is the handler for "/api/v1/" resource
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the REST API Server!")
}

// allcourses is the handler for "/api/v1/courses" resource
func allcourses(w http.ResponseWriter, r *http.Request) {

	// Use mysql as driverName and a valid DSN as dataSourceName:
	// user set up password that can access this db connection
	// db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:58710)/courseDB")
	// db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:58710)/courseDB")
	db, err := sql.Open("mysql", cfg.FormatDSN())

	// handle error
	if err != nil {
		panic(err.Error())
	}
	// defer the close till after the main function has finished executing
	defer db.Close()
	//	fmt.Println("Database opened")

	fmt.Fprintf(w, "List of all courses\n")

	//var bufferMap map[string]interface{}
	bufferMap := GetRecords(db)
	//	fmt.Println("BufferMap :", bufferMap)

	// map assertion to interface to string
	for k, v := range bufferMap {
		fmt.Fprintln(w, k, v.(string))
	}

	// returns all the courses in JSON
	json.NewEncoder(w).Encode(bufferMap)
}

// course() is the hanlder for "/api/v1/courses/{courseid}" resource
func course(w http.ResponseWriter, r *http.Request) {

	// vakidate key for parameter key-value
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Use mysql as driverName and a valid DSN as dataSourceName:
	// user set up password that can access this db connection
	//db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:58710)/courseDB")
	//db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:58710)/courseDB")
	db, err := sql.Open("mysql", cfg.FormatDSN())

	// handle error
	if err != nil {
		panic(err.Error())
	}
	// defer the close till after the main function has finished executing
	defer db.Close()
	//	fmt.Println("Database opened")

	// mux.Vars(r) is the variable immediately after the URL
	params := mux.Vars(r)
	//fmt.Println("parameter =", params)

	// Get does not have a body so only header
	if r.Method == "GET" {

		// check if there is a row for this record with the ID
		bufferMap, err := GetOneRecord(db, params["courseid"])

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
			return
		}

		//map assertion to interface to map[string]string
		for k, v := range bufferMap {
			fmt.Fprintln(w, k, v.(string))
		}

		// return the specific course in Json
		json.NewEncoder(w).Encode(bufferMap[params["courseid"]])
	}

	// Delete may have a body but not encouraged, safest not to use
	if r.Method == "DELETE" {
		_, err := GetOneRecord(db, params["courseid"])

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
			return
		}

		DeleteRecord(db, params["courseid"])
		// 	delete(courses, params["courseid"])
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - Course deleted: " + params["courseid"]))
	}

	// check for json application
	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new course
		if r.Method == "POST" { // check request method
			// read the string sent to the service
			var newCourse courseInfo
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// parse JSON to object data structure
				json.Unmarshal(reqBody, &newCourse)
				if newCourse.Title == "" { // empty title
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply course " + "information " + "in JSON format"))
					return
				} // check if course exists; add only if // course does not exist

				// check if there is a row for this record with the ID
				_, err := GetOneRecord(db, params["courseid"])

				fmt.Println("Title", newCourse.Title)
				if err != nil {
					// Row not found, so add new row (new record)
					if err == errEmptyRow {
						InsertRecord(db, params["courseid"], newCourse.Title)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Course added: " + params["courseid"] + " Title: " + newCourse.Title))
					} else {
						// some sql error if any such error
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("500 - Internal Server Error"))
					}
				} else {
					w.WriteHeader(http.StatusConflict) // course key already exist
					w.Write([]byte("409 - Duplicate course ID"))
				}
			} else {
				// Problem with the body from response
				w.WriteHeader(http.StatusUnprocessableEntity) // error
				w.Write([]byte("422 - Please supply course information " + "in JSON format"))
			}
		}

		//---PUT is for creating or updating exiting course---
		if r.Method == "PUT" {
			var newCourse courseInfo
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// parse JSON to object data structure
				json.Unmarshal(reqBody, &newCourse)
				if newCourse.Title == "" { // empty title in body
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply course " + " information " + "in JSON format"))
					return
				} // check if course exists; add only if // course does not exist

				// check if there is a row for this record with the ID
				_, err := GetOneRecord(db, params["courseid"])

				fmt.Println("Title", newCourse.Title)

				if err != nil {
					// Row not found, so creat new row
					// 	courses[params["courseid"]] = newCourse // create the key-value
					if err == errEmptyRow {
						InsertRecord(db, params["courseid"], newCourse.Title)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Course added: " + params["courseid"]))
					} else {
						// some sql error if any such error
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("500 - Internal Server Error"))
					}
				} else {
					// Edit row if row exist
					EditRecord(db, params["courseid"], newCourse.Title)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Course updated: " + params["courseid"] + " Title: " + newCourse.Title))
				}

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity) // error
				w.Write([]byte("422 - Please supply " + "course information " + "in JSON format"))
			}
		}
	}
}
*/

//
//
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