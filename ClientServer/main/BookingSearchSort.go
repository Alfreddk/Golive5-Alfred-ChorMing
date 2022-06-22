package main

import (
	"errors"
	"fmt"
)

// Type for receiver use
// type bookingsType []booking

// The same keys are used for search and sort to ensure consistent results
// between search and sort
type bookingKeyType int

// const is search or sort key used for search or sorting
//bookName bookingKeyType = iota  	// for sort/search in order of booking name
//bookDateAndTimeSlot				// for sort/search in order of Date and TimeSlot
//bookDate							// for sort/search in order of Date
//bookTimeSlot						// for sort/search in order of TimeSlot
//bookVenueName				    	// for sort/search in order of venue name
//bookVenueTypeName             	// for sort/search in order of venue type name
//bookDateTimeSlotVenueName     	// for sort/search in order of priority date, timeslot and venue name
//bookDateTimeSlotVenueTypeName 	// for sort/search in order of priority date, timeslot and venue type name
//bookVenueCapacity			 		// for sort/search in order of venue capacity
const (
	bookName         bookingKeyType = iota
	bookFullyMatched                // only for search
	bookDateAndTimeSlotAndVenueTypeNum
	bookDateAndTimeSlot
	bookDate
	bookTimeSlot
	bookVenueTypeName
	bookVenueName
	bookDateTimeSlotVenueName
	bookDateTimeSlotVenueTypeName
	bookVenueCapacity
)

// iterativeSeqUnsortedSearch is the iterative search using sequential search algorithm for *[]booking
// return the position of index (0 based)
// return -1, if not found
// Sequentual Search for unsorted search
// This algorithm applies for both sorted and unsorted array, to find the first occurance that matches but
// for sorted array, the first occurance of an element found is also the first element in the sorted array to have this value
func iterativeSeqUnsortedSearch(n int, arr *[]booking, target booking, key bookingKeyType) int {
	for i := 0; i < n; i++ {
		if (*arr)[i].isEqual(target, key) { // return when target is found
			//			fmt.Println("Found  ")
			return i
		}
	}
	//	fmt.Println("Not Found  ")
	return -1
}

// iterativeSeqUnsortedSearch is the iterative search using sequential search algorithm for *[]* booking
// return the position of index (0 based)
// return -1, if not found
// Sequentual Search for unsorted search
// This algorithm applies for both sorted and unsorted array, to find the first occurance that matches but
// for sorted array, the first occurance of an element found is also the first element in the sorted array to have this value
func iterativeSeqUnsortedSearch2(n int, arr *[]*booking, target booking, key bookingKeyType) int {
	ptrs := *arr
	for i := 0; i < n; i++ {
		if (*ptrs[i]).isEqual(target, key) { // return when target is found
			//			fmt.Println("Found  ")
			return i
		}
	}
	//	fmt.Println("Not Found  ")
	return -1
}

// iterativeBinSearch is the iterative search for sorted array using binary search algorithm for *[] booking
// return the position of index (0 based) of the first occurance
// return -1, if not found
func iterativeBinSearch(n int, arr *[]booking, target booking, key bookingKeyType) int {
	first := 0
	last := n - 1

	for first <= last {
		mid := (first + last) / 2

		if (*arr)[mid].isEqual(target, key) {
			return mid
		} else if (*arr)[mid].isLargerThan(target, key) {
			last = mid - 1 // (tested value) > target, so answer in the lower half, so shift down the upper bound
		} else {
			first = mid + 1 // (tested value) < target, so answer in the upper hald, so shift up the lower bound
		}
	}
	return -1
}

// iterativeBinSearch2 is the iterative search for sorted array using binary search algorithm for *[]* booking
// return the position of index (0 based) of the first occurance
// return -1, if not found
func iterativeBinSearch2(n int, arr *[]*booking, target booking, key bookingKeyType) int {
	ptrs := *arr
	first := 0
	last := n

	for first <= last {
		mid := (first + last) / 2
		//		fmt.Printf("First = %d, Mid = %d, Last = %d\n", first, mid, last)

		if (*ptrs[mid]).isEqual(target, key) {
			return mid
		} else if (ptrs[mid]).isLargerThan(target, key) {
			last = mid - 1 // (tested value) > target, so answer in the lower half, so shift down the upper bound
		} else {
			first = mid + 1 // (tested value) < target, so answer in the upper hald, so shift up the lower bound
		}
	}
	return -1
}

