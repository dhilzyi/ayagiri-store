package domain

import "time"

type ProductResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
	Price       int32  `json:"price"`
	CategoryID  int32  `json:"category_id"`
	Discount    int32  `json:"discount"`
}

type ProductResponseAdmin struct {
	ProductResponse
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CategoryName string    `json:"category_name"`
}
