// Lab 3 This is the server implementation for REST API
package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

/*
type Course struct {
	// map this type to the record in the table
	ID    string
	Title string
}
*/

var errEmptyRow = errors.New("sql: Empty Row")

func sqlGetAllItems(db *sql.DB) []Item {
	rows, err := db.Query("Select * FROM Items")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.Name, &item.Description, &item.HideGiven, &item.HideGotten, &item.HideWithdrawn, &item.GiverID, &item.GetterID, &item.State, &item.Date)
		if err != nil {
			panic(err.Error())
		}
		items = append(items, item)
	}
	//fmt.Println(items)
	// .. to log successful call.

	return items
}

//func AddNewItem(db *sql.DB, name string, desc string, hideGiven int, hideGotten int, hideWithdrawn int, giverID string, getterID string, state int, date string) {
func sqlAddNewItem(db *sql.DB, item Item) {
	query := fmt.Sprintf("INSERT INTO Items (Name, Description, HideGiven, HideGotten, HideWithdrawn, GiverID, GetterID, State, Date) VALUES ('%s', '%s', %v, %v, %v, %v, %v, %v, '%s')",
		item.Name, item.Description, item.HideGiven, item.HideGotten, item.HideWithdrawn, item.GiverID, item.GetterID, item.State, item.Date)

	row, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer row.Close()

	//fmt.Println("Insert Successful")
	// .. to log successful call.
}

func sqlEditItem(db *sql.DB, item Item) {
	query := fmt.Sprintf("UPDATE Items SET HideGiven = %v, HideGotten = %v, HideWithdrawn = %v, GetterID = '%s', State = %v WHERE ID = '%s')",
		item.HideGiven, item.HideGotten, item.HideWithdrawn, item.GetterID, item.State)

	row, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer row.Close()

	//fmt.Println("Edit Successful")
	// .. to log successful call.
}

/*
// GetRecords gets all the rows of the current table and return as a slice of map
func GetRecords(db *sql.DB) map[string]interface{} {
	// query to get all records of table (persons) of my_db
	rows, err := db.Query("Select * FROM Courses")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	//var course map[string]string
	var courseMap = make(map[string]interface{})

	for rows.Next() {
		// map this type to the record in the table
		var course Course
		err = rows.Scan(&course.ID, &course.Title)
		if err != nil {
			panic(err.Error())
		}
		courseMap[course.ID] = course.Title
	}
	fmt.Println(courseMap)
	return courseMap

}

// GetOneRecord checks if there is a existence of a record based on the ID primary key
// If there is a record, it returns a map of the record key:title pair
// error = nil, there is a record
// error = emptyRow, there is no record
func GetOneRecord(db *sql.DB, id string) (map[string]interface{}, error) {

	// This should not be done this way to avaoid sql injection risk
	// see https://go.dev/doc/database/sql-injection
	//	query := fmt.Sprintf("Select * FROM Courses where ID='%s'", id)

	row, err := db.Query("Select * FROM Courses where ID=?", id)
	if err != nil {
		panic(err.Error())
	}
	defer row.Close()

	//var course map[string]string
	courseMap := make(map[string]interface{})

	if row.Next() {
		var course Course
		err = row.Scan(&course.ID, &course.Title)
		if err != nil {
			panic(err.Error())
		}
		courseMap[course.ID] = course.Title
		fmt.Println("Course:", courseMap)
		return courseMap, nil
	} else {
		return courseMap, errEmptyRow
	}
}

// DeleteRecord deletes a record from the current table using the ID primary key
func DeleteRecord(db *sql.DB, ID string) {
	// create the sql query to delete with primary key
	// Note deleting a non-existent record is considered as deleted, so will always passed
	//query := fmt.Sprintf("DELETE FROM Courses WHERE ID='%s'", ID)
	//row, err := db.Query(query)
	row, err := db.Query("DELETE FROM Courses WHERE ID=?", ID)

	if err != nil {
		panic(err.Error())
	}
	defer row.Close()

	fmt.Println("Delete Successful")
}

// EditRecord edits the record of the current table based on the primary key ID with title
func EditRecord(db *sql.DB, ID string, title string) {
	// create the sql query to update record
	// query := fmt.Sprintf("UPDATE Courses SET Title='%s' WHERE ID='%s'", title, ID)
	// row, err := db.Query(query)
	row, err := db.Query("UPDATE Courses SET Title=? WHERE ID=?", title, ID)
	if err != nil {
		panic(err.Error())
	}
	defer row.Close()
	fmt.Println("Edit Successful")
}

// InsertRecord instead a row record into the current table based on the primary key and title
func InsertRecord(db *sql.DB, ID string, title string) {
	// create the sql query to insert record
	// note the quote for string
	// query := fmt.Sprintf("INSERT INTO Courses VALUES ('%s', '%s')", ID, title)
	// _, err := db.Query(query)
	row, err := db.Query("INSERT INTO Courses VALUES (?, ?)", ID, title)
	if err != nil {
		panic(err.Error())
	}
	defer row.Close()
	fmt.Println("Insert Successful")
}
*/
