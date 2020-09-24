package model

type Branch struct {
	ID uint `json:"id"`

	Name           string `json:"name"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phoneNumber"`
	OrganizationID uint   `json:"organizationId"`

	Organization *Organization `json:"organization"`
}

type Organization struct {
	ID uint `json:"id"`

	Name     string    `json:"name"`
	Branches []*Branch `json:"branches"`
}
