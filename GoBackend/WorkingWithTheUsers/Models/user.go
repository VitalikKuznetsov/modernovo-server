package models

type User struct {
	Email       string
	Password    string
	Name        string
	PhoneNumber string
}

type UserRegOrLog struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InfoForUser struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phonenumber"`
}

type Session struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type Favorite struct {
	ID        int    `json:"id"`
	UserEmail string `json:"user_email"`
	ProductID int    `json:"product_id"`
}

type FavoriteRequest struct {
	ProductID int `json:"product_id"`
}

type CartItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	ImageURL  string  `json:"image_url"`
}

type CartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CartResponse struct {
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	ItemCount int        `json:"item_count"`
}

type AdminProduct struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	ImageURL    string   `json:"image_url"`
	ImageURLs   []string `json:"image_urls"`
}

type AdminProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	ImageURL    string   `json:"image_url"`
	ImageURLs   []string `json:"image_urls"`
}
