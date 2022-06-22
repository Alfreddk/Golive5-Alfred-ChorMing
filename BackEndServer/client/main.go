// Note this needs to be run with lab3 since lab3 is the server and this is the client
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/joho/godotenv"
)

var baseURL string
var key string

// init() initialises the system
// Set up the environment for client system
// For this exercise the same .env file is used for both server and client,
// in actual deployment, this should be separate .env file
func init() {

	// set path for the env file
	envFile := path.Join("..", "config", ".env")

	//err := godotenv.Load(".env")
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalln("Error loading .env file: ", err)
	}

	// getting env variables SITE_TITLE and DB_HOST
	siteTitle := os.Getenv("CLIENT_TITLE")
	clientHost := os.Getenv("CLIENT_HOST")
	clientPort := os.Getenv("CLIENT_PORT")
	key = os.Getenv("CLIENT_URLKEY")

	// Create base URL from environment variable
	baseURL = fmt.Sprintf("http://%s:%s/api/v1/courses", clientHost, clientPort)

	fmt.Println(baseURL)

	fmt.Printf("%s = %s\n", "Site Title", siteTitle)
	fmt.Printf("Use https:// %s:%s\n", clientHost, clientPort)
}

// main function
func main() {
	// run menu use channel with 1 value to block till done, to prevent main from exiting until menu is done
	done := make(chan bool, 1)
	// run menu concurrently
	go menu(done)
	// block until menu is done
	<-done
}

//getCourse gets the course via REST API
func getCourse(code string) {
	url := baseURL
	if code != "" {
		url = baseURL + "/" + code + "?key=" + key
	}
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

//addCourse adds the course via REST API
func addCourse(code string, jsonData map[string]string) {
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post(baseURL+"/"+code+"?key="+key, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

//updateCourse updates the course via REST API
func updateCourse(code string, jsonData map[string]string) {
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest(http.MethodPut, baseURL+"/"+code+"?key="+key, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

//deleteCourse deletes a course via REST API
func deleteCourse(code string) {
	request, err := http.NewRequest(http.MethodDelete, baseURL+"/"+code+"?key="+key, nil)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

//menu is a console menu to manange the input and firing of the REST APIs
func menu(done chan bool) {
	operationMenu := []string{
		"Course Management System",
		"==========================",
		"1.  List Courses",
		"2.  Get a course",
		"3.  Add Course",
		"4.  Update Course",
		"5.  Delete Course",
		"x.  To Exit",
		"Select your choice:",
	}

	for {
		setColor(Red)
		for _, item := range operationMenu {
			fmt.Println(item)
		}
		setColor(Reset)

		index, err := get1KeyMatch("123456xX")

		if err == nil && index >= 0 {
			switch index {
			case 0: //1
				getCourse("")
			case 1: //2
				menuGetCourse()
			case 2: //2
				menuAddCourse()
			case 3: //3
				menuUpdateCourse()
			case 4:
				menuDeleteCourse()
			case 5:
			case 6:
				// break only work for switch but not applied to for loop
			default: // not needed
				// break only work for switch but not applied to for loop
			}
			//			fmt.Printf("Index = %d", index)
			if index >= 5 && index <= 6 {
				break // use break here to break the loop
			}
		}
	}
	done <- true
}

// menuGetCourse list a specified course
func menuGetCourse() {
	fmt.Println("Enter the course ID")
	courseId, err := readString()
	if err != nil {
		fmt.Println(err)
		return
	}
	getCourse(courseId)
}

// menuAddCourse adds a course to the list of courses
func menuAddCourse() {
	fmt.Println("Enter the course ID")
	courseId, err1 := readString()
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	fmt.Println("Enter the title of the course")
	title, err2 := readString()
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	jsonData := map[string]string{}
	jsonData["title"] = title
	addCourse(courseId, jsonData)
}

// menuUpdateCourse update the course title with the respective ID
func menuUpdateCourse() {
	fmt.Println("Enter the course ID")
	courseId, err1 := readString()
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	fmt.Println("Enter the title of the course")
	title, err2 := readString()
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	jsonData := map[string]string{}
	jsonData["title"] = title
	updateCourse(courseId, jsonData)
}

// menuDeleteCourse deletes the course with the respective ID
func menuDeleteCourse() {
	fmt.Println("Enter the course ID")
	courseId, err1 := readString()
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	deleteCourse(courseId)
}
