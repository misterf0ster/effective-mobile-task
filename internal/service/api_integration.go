package service

import (
	m "effective-mobile-task/internal/model"
	"effective-mobile-task/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Для горутин
type apiResult struct {
	data m.APIResponse
	err  error
}

func APIRespData(name string) (*int, *string, *string, error) {
	agifyURL := os.Getenv("AGIFY_API_URL") + "?name=" + name
	genderizeURL := os.Getenv("GENDERIZE_API_URL") + "?name=" + name
	nationalizeURL := os.Getenv("NATIONALIZE_API_URL") + "?name=" + name

	results := make(chan apiResult, 3)

	// Функция для выполнения HTTP-запроса
	fetchAPI := func(url string) {
		resp, err := http.Get(url)
		if err != nil {
			results <- apiResult{err: err}
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- apiResult{err: err}
			return
		}

		var data m.APIResponse
		if err := json.Unmarshal(body, &data); err != nil {
			results <- apiResult{err: err}
			return
		}
		results <- apiResult{data: data}
	}

	// Запуск горутин запросов
	go fetchAPI(agifyURL)
	go fetchAPI(genderizeURL)
	go fetchAPI(nationalizeURL)

	// Сбор результатов
	var age *int
	var gender *string
	var nationality *string
	for i := 0; i < 3; i++ {
		result := <-results
		if result.err != nil {
			logger.LogError("API error: %v", result.err)
			continue
		}
		if result.data.Age != nil {
			logger.LogInfo("Fetched age from agify.io")
			age = result.data.Age
		}
		if result.data.Gender != nil {
			logger.LogInfo("Fetched gender from genderize.io")
			gender = result.data.Gender
		}
		if len(result.data.Country) > 0 {
			logger.LogInfo("Fetched nationality from nationalize.io")
			nationality = &result.data.Country[0].CountryID
		}
	}

	if age == nil && gender == nil && nationality == nil {
		logger.LogInfo("no data fetched from APIs")
		return nil, nil, nil, fmt.Errorf("no data fetched from APIs")
	}
	return age, gender, nationality, nil
}