// recursiveBinSearch is the recursive search for sorted array using binary search algorithm for *[] booking
// first - first element to be searched
// last  - last element to be searched
// return index >= 0 if found
// return index -1 if not found
//  This algorithm finds the first occurance of an element (Not necessarily the first element in the sorted array)
func recursiveBinSearch(n int, arr *[]booking, target booking, key bookingKeyType, first *int, last *int) int {

	if *first > *last { // base case when first > last element
		return -1
	} else {
		mid := (*first + *last) / 2           // find the middle element
		if (*arr)[mid].isEqual(target, key) { // look for match case for early termination
			return mid
		} else if (*arr)[mid].isLargerThan(target, key) { // check value to decide lower pointer be increased or upper pointer be lowered
			*last = mid - 1 // tested value > target mean the answer is in the lower half, so move down upper bound
		} else {
			*first = mid + 1 // test valued < taget, then answer must be in the upper half, so move up the lower bound
		}
		return recursiveBinSearch(n, arr, target, key, first, last)
	}
}

// recursiveBin is the wrapper for recursive search recursiveBinSearch
// first - first element to be searched
// last  - last element to be searched
// return index >= 0 if found
// return index -1 if not found
//  This algorithm finds the first occurance of an element (Not necessarily the first element in the sorted array)
func recursiveBin(n int, arr *[]booking, target booking, key bookingKeyType) int {

	// n := len(*arr)
	first := 0
	last := n
	var index int
	index = recursiveBinSearch(n, arr, target, key, &first, &last)
	return index
}

// recursiveBinSearch2 is the recursive search for sorted array using binary search algorithm for *[]* booking
// first - first element to be searched
// last  - last element to be searched
// return index >= 0 if found
// return index -1 if not found
//  This algorithm finds the first occurance of an element (Not necessarily the first element in the sorted array)
func recursiveBinSearch2(n int, arr *[]*booking, target booking, key bookingKeyType, first *int, last *int) int {
	ptrs := *arr
	if *first > *last { // base case when first > last element
		return -1
	} else {
		mid := (*first + *last) / 2            // find the middle element
		if (*ptrs[mid]).isEqual(target, key) { // look for match case for early termination
			return mid
		} else if (*ptrs[mid]).isLargerThan(target, key) { // check value to decide lower pointer be increased or upper pointer be lowered
			*last = mid - 1 // tested value > target mean the answer is in the lower half, so move down upper bound
		} else {
			*first = mid + 1 // test valued < taget, then answer must be in the upper half, so move up the lower bound
		}
		return recursiveBinSearch2(n, arr, target, key, first, last)
	}
}

// recursiveBin2 is the wrapper for recursive search recursiveBinSearch2
// first - first element to be searched
// last  - last element to be searched
// return index >= 0 if found
// return index -1 if not found
//  This algorithm finds the first occurance of an element (Not necessarily the first element in the sorted array)
func recursiveBin2(n int, arr *[]*booking, target booking, key bookingKeyType) int {

	// n := len(*arr)
	first := 0
	last := n
	var index int
	index = recursiveBinSearch2(n, arr, target, key, &first, &last)
	return index
}

// inserttionSortAscending is the insertion sort for *[]booking sorted according to sortkey
func insertionSortAscending(arr *[]booking, sortKey bookingKeyType) {
	var i int
	// repeat this insertion for every considered element from 1 to N
	for i = 1; i < len(*arr); i++ { // "i" item is the element in the unsorted list considered for insertion
		// i is the position (start from 1) of the element (in unsorted list) considered for insertion into the sorted list
		// element < i are in the sorted list (start from 0)
		data := (*arr)[i] // "i" item position is copied first and the position is freed up for later use in the shift
		// 2. shift (consider element i) into the sorted list to a positon above the element with value below the considered element i
		var pos int // current position considered for shifting in sorted list
		// for every elements within the sorted list (element i-1 to element 0)
		// shift whenever the sorted list elements (on the right (pos)-->(pos+1) are larger than the i element (unsorted)
		// exit when a lower element is found or the end of the search is reached
		for pos = i - 1; pos >= 0 && (*arr)[pos].isLargerThan(data, sortKey); pos-- {
			// (sorted list) element with value higher than (considered element i) are shift right 1 position to make space for insertion later
			(*arr)[pos+1] = (*arr)[pos]
		}
		(*arr)[pos+1] = data // cur is the right most element (sorted list) that is lower than the (considered element i), so insert at cur+1
	}
}

