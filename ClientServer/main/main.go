// Package main is the code for Venue Booking Management System
// Created for Go in Action 2
// By Tan Chor Ming
// 18 May 2022
// Development Time : 170 hours (Go in Action 1 and 2)
package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
type user struct {
	Username  string //username
	Password  []byte //password hash
	Name      string
	Address   string
	Postal    string
	Telephone string
	LastLogin string
}*/

var tpl *template.Template

var mapSessions = map[string]string{}
var mapDeletedUser = map[string]string{}
var mapDeletedSession = map[string]string{}

var userLastVisit = map[string]string{}

// create an empty linked list for booking
var myList = &linkedList{nil, 0}

// init is the system initialisation
var vrsHost, vrsPort string
var adminSubName string
var errLogDir string

// init() initialises the system
// Set up the environment
// Set up the logger
// Set up database for booking and Properties
// Initialises admin password
func init() {

	// set path for the env file
	envFile := path.Join("..", "config", ".env")

	//err := godotenv.Load(".env")
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalln("Error loading .env file: ", err)
	}

	// getting env variables SITE_TITLE and DB_HOST
	siteTitle := os.Getenv("SITE_TITLE")
	vrsHost = os.Getenv("VRS_HOST")
	vrsPort = os.Getenv("VRS_PORT")
	// Error log relative dir
	errLogDir = os.Getenv("VRS_ERR_LOG_DIR")
	//fmt.Println("Error Dir", errLogDir)

	// system allow for 3 admin users
	adminSubName = os.Getenv("VRS_ADMIN_SUBNAME")
	adminName1 := os.Getenv("VRS_USERNAME1")
	password1 := os.Getenv("VRS_PASSWORD1")
	adminName2 := os.Getenv("VRS_USERNAME2")
	password2 := os.Getenv("VRS_PASSWORD2")
	adminName3 := os.Getenv("VRS_USERNAME3")
	password3 := os.Getenv("VRS_PASSWORD3")

	fmt.Printf("%s = %s\n", "Site Title", siteTitle)
	fmt.Printf("Use https:// %s:%s\n", vrsHost, vrsPort)
	// fmt.Printf("godotenv : %s = %s \n", "VRS Host", vrsHost)
	// fmt.Printf("godotenv : %s = %s \n", "VRS Port", vrsPort)
	// fmt.Printf("godotenv : %s = %s \n", "VRS Username", username)
	// fmt.Printf("godotenv : %s = %s \n", "VRS Password", password)

	// Set up Logging
	warningsFile := path.Join("..", errLogDir, "warnings.log")
	errorsFile := path.Join("..", errLogDir, "errors.log")

	// check for existence of file in the directory
	// adding this part will prevent automatic creation of file.
	_, err = os.Stat(warningsFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("%s not found", warningsFile)
			return
		}
	}
	_, err = os.Stat(errorsFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("%s not found", errorsFile)
			return
		}
	}

	// Open a file for warnings.
	warningsFH, err := os.OpenFile(warningsFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open %s", warningsFile)
		return
	}
	//  Note init() will execute defer function so has to comment off
	//	defer warningsFH.Close()

	// Open a file for errors.
	errorsFH, err := os.OpenFile(errorsFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open %s", errorsFile)
		return
	}
	//  Note init() will execute defer function so has to comment off
	//	defer errorsFH.Close()

	// Create a multi writer for errors.
	multi := io.MultiWriter(errorsFH, os.Stderr)

	// Init the log package for each message type.
	initLog(ioutil.Discard, os.Stdout, warningsFH, multi)

	// Test each log type.
	// Trace.Println("Test Trace")
	// Warning.Println("Test Warining")
	// Error.Println("Test Error")
	// Info.Println("Venue Booking System Initialisation Completed")

	// application initialisation
	//propertyInit() // create the property features
	//bookingsInit() // create some initial bookings
	// template initialisation
	tpl = template.Must(template.ParseGlob("templates/*"))
	bPassword1, _ := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	mapUsers[adminName1] = User{"", adminName1, string(bPassword1), "Staff1", "", "", "", "", ""}
	bPassword2, _ := bcrypt.GenerateFromPassword([]byte(password2), bcrypt.DefaultCost)
	mapUsers[adminName2] = User{"", adminName2, string(bPassword2), "Staff2", "", "", "", "", ""}
	bPassword3, _ := bcrypt.GenerateFromPassword([]byte(password3), bcrypt.DefaultCost)
	mapUsers[adminName3] = User{"", adminName3, string(bPassword3), "Staff3", "", "", "", "", ""}

	userInit()

	// Initialise the business logic
	bizInit()

}

