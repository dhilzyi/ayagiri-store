package domain

type ProductResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
	Price       int32  `json:"price"`
	CategoryID  int32  `json:"category_id"`
	Discount    int32  `json:"discount"`
}
