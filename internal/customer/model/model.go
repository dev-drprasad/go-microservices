package model

import "time"

type Customer struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string `json:"name"`
	Address     string `json:"address"`
	Zipcode     string `json:"zipcode"`
	PhoneNumber string `json:"phoneNumber"`
}
