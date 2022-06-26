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
var mapSessionSearchedList = map[string][]Item{}
var mapSessionMyTrayList = map[string][]Item{}

var mapSessionPreviousMenu = map[string]string{}

// postRedirectionNameDescription
// posting redirection from the
// searchItem()  --> /showSearchList  (after name and description is obtained)
// toGiveItem()  --> /giveItem (after name and description is obtained)
// To get the input for name and description
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
			fmt.Println("toGiveItem")
			http.Redirect(res, req, "/giveItem", http.StatusSeeOther)
			return
		}
	}
}

// postRedirectSortSelect post http to get parameter (select option) for the following
// displaySelect() --> /displayList to display the list of sorted items (after the picklist is obtained)
func postRedirectSortSelect(uuid string, res http.ResponseWriter,
	req *http.Request) {
	Trace.Println("postRedirectSortSelect")
	// Get selection from requester (browser)
	if req.Method == http.MethodPost {
		// get form values from browser
		sortBy := req.FormValue("select")
		mapSessionSort[uuid] = sortBy
		//		fmt.Printf("Sort by: %v, type: %T\n", sortBy, sortBy)

		// selective redirect to based on last menu
		switch mapSessionPreviousMenu[uuid] {

		case "displaySelect":
			// redirect if needed
			http.Redirect(res, req, "/displayList", http.StatusSeeOther)
			return
		}
	}
}

// postRedirectPickItems
// Post redirection after the picklist of obtained based on submit button input actions
// Withdraw Item --> /withdrawItem
// View Giver Details  --> /viewGiverDetails
// View Getter Details --> /viewGetterDetails
// Remove From Tray --> /removeFromMyTray
// Get Items --> /getListedItems
func postRedirectPickItems(uuid string, res http.ResponseWriter,
	req *http.Request) {
	Trace.Println("postRedirectPickItems")
	// Get selection from requester (browser)
	if req.Method == http.MethodPost {

		req.ParseForm()
		selected := req.Form["selected"]
		submit := req.Form["submit"] // get the submit button
		fmt.Println("Submit =", submit[0])

		mapSessionSelect[uuid] = make([]string, len(selected))
		mapSessionSelect[uuid] = selected
		fmt.Printf("Selected Item: %v, Type: %T\n", selected, selected)
		fmt.Println("MapSessionSelect", mapSessionSelect[uuid])

		//selective redirect to based on submit button
		switch submit[0] {
		case "Withdraw Item":
			http.Redirect(res, req, "/withdrawItem", http.StatusSeeOther)
			return

		case "View Giver Details":
			http.Redirect(res, req, "/viewGiverDetails", http.StatusSeeOther)
			return

		case "View Getter Details":
			http.Redirect(res, req, "/viewGetterDetails", http.StatusSeeOther)
			return

		// common for Given Items, Gotten items, Withdrawn Items
		case "Remove From Tray":
			fmt.Println("Tray-myWithdrawn")
			http.Redirect(res, req, "/removeFromMyTray", http.StatusSeeOther)
			return

		case "Get Items":
			fmt.Println("Search List")
			http.Redirect(res, req, "/getListedItems", http.StatusSeeOther)
			return
		}
	}
}

// seacrhItem
// called from /searchItem from index to display page for name and description input
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

// showSearchList (redirected handler)
// shows the list of searched items
// called from /showSearchList after the name and description is submitted
// and display the list as a picklist (by displaying the picklist)
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

	searchedItems, err := bizListSearchItems(name, itemDescription, searchLogic)
	mapSessionSearchedList[myCookie.Value] = make([]Item, len(searchedItems))
	mapSessionSearchedList[myCookie.Value] = searchedItems

	if err != nil {
		http.Error(res, "Error in ShowSearchList ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}

	var list []string

	list = showIdNameDescriptionDate2String(mapSessionSearchedList[myCookie.Value])

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

// toGiveItem handler
// called from index page from /toGiveItem and proceed to get name and description for an item to give away
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

	tpl.ExecuteTemplate(res, "getNameDescription.gohtml", menu)
}

// displaySelect handler
// called from index page /displaySelect to display the sort options for display
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

	tpl.ExecuteTemplate(res, "selectBy.gohtml", listMenu)
}

// myTrayToGive (redirected handler)
// called from index to display the myTray list of items for "To Give" items
func myTrayToGive(res http.ResponseWriter, req *http.Request) {
	myTray(res, req, "myTrayToGive")
}

// myTrayGiven (redirected handler)
// called from index to display the myTray list of items for "Given" items
func myTrayGiven(res http.ResponseWriter, req *http.Request) {
	myTray(res, req, "myTrayGiven")
}

// myTrayGotten (redirected handler)
// called from to display the myTray list of items for "Gotten" items
func myTrayGotten(res http.ResponseWriter, req *http.Request) {
	myTray(res, req, "myTrayGotten")
}

