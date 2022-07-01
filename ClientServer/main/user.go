package main

import (
	"log"
	"sync"
)

var mutex2 sync.Mutex
var mapUsers = map[string]User{}

// userInit performs getAllUsers() function call to retreive all users records and maps these records to mapUsers on runtime memory.
func userInit() {

	mutex2.Lock()
	defer mutex2.Unlock()

	users, err := getAllUsers()
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range users {
		mapUsers[v.Username] = v
		userLastVisit[v.Username] = v.LastLogin
	}

}