// inserttionSortAscending is the insertion sort for *[]*booking sorted according to sortkey
func insertionSortAscending2(arr *[]*booking, sortKey bookingKeyType) {
	var i int
	ptrs := *arr
	// repeat this insertion for every considered element from 1 to N
	for i = 1; i < len(ptrs); i++ { // "i" item is the element in the unsorted list considered for insertion
		// i is the position (start from 1) of the element (in unsorted list) considered for insertion into the sorted list
		// element < i are in the sorted list (start from 0)
		dataPtr := ptrs[i] // "i" item position is copied first and the position is freed up for later use in the shift
		// 2. shift (consider element i) into the sorted list to a positon above the element with value below the considered element i
		var pos int // current position considered for shifting in sorted list
		// for every elements within the sorted list (element i-1 to element 0)
		// shift whenever the sorted list elements (on the right (pos)-->(pos+1) are larger than the i element (unsorted)
		// exit when a lower element is found or the end of the search is reached
		//		for pos = i - 1; pos >= 0 && isLargerThan2(ptrs[pos], data, sortKey); pos-- {
		for pos = i - 1; pos >= 0 && (*ptrs[pos]).isLargerThan((*dataPtr), sortKey); pos-- {
			// (sorted list) element with value higher than (considered element i) are shift right 1 position to make space for insertion later
			ptrs[pos+1] = ptrs[pos]
		}
		ptrs[pos+1] = dataPtr // cur is the right most element (sorted list) that is lower than the (considered element i), so insert at cur+1
	}
}

// isEqual implements the method for comparision used for search algorithm
func (b booking) isEqual(target booking, searchKey bookingKeyType) bool {

	// 	switch searchKey {

	// 	case bookFullyMatched:
	// 		if b == target {
	// 			return true
	// 		}
	// 	case bookName:
	// 		if strings.Compare(b.name, target.name) == 0 { // 1 for greater
	// 			//			fmt.Printf("N %s > %s\n", b.name, data.name)
	// 			return true
	// 		}

	// 	case bookDateAndTimeSlotAndVenueTypeNum:
	// 		if b.date == target.date && b.timeSlot == target.timeSlot && b.venueType == target.venueType {
	// 			fmt.Println("bookDateAndTimeSlotAndVenueTypeNum")
	// 			return true
	// 		}

	// 	case bookDateAndTimeSlot:
	// 		if b.date == target.date && b.timeSlot == target.timeSlot {
	// 			// fmt.Println(b.date)
	// 			// fmt.Println(target.date)
	// 			// fmt.Println(b.timeSlot)
	// 			// fmt.Println(target.timeSlot)
	// 			return true
	// 		}

	// 	case bookDate:
	// 		if b.date == target.date {
	// 			return true
	// 		}

	// 	case bookTimeSlot:
	// 		if b.timeSlot == target.timeSlot {
	// 			return true
	// 		}

	// case bookVenueTypeAndNum:
	// 	if b.venueType == target.venueType && b.venueNum == target.venueNum {
	// 		return true
	// 	}

	// case bookVenueTypeName:
	// 	if strings.Compare(venueTypeList[b.venueType], venueTypeList[target.venueType]) == 0 {
	// 		return true
	// 	}
	// case bookVenueName:
	// 	//	if b.venueNum == target.venueNum {
	// 	if strings.Compare(property[b.venueType][b.venueNum].name, property[target.venueType][target.venueNum].name) == 0 { // for name equality venue type and num must match
	// 		return true
	// 	}
	// case bookDateTimeSlotVenueTypeName:
	// 	if b.date == target.date && b.timeSlot == target.timeSlot &&
	// 		strings.Compare(venueTypeList[b.venueType], venueTypeList[target.venueType]) == 0 {
	// 		return true
	// 	}
	// case bookDateTimeSlotVenueName:
	// 	if b.date == target.date && b.timeSlot == target.timeSlot &&
	// 		strings.Compare(property[b.venueType][b.venueNum].name, property[target.venueType][target.venueNum].name) == 0 {
	// 		return true
	// 	}
	// case bookVenueCapacity:
	// 	if property[b.venueType][b.venueNum].capacity == property[target.venueType][target.venueNum].capacity {
	// 		return true
	// 	}
	// }
	return false
}

