package models

type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	ImageURL    string   `json:"ImageUrl"`
	ImageURLs   []string `json:"image_urls"`
}

type ProductList struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
}

type ProductDetail struct {
	Product
	AdditionalInfo string `json:"additional_info,omitempty"`
}
