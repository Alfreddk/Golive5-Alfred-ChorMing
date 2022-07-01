// Package main is the code for GiveNGet System
// Created for Project Go Live
// By Tan Chor Ming & Alfred Wung
// 26 June 2022
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

var tpl *template.Template

var mapSessions = map[string]string{}
var mapDeletedUser = map[string]string{}
var mapDeletedSession = map[string]string{}

var userLastVisit = map[string]string{}

var vrsHost, vrsPort string

var errLogDir string

var backendHost string
var backendPort string
var urlKey string

// init() initialises the system
// Set up the environment
// Set up the logger
// Set up local database on runtime memory
func init() {

	// set path for the env file
	envFile := path.Join("..", "config", ".env")

	//err := godotenv.Load(".env")
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalln("Error loading .env file: ", err)
	}

	// getting env variables SITE_TITLE , VRS_HOST and VRS_PORT
	siteTitle := os.Getenv("SITE_TITLE")
	vrsHost = os.Getenv("VRS_HOST")
	vrsPort = os.Getenv("VRS_PORT")

	// getting env variables BACKEND_HOST, BACKEND_PORT and CLIENT_URLKEY
	backendHost = os.Getenv("BACKEND_HOST")
	backendPort = os.Getenv("BACKEND_PORT")
	urlKey = os.Getenv("CLIENT_URLKEY")

	// Error log relative dir
	errLogDir = os.Getenv("VRS_ERR_LOG_DIR")

	fmt.Printf("%s = %s\n", "Site Title", siteTitle)
	fmt.Printf("Listening on https://%s:%s\n", vrsHost, vrsPort)

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

	// template initialisation
	tpl = template.Must(template.ParseGlob("templates/*"))

	// initialise user and business logic for application
	go userInit()
	go bizInit()

}

// main package starts here
func main() {
	// Test each log type.
	// Trace.Println("Test Trace")
	// Warning.Println("Test Warining")
	// Error.Println("Test Error")
	// Info.Println("Venue Booking System Initialisation Completed")

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

	r.HandleFunc("/searchItem", searchItem)
	r.HandleFunc("/showSearchList", showSearchList)
	r.HandleFunc("/toGiveItem", toGiveItem)
	r.HandleFunc("/giveItem", giveItem)
	r.HandleFunc("/displaySelect", displaySelect)
	r.HandleFunc("/displayList", displayList)
	r.HandleFunc("/getListedItems", getListedItems)
	r.HandleFunc("/myTrayToGive", myTrayToGive)
	r.HandleFunc("/myTrayGiven", myTrayGiven)
	r.HandleFunc("/myTrayGotten", myTrayGotten)
	r.HandleFunc("/myTrayWithdrawn", myTrayWithdrawn)
	r.HandleFunc("/removeFromMyTray", removeFromMyTray)
	r.HandleFunc("/withdrawItem", withdrawItem)
	r.HandleFunc("/viewGiverDetails", viewGiverDetails)
	r.HandleFunc("/viewGetterDetails", viewGetterDetails)

	//http.ListenAndServe(vrsHost+":"+vrsPort, nil)
	// Using Open SSL certification and key for development and testing only
	http.ListenAndServeTLS(vrsHost+":"+vrsPort, "../OpenSSL/cert.pem", "../OpenSSL/key.pem", nil)
}

// index is the Index handler
func index(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)

	// update last visit status to client
	myUser.LastLogin = userLastVisit[myUser.Username]

	Trace.Println("Index")

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

			mutex2.Lock()
			defer mutex2.Unlock()

			err = addNewUser(myUser)
			if err != nil {
				Trace.Println(err)
			}

			users, err := getAllUsers()
			if err != nil {
				Trace.Println(err)
			}

			for _, v := range users {
				mapUsers[v.Username] = v
			}

			fmt.Println(mapUsers)

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

	// This is executed for first entry when http post hasn't happen yet
	tpl.ExecuteTemplate(res, "signup.gohtml", myUser)
}

// login handler
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
	updateLastVist(myCookie.Value, fmt.Sprintf("Passed : %s", time.Now().Format("Jan-02-2006, 3:04:05 pm")))
	if req.Method == http.MethodPost {
		// get form values from browser
		userName := req.FormValue("userName")
		fmt.Println("User name :", userName)

		for _, v := range Items {
			if v.GiverUsername == userName && v.State == stateToGive {
				v.State = stateInvalid // change all "togive" state items listed by this user to "invalid" state on runtime memory Items slice.
				err := editItem(v)     // change all "togive" state items listed by this user to "invalid" state on backend mysql database.
				if err != nil {
					Trace.Println(err)
				}
			}
		}

		// delete user from backend server mysql database.
		var myUser User
		myUser = mapUsers[userName]
		mutex2.Lock()
		defer mutex2.Unlock()

		err := delUser(myUser)
		if err != nil {
			Trace.Println(err)
		}

		// delete user from runtime mapUsers.
		delete(mapUsers, userName)

		// redirect to browse Venue List
		http.Redirect(res, req, "/showDeletedUser", http.StatusSeeOther) // alfred 24.06.2022: No longer tracking deleted users. this page is no longer relevant. CM to advise.
		return
	}

	tpl.ExecuteTemplate(res, "deleteUser.gohtml", mapUsers)
}

// showDeletedUser handler
func showDeletedUser(res http.ResponseWriter, req *http.Request) { // alfred 24.06.2022: No longer tracking deleted users. this func is no longer relevant. CM to advise.
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
	updateLastVist(myCookie.Value, fmt.Sprintf("Passed : %s", time.Now().Format("Jan-02-2006, 3:04:05 pm")))

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

		var myUser User
		mutex2.Lock()
		defer mutex2.Unlock()

		myUser = mapUsers[username]
		myUser.LastLogin = status
		mapUsers[username] = myUser // Update user's lastlogin record on runtime memory.

		err := editUser(myUser) // Update user's lastlogin record to backend server.
		if err != nil {
			Trace.Println(err)
		}

	}

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

	delete(mapSessionName, uuid)
	delete(mapSessionItemName, uuid)
	delete(mapSessionItemDescription, uuid)
	delete(mapSessionSearchLogic, uuid)
	delete(mapSessionSearch, uuid)
	delete(mapSessionSelect, uuid)
	delete(mapSessionSort, uuid)
	delete(mapSessionSearchedList, uuid)
	delete(mapSessionMyTrayList, uuid)
	delete(mapSessionPreviousMenu, uuid)
}
