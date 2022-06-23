package main

// alfred 23.06.2022: Item struct needs to be the same struct that the backend server uses.
type Item struct {
	ID            string `json:"ID"`
	Name          string `json:"Name"`
	Description   string `json:"Description"`
	HideGiven     int    `json:"HideGiven"`
	HideGotten    int    `json:"HideGotten"`
	HideWithdrawn int    `json:"HideWithdrawn"`
	GiverID       string `json:"GiverID"`
	GetterID      string `json:"GetterID"`
	State         int    `json:"State"`
	Date          string `json:"Date"`
}
