package model

//Входной JSON
type PersonRequest struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
}

//Ответ
type APIResponse struct {
	Age     *int    `json:"age"`
	Gender  *string `json:"gender"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
	Count int `json:"count"`
}
