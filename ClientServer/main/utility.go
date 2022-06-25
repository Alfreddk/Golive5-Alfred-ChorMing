package main

import (
	"fmt"
	"strconv"
)

/* // alfred 24.06.2022: not used.
// mapKeyPresent implement the method to check the the presence of the key for
// monthMapType
func (m monthMapType) mapKeyPresent(k int) (ok bool) {
	_, ok = m[k]
	return
}

// mapKeyPresent implement the method to check the the presence of the key for
// timeSlotType
func (m timeSlotType) mapKeyPresent(k int) (ok bool) {
	_, ok = m[k]
	return
}

// getDateStr return the name string of the date
func getDateStr(date dateType) string {
	return fmt.Sprintf("%d, %s %d", date.day, nameOfMonth[date.month], date.year)
}

// getTimeSlotStr return the name string of the timeSloNum
func getTimeSlotStr(timeSlotNum int) string {
	return fmt.Sprintf("<%s>", timeSlot[timeSlotNum])
}

// onvertYYYYMMDD_toDateType convert html YYYYMMDD date format to our dateType format
func convertYYYYMMDD_toDateType(dateIn string) dateType {
	var date dateType
	if dateIn == "" {
		date = invalidDate
	} else {
		date.year, _ = strconv.Atoi(dateIn[:4])
		date.month, _ = strconv.Atoi(dateIn[5:7])
		date.day, _ = strconv.Atoi(dateIn[8:10])
		if date.year <= 0 || date.month < 1 || date.month > 12 || date.day < 1 || date.day > 31 {
			date = invalidDate
		}
	}
	return date
}

// convertDateType_toYYYYMMDD convert dateType format to html YYYYMMDD date format
func convertDateType_toYYYYMMDD(date dateType) string {

	year := strconv.Itoa(date.year)
	month := strconv.Itoa(date.month)
	day := strconv.Itoa(date.day)
	for len(year) < 4 {
		year = "0" + year
	}
	for len(month) < 2 {
		month = "0" + month
	}
	for len(day) < 2 {
		day = "0" + day
	}
	yyyyMMdd := year + "-" + month + "-" + day
	return yyyyMMdd
}


//ShowIndex displays the index of the value passed if it is positive
// display index value if index > 0, otherwise displays "Not Found"
func showIndex(index int) {
	if index < 0 {
		fmt.Println("Not Found")
	} else {
		fmt.Printf("Index : %d\n", index)
	}
}

// GetDefValue checks integer for unInitialised value.
// Returns the default initialValue, if the current value is the unInitialised value, otherwise the current value
func getDefValue(currentValue int, unInitialisedValue int, initialValue int) int {
	if currentValue == unInitialisedValue {
		return initialValue
	} else {
		return currentValue
	}
}
*/

// sanitizeAtoi convert from Ascii to Integer.
// Checks that input string number is within "first" and "last" integer converted string value,
// returns the "first" int if value is out of range, otherwise, the string value is converted to integer
func sanitizeAtoi(input string, first int, last int) int {

	// check that input is witjhn expected range
	if input >= strconv.Itoa(first) && input <= strconv.Itoa(last) {
		value, _ := strconv.Atoi(input)
		return value
	} else {
		return first
	}
}

/* // alfred 24.06.2022: not used.
// Convert items data to string slices
func convertToString(data []itemType) []string {
	listdata := []string{}

	for _, v := range data {
		listdata = append(listdata, v.Id+" "+v.Name+" "+v.Description)
	}

	return listdata
}
*/

// ConvertItems2String converts items data to string slices.
func convertItems2String(data []Item) []string {

	fmt.Println("convertItems2String")
	listdata := []string{}
	for _, v := range data {
		hideStatus := fmt.Sprintf("%v:%v:%v", v.HideGiven, v.HideGotten, v.HideWithdrawn)
		//state := fmt.Sprintf("%d", v.State)
		listdata = append(listdata, "ID: "+v.ID+",  Name: "+v.Name+",  Summary: "+v.Description+
			",  State: "+itemState[v.State]+",  Hide: "+hideStatus+",  GiverID: "+v.GiverUsername+
			",  GetterID: "+v.GetterUsername+",  Date: "+v.Date)
	}
	return listdata
}

