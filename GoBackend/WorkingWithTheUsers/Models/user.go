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
