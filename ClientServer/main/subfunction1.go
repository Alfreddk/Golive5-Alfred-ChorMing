package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// general reply message to server
type pageType struct {
	PageType int // Type of message
	//	Err      int      // use for internal processing status
	Msg []string // List of messages
}

type menuNameDescType struct {
	PageType       int
	DefName        string
	DefDescription string
}

// Menu for List
type menuListType struct {
	PageType  int
	DefSelect int
	List      []string
}

type messageType struct {
	PageType int
	Msg      []string
}

// for testing only
type itemType struct {
	Id          string
	Name        string
	Description string
}

var mutex1 sync.Mutex
var mutex2 sync.Mutex

// session variables
var mapSessionName = map[string]string{}
var mapSessionItemName = map[string]string{}
var mapSessionItemDescription = map[string]string{}
var mapSessionSearchLogic = map[string]string{}
var mapSessionSearch = map[string]string{}
var mapSessionSelect = map[string][]string{}
var mapSessionSort = map[string]string{}

var mapSessionPreviousMenu = map[string]string{}

func postRedirectNameDescription(uuid string, res http.ResponseWriter,
	req *http.Request) {

	if req.Method == http.MethodPost {

		name := req.FormValue("name")
		name = strings.TrimSpace(name)
		mapSessionItemName[uuid] = name

		description := req.FormValue("description")
		description = strings.TrimSpace(description)
		mapSessionItemDescription[uuid] = strings.TrimSpace(description)

		searchLogic := req.FormValue("searchLogic")
		searchLogic = strings.TrimSpace(searchLogic)
		mapSessionSearchLogic[uuid] = strings.TrimSpace(searchLogic)

		fmt.Println("Name :", name)
		fmt.Println("Description :", description)
		fmt.Println("Search Logic :", searchLogic)

		// selective redirect to based on last menu
		switch mapSessionPreviousMenu[uuid] {

		case "searchItem":
			fmt.Println("searchItem")
			http.Redirect(res, req, "/showSearchList", http.StatusSeeOther)
			return
		case "toGiveItem":
			fmt.Println("searchItem")
			http.Redirect(res, req, "/giveItem", http.StatusSeeOther)
			return
		}
	}
}

// postRedirectSortSelect post http to get parameter (select option) for
func postRedirectSortSelect(uuid string, res http.ResponseWriter,
	req *http.Request) {
	Trace.Println("postRedirectSortSelect")
	// Get selection from requester (browser)
	if req.Method == http.MethodPost {
		// get form values from browser
		sortBy := req.FormValue("select")
		mapSessionSort[uuid] = sortBy
		fmt.Printf("Sort by: %v, type: %T\n", sortBy, sortBy)

		// selective redirect to based on last menu
		switch mapSessionPreviousMenu[uuid] {

		case "displaySelect":
			// redirect if needed
			http.Redirect(res, req, "/displayList", http.StatusSeeOther)
			return
		}
	}
}

func postRedirectPickItems(uuid string, res http.ResponseWriter,
	req *http.Request) {
	Trace.Println("postRedirectPickItems")
	// Get selection from requester (browser)
	if req.Method == http.MethodPost {
		// get form values from browser
		// selected := req.FormValue("selected")
		// mapSessionSelect[uuid] = selected
		req.ParseForm()
		selected := req.Form["selected"]
		mapSessionSelect[uuid] = make([]string, len(selected))
		mapSessionSelect[uuid] = selected
		fmt.Printf("Selected Item: %v, Type: %T\n", selected, selected)
		fmt.Println("MapSessionSelect", mapSessionSelect[uuid])

		// selective redirect to based on last menu
		switch mapSessionPreviousMenu[uuid] {
		case "myTrayToGive":
			http.Redirect(res, req, "/withdrawItem", http.StatusSeeOther)
			return

		case "myTrayGiven":
			fmt.Println("Tray-myTrayGiven")
			http.Redirect(res, req, "/removeFromMyTray", http.StatusSeeOther)
			return

		case "myTrayGotten":
			fmt.Println("Tray-myTrayGotten")
			http.Redirect(res, req, "/removeFromMyTray", http.StatusSeeOther)
			return

		case "myTrayWithdrawn":
			fmt.Println("Tray-myWithdrawn")
			http.Redirect(res, req, "/removeFromMyTray", http.StatusSeeOther)
			return

		case "showSearchList":
			fmt.Println("Search List")
			http.Redirect(res, req, "/getListedItems", http.StatusSeeOther)
			return
		}
	}
}

