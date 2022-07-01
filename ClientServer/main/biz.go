package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Use mutex1 to lock local DB and remote DB to ensure atomic operation
var mutex1 sync.Mutex

// Items stores the items record at package level on local database on runtime memory.
var Items []Item

var cfg mysql.Config // configuration for DSN.

// bizInit initialisation for business logic.
func bizInit() {

	bizItemListInit(mutex1)
}

// bizItemListInit Iniitialises items record on local database on runtime memory.
func bizItemListInit(mu sync.Mutex) {

	mu.Lock()
	defer mu.Unlock()

	// Initialise items record onto a slice.
	items, err := getAllItems()
	if err != nil {
		log.Fatalln(err)
	}

	// Deep copy item slice onto var Item.
	Items = make([]Item, len(items))
	copy(Items, items)

}

// bizListSearchItems searchs for and returns a list of items that has item name OR/AND Description. If item name OR/AND description were left blank, it returns a complete list of available items.
func bizListSearchItems(name string, description string, searchLogic string) ([]Item, error) {

	var foundList []Item

	// Convert to lower case for search.
	name = strings.ToLower(name)
	description = strings.ToLower(description)

	// list all available items if search entry is empty.
	if len(name)+len(description) == 0 {
		for _, v := range Items {
			if v.State == stateToGive {
				foundList = append(foundList, v)
			}
		}
	}

	if searchLogic == "OR" {
		for _, v := range Items {
			if v.State == stateToGive {
				if (v.State == stateToGive && len(name) > 0 && strings.Contains(strings.ToLower(v.Name), name)) ||
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

	return foundList, nil
}

// bizGetListedItems
// Get a list of selected items (From searchedList)
// selected item will have item.getter = userID, and state changed to stateGiven
func bizGetListedItems(uuid string, selectedItem []string) ([]string, error) {

	// item list is in this map mapSessionSearchedList[uuid]

	// pick up the selected items, only display item ID, name and description
	userID := mapSessions[uuid]
	fmt.Println("UserID :", userID)
	var item string
	var err error
	fmt.Println("Searched List", mapSessionSearchedList[uuid])

	var count int = 0

	var msg []string
	for _, v := range selectedItem {
		intVar, _ := strconv.Atoi(v) // use this to get the integer value of the index

		err = bizSetItemStateToGiven(mutex1, userID, mapSessionSearchedList[uuid][intVar].ID)

		fmt.Println("Error", err)

		// Reform the item message if there is error
		if err != nil {
			if err == ErrorItemAlreadyGiven {
				item = fmt.Sprintf("Item:%d, ID: %s, Name: %s, Sorry: %s", intVar+1, mapSessionSearchedList[uuid][intVar].ID,
					mapSessionSearchedList[uuid][intVar].Name, err)
			} else if err == ErrorItemAlreadyWithdrawn {
				item = fmt.Sprintf("Item:%d, ID: %s, Name: %s, Sorry: %s", intVar+1, mapSessionSearchedList[uuid][intVar].ID,
					mapSessionSearchedList[uuid][intVar].Name, err)
			} else {
				return msg, err
			}
		} else {
			item = fmt.Sprintf("Item:%d, ID: %s, Name: %s, Description: %s", intVar+1, mapSessionSearchedList[uuid][intVar].ID,
				mapSessionSearchedList[uuid][intVar].Name, mapSessionSearchedList[uuid][intVar].Description)
			count++
		}
		msg = append(msg, item)
	}

	// move msg to 1 position to the right
	msg = append(msg[:1], msg[0:]...)
	//	var msg []string
	msg[0] = fmt.Sprintf("Number of items Gotten = %d", count)
	//	msg = append(msg, num)

	return msg, nil
}

var ErrorItemAlreadyGiven error = errors.New("Item is already Taken")
var ErrorItemAlreadyWithdrawn error = errors.New("Item is already withdrawn")
var ErrorItemStateInvalid error = errors.New("Item is in Invalid State")

// bizSetItemStateToGiven
// Update the state of the item in slice and in SQL DB.
func bizSetItemStateToGiven(mu sync.Mutex, userID string, id string) error {
	//	fmt.Println("ID:", id)
	mu.Lock()
	defer mu.Unlock()

	var count int = 0
	for i, v := range Items {
		if v.ID == id { // search for ID to match item
			if Items[i].State == stateToGive {
				Items[i].State = stateGiven // use index to change the state directly local DB
				Items[i].GetterUsername = userID
				err := editItem(Items[i]) // update remote DB with the change
				count++
				if err != nil {
					fmt.Println("Error in bizSetItemStateToGiven", err)
					return err
				}
			} else if Items[i].State == stateGiven {
				return ErrorItemAlreadyGiven
			} else if Items[i].State == stateWithdrawn {
				return ErrorItemAlreadyWithdrawn
			} else {
				return ErrorItemStateInvalid
			}
			break // match found, so can break
		}
	}
	return nil
}

// bizWithdrawItems
// withdraw a list of selected items
// items is not displayed list
// selected is the selected items
func bizWithdrawItems(mu sync.Mutex, items []Item, selectedItem []string) ([]string, error) {

	var withdrawList []Item
	withdrawList = make([]Item, len(selectedItem))

	fmt.Println("Items : ", items)
	// get the selected list to be withdrawn
	if len(selectedItem) > 0 {
		fmt.Println("selectedItem", selectedItem)
		for i, v := range selectedItem {
			intVar, _ := strconv.Atoi(v)
			withdrawList[i] = items[intVar]
		}
	}

	fmt.Println("withdraw List =", withdrawList)
	// Set the state to withdrawn for the selected items
	mu.Lock()
	defer mu.Unlock()

	var msg []string
	var item string
	var count int = 0
	for _, v := range withdrawList {

		// State change to local DB
		err := setStateWithdraw(v.ID)
		if err == nil {
			// state change to SQL DB to stateWithdrawn
			v.State = stateWithdrawn

			err := editItem(v)
			if err != nil {
				Trace.Println(err)
				return msg, err
			}
			count++
			item = fmt.Sprintf("item ID: %v, is withdrawn", v.ID)
		} else {
			item = fmt.Sprintf("item ID:%v Withdrawal Error: %s", v.ID, err.Error())
		}
		msg = append(msg, item)
	}

	fmt.Println("withdraw:", withdrawList)

	// move msg to 1 position to the right
	msg = append(msg[:1], msg[0:]...)
	msg[0] = fmt.Sprintf("Number of items Withdrawn = %d", count)

	return msg, nil

}

// bizGiveItem allows user to list an item by providing a item name and description.
// Item listed will be added to backend server mysql database
func bizGiveItem(mu sync.Mutex, name string, description string, username string) ([]string, error) {

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")

	item := Item{"", name, description, 0, 0, 0, username, "", 0, date}

	mu.Lock()
	defer mu.Unlock()

	var msg []string
	err := addNewItem(item) // add item to items table in mysql
	if err != nil {
		fmt.Println(err)
		return msg, err
	}

	Items, err = getAllItems() // in order to get item ID. pull out all items from items table in mysql again to update/overwrite items (all items slice).
	if err != nil {
		Trace.Println(err)
	}

	msg = append(msg, "Item : "+name+", "+description+"  ===> To-Give Tray")

	return msg, nil

}

// bizRemoveFromTray
// Make these Item from Tray not visible in the Tray
func bizRemoveFromTray(mu sync.Mutex, items []Item, selectedList []string, tray string) ([]string, error) {
	fmt.Println("Tray", tray)

	fmt.Println("Here!!")
	fmt.Println("Item to be withdrawn", items) // alfred 25.06.2022: item to be withdrawn?
	fmt.Println("Item index", selectedList)

	var msg []string
	var num string
	var hideList []Item // final selected list to hide
	if len(selectedList) == 0 {
		num = "Nothing to Remove"
	} else {
		// get the items list for hiding
		for _, v := range selectedList {
			intVar, _ := strconv.Atoi(v)
			hideList = append(hideList, items[intVar])
		}

		mu.Lock()
		defer mu.Unlock()

		// Set up hide flag in local db
		for _, v := range hideList {
			hideItem(tray, v.ID)
			switch tray {
			case "myTrayGiven":
				v.HideGiven = 1
			case "myTrayGotten":
				v.HideGotten = 1
			case "myTrayWithdrawn":
				v.HideWithdrawn = 1
			}

			// update SQL DB
			err := editItem(v)
			if err != nil {
				fmt.Println("Error :", err)
				return msg, err
			}
		}
		num = fmt.Sprintf("Number of items removed from Tray = %d", len(selectedList))
	}
	msg = append(msg, num)

	return msg, nil
}

// bizGetItemWithGiverDetails
// Get the Giver's contact details for each item in the selected list and form a message slice
func bizGetItemWithGiverDetails(items []Item, selectedList []string) ([]string, error) {
	var msg []string

	fmt.Println("Items", items)
	fmt.Println("Select", selectedList)

	if len(selectedList) == 0 {
		msg = append(msg, "Nothing Selected")
	} else {

		for _, v := range selectedList {
			intVar, _ := strconv.Atoi(v)
			item := items[intVar]

			formGiverGetterDetails(&msg, "Giver", item.GiverUsername, item)

			msg = append(msg, "\n")
		}
	}
	return msg, nil
}

// bizGetItemWithGetterDetails
// Get the Getter's contact details for each item in the selected list and form a message slice
func bizGetItemWithGetterDetails(items []Item, selectedList []string) ([]string, error) {
	var msg []string

	fmt.Println("Items", items)
	fmt.Println("Select", selectedList)

	if len(selectedList) == 0 {
		msg = append(msg, "Nothing Selected")
	} else {

		for _, v := range selectedList {
			intVar, _ := strconv.Atoi(v)
			item := items[intVar]
			formGiverGetterDetails(&msg, "Getter", item.GetterUsername, item)

			msg = append(msg, "\n")
		}
	}

	return msg, nil
}

// bizGetSortedList
// Sorted data in slice of global item based on the sorted key and result in a slice of string
// The slice message is formed based on the sort choice
func bizGetSortedList(sortBy string) ([]string, error) {
	fmt.Println("Sort By:", sortBy)

	// make a copy of the items list before sort
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
		})
		msg = convertGiverIDFirst2String(items)
	case "5":
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].GetterUsername < items[j].GetterUsername
		})
		msg = convertGetterIDFirst2String(items)
	}

	return msg, nil
}

// bizMyTrayItems collects the list of items pertaining to the user's MyTray based on subtray (tray)
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

// HiteItem hides the item for the local db that match the id
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

var ErrorFailedWithdrwal error = errors.New("Item failed withdrawal")
var ErrorItemIDNotFound error = errors.New("Item ID cannot be found")

// setStateWithdraw sets the state of the item that match ID to stateWithdrawn
func setStateWithdraw(id string) error {

	var err error
	for i, v := range Items {
		// found item
		if v.ID == id {
			if Items[i].State == stateToGive {
				Items[i].State = stateWithdrawn
				err = nil
			} else {
				err = ErrorFailedWithdrwal
			}
			return err
		}
	}
	return ErrorItemIDNotFound
}
