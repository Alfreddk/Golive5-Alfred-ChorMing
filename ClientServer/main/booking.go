package main

// This file contains struct and functions pertaining to bookings
type booking struct {
	name      string
	date      dateType
	timeSlot  int
	venueType int
	venueNum  int
}

// // Use for date
type dateType struct {
	year  int
	month int
	day   int
}

// // Type definitions for Receiver
type monthMapType map[int]string
type timeSlotType map[int]string

// // create maps first and allow access within package
// // Initialise Map that are not to be changed
// // Map the number of days of a month

var daysOfMonth = map[int]int{1: 31, 2: 28, 3: 31, 4: 30, 5: 31, 6: 30, 7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31}
var nameOfMonth = map[int]string{1: "Jan", 2: "Feb", 3: "Mar", 4: "Apr", 5: "May", 6: "Jun", 7: "Jul", 8: "Aug", 9: "Oct", 11: "Nov", 12: "Dec"}

var timeSlot = map[int]string{
	1: "8:00 - 9:00", 2: "9:00 - 10:00", 3: "10:00 - 11:00", 4: "11:00 - 12:00",
	5: "12:00 - 13:00", 6: "13:00 - 14:00", 7: "14:00 - 15:00", 8: "15:00 - 16:00",
	9: "16:00 - 17:00", 10: "17:00 - 18:00", 11: "18:00 - 19:00", 12: "19:00 - 20:00",
	13: "20:00 - 21:00", 14: "21:00 - 22:00", 15: "22:00 - 23:00", 16: "23:00 - 00:00"}

var invalidDate = dateType{0, 0, 0}
