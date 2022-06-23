package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var items []Item

// initialisation for business logic
func bizInit() {
	items, err := getAllItems()
	if err != nil {
		fmt.Println(err)
		// log error
	}

	fmt.Println(items)
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
	var msg []string
	msg = append(msg, "Item Given :"+name+", "+description+" is moved to To-Give Tray") // One one item

	//var test []string
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

	return items, fmt.Errorf("Error: resp.StatusCode is not 200 - %v", err)
}
