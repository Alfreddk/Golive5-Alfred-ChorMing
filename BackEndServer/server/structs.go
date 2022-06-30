package server

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

// User defines the details of users.
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