// main package starts here
func main() {
	// Test each log type.
	// Trace.Println("Test Trace")
	// Warning.Println("Test Warining")
	// Error.Println("Test Error")
	// Info.Println("Venue Booking System Initialisation Completed")

	// Experiment with gorilla/mux
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/deleteSession", deleteSession)
	r.HandleFunc("/deleteUser", deleteUser)
	r.HandleFunc("/showDeletedUser", showDeletedUser)
	r.HandleFunc("/showDeletedSession", showDeletedSession)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/", r)

	// handler for various URL signature with net/http
	// http.HandleFunc("/listItem", listItem)
	// http.HandleFunc("/toGiveItem", toGiveItem)
	// http.HandleFunc("/userTray", userTray)

	// Switch to Gorilla Mux
	r.HandleFunc("/searchItem", searchItem)
	r.HandleFunc("/showSearchList", showSearchList)
	r.HandleFunc("/toGiveItem", toGiveItem)
	r.HandleFunc("/giveItem", giveItem)
	//r.HandleFunc("/userTray", userTray)
	r.HandleFunc("/displaySelect", displaySelect)
	r.HandleFunc("/displayList", displayList)
	r.HandleFunc("/getListedItems", getListedItems)
	r.HandleFunc("/myTrayToGive", myTrayToGive)
	r.HandleFunc("/myTrayGiven", myTrayGiven)
	r.HandleFunc("/myTrayGotten", myTrayGotten)
	r.HandleFunc("/myTrayWithdrawn", myTrayWithdrawn)
	r.HandleFunc("/removeFromMyTray", removeFromMyTray)
	r.HandleFunc("/withdrawItem", withdrawItem)

	http.ListenAndServe(vrsHost+":"+vrsPort, nil)
	// Using Open SSL certification and key for development and testing only
	//http.ListenAndServeTLS(vrsHost+":"+vrsPort, "../OpenSSL/cert.pem", "../OpenSSL/key.pem", nil)
}

// index is the Index handler
func index(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req) // alfred 24.06.2022: need to look into this. does this pulls the latest LastLogin?

	// update last visit status to client
	myUser.LastLogin = userLastVisit[myUser.Username] // alfred 24.06.2022: need to look into this. do we need this or is it already updated?

	Trace.Println("Index")

	/*
		regexStr := "^" + adminSubName + "$"
		//regex := regexp.MustCompile(`(^admin\d$)`)
		regex := regexp.MustCompile(regexStr)


			if regex.MatchString(myUser.Username) {
				//admin page
				Trace.Println("Index Admin Page")
				tpl.ExecuteTemplate(res, "index_admin.gohtml", myUser)
			} else {
				// non admin page
				Trace.Println("Index Non Admin Page")
				tpl.ExecuteTemplate(res, "index.gohtml", myUser)
			}
	*/

	if myUser.Role == "admin" {
		//admin page
		Trace.Println("Index Admin Page")
		tpl.ExecuteTemplate(res, "index_admin.gohtml", myUser)
	} else {
		// non admin page
		Trace.Println("Index Non Admin Page")
		tpl.ExecuteTemplate(res, "index.gohtml", myUser)
	}

}

