package model

type Branch struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	// Organization *Organization `json:"organization"`
}

type Organization struct {
	Name     string    `json:"name"`
	Branches []*Branch `json:"branches"`
}
