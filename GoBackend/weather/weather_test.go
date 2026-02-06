package weather

import (
	"testing"
)

func TestGetWeatherData(t *testing.T) {
	// Тест с валидным городом
	weatherData, err := GetWeatherData("Moscow")
	if err != nil {
		t.Errorf("Failed to get weather data for Moscow: %v", err)
	}

	// Проверяем наличие ожидаемых полей
	if city, ok := weatherData["city"]; !ok || city == "" {
		t.Error("City field is missing or empty")
	}

	if temp, ok := weatherData["temperature"]; !ok {
		t.Error("Temperature field is missing")
	} else {
		// Температура должна быть в разумных пределах
		if tempFloat, ok := temp.(float64); ok {
			if tempFloat < -50 || tempFloat > 50 {
				t.Errorf("Temperature %f seems out of reasonable range", tempFloat)
			}
		}
	}
}

func TestGetWeatherDataInvalidCity(t *testing.T) {
	// Тест с несуществующим городом
	_, err := GetWeatherData("NonexistentCity12345")
	if err == nil {
		t.Error("Expected error for nonexistent city, got nil")
	}
}

func TestGetWeatherDataEmptyCity(t *testing.T) {
	// Тест с пустым городом
	_, err := GetWeatherData("")
	if err == nil {
		t.Error("Expected error for empty city, got nil")
	}
}
