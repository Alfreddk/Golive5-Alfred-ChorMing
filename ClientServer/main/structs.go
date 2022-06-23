package main

// alfred 23.06.2022: Item struct needs to be the same struct that the backend server uses.
type Item struct {
	ID             string `json:"ID"`
	Name           string `json:"Name"`
	Description    string `json:"Description"`
	HideGiven      int    `json:"HideGiven"`
	HideGotten     int    `json:"HideGotten"`
	HideWithdrawn  int    `json:"HideWithdrawn"`
	GiverUsername  string `json:"GiverUsername"`
	GetterUsername string `json:"GetterUsername"`
	State          int    `json:"State"`
	Date           string `json:"Date"`
}

// alfred 23.06.2022: User struct needs to be the same struct that the backend server uses.
type User struct {
	ID        string `json:"ID"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	Name      string `json:"Name"`
	Address   string `json:"Address"`
	Postal    string `json:"Postal"`
	Telephone string `json:"Telephone"`
	LastLogin string `json:"LastLogin"`
}
