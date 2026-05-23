package api

type ProductRequest struct {
	Name        string
	EnglishName string `json:"english_name"`
	Price       int32
	Discount    int32
	CategoryID  int32 `json:"category_id"`
}

type CategoryRequest struct {
	Name        string
	EnglishName string `json:"english_name"`
}

type CategoryResponse struct {
	ID          int `json:"id"`
	Name        string
	EnglishName string `json:"english_name"`
}
