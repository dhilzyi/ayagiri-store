package api

type Order struct {
	Product []Product
	Table   Table
}

type PayloadProductParams struct {
	Name       string
	RomajiName string `json:"romaji_name"`
	Price      uint32
	Discount   int32
	CategoryID int32 `json:"category_id"`
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