// myTrayWithdrawn (redirected handler)
// called from index to display the myTray list of items for "Withdrawn" items
func myTrayWithdrawn(res http.ResponseWriter, req *http.Request) {
	myTray(res, req, "myTrayWithdrawn")
}

// myTray
// common functions for
// myTrayToGive, myTrayGiven, myTrayGotten, myTrayWithdrawn
func myTray(res http.ResponseWriter, req *http.Request, tray string) {

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

	var err error
	// capture my tray list of items
	mapSessionMyTrayList[myCookie.Value], err = bizMyTrayItems(mapSessions[myCookie.Value], tray)
	if err != nil {
		http.Error(res, "Error in bizMyTrayItems ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}

	// convert to [] string
	listMenu.List = showIdNameDescriptionDate2String(mapSessionMyTrayList[myCookie.Value])

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

	tpl.ExecuteTemplate(res, "showPickList.gohtml", listMenu)
}

// withdrawItem (redirected handler)
// called from postRedirectPickItems (picked list) to withdrawn the picked items
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
	// to set selected items (from MyTraylList) to "withdraw" state
	msg, err := bizWithdrawItems(mapSessionMyTrayList[myCookie.Value], selectedList)
	if err != nil {
		http.Error(res, "Error in bizWithdrawItems ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}
	/******************************/

	showMessages(res, req, 2, msg)
}

// getListedItems (redirected handler)
// called from postRedirectPickItems (picked list) to get the picked items and put them into "Given" state
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

	msg, err := bizGetListedItems(myCookie.Value, selectedList)
	if err != nil {
		http.Error(res, "Error in bizGetListedItems ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}
	/******************************/

	showMessages(res, req, 4, msg)
}

// giveItem (redirected handler)
//postRedirectNameDescription after the name and description are entered to put items to "ToGive" state
func giveItem(res http.ResponseWriter, req *http.Request) {

	// precautionary - to handle direct URL pattern for signup
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	name := mapSessionItemName[myCookie.Value]
	description := mapSessionItemDescription[myCookie.Value]
	username := mapSessions[myCookie.Value]
	lastmenu := mapSessionPreviousMenu[myCookie.Value]

	fmt.Println("Last Menu :", lastmenu)
	// fmt.Println("Name :", name)
	// fmt.Println("Description :", description)

	/******************************/
	// Process the item for listing, change item to "togive" state
	msg, err := bizGiveItem(name, description, username)
	if err != nil {
		http.Error(res, "Error in bizGiveItem ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}
	/******************************/

	showMessages(res, req, 5, msg)
}

// removeFromTray (redirected handler)
// called from postRedirectPickItems to process of pick list items to hide from display from tray
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

	msg, err := bizRemoveFromTray(mapSessionMyTrayList[myCookie.Value], selectedList, tray)
	if err != nil {
		http.Error(res, "Error in bizRemoveFromTray ", http.StatusInternalServerError)
		fmt.Println("Error :", err)
		return
	}
	/******************************/
	showMessages(res, req, 3, msg)
}

// ShowMeessages send a list of messages to client for display
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

// displayList (redirected hanlder)
// called from postRedirectSortSelect after sorted key is entered to as to get the sorted list of items
// to be displayed
func displayList(res http.ResponseWriter, req *http.Request) {

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
		http.Error(res, "Server Error", http.StatusInternalServerError)
		fmt.Println("Error in bizGetSortedList :", err)
		return
	}

	showMessages(res, req, 1, msg)
}

// viewGiverDetails (redirected handler)
// called from postRedirectPickItems after the picklist to get Giver's contact details for each item
// and display them
func viewGiverDetails(res http.ResponseWriter, req *http.Request) {

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	selectedList := mapSessionSelect[myCookie.Value]

	fmt.Println("Selected :", selectedList)
	// Get the items picked and Giver's details
	msg1, err := bizGetItemWithGiverDetails(mapSessionMyTrayList[myCookie.Value], selectedList)
	if err != nil {
		http.Error(res, "Server Error", http.StatusInternalServerError)
		fmt.Println("Error in bizGetItemWithGiverDetails:", err)
		return
	}
	fmt.Println("viewGiverDetails")

	showMessages(res, req, 6, msg1)
}

// viewGetterDetails (redirected handler)
// called from postRedirectPickItems after the picklist get the Giver's contact details for each item
// and display them
func viewGetterDetails(res http.ResponseWriter, req *http.Request) {

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	selectedList := mapSessionSelect[myCookie.Value]
	fmt.Println("Selected :", selectedList)
	msg1, err := bizGetItemWithGetterDetails(mapSessionMyTrayList[myCookie.Value], selectedList)
	if err != nil {
		http.Error(res, "Server Error", http.StatusInternalServerError)
		fmt.Println("Error in bizGetItemWithGetterDetails :", err)
		return
	}
	fmt.Println("viewGetterDetails")

	showMessages(res, req, 7, msg1)
}
