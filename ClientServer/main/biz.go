package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var items []Item

//var sqlDBConnection string
var cfg mysql.Config // configuration for DSN

// initialisation for business logic
func bizInit() {

	items, err := getAllItems()
	if err != nil {
		fmt.Println(err)
		// log error
	}

	fmt.Println(items)

	//	bizSqlInit()  // This is for SQL initialisation if SQL if front/back server is combined
	//bizItemListInit()
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
	// date format using ISO8601
	item1 := Item{ID: "0", Name: "Clothes", Description: "A box of 10 shirts", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-11T23:28:56.782Z"}
	item2 := Item{ID: "1", Name: "Clothes", Description: "A box of 20 shirts", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-10T23:28:56.782Z"}
	item3 := Item{ID: "2", Name: "Saw", Description: "A 10 inch saw", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-09T23:28:56.782Z"}
	item4 := Item{ID: "3", Name: "Computer System", Description: "A Intel Dual Core Computer and monitor", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-08T23:28:56.782Z"}
	item5 := Item{ID: "4", Name: "Calculator", Description: "A scientific calculator", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-07T23:28:56.782Z"}
	item6 := Item{ID: "5", Name: "Monitor", Description: "Dell Model P2412, 24 inch Monitor", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-06T23:28:56.782Z"}
	item7 := Item{ID: "6", Name: "Monitor", Description: "LG Model XYZ, 24 inches", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-05T23:28:56.782Z"}
	item8 := Item{ID: "7", Name: "Clothes", Description: "A box of 20 shirts", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-04T23:28:56.782Z"}
	item9 := Item{ID: "8", Name: "Clothes", Description: "A box of 10 shorts", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-03T23:28:56.782Z"}
	item10 := Item{ID: "9", Name: "Bed Sheets", Description: "3 Queen size bed sheet", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-02-02T23:28:56.782Z"}
	item11 := Item{ID: "10", Name: "Bed Sheets", Description: "3 king and 2 super singles bedsheets", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-02T23:28:56.782Z"}
	item12 := Item{ID: "11", Name: "Shoes", Description: "2 pair of size 12 shoes for men", HideGiven: 0, HideGotten: 0, HideWithdrawn: 0,
		GiverID: "GiverID", GetterID: "GetterID", State: 0, Date: "2022-01-01T23:28:56.782Z"}
	items := []Item{item1, item2, item3, item4, item5, item6, item7, item8, item9, item10, item11, item12}

	fmt.Println("List of Items")
	for i, v := range items {
		fmt.Printf("item %d:, ID: %s, Name: %s, Description: %s, HideGiven: %d, HideGotten: %d, HideWithdrawn: %d, GiverID: %s, GetterID: %s, State: %d, Date: %s\n",
			i, v.ID, v.Name, v.Description, v.HideGiven, v.HideGotten, v.HideWithdrawn, v.GiverID, v.GetterID, v.State, v.Date)
	}
	//fmt.Println("List of Items", items)

}

// bizListSearchItems - searchs for a list of item that has name OR/AND itemDescription
func bizListSearchItems(name string, itemDescription string, searchLogic string) ([]string, error) {
	// test data
	mr1 := itemType{Id: "MR1", Name: "Clothes", Description: "A box of 10 shirts"}
	mr2 := itemType{Id: "MR2", Name: "Clothes", Description: "A box of 20 shirts"}
	mr3 := itemType{Id: "MR3", Name: "Saw", Description: "A 10 inch saw"}
	mr4 := itemType{Id: "MR4", Name: "Computer", Description: "A Intel Computer and monitor"}
	mr5 := itemType{Id: "MR5", Name: "Calculator", Description: "A scientific calculator"}
	mr6 := itemType{Id: "MR6", Name: "Monitor", Description: "Dell Model 123"}
	mr7 := itemType{Id: "MR7", Name: "Monitor", Description: "LG Model XYZ, 24 inche"}
	mr8 := itemType{Id: "MR8", Name: "Clothes", Description: "A box of 10 shorts"}
	mr9 := itemType{Id: "MR9", Name: "Bed Sheets", Description: "3 Queen size bed sheet"}
	mr10 := itemType{Id: "MR10", Name: "Shoe", Description: "A pair of size 10 shoes for men"}
	mr11 := itemType{Id: "MR11", Name: "Bed Sheets", Description: "3 king size bed sheet"}
	mr12 := itemType{Id: "MR12", Name: "Shoe", Description: "A pair of size 12 shoes for men"}
	list := []itemType{mr1, mr2, mr3, mr4, mr5, mr6, mr7, mr8, mr9, mr10, mr11, mr12}

	strList := convertToString(list)
	fmt.Println("String List", strList)

	return strList, nil
}

// Get a list of selected items
func bizGetListedItems(selectedItem []string) ([]string, error) {

	var msg []string
	num := fmt.Sprintf("Number of items Gotten = %d", len(selectedItem))
	msg = append(msg, num)
	// test data
	for _, v := range selectedItem {
		item := fmt.Sprintf("item %v, moved to Gotten Tray", v)
		msg = append(msg, item)
	}

	//var test []string
	return msg, nil

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

// Give an item for listing
func bizGiveItem(name string, description string) ([]string, error) {

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")

	item := Item{"", name, description, 0, 0, 0, "1", "0", 0, date} //GiverID hardcoded for testing purpose..

	err := addNewItem(item) // add item to items table in mysql
	if err != nil {
		fmt.Println(err)
		// log error
	}

	items, err = getAllItems() // in order to get item ID. pull out all items from items table in mysql again to update/overwrite items (all items slice).
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

	mr12 := itemType{Id: "MR1", Name: "Clothes", Description: "A box of 10 shirts"}
	mr11 := itemType{Id: "MR2", Name: "Clothes", Description: "A box of 20 shirts"}
	mr10 := itemType{Id: "MR3", Name: "Saw", Description: "A 10 inch saw"}
	mr9 := itemType{Id: "MR4", Name: "Computer", Description: "A Intel Computer and monitor"}
	mr8 := itemType{Id: "MR5", Name: "Calculator", Description: "A scientific calculator"}
	mr7 := itemType{Id: "MR6", Name: "Monitor", Description: "Dell Model 123"}
	mr6 := itemType{Id: "MR7", Name: "Monitor", Description: "LG Model XYZ, 24 inche"}
	mr5 := itemType{Id: "MR8", Name: "Clothes", Description: "A box of 10 shorts"}
	mr4 := itemType{Id: "MR9", Name: "Bed Sheets", Description: "3 Queen size bed sheet"}
	mr3 := itemType{Id: "MR10", Name: "Shoe", Description: "A pair of size 10 shoes for men"}
	mr2 := itemType{Id: "MR11", Name: "Bed Sheets", Description: "3 king size bed sheet"}
	mr1 := itemType{Id: "MR12", Name: "Shoe", Description: "A pair of size 12 shoes for men"}
	list := []itemType{mr1, mr2, mr3, mr4, mr5, mr6, mr7, mr8, mr9, mr10, mr11, mr12}
	//fmt.Println(list)
	// selected soted selection
	msg := convertToString(list)

	//	var test []string
	return msg, nil
}

func bizMyTrayItems(tray string) ([]string, error) {
	fmt.Println("Tray", tray)
	// test data
	mr1 := itemType{Id: "MR1", Name: "Clothes", Description: "A box of 10 shirts"}
	mr2 := itemType{Id: "MR2", Name: "Clothes", Description: "A box of 20 shirts"}
	mr3 := itemType{Id: "MR3", Name: "Saw", Description: "A 10 inch saw"}
	mr4 := itemType{Id: "MR4", Name: "Computer", Description: "A Intel Computer and monitor"}
	mr5 := itemType{Id: "MR5", Name: "Calculator", Description: "A scientific calculator"}
	mr6 := itemType{Id: "MR6", Name: "Monitor", Description: "Dell Model 123"}
	mr7 := itemType{Id: "MR7", Name: "Monitor", Description: "LG Model XYZ, 24 inche"}
	mr8 := itemType{Id: "MR8", Name: "Clothes", Description: "A box of 10 shorts"}
	mr9 := itemType{Id: "MR9", Name: "Bed Sheets", Description: "3 Queen size bed sheet"}
	mr10 := itemType{Id: "MR10", Name: "Shoe", Description: "A pair of size 10 shoes for men"}
	mr11 := itemType{Id: "MR11", Name: "Bed Sheets", Description: "3 king size bed sheet"}
	mr12 := itemType{Id: "MR12", Name: "Shoe", Description: "A pair of size 12 shoes for men"}
	list := []itemType{mr1, mr2, mr3, mr4, mr5, mr6, mr7, mr8, mr9, mr10, mr11, mr12}
	//fmt.Println(list)
	strList := convertToString(list)

	//var test []string
	return strList, nil
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

func editItem(item Item) error {

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