// signup handler
func signup(res http.ResponseWriter, req *http.Request) {
	// precautionary - to handle direct URL pattern for signup
	Trace.Println("Sign up")
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// regex to enforce 1 Alpha + 4 AlphaNumeric char min for username
	regex1 := regexp.MustCompile(`(^[A-Za-z])+([A-Za-z0-9]){4,}`)
	// regex to enforce 1 Alpha + 4 AlphaNumeric char min for password
	regex2 := regexp.MustCompile(`([A-Za-z0-9]){6,}`)
	// for future use
	//regex2 := regexp.MustCompile('^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9])(?=.*[!@#$%^&*()_+,.\\\/;':"-]).{6,}$')
	// First name enforce 2 alpha characters
	regex3 := regexp.MustCompile(`([a-zA-Z]){2,}(\x20)*([a-zA-Z])*`)

	var myUser User
	// process if there is a http post form submission
	if req.Method == http.MethodPost {
		// get form values from browser
		username := req.FormValue("username")
		password := req.FormValue("password")
		name := req.FormValue("name")
		address := req.FormValue("address")
		postalcode := req.FormValue("postalcode")
		telnumber := req.FormValue("telnumber")

		// trim leading and trailing spaces
		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)
		name = strings.TrimSpace(name)
		address = strings.TrimSpace(address)
		postalcode = strings.TrimSpace(postalcode)
		telnumber = strings.TrimSpace(telnumber)

		/*
			// Use this to send to backend server
			fmt.Println("UserName :", username)
			fmt.Println("Password :", password)
			fmt.Println("Name :", name)
			fmt.Println("Address :", address)
			fmt.Println("Postal Code :", postalcode)
			fmt.Println("Phone Number :", telnumber)
		*/

		// User name , password or firstname does not meet policy

		if !regex1.MatchString(username) {
			Trace.Println("Signup failed username")
			http.Error(res, "Username (Please begin with Alphabet and be at least 5 characters)", http.StatusForbidden)
			return
		}

		if !regex2.MatchString(password) {
			Trace.Println("Signup failed password")
			http.Error(res, "Password (Please enter at least 6 Aplha-numeric characters)", http.StatusForbidden)
			return
		}

		if !regex3.MatchString(name) {
			Trace.Println("Signup failed firstname")
			http.Error(res, "First name (Please enter at least 2 Alphabets)", http.StatusForbidden)
			return
		}

		if username != "" {
			// check if username exist
			if _, ok := mapUsers[username]; ok {
				userLastVisit[username] = fmt.Sprintf("Failed : %s", time.Now().Format("Jan-02-2006, 3:04:05 pm")) // alfred 24.06.2022: need to look at this...
				Trace.Println("Signup failed duplicate username")
				http.Error(res, "Username already taken", http.StatusForbidden)
				return
			}

			// generate bcrypt password Hash from the password
			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				Error.Println("Signup bcrypt password generation failure")
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			// initialise myUser
			myUser = User{"", username, string(bPassword), name, address, postalcode, telnumber, "user", ""}

			err = addNewUser(myUser)
			if err != nil {
				fmt.Println(err)
				// log error
			}

			users, err := getAllUsers()
			if err != nil {
				fmt.Println(err)
				// log error
			}

			for _, v := range users {
				mapUsers[v.Username] = v
			}

			fmt.Println(mapUsers)

			//mapUsers[username] = myUser
			userLastVisit[username] = "None"

			// create session ID using cookie
			id := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:    "myCookie",
				Value:   id.String(),
				Expires: time.Now().Add(10 * time.Minute),
			}

			Trace.Println("Signup sucessful")
			// prepare cookie in http header for cookie
			http.SetCookie(res, myCookie)
			// keep session ID with username for future login validation
			mapSessions[myCookie.Value] = username

		}
		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	//	fmt.Println("For First Time http Post Request")
	// This is executed for first entry when http post hasn't happen yet
	tpl.ExecuteTemplate(res, "signup.gohtml", myUser)
}

// Login handler
func login(res http.ResponseWriter, req *http.Request) {
	// disallow concurrent login
	Trace.Println("Login")
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// regex to enforce 1 Alpha + 4 AlphaNumeric char min for username
	regex1 := regexp.MustCompile(`(^[A-Za-z])+([A-Za-z0-9]){4,}`)
	// regex to enforce 1 Alpha + 6 AlphaNumeric char min for password
	regex2 := regexp.MustCompile(`([A-Za-z0-9]){6,}`)

	// process form submission if there is a http post
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")

		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)

		// User name , password or firstname does not meet policy
		if !regex1.MatchString(username) {
			Trace.Println("Login username did not meet policy")
			http.Error(res, "Username (Please begin with Alphabet and be at least 5 characters long)", http.StatusForbidden)
			return
		}

		if !regex2.MatchString(password) {
			Trace.Println("Login password did not meet policy")
			http.Error(res, "Password (Please enter at least 6 alpha-numeric characters)", http.StatusForbidden)
			return
		}

		// check if user exist with username
		myUser, ok := mapUsers[username]
		if !ok {
			Trace.Println("Login failure, invalid user")
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}
		// Username verified, checking Matching of password entered
		err := bcrypt.CompareHashAndPassword([]byte(myUser.Password), []byte(password))
		if err != nil {
			userLastVisit[username] = fmt.Sprintf("Failed : %s", time.Now().Format("Jan-02-2006, 3:04:05 pm"))

			Trace.Println("Login failure, invalid password")
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:    "myCookie",
			Value:   id.String(),
			Expires: time.Now().Add(10 * time.Minute),
		}
		http.SetCookie(res, myCookie)
		mapSessions[myCookie.Value] = username

		Trace.Println("Login sucessful")

		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	// This is executed for first entry when http post hasn't happen yet
	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

