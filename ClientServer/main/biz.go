package main

import (
	"fmt"
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

}

// bizListSearchItems - searchs for a list of item that has name OR/AND itemDescription
func bizListSearchItems(name string, description string, searchLogic string) ([]Item, error) {

	var foundList []Item

	// convert to lower case for search
	name = strings.ToLower(name)
	description = strings.ToLower(description)

	// list 20 items if search entry is empty
	if len(name)+len(description) == 0 {
		for i, v := range Items {
			if v.State == stateToGive {
				foundList = append(foundList, v)
			}
			if i == 20 {
				break
			}
		}
	}

	if searchLogic == "OR" {
		for _, v := range Items {
			if v.State == stateToGive {
				if (len(name) > 0 && strings.Contains(strings.ToLower(v.Name), name)) ||
					len(description) > 0 && strings.Contains(strings.ToLower(v.Description), description) {
					foundList = append(foundList, v)
				}
			}
		}
	}
	if searchLogic == "AND" {
		for _, v := range Items {
			if v.State == stateToGive {
				if (len(name) > 0 && strings.Contains(strings.ToLower(v.Name), name)) &&
					len(description) > 0 && strings.Contains(strings.ToLower(v.Description), description) {
					foundList = append(foundList, v)
				}
			}
		}
	}
	//strList := convertItems2String(foundList)
	//fmt.Println("String List", strList)

	return foundList, nil
}

// Get a list of selected items
// selected item will have item.getter = userID, and state changed to stateGiven
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

		bizSetItemStateToGiven(userID, mapSessionSearchedList[uuid][intVar].ID)
		msg = append(msg, item)
	}

	//var test []string
	return msg, nil
}

// Update the state of the item in slice and in SQL DB
// SQL DB need API, pending implementation
func bizSetItemStateToGiven(userID string, id string) {
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
// items is not displayed list
// selected is the selected items
func bizWithdrawItems(items []Item, selectedItem []string) ([]string, error) {
	var msg []string
	var withdrawList []Item
	withdrawList = make([]Item, len(selectedItem))

	fmt.Println("Items : ", items)
	// get the selected list to be withdrawn
	if len(selectedItem) > 0 {
		fmt.Println("selectedItem", selectedItem)
		for _, v := range selectedItem {
			intVar, _ := strconv.Atoi(v)
			withdrawList = append(withdrawList, items[intVar])
		}
	}

	// Set the state to withdrawn for the selected items
	for _, v := range withdrawList {

		// State change to local DB
		setStateWithdraw(v.ID) // set item to withdrawn

		// state change to SQL DB to stateWithdrawn
		v.State = stateWithdrawn
		err := addNewItem(v) // add item to items table in mysql
		if err != nil {
			fmt.Println(err)
			// log error
		}
	}

	fmt.Println("withdraw:", withdrawList)

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

// Give an item for listing
func bizGiveItem(name string, description string, username string) ([]string, error) {

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")

	item := Item{"", name, description, 0, 0, 0, username, "", 0, date} //GiverUsername hardcoded for testing purpose..

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

// Make these Item from Tray not visible in the Tray
func bizRemoveFromTray(items []Item, selectedList []string, tray string) ([]string, error) {
	fmt.Println("Tray", tray)

	fmt.Println("Here!!")
	fmt.Println("Item to be withdrawn", items)
	fmt.Println("Item index", selectedList)

	var msg []string
	var num string
	var hideList []Item // final selected list to hide
	if len(selectedList) == 0 {
		num = "Nothing to Withdraw"
	} else {
		// get the items list for hiding
		for _, v := range selectedList {
			intVar, _ := strconv.Atoi(v)
			hideList = append(hideList, items[intVar])
		}
		// Set up hide flag in local db
		for _, v := range hideList {
			hideItem(tray, v.ID)
			switch tray {
			case "myTrayGiven":
				v.HideGiven = 1
			case "myTrayGotten":
				//v.HideGottem = 1
			case "myTrayWithdrawn":
				v.HideWithdrawn = 1
			}
			// update SQL DB
			err := editItem(v)
			if err != nil {
				fmt.Println("Error :", err)
			}
		}
		num = fmt.Sprintf("Number of items removed from Tray = %d", len(selectedList))
	}
	msg = append(msg, num)

	fmt.Println("tray Type:", tray)
	fmt.Println("hide List", hideList)

	fmt.Println("Items ", msg)
	// var test []string
	return msg, nil
}

// Get a list of sorted data based on the sorted key
func bizGetSortedList(sortBy string) ([]string, error) {
	fmt.Println("Sort By:", sortBy)

	// make a copy of the list before sort
	items := make([]Item, len(Items))
	copy(items, Items) // deep copy

	var msg []string
	switch sortBy {
	case "0":
		msg = convertItems2String(items)
	case "1":
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Name < items[j].Name
		})
		msg = convertNameFirst2String(items)
	case "2":
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].State < items[j].State
		})
		msg = convertStateFirst2String(items)
	case "3":
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Date < items[j].Date
		})
		msg = convertDateFirst2String(items)
	case "4":
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].GiverUsername < items[j].GiverUsername
		}) //alfred 23.06.2022: ChorMing you need to relook into this. Changed from ID to username.
		msg = convertGiverIDFirst2String(items)
	case "5":
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].GetterUsername < items[j].GetterUsername
		}) //alfred 23.06.2022: ChorMing you need to relook into this. Changed from ID to username.
		msg = convertGetterIDFirst2String(items)
	}

	//	var test []string
	return msg, nil
}

// collect a list of items for myTray
// myTrayToGive - Item given but not received
// myTrayGiven - Item Received by Given
// myTrayGotten - Not used
// myTrayWithdrawn - Item to give but withdrawn before any taker
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

// hide the item for the local db that match the id
func hideItem(tray string, id string) {

	switch tray {
	case "myTrayGiven":
		for i, v := range Items {
			if v.ID == id {
				Items[i].HideGiven = 1
				break
			}
		}
	case "myTrayGotten":
		for i, v := range Items {
			if v.ID == id {
				Items[i].HideGotten = 1
				break
			}
		}
	case "myTrayWithdrawn":
		for i, v := range Items {
			if v.ID == id {
				Items[i].HideWithdrawn = 1
				break
			}
		}
	}
}

func setStateWithdraw(id string) {

	for i, v := range Items {
		if v.ID == id {
			Items[i].State = stateWithdrawn
			break
		}
	}
}