// isEqual implements the method for comparision used for sorting algorithm
func (b booking) isLargerThan(data booking, sortKey bookingKeyType) bool {

	// switch sortKey {

	// case bookName:
	// 	if strings.Compare(b.name, data.name) == 1 { // 1 for greater
	// 		//			fmt.Printf("N %s > %s\n", b.name, data.name)
	// 		return true
	// 	}

	// case bookDateAndTimeSlotAndVenueTypeNum: // priority of sort order Year, month, date, time slot
	// 	if b.date.year > data.date.year {
	// 		return true
	// 	}
	// 	if b.date.year == data.date.year && b.date.month > data.date.month {
	// 		return true
	// 	}
	// 	if b.date.year == data.date.year && b.date.month == data.date.month && b.date.day > data.date.day {
	// 		return true
	// 	}
	// 	if b.date.year == data.date.year && b.date.month == data.date.month && b.date.day == data.date.day && b.timeSlot > data.timeSlot {
	// 		return true
	// 	}
	// 	if b.date.year == data.date.year && b.date.month == data.date.month &&
	// 		b.date.day == data.date.day && b.timeSlot == data.timeSlot &&
	// 		b.venueType > data.venueType {
	// 		return true
	// 	}

	// case bookDateAndTimeSlot: // priority of sort order Year, month, date, time slot
	// 	if b.date.year > data.date.year {
	// 		return true
	// 	}
	// 	if b.date.year == data.date.year && b.date.month > data.date.month {
	// 		return true
	// 	}
	// 	if b.date.year == data.date.year && b.date.month == data.date.month && b.date.day > data.date.day {
	// 		return true
	// 	}
	// 	if b.date.year == data.date.year && b.date.month == data.date.month && b.date.day == data.date.day && b.timeSlot > data.timeSlot {
	// 		return true
	// 	}

	// case bookDate: // priority of sort order Year, month, date
	// 	if b.date.year > data.date.year {
	// 		return true
	// 	}
	// 	if b.date.year == data.date.year {
	// 		if b.date.month > data.date.month {
	// 			return true
	// 		}
	// 		if b.date.month == data.date.month {
	// 			if b.date.day > data.date.day {
	// 				return true
	// 			}
	// 		}
	// 	}

	// case bookTimeSlot:
	// 	if b.timeSlot > data.timeSlot {
	// 		return true
	// 	}

	// // case bookVenueTypeAndNum: // priority of sort order Venue Type ,venue num
	// // 	if b.venueType > data.venueType {
	// // 		return true
	// // 	}
	// // 	if b.venueType == data.venueType && b.venueNum > data.venueNum {
	// // 		return true
	// // 	}

	// case bookVenueTypeName:
	// 	if strings.Compare(venueTypeList[b.venueType], venueTypeList[data.venueType]) == 1 {
	// 		return true
	// 	}

	// case bookVenueName:
	// 	if strings.Compare(property[b.venueType][b.venueNum].name, property[data.venueType][data.venueNum].name) == 1 {
	// 		return true
	// 	}
	// case bookDateTimeSlotVenueTypeName: // priority of sort Date, TimeSlot, Venue Type

	// 	if b.date.year > data.date.year {
	// 		return true
	// 	}
	// 	if b.date.year == data.date.year {
	// 		if b.date.month > data.date.month {
	// 			return true
	// 		}
	// 		if b.date.month == data.date.month {
	// 			if b.date.day > data.date.day {
	// 				return true
	// 			}
	// 			if b.date.day == data.date.day {
	// 				if strings.Compare(venueTypeList[b.venueType], venueTypeList[data.venueType]) == 1 { // 1 means greater in alphabetical order
	// 					return true
	// 				}
	// 			}
	// 		}
	// 	}

	// case bookDateTimeSlotVenueName: // priority of sort Date, TimeSlot, Venue Type, Venue Num
	// 	if b.date.year > data.date.year {
	// 		return true
	// 	}
	// 	if b.date.year == data.date.year {
	// 		if b.date.month > data.date.month {
	// 			return true
	// 		}
	// 		if b.date.month == data.date.month {
	// 			if b.date.day > data.date.day {
	// 				return true
	// 			}
	// 			if b.date.day == data.date.day {
	// 				if strings.Compare(venueTypeList[b.venueType], venueTypeList[data.venueType]) == 1 {
	// 					return true
	// 				}
	// 				if strings.Compare(venueTypeList[b.venueType], venueTypeList[data.venueType]) == 0 {
	// 					if strings.Compare(property[b.venueType][b.venueNum].name, property[data.venueType][data.venueNum].name) == 1 {
	// 						return true
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// case bookVenueCapacity:
	// 	if property[b.venueType][b.venueNum].capacity > property[data.venueType][data.venueNum].capacity {
	// 		return true
	// 	}
	// }
	return false
}

