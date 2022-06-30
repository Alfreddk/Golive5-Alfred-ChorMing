package main

// Item defines the details of items.
// It also serves for JSON marshalling and unmarshalling.
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

// Item defines the details of items.
// It also serves for JSON marshalling and unmarshalling.
type User struct {
	ID        string `json:"ID"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	Name      string `json:"Name"`
	Address   string `json:"Address"`
	Postal    string `json:"Postal"`
	Telephone string `json:"Telephone"`
	Role      string `json:"Role"`
	LastLogin string `json:"LastLogin"`
}

// ToGive - item given but not received yet
// Given - items given and has a receiver
// Gotten - items given but not received yet  (Not used at the moment)
// Withdrawn - items given, not received yet and withdrawn
const (
	stateToGive = iota
	stateGiven
	stateWithdrawn
	stateInvalid
)

type state map[int]string

var itemState = state{
	stateToGive:    "ToGive",
	stateGiven:     "Given",
	stateWithdrawn: "Withdrawn",
	stateInvalid:   "Invalid",
}
