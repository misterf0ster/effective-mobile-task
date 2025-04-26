package model

type AgifyResponse struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type GenderizeResponse struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

type NationalizeResponse struct {
	Name        string `json:"name"`
	Nationality []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"nationality"`
}
