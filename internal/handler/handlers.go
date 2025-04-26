package handler

import (
	"context"
	m "effective-mobile-task/internal/model"
	s "effective-mobile-task/internal/service"
	psql "effective-mobile-task/internal/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	DB *psql.DB
}

func CreateUserHandler(db *psql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// POST
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req m.PersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Debug: Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: name and surname are required"})
		return
	}

	log.Printf("Debug: Received POST with name=%s, surname=%s, patronymic=%s", req.Name, req.Surname, req.Patronymic)

	// Данные от 3 API
	age, gender, nationality, err := s.APIRespData(req.Name)
	if err != nil {
		log.Printf("Debug: Failed to enrich data: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch data from external APIs"})
		return
	}

	// Подготовка записи для БД
	person := m.User{
		Name:        req.Name,
		Surname:     req.Surname,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}
	if req.Patronymic != "" {
		person.Patronymic = &req.Patronymic
	}

	// Сохранение в БД
	var personID int
	err = h.DB.Psql.QueryRow(context.Background(), `
		INSERT INTO Users (name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality,
	).Scan(&personID)
	if err != nil {
		log.Printf("Debug: Failed to save to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database"})
		return
	}

	log.Printf("Info: Saved person with id=%d", personID)
	person.ID = personID

	c.JSON(http.StatusCreated, person)
}

// GET
func (h *UserHandler) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Status":  "Successfuel",
		"Message": "pong",
	})
}

// PUT
func (h *UserHandler) PutUser(c *gin.Context) {

}

// DELETE
func (h *UserHandler) DeleteUser(c *gin.Context) {

}
