package database

import (
	"testing"
)

func TestEmailValidation(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		isValid bool
	}{
		{"Valid email", "test@example.com", true},
		{"Another valid", "user.name+tag@domain.co.uk", true},
		{"No @ symbol", "invalid-email", false},
		{"No domain", "test@", false},
		{"No username", "@domain.com", false},
		{"Empty", "", false},
		{"With spaces", "test @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Простая проверка наличия @ и точки в домене
			hasAt := false
			hasDotAfterAt := false
			for i, c := range tt.email {
				if c == '@' {
					hasAt = true
					// Проверяем есть ли точка после @
					for j := i + 1; j < len(tt.email); j++ {
						if tt.email[j] == '.' {
							hasDotAfterAt = true
							break
						}
					}
				}
			}

			isValid := hasAt && hasDotAfterAt && len(tt.email) > 3

			if isValid != tt.isValid {
				t.Errorf("Email validation failed for %s: expected %v, got %v",
					tt.email, tt.isValid, isValid)
			}
		})
	}
}

func TestPasswordValidation(t *testing.T) {
	tests := []struct {
		name     string
		password string
		isValid  bool
	}{
		{"Valid password", "password123", true},
		{"Short password", "123", false},
		{"Empty password", "", false},
		{"Long password", "verylongpasswordthatshouldbevalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := len(tt.password) >= 6
			if isValid != tt.isValid {
				t.Errorf("Password validation failed: expected %v, got %v",
					tt.isValid, isValid)
			}
		})
	}
}

func TestTokenGeneration(t *testing.T) {
	// Тест генерации токена
	token1 := generateToken()
	token2 := generateToken()

	if token1 == token2 {
		t.Error("Generated tokens should be unique")
	}

	if len(token1) != 64 { // 32 байта в hex = 64 символа
		t.Errorf("Token length should be 64 chars, got %d", len(token1))
	}
}

func TestCartQuantityValidation(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		isValid  bool
	}{
		{"Valid quantity", 1, true},
		{"Valid quantity 5", 5, true},
		{"Zero quantity", 0, false},
		{"Negative quantity", -1, false},
		{"Large quantity", 999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.quantity > 0
			if isValid != tt.isValid {
				t.Errorf("Quantity validation failed: expected %v, got %v",
					tt.isValid, isValid)
			}
		})
	}
}
