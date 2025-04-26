package model

type User struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Gender      string  `json:"gender"`
	Age         int     `json:"age"`
	Nationality string  `json:"nationality"`
}