// logout handler
func logout(res http.ResponseWriter, req *http.Request) {
	// precautionary - to handle direct URL pattern for logout
	Trace.Println("Logout")

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	updateLastVist(myCookie.Value, fmt.Sprintf("Passed : %s", time.Now().Format("Jan-02-2006, 3:04:05 pm")))

	// delete the session
	//delete(mapSessions, myCookie.Value)
	cleanupSession(myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	Trace.Println("Logout sucessful")

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// delete user handler
func deleteUser(res http.ResponseWriter, req *http.Request) {
	Trace.Println("deleteUser")
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	updateLastVist(myCookie.Value, fmt.Sprintf("Passed : %s", time.Now().Format("Jan-02-2006, 3:04:05 pm"))) // alfred 24.06.2022: question to CM, why is this here?
	if req.Method == http.MethodPost {
		// get form values from browser
		userName := req.FormValue("userName")
		fmt.Println("User name :", userName)
		// to delete the user in next state
		mapDeletedUser[myCookie.Value] = userName
		//delete(mapUsers, userName)
		// test this concept
		//		tpl.ExecuteTemplate(res, "showUser.gohtml", userName)
		// redirect to browse Venue List
		http.Redirect(res, req, "/showDeletedUser", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "deleteUser.gohtml", mapUsers)
}

// showDeletedUser handler
func showDeletedUser(res http.ResponseWriter, req *http.Request) {
	Trace.Println("showDeleteUser")
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	Trace.Println("delete User successful")

	// retrieve the entered data
	Result := mapDeletedUser[myCookie.Value]
	delete(mapUsers, Result)

	fmt.Println("Result :", Result)
	tpl.ExecuteTemplate(res, "showDeletedUser.gohtml", Result)
}

// deleteSession handler
func deleteSession(res http.ResponseWriter, req *http.Request) {
	Trace.Println("deleteSession")
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	updateLastVist(myCookie.Value, fmt.Sprintf("Passed : %s", time.Now().Format("Jan-02-2006, 3:04:05 pm"))) // alfred 24.06.2022: question to CM, why is this here?

	if req.Method == http.MethodPost {
		// get form values from browser
		session := req.FormValue("session")
		fmt.Println("Session :", session)
		// set to delete the session on next state
		mapDeletedSession[myCookie.Value] = session

		http.Redirect(res, req, "/showDeletedSession", http.StatusSeeOther)
		return
	}
	Trace.Println("Show session to be deleted")
	tpl.ExecuteTemplate(res, "deleteSession.gohtml", mapSessions)
}

// showDeletedSession Handler
func showDeletedSession(res http.ResponseWriter, req *http.Request) {
	Trace.Println("showDeletedSession")
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	// retrieve the entered data
	Result := mapDeletedSession[myCookie.Value]
	// delete all sessions variables
	cleanupSession(Result)
	fmt.Println("Result :", Result)
	Trace.Println("Delete Session sucessful")
	tpl.ExecuteTemplate(res, "showDeletedSession.gohtml", Result)
}

// getUser validate user's cookie with server's user session presence
// generate a new cookie for user with no server session presence
func getUser(res http.ResponseWriter, req *http.Request) User {
	Trace.Println("getUser")
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil { // create new one if there is none
		id := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:    "myCookie",
			Value:   id.String(),
			Expires: time.Now().Add(10 * time.Minute),
		}
		Trace.Println("New Cookie for User")
	}
	http.SetCookie(res, myCookie)

	// if the user exists already, get user
	var myUser User
	// get the user identity from record
	if username, ok := mapSessions[myCookie.Value]; ok { //retrieve of the session
		myUser = mapUsers[username] // verify the user
		Trace.Println("getUser identity")
	}
	Trace.Println("getUser completed")

	return myUser
}

// updateLastVist updates the status of userLastVisit[username] map with status and time stamp
func updateLastVist(uuid string, status string) {
	Trace.Println("updateLastVisit")

	if username, ok := mapSessions[uuid]; ok { //retrieve of the session
		userLastVisit[username] = status

		/*
			var myUser User

			myUser = mapUsers[username]
			myUser.LastLogin = status
			mapUsers[username] = myUser // alfred 24.06.2022: update user's lastlogin record on runtime memory.

			err := editUser(myUser) // alfred 24.06.2022: update user's lastlogin record to backend server.
			if err != nil {
				fmt.Println(err)
				// log error
			}
		*/

	}

	// else {
	// 	userLastVisit[username] = "None"
	// }

}

// alreadyLoggedIn check if the user is already logged in
// returns true if already logged in or with a session presence
func alreadyLoggedIn(req *http.Request) bool {
	Trace.Println("alreadyLoggedIn")
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		Trace.Println("mycookie failed validation")
		return false
	}
	Trace.Println("mycookie validated")
	username := mapSessions[myCookie.Value]
	_, ok := mapUsers[username] // check that session user against stored User
	return ok
}

// cleanupSession remove session and all session variables
func cleanupSession(uuid string) {
	Trace.Println("cleanupSession")
	delete(mapSessions, uuid)
	delete(mapDeletedUser, uuid)
	delete(mapDeletedSession, uuid)

	// delete(mapSessionSearch, uuid)
	// delete(mapSessionVenueType, uuid)
	// delete(mapSessionVenueTypeNum, uuid)
	// delete(mapSessionDate, uuid)
	// delete(mapSessionTimeSlot, uuid)
	// delete(mapSessionSearchFoundBookings, uuid)
	// delete(mapSessionRecordToEdit, uuid)
	// delete(mapSessionFoundBookingPtr, uuid)
	// delete(mapSessionPreviousMenu, uuid)
}
