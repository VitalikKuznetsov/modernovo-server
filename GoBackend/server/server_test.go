package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCorsMiddleware(t *testing.T) {
	// Создаем тестовый хэндлер
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Оборачиваем в CORS middleware
	handler := corsMiddleware(testHandler)

	// Тест OPTIONS запроса
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("OPTIONS request should return 204, got %d", w.Code)
	}

	// Проверяем CORS заголовки
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization",
	}

	for header, expectedValue := range expectedHeaders {
		actualValue := w.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("Header %s: expected %s, got %s", header, expectedValue, actualValue)
		}
	}

	// Тест обычного GET запроса
	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET request should return 200, got %d", w.Code)
	}
}

func TestRouteSetup(t *testing.T) {
	// Проверяем, что все основные маршруты определены
	router := mux.NewRouter()
	staticPath := "./test_static"

	SetupRoutes(router, staticPath)

	// Проверяем несколько ключевых маршрутов
	routesToCheck := []struct {
		method string
		path   string
	}{
		{"GET", "/api/weather"},
		{"POST", "/api/register"},
		{"POST", "/api/login"},
		{"GET", "/api/products"},
		{"GET", "/api/products/1"},
	}

	for _, route := range routesToCheck {
		req := httptest.NewRequest(route.method, route.path, nil)
		var match mux.RouteMatch
		if !router.Match(req, &match) {
			t.Errorf("Route %s %s not found", route.method, route.path)
		}
	}
}
