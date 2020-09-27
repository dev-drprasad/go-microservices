package model

import (
	customermodel "gomicroservices/internal/customer/model"
	"time"
)

type Status string

const (
	StatusNew         Status = "new"
	StatusPreparation Status = "preparation"
	StatusReady       Status = "ready"
	StatusDelivered   Status = "delivered"
)

type Order struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Status     Status `json:"status"`
	CustomerID uint   `json:"customerId"`

	Customer *customermodel.Customer `json:"customer"`
	Products []*OrderProduct         `json:"products"`
}

type OrderProduct struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	OrderID   uint    `json:"orderId"`
	ProductID uint    `json:"productId"`
	UnitPrice float64 `json:"unitPrice"`
	Quantity  uint    `json:"quantity"`
}