func searchItem(res http.ResponseWriter, req *http.Request) {
	Trace.Println("SearchItem")
	// precautionary - to handle direct URL pattern for signup
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	// update last visit
	updateLastVist(myCookie.Value, fmt.Sprintf("Passed : %s", time.Now().Format("Jan-02-2006, 3:04:05 pm")))
	mapSessionPreviousMenu[myCookie.Value] = "searchItem"
	// fmt.Println(myCookie.Value)

	postRedirectNameDescription(myCookie.Value, res, req)

	var menu menuNameDescType
	menu.PageType = 2
	// use last value as default
	menu.DefName = mapSessionItemName[myCookie.Value]
	menu.DefDescription = mapSessionItemDescription[myCookie.Value]

	// fmt.Println("For First Time http Post Request")
	// This is executed for first entry when http post hasn't happen yet
	tpl.ExecuteTemplate(res, "getNameDescription.gohtml", menu)
}

// Shows the list of searched items
func showSearchList(res http.ResponseWriter, req *http.Request) {
	Trace.Println("showSearchList")
	// precautionary - to handle direct URL pattern for signup
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	mapSessionPreviousMenu[myCookie.Value] = "showSearchList"
	// fmt.Println(myCookie.Value)
	postRedirectPickItems(myCookie.Value, res, req)

	// Information to be used for search of DB
	name := mapSessionItemName[myCookie.Value]
	itemDescription := mapSessionItemDescription[myCookie.Value]
	searchLogic := mapSessionSearchLogic[myCookie.Value]

	fmt.Println("Name :", name)
	fmt.Println("Description :", itemDescription)
	fmt.Println("Search Logic :", searchLogic)

	if len(name)+len(itemDescription) == 0 {
		fmt.Println("There is no info on search criteria - use a default listing or no listing?")
	}
	/****************************************/
	// Process the list of item for display
	//test data
	// mr1 := itemType{Id: "MR1", Name: "Clothes", Description: "A box of 10 shirts"}
	// mr2 := itemType{Id: "MR2", Name: "Clothes", Description: "A box of 20 shirts"}
	// mr3 := itemType{Id: "MR3", Name: "Saw", Description: "A 10 inch saw"}
	// mr4 := itemType{Id: "MR4", Name: "Computer", Description: "A Intel Computer and monitor"}
	// mr5 := itemType{Id: "MR5", Name: "Calculator", Description: "A scientific calculator"}
	// mr6 := itemType{Id: "MR6", Name: "Monitor", Description: "Dell Model 123"}
	// mr7 := itemType{Id: "MR7", Name: "Monitor", Description: "LG Model XYZ, 24 inche"}
	// mr8 := itemType{Id: "MR8", Name: "Clothes", Description: "A box of 10 shorts"}
	// mr9 := itemType{Id: "MR9", Name: "Bed Sheets", Description: "3 Queen size bed sheet"}
	// mr10 := itemType{Id: "MR10", Name: "Shoe", Description: "A pair of size 10 shoes for men"}
	// mr11 := itemType{Id: "MR11", Name: "Bed Sheets", Description: "3 king size bed sheet"}
	// mr12 := itemType{Id: "MR12", Name: "Shoe", Description: "A pair of size 12 shoes for men"}
	// list := []itemType{mr1, mr2, mr3, mr4, mr5, mr6, mr7, mr8, mr9, mr10, mr11, mr12}
	// listMenu.List = list

	list, err := bizListSearchItems(name, itemDescription, searchLogic)
	if err != nil {
		http.Error(res, "Error in ShowSearchList ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}

	fmt.Println("ItemList : ", list)

	/****************************************/

	var listMenu menuListType
	// menu for search for browser to select search option before editing
	// listMenu.List = convertToString(list)
	listMenu.List = list
	listMenu.PageType = 5

	// fmt.Println("For First Time http Post Request")
	// This is executed for first entry when http post hasn't happen yet
	tpl.ExecuteTemplate(res, "showPickList.gohtml", listMenu)
}

func toGiveItem(res http.ResponseWriter, req *http.Request) {
	Trace.Println("toGiveItem")
	// precautionary - to handle direct URL pattern for signup
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	// update last visit
	updateLastVist(myCookie.Value, fmt.Sprintf("Passed : %s", time.Now().Format("Jan-02-2006, 3:04:05 pm")))
	mapSessionPreviousMenu[myCookie.Value] = "toGiveItem"
	// fmt.Println(myCookie.Value)

	postRedirectNameDescription(myCookie.Value, res, req)

	var menu menuNameDescType
	menu.PageType = 1
	// use last value as default
	menu.DefName = mapSessionItemName[myCookie.Value]
	menu.DefDescription = mapSessionItemDescription[myCookie.Value]

	// fmt.Println("For First Time http Post Request")
	// This is executed for first entry when http post hasn't happen yet
	tpl.ExecuteTemplate(res, "getNameDescription.gohtml", menu)
}

// displaySelect handler
func displaySelect(res http.ResponseWriter, req *http.Request) {
	Trace.Println("displaySelect")
	// precautionary - to handle direct URL pattern for signup
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	// update last visit
	updateLastVist(myCookie.Value, fmt.Sprintf("Passed : %s", time.Now().Format("Jan-02-2006, 3:04:05 pm")))
	mapSessionPreviousMenu[myCookie.Value] = "displaySelect"
	// fmt.Println(myCookie.Value)

	postRedirectSortSelect(myCookie.Value, res, req) // go next to searchparameter

	var listMenu menuListType
	// menu for search for browser to select search option before editing
	listMenu.List = []string{
		"Unsorted Listing Order",
		"Name",
		"Item State",
		"Date",
		"Giver ID",
		"Getter ID",
	}

	listMenu.PageType = 2
	// get last value as default
	listMenu.DefSelect = sanitizeAtoi(mapSessionSort[myCookie.Value], 0, len(listMenu.List)-1)

	// fmt.Println("For First Time http Post Request")
	// This is executed for first entry when http post hasn't happen yet
	tpl.ExecuteTemplate(res, "selectBy.gohtml", listMenu)
}

func myTrayToGive(res http.ResponseWriter, req *http.Request) {
	myTray(res, req, "myTrayToGive")
}

func myTrayGiven(res http.ResponseWriter, req *http.Request) {
	myTray(res, req, "myTrayGiven")
}

func myTrayGotten(res http.ResponseWriter, req *http.Request) {
	myTray(res, req, "myTrayGotten")
}

func myTrayWithdrawn(res http.ResponseWriter, req *http.Request) {
	myTray(res, req, "myTrayWithdrawn")
}

func myTray(res http.ResponseWriter, req *http.Request, tray string) {
	// test data
	// mr1 := itemType{Id: "MR1", Name: "Clothes", Description: "A box of 10 shirts"}
	// mr2 := itemType{Id: "MR2", Name: "Clothes", Description: "A box of 20 shirts"}
	// mr3 := itemType{Id: "MR3", Name: "Saw", Description: "A 10 inch saw"}
	// mr4 := itemType{Id: "MR4", Name: "Computer", Description: "A Intel Computer and monitor"}
	// mr5 := itemType{Id: "MR5", Name: "Calculator", Description: "A scientific calculator"}
	// mr6 := itemType{Id: "MR6", Name: "Monitor", Description: "Dell Model 123"}
	// mr7 := itemType{Id: "MR7", Name: "Monitor", Description: "LG Model XYZ, 24 inche"}
	// mr8 := itemType{Id: "MR8", Name: "Clothes", Description: "A box of 10 shorts"}
	// mr9 := itemType{Id: "MR9", Name: "Bed Sheets", Description: "3 Queen size bed sheet"}
	// mr10 := itemType{Id: "MR10", Name: "Shoe", Description: "A pair of size 10 shoes for men"}
	// mr11 := itemType{Id: "MR11", Name: "Bed Sheets", Description: "3 king size bed sheet"}
	// mr12 := itemType{Id: "MR12", Name: "Shoe", Description: "A pair of size 12 shoes for men"}
	// list := []itemType{mr1, mr2, mr3, mr4, mr5, mr6, mr7, mr8, mr9, mr10, mr11, mr12}
	//fmt.Println(list)

	// precautionary - to handle direct URL pattern for signup
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	// update last visit
	mapSessionPreviousMenu[myCookie.Value] = tray
	postRedirectPickItems(myCookie.Value, res, req)

	var listMenu menuListType

	// List Items base on the tray type

	var err error
	listMenu.List, err = bizMyTrayItems(tray)
	if err != nil {
		http.Error(res, "Error in bizMyTrayItems ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}

	// listMenu.List = convertToString(list)
	fmt.Println("Tray Selected Items:", tray)
	// **********************

	switch tray {
	case "myTrayToGive":
		listMenu.PageType = 1
	case "myTrayGiven":
		listMenu.PageType = 2
	case "myTrayGotten":
		listMenu.PageType = 3
	case "myTrayWithdrawn":
		listMenu.PageType = 4
	}

	// fmt.Println("For First Time http Post Request")
	// This is executed for first entry when http post hasn't happen yet
	tpl.ExecuteTemplate(res, "showPickList.gohtml", listMenu)
}

func withdrawItem(res http.ResponseWriter, req *http.Request) {
	// precautionary - to handle direct URL pattern for signup
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	selectedList := mapSessionSelect[myCookie.Value]
	fmt.Println("selected for withdrawal", selectedList)
	/******************************/
	// withdrawal from the displayed
	msg, err := bizWithdrawItems(selectedList)
	if err != nil {
		http.Error(res, "Error in bizWithdrawItems ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}
	/******************************/

	// var msg []string
	// msg = append(msg, "Item1 withdrawn")
	// msg = append(msg, "Item2 withdrawn")

	showMessages(res, req, 2, msg)
}

func getListedItems(res http.ResponseWriter, req *http.Request) {

	// precautionary - to handle direct URL pattern for signup
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	lastmenu := mapSessionPreviousMenu[myCookie.Value]
	selectedList := mapSessionSelect[myCookie.Value]
	fmt.Println("Last Menu :", lastmenu)
	fmt.Println("selected to get :", selectedList)
	/******************************/
	// Process the list of items picked by getter

	msg, err := bizGetListedItems(selectedList)
	if err != nil {
		http.Error(res, "Error in bizGetListedItems ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}
	/******************************/

	//msg := fmt.Sprintf("Items are moved from listing into Gotten Tray")
	// var msg []string
	// msg = append(msg, "Item1 Gotten")
	// msg = append(msg, "Item2 Gotten")

	showMessages(res, req, 4, msg)
}

func giveItem(res http.ResponseWriter, req *http.Request) {

	// precautionary - to handle direct URL pattern for signup
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	name := mapSessionItemName[myCookie.Value]
	description := mapSessionItemDescription[myCookie.Value]
	lastmenu := mapSessionPreviousMenu[myCookie.Value]

	fmt.Println("Last Menu :", lastmenu)
	fmt.Println("Name :", name)
	fmt.Println("Description :", description)
	/******************************/
	// Process the item for listing, change item to "togive" state

	msg, err := bizGiveItem(name, description)
	if err != nil {
		http.Error(res, "Error in bizGiveItem ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}

	/******************************/
	// var msg []string
	// msg = append(msg, "Item Given :"+name+", "+description+" is given") // One one item

	showMessages(res, req, 5, msg)
}

func removeFromMyTray(res http.ResponseWriter, req *http.Request) {

	// precautionary - to handle direct URL pattern for signup
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	tray := mapSessionPreviousMenu[myCookie.Value]
	selectedList := mapSessionSelect[myCookie.Value]
	fmt.Println("Tray to Remove from View", tray)
	fmt.Println("selected for withdrawal", selectedList)
	/******************************/
	// Undisplayed the picked list from Tray (Given Items, Gotten Items, Withdrawn Items)
	// This common for all 3 group, so check the tray to know which group

	/******************************/
	msg, err := bizRemoveFromTray(selectedList, tray)
	if err != nil {
		http.Error(res, "Error in bizRemoveFromTray ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}
	// var msg []string
	// msg1 := fmt.Sprintf("list of items are removed from %s", tray)
	// msg = append(msg, msg1)
	// msg = append(msg, "Item1 removed from Tray")
	// msg = append(msg, "Item2 removed from Tray")

	showMessages(res, req, 3, msg)
}

// Send a fix message
func showMessage(res http.ResponseWriter, req *http.Request, pageType int, message string) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	var msg messageType
	msg.PageType = pageType
	msg.Msg = append(msg.Msg, message)

	tpl.ExecuteTemplate(res, "showMessage.gohtml", msg)
}

// Send a list of messages
func showMessages(res http.ResponseWriter, req *http.Request, pageType int, message []string) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	var msg messageType
	msg.PageType = pageType
	msg.Msg = make([]string, len(message))
	copy(msg.Msg, message)

	tpl.ExecuteTemplate(res, "showMessage.gohtml", msg)
}

func displayList(res http.ResponseWriter, req *http.Request) {
	// display the list
	// test data
	// mr1 := itemType{Id: "MR1", Name: "Clothes", Description: "A box of 10 shirts"}
	// mr2 := itemType{Id: "MR2", Name: "Clothes", Description: "A box of 20 shirts"}
	// mr3 := itemType{Id: "MR3", Name: "Saw", Description: "A 10 inch saw"}
	// mr4 := itemType{Id: "MR4", Name: "Computer", Description: "A Intel Computer and monitor"}
	// mr5 := itemType{Id: "MR5", Name: "Calculator", Description: "A scientific calculator"}
	// mr6 := itemType{Id: "MR6", Name: "Monitor", Description: "Dell Model 123"}
	// mr7 := itemType{Id: "MR7", Name: "Monitor", Description: "LG Model XYZ, 24 inche"}
	// mr8 := itemType{Id: "MR8", Name: "Clothes", Description: "A box of 10 shorts"}
	// mr9 := itemType{Id: "MR9", Name: "Bed Sheets", Description: "3 Queen size bed sheet"}
	// mr10 := itemType{Id: "MR10", Name: "Shoe", Description: "A pair of size 10 shoes for men"}
	// mr11 := itemType{Id: "MR11", Name: "Bed Sheets", Description: "3 king size bed sheet"}
	// mr12 := itemType{Id: "MR12", Name: "Shoe", Description: "A pair of size 12 shoes for men"}
	// list := []itemType{mr1, mr2, mr3, mr4, mr5, mr6, mr7, mr8, mr9, mr10, mr11, mr12}
	// //fmt.Println(list)
	// // selected soted selection
	// msg := convertToString(list)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")

	sortBy := mapSessionSort[myCookie.Value]
	fmt.Println("Sort Selected", sortBy)

	// ***********************
	// // Get the sorted list here !!
	msg, err := bizGetSortedList(sortBy)
	if err != nil {
		http.Error(res, "Error in bizGetSortedList", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}

	showMessages(res, req, 1, msg)
}
