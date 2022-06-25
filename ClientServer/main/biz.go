package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var Items []Item

//var sqlDBConnection string
var cfg mysql.Config // configuration for DSN

// initialisation for business logic
func bizInit() {

	//	bizSqlInit()  // This is for SQL initialisation if SQL if front/back server is combined
	bizItemListInit()
}

// Initialise the SQL server connection from main init
func bizSqlInit() {
	// SQL DB Data Source Name config
	cfg = mysql.Config{
		User:   os.Getenv("SQL_USER"),
		Passwd: os.Getenv("SQL_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("SQL_ADDR"),
		DBName: os.Getenv("SQL_DB"),
	}
}

// Iniitialises item for testing purpose
func bizItemListInit() {

	// This is the place to initialise the package slice of items
	items, err := getAllItems()
	if err != nil {
		fmt.Println(err)
		// log error
	}

	Items = make([]Item, len(items))
	copy(Items, items)

	fmt.Println("List of Items")
	for i, v := range items {
		fmt.Printf("item %d:, ID: %s, Name: %s, Description: %s, HideGiven: %d, HideGotten: %d, HideWithdrawn: %d, GiverUsername: %s, GetterUsername: %s, State: %d, Date: %s\n",
			i, v.ID, v.Name, v.Description, v.HideGiven, v.HideGotten, v.HideWithdrawn, v.GiverUsername, v.GetterUsername, v.State, v.Date)
	}
	//fmt.Println("List of Items", items)

}

// bizListSearchItems searchs for and returns a list of items that has item name OR/AND Description. If item name OR/AND description were left blank, it returns a complete list of available items.
func bizListSearchItems(name string, description string, searchLogic string) ([]Item, error) {

	var foundList []Item

	// convert to lower case for search
	name = strings.ToLower(name)
	description = strings.ToLower(description)

	if searchLogic == "OR" {
		for _, v := range Items {
			if v.State == stateToGive {
				if len(name) == 0 {
					foundList = append(foundList, v)
				} else if (v.State == stateToGive && len(name) > 0 && strings.Contains(strings.ToLower(v.Name), name)) ||
					len(description) > 0 && strings.Contains(strings.ToLower(v.Description), description) {
					foundList = append(foundList, v)
				}
				/*
					if (v.State == stateToGive && len(name) > 0 && strings.Contains(strings.ToLower(v.Name), name)) ||
						len(description) > 0 && strings.Contains(strings.ToLower(v.Description), description) {
						foundList = append(foundList, v)
					}
				*/
			}
		}
	}
	if searchLogic == "AND" {
		for _, v := range Items {
			if v.State == stateToGive {
				if len(name) == 0 && len(description) == 0 {
					foundList = append(foundList, v)
				} else if (len(name) > 0 && strings.Contains(strings.ToLower(v.Name), name)) &&
					len(description) > 0 && strings.Contains(strings.ToLower(v.Description), description) {
					foundList = append(foundList, v)
				}
				/*
					if (len(name) > 0 && strings.Contains(strings.ToLower(v.Name), name)) &&
						len(description) > 0 && strings.Contains(strings.ToLower(v.Description), description) {
						foundList = append(foundList, v)
					}
				*/
			}
		}
	}
	strList := convertItems2String(foundList)
	fmt.Println("String List", strList)

	return foundList, nil
}

// Get a list of selected items
func bizGetListedItems(uuid string, selectedItem []string) ([]string, error) {

	// item list is in this map mapSessionSearchedList[uuid]
	var msg []string
	num := fmt.Sprintf("Number of items Gotten = %d", len(selectedItem))
	msg = append(msg, num)
	// test data
	// pick up the selected items, only display ID and name
	// Need also to set the flag for the database
	userID := mapSessions[uuid]
	for _, v := range selectedItem {
		intVar, _ := strconv.Atoi(v) // use this to get the integer value of the index
		item := fmt.Sprintf("Item:%d, ID: %s, Name: %s, Description: %s", intVar+1, mapSessionSearchedList[uuid][intVar].ID,
			mapSessionSearchedList[uuid][intVar].Name, mapSessionSearchedList[uuid][intVar].Description)
		bizUpdateItemState(userID, mapSessionSearchedList[uuid][intVar].ID)
		msg = append(msg, item)
	}

	//var test []string
	return msg, nil
}

// Update the state of the item in slice and in SQL DB
// SQL DB need API, pending implementation
func bizUpdateItemState(userID string, id string) {
	//	fmt.Println("ID:", id)
	for i, v := range Items {
		if v.ID == id { // search for ID to match item
			Items[i].State = stateGiven // use index to change the state directly local DB
			Items[i].GetterUsername = userID
			err := editItem(Items[i]) // update remote DB with the change
			if err != nil {
				fmt.Println("Error", err)
			}
			break // match found, so can break
		}
	}
}

// withdraw a list of selected items
func bizWithdrawItems(selectedItem []string) ([]string, error) {
	var msg []string
	num := fmt.Sprintf("Number of items Withdrawn = %d", len(selectedItem))
	msg = append(msg, num)
	// test data
	for _, v := range selectedItem {
		item := fmt.Sprintf("item %v, is withdrawn", v)
		msg = append(msg, item)
	}

	//var test []string
	return msg, nil

}

// bizGiveItem allows user to list an item by providing a item name and description.
// Item listed will be added to backend server mysql database
func bizGiveItem(name string, description string, username string) ([]string, error) {

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")

	item := Item{"", name, description, 0, 0, 0, username, "", 0, date}

	err := addNewItem(item) // add item to items table in mysql
	if err != nil {
		fmt.Println(err)
		// log error
	}

	Items, err = getAllItems() // in order to get item ID. pull out all items from items table in mysql again to update/overwrite items (all items slice).
	if err != nil {
		fmt.Println(err)
		// log error
	}

	var msg []string
	msg = append(msg, "Item Given :"+name+", "+description+" is moved to To-Give Tray")

	return msg, nil

}

// Make thest Item from Tray not visible in the Tray
func bizRemoveFromTray(selectedList []string, tray string) ([]string, error) {
	fmt.Println("Tray", tray)

	var msg []string
	num := fmt.Sprintf("Number of items Withdrawn = %d", len(selectedList))
	msg = append(msg, num)

	switch tray {
	case "myTrayGiven":
		// test data
		for _, v := range selectedList {
			item := fmt.Sprintf("item %v, removed from Given Tray", v)
			msg = append(msg, item)
		}

	case "myTrayGotten":
		// test data
		for _, v := range selectedList {
			item := fmt.Sprintf("item %v, removed from Gotten Tray", v)
			msg = append(msg, item)
		}

	case "myTrayWithdrawn":
		// test data
		for _, v := range selectedList {
			item := fmt.Sprintf("item %v, removed from Withdrawn Tray", v)
			msg = append(msg, item)
		}
	}

	// var test []string
	return msg, nil
}

// Get a list of sorted data
func bizGetSortedList(sortBy string) ([]string, error) {
	fmt.Println("Sort By:", sortBy)

	// make an empty list before sort
	items := make([]Item, len(Items))
	copy(items, Items) // deep copy

	switch sortBy {
	case "0":
		// unsorted list
	case "1":
		sort.SliceStable(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	case "2":
		sort.SliceStable(items, func(i, j int) bool { return items[i].State < items[j].State })
	case "3":
		sort.SliceStable(items, func(i, j int) bool { return items[i].Date < items[j].Date })
	case "4":
		sort.SliceStable(items, func(i, j int) bool { return items[i].GiverUsername < items[j].GiverUsername }) //alfred 23.06.2022: ChorMing you need to relook into this. Changed from ID to username.
	case "5":
		sort.SliceStable(items, func(i, j int) bool { return items[i].GetterUsername < items[j].GetterUsername }) //alfred 23.06.2022: ChorMing you need to relook into this. Changed from ID to username.
	}

	msg := convertItems2String(items)

	//	var test []string
	return msg, nil
}

func bizMyTrayItems(userID string, tray string) ([]Item, error) {

	var trayList []Item

	switch tray {
	case "myTrayToGive":
		// search for item with State=stateToGive
		for _, v := range Items {
			if v.State == stateToGive && v.GiverUsername == userID {
				trayList = append(trayList, v)
			}
		}

	case "myTrayGiven":
		// search for items with State=stateGiven
		for _, v := range Items {
			if v.State == stateGiven && v.GiverUsername == userID && v.HideGiven == 0 {
				trayList = append(trayList, v)
			}
		}

	case "myTrayGotten":
		// search for items with State=stateGiven
		for _, v := range Items {
			if v.State == stateGiven && v.GetterUsername == userID && v.HideGotten == 0 {
				trayList = append(trayList, v)
			}
		}

	case "myTrayWithdrawn":
		// search for items with State=stateWithdrawn

		for _, v := range Items {
			if v.State == stateWithdrawn && v.GiverUsername == userID && v.HideWithdrawn == 0 {
				trayList = append(trayList, v)
			}
		}
	}

	return trayList, nil
}

func getAllItems() (items []Item, err error) {

	items = []Item{}

	backendURL := "http://127.0.0.1:5000/api/v1/allitems/?key=2c78afaf-97da-4816-bbee-9ad239abb296"

	resp, err := http.Get(backendURL)
	if err != nil {
		return items, fmt.Errorf("Error: POST request - %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		err := json.Unmarshal(respData, &items)
		if err != nil {
			return items, fmt.Errorf("Error: JSON unmarshaling session - %v", err)
		}

		return items, nil
	}

	return items, errors.New("Error: resp.StatusCode is not 200")
}

func addNewItem(item Item) error {

	backendURL := "http://127.0.0.1:5000/api/v1/addnewitem/?key=2c78afaf-97da-4816-bbee-9ad239abb296"

	jsonData, err := json.Marshal(item)
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

func editItem(item Item) error { // alfred 23.06.2022: not tested...

	backendURL := "http://127.0.0.1:5000/api/v1/edititem/?key=2c78afaf-97da-4816-bbee-9ad239abb296"

	jsonData, err := json.Marshal(item)
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
