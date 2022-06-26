//
package server

import (
	"database/sql"

	"BackEndServer/logger"

	_ "github.com/go-sql-driver/mysql"
)

// sqlGetAllItems execute a query to database to retrieve all items.
// It returns all items as []Item data type.
func sqlGetAllItems(db *sql.DB) []Item {

	rows, err := db.Query("Select * FROM Items")
	if err != nil {
		logger.Trace.Fatalln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.Name, &item.Description, &item.HideGiven, &item.HideGotten, &item.HideWithdrawn, &item.GiverUsername, &item.GetterUsername, &item.State, &item.Date)
		if err != nil {
			logger.Trace.Fatalln(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		}
		items = append(items, item)
	}

	logger.Info.Println("All items successfully retreived from Database.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

	return items
}

// sqlAddNewItem execute a query to database to add a new item.
func sqlAddNewItem(db *sql.DB, item Item) {

	row, err := db.Query("INSERT INTO Items (Name, Description, HideGiven, HideGotten, HideWithdrawn, GiverUsername, GetterUsername, State, Date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		item.Name, item.Description, item.HideGiven, item.HideGotten, item.HideWithdrawn, item.GiverUsername, item.GetterUsername, item.State, item.Date)
	if err != nil {
		logger.Trace.Fatalln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer row.Close()

	logger.Info.Println("Item successfully inserted into Database.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)
}

func sqlEditItem(db *sql.DB, item Item) {

	row, err := db.Query("UPDATE Items SET HideGiven = ?, HideGotten = ?, HideWithdrawn = ?, GetterUsername = ?, State = ? WHERE ID = ?",
		item.HideGiven, item.HideGotten, item.HideWithdrawn, item.GetterUsername, item.State, item.ID)
	if err != nil {
		logger.Trace.Fatalln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer row.Close()

	logger.Info.Println("Item successfully updated in Database.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)
}

func sqlGetAllUsers(db *sql.DB) []User {
	rows, err := db.Query("Select * FROM Users")
	if err != nil {
		logger.Trace.Fatalln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.Address, &user.Postal, &user.Telephone, &user.Role, &user.LastLogin)
		if err != nil {
			logger.Trace.Fatalln(err)
			logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
		}
		users = append(users, user)
	}

	logger.Info.Println("All users successfully retreived from Database.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)

	return users
}

func sqlAddNewUser(db *sql.DB, user User) {

	row, err := db.Query("INSERT INTO Users (Username, Password, Name, Address, Postal, Telephone, Role, LastLogin) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.Username, user.Password, user.Name, user.Address, user.Postal, user.Telephone, user.Role, user.LastLogin)
	if err != nil {
		logger.Trace.Fatalln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer row.Close()

	logger.Info.Println("User successfully inserted into Database.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)
}

func sqlEditUser(db *sql.DB, user User) {

	row, err := db.Query("UPDATE Users SET LastLogin = ? WHERE ID = ?", user.LastLogin, user.ID)
	if err != nil {
		logger.Trace.Fatalln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer row.Close()

	logger.Info.Println("User successfully updated in Database.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)
}

func sqlDeleteUser(db *sql.DB, user User) {

	row, err := db.Query("DELETE FROM Users WHERE ID = ?", user.ID)
	if err != nil {
		logger.Trace.Fatalln(err)
		logger.LogHashing(logger.TraceLogFile, logger.TraceLogHashFile)
	}
	defer row.Close()

	logger.Info.Println("User successfully deleted from Database.")
	logger.LogHashing(logger.InfoLogFile, logger.InfoLogHashFile)
}