// bookingFindMoreMatch find more match of a sorted array from the staring index
// Operate on sorted array *[] booking to find N items of the same target
// return count, err
// if no err !- nil, count is the number of match to the target

func bookingFindMoreMatch(b *[]booking, index int, target booking, key bookingKeyType) (int, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Trapped panic: %s (%T)\n", err, err)
		}
	}()

	var err error

	if index < 0 {
		err = IndexOutOfRange
		return -1, err
	}
	count := 0
	for i := 0; i < len(*b); i++ {
		if (*b)[i].isEqual(target, key) {
			count++
		}
	}
	return count, nil
}

// bookingFindMoreMatch2 find more match of a sorted array from the staring index
// Operate on sorted array *[]* booking to find N items of the same target
// return count, err
// if no err !- nil, count is the number of match to the target

func bookingFindMoreMatch2(b *[]*booking, index int, target booking, key bookingKeyType) (int, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Trapped panic: %s (%T)\n", err, err)
		}
	}()

	var err error

	if index < 0 {
		err = IndexOutOfRange
		return -1, err
	}
	count := 0

	for i := 0; i < len(*b); i++ {
		if (*(*b)[i]).isEqual(target, key) {
			count++
		}
	}
	return count, nil
}

// bookingFindFirstMatch find the first index match of a sorted array from the index (between first and last)
// Operate on sorted array *[] booking to find the start of index that match the target of the given index
// return index, err
// if index = -1, err show the error type
// if index >= 0, index is the index of the first match to the target
var SearchElementMismatch = errors.New("Search Element Mismatch")

func bookingFindFirstMatch2(b *[]*booking, index int, target booking, key bookingKeyType) (int, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Trapped panic: %s (%T)\n", err, err)
		}
	}()

	var err error

	if index < 0 {
		err = IndexOutOfRange
		return -1, err
	}

	var i int
	for i = index; i >= 0; {
		if (*(*b)[i]).isEqual(target, key) {
			i-- // search the lower bound
			continue
		}
		if i == index {
			err = SearchElementMismatch
			return -1, err
		}
		return i + 1, nil // exit on first occurance of not equal
	}
	return i + 1, nil
}

// bookingFindFirstMatch find the first index match of a sorted array from the staring index
// Operate on sorted array *[] booking to find the start of index that match the target of the given index
// if index = -1, err show the error type
// if index >= 0, index is the index of the last match to the target
func bookingFindLastMatch2(b *[]*booking, index int, target booking, key bookingKeyType) (int, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Trapped panic: %s (%T)\n", err, err)
		}
	}()

	var err error

	if index < 0 {
		err = IndexOutOfRange
		return -1, err
	}

	var i int
	for i = index; i < len(*b); {
		if (*(*b)[i]).isEqual(target, key) {
			i++ // searching the upper bound
			continue
		}
		if i == index {
			err = SearchElementMismatch
			return -1, err
		}
		return i - 1, nil // exit on first occurance of not equal
	}
	return i - 1, nil
}

var NoMatchFound = errors.New("No Match found in Search")

// iterativeBinSortedFindFirst2 applies to sorted array to use iterative binary search to locate the first element that matches
// return -1, err if element is not found
// return index of first element, nil if index is found
func iterativeBinSortedFindFirst2(b *[]*booking, target booking, key bookingKeyType) (int, error) {

	var err error = nil
	var index int
	index = iterativeBinSearch2(len(*b), b, target, key) // Binary search to find first occurance of target (not necessary the first in a sorted array)
	showIndex(index)
	if index < 0 {
		err = NoMatchFound
		return -1, err
	}
	index, err = bookingFindFirstMatch2(b, index, target, key) // locate first element of the target in the sort array
	if err != nil {
		return -1, err
	}
	showIndex(index)
	return index, err
}

// recursiveBinSortedFindFirst2 applies to sorted array to use recursive binary search to locate the first element that matches
// return -1, err if element is not found
// return index of first element, nil if index is found
func recursiveBinSortedFindFirst2(b *[]*booking, target booking, key bookingKeyType) (int, error) {

	var err error = nil
	var index int
	index = recursiveBin2(len(*b), b, target, key) // Binary search to find first occurance of target (not necessary the first in a sorted array)
	showIndex(index)
	if index < 0 {
		err = NoMatchFound
		return -1, err
	}
	index, err = bookingFindFirstMatch2(b, index, target, key) // locate first element of the target in the sort array
	if err != nil {
		return -1, err
	}
	showIndex(index)
	return index, err
}
