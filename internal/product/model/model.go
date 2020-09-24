package model

type Category struct {
	ID uint `json:"id"`

	Name string `json:"name"`
}

type Brand struct {
	ID uint `json:"id"`

	Name string `json:"name"`
}

type Product struct {
	ID uint `json:"id"`

	Name       string   `json:"name"`
	Cost       float64  `json:"cost"`
	SellPrice  float64  `json:"sellPrice"`
	BrandID    uint     `json:"brandId"`
	CategoryID uint     `json:"categoryId"`
	ImageURLs  []string `json:"imageUrls"`
	Stock      uint     `json:"stock"`
}
