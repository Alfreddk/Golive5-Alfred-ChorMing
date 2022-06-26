package main

import (
	"log"
)

var mapUsers = map[string]User{}

// userInit performs getAllUsers() function call to retreive all users records and maps these records to mapUsers on runtime memory.
func userInit() {
	users, err := getAllUsers()
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range users {
		mapUsers[v.Username] = v
		userLastVisit[v.Username] = v.LastLogin
	}

}
