package api

type Order struct {
	Product []Product
	Table   Table
}

type PayloadProductParams struct {
	Name     string
	Price    uint32
	Discount int32
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
