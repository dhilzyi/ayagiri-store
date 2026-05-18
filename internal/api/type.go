package api

type Order struct {
	Product []Product
	Table   Table
}

type ProductRequest struct {
	Name       string
	RomajiName string `json:"romaji_name"`
	Price      int32
	Discount   *int32
	CategoryID *int32 `json:"category_id"`
}

type PayloadCategoryParams struct {
	Name       string
	RomajiName string `json:"romaji_name"`
}

type Product struct {
	ID     int
	Name   string
	Price  int
	Amount int
}

type Table struct {
	ID int
}

type ProductResponse struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	RomajiName string `json:"romaji_name"`
	Price      int32  `json:"price"`
	CategoryID int32  `json:"category_id"`
	Discount   int32  `json:"discount"`
}
