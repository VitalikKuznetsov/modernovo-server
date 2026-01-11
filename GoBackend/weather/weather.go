package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetWeatherData(city string) (map[string]interface{}, error) {
	apiKey := "a635da31f15f34aa8aaf69a39a3b6570"

	// Используйте API 2.5, который работает с вашим ключом
	url := fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ru",
		city,
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Проверяем на ошибки API
	var apiError struct {
		Cod     int    `json:"cod"`
		Message string `json:"message"`
	}

	// Сначала пробуем распарсить как ошибку
	if err = json.Unmarshal(body, &apiError); err == nil && apiError.Cod != 200 {
		return nil, fmt.Errorf("API error %d: %s", apiError.Cod, apiError.Message)
	}

	// Если не ошибка, парсим как погоду
	var weatherData struct {
		Name string `json:"name"`
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Cod int `json:"cod"`
	}

	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, err
	}

	if weatherData.Cod != 200 {
		return nil, fmt.Errorf("API error %d", weatherData.Cod)
	}

	result := map[string]interface{}{
		"city":        weatherData.Name,
		"temperature": weatherData.Main.Temp,
	}

	return result, nil
}