// Convert items data to string slices
func convertNameFirst2String(data []Item) []string {

	listdata := []string{}
	for _, v := range data {
		hideStatus := fmt.Sprintf("%v:%v:%v", v.HideGiven, v.HideGotten, v.HideWithdrawn)
		//state := fmt.Sprintf("%d", v.State)
		listdata = append(listdata, "Name: "+v.Name+",  ID: "+v.ID+",  Summary: "+v.Description+
			",  State: "+itemState[v.State]+",  Hide: "+hideStatus+",  GiverID: "+v.GiverUsername+
			",  GetterID: "+v.GetterUsername+",  Date: "+v.Date)
	}
	return listdata
}

// Convert items data to string slices
func convertStateFirst2String(data []Item) []string {

	fmt.Println("convertItems2String")
	listdata := []string{}
	for _, v := range data {
		hideStatus := fmt.Sprintf("%v:%v:%v", v.HideGiven, v.HideGotten, v.HideWithdrawn)
		//state := fmt.Sprintf("%d", v.State)
		listdata = append(listdata, "State: "+itemState[v.State]+",  ID: "+v.ID+",  Name: "+v.Name+
			",  Summary: "+v.Description+",  Hide: "+hideStatus+",  GiverID: "+v.GiverUsername+
			",  GetterID: "+v.GetterUsername+",  Date: "+v.Date)
	}
	return listdata
}

// Convert items data to string slices
func convertDateFirst2String(data []Item) []string {

	fmt.Println("convertItems2String")
	listdata := []string{}
	for _, v := range data {
		hideStatus := fmt.Sprintf("%v:%v:%v", v.HideGiven, v.HideGotten, v.HideWithdrawn)
		//state := fmt.Sprintf("%d", v.State)
		listdata = append(listdata, "Date: "+v.Date+",  ID: "+v.ID+",  Name: "+v.Name+",  Summary: "+v.Description+
			",  State: "+itemState[v.State]+",  Hide: "+hideStatus+",  GiverID: "+v.GiverUsername+
			",  GetterID: "+v.GetterUsername)
	}
	return listdata
}

// Convert items data to string slices
func convertGiverIDFirst2String(data []Item) []string {

	fmt.Println("convertItems2String")
	listdata := []string{}
	for _, v := range data {
		hideStatus := fmt.Sprintf("%v:%v:%v", v.HideGiven, v.HideGotten, v.HideWithdrawn)
		//state := fmt.Sprintf("%d", v.State)
		listdata = append(listdata, "GiverID: "+v.GiverUsername+",  ID: "+v.ID+",  Name: "+v.Name+",  Summary: "+v.Description+
			",  State: "+itemState[v.State]+", Hide: "+hideStatus+
			",  GetterID: "+v.GetterUsername+",  Date: "+v.Date)
	}
	return listdata
}

// Convert items data to string slices
func convertGetterIDFirst2String(data []Item) []string {

	fmt.Println("convertItems2String")
	listdata := []string{}
	for _, v := range data {
		hideStatus := fmt.Sprintf("%v:%v:%v", v.HideGiven, v.HideGotten, v.HideWithdrawn)
		//state := fmt.Sprintf("%d", v.State)
		listdata = append(listdata, "GetterID: "+v.GetterUsername+",  ID: "+v.ID+",  Name: "+v.Name+",  Summary: "+v.Description+
			",  State: "+itemState[v.State]+",  Hide: "+hideStatus+",  GiverID: "+v.GiverUsername+
			",  Date: "+v.Date)
	}
	return listdata
}

// Convert items data to string slices
func showIdNameDescriptionDate2String(data []Item) []string {

	fmt.Println("showIdNameDescriptionDate2String")
	listdata := []string{}
	for _, v := range data {
		listdata = append(listdata, "ID: "+v.ID+",  Name: "+v.Name+",  Summary: "+v.Description+",  Given On: "+v.Date)
	}
	return listdata
}

func formGiverGetterDetails(msg *[]string, role string, userID string, item Item) {

	s := *msg
	//var msg []string
	var str1, str2, str3, str4 string
	str1 = "ID: " + item.ID + ",  Name: " + item.Name + ",  Summary: " + item.Description
	str2 = role + " Contact:"
	str3 = "Name: " + mapUsers[userID].Name + ",  Phone:  " + mapUsers[userID].Telephone
	str4 = "Address:  " + mapUsers[userID].Address + ", S" + mapUsers[userID].Postal
	fmt.Println("str", str1)
	fmt.Println("str2 =", str2)
	fmt.Println("str2 =", str3)
	fmt.Println("str2 =", str4)
	*msg = append(s, str1, str2, str3, str4)
	//	return msg
}
