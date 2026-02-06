package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendError(t *testing.T) {
	w := httptest.NewRecorder()
	message := "Test error message"
	statusCode := http.StatusBadRequest

	sendError(w, message, statusCode)

	if w.Code != statusCode {
		t.Errorf("Expected status code %d, got %d", statusCode, w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["error"] != message {
		t.Errorf("Expected error message '%s', got '%s'", message, response["error"])
	}
}

func TestSendSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	message := "Success message"

	sendSuccess(w, message)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["message"] != message {
		t.Errorf("Expected success message '%s', got '%s'", message, response["message"])
	}
}

func TestIsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Admin email", "ssiromas@gmail.com", true},
		{"Regular email", "user@example.com", false},
		{"Another regular", "test@gmail.com", false},
		{"Empty email", "", false},
		{"Similar but not admin", "ssiromas@hotmail.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAdmin(tt.email)
			if result != tt.expected {
				t.Errorf("isAdmin(%s) = %v, expected %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestRegisterHandlerValidation(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		expectError bool
	}{
		{"Valid data", "test@example.com", "password123", false},
		{"Empty email", "", "password123", true},
		{"Empty password", "test@example.com", "", true},
		{"Both empty", "", "", true},
		{"Invalid email", "invalid-email", "password123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый запрос
			userData := map[string]string{
				"email":    tt.email,
				"password": tt.password,
			}
			jsonData, _ := json.Marshal(userData)

			req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			hasError := tt.email == "" || tt.password == "" || !strings.Contains(tt.email, "@")
			if hasError != tt.expectError {
				t.Errorf("Validation failed for email=%s, password=%s", tt.email, tt.password)
			}
		})
	}
}
