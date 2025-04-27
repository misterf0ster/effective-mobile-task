package handler

import (
	"context"
	m "effective-mobile-task/internal/model"
	s "effective-mobile-task/internal/service"
	psql "effective-mobile-task/internal/storage"
	"effective-mobile-task/pkg/logger"
	"net/http"
	"time"

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
		logger.Log.Printf("Debug: Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid input: name and surname are required",
			"Status":  "Error",
		})
		return
	}

	logger.Log.Printf("Debug: Received POST with name=%s, surname=%s, patronymic=%s", req.Name, req.Surname, req.Patronymic)

	// Данные от 3 API
	age, gender, nationality, err := s.APIRespData(req.Name)
	if err != nil {
		logger.Log.Printf("Debug: Failed to enrich data: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"Message": "Failed to fetch data from external APIs",
			"Status":  "Error",
		})
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
		logger.Log.Printf("Debug: Failed to save to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Failed to save to database",
			"Status":  "Error",
		})
		return
	}

	logger.Log.Printf("Info: Saved person with id=%d", personID)
	person.ID = personID

	c.JSON(http.StatusCreated, person)
}

// GET
func (h *UserHandler) GetUsers(c *gin.Context) {
	sql := "SELECT id, name, surname, patronymic, age, gender, nationality FROM users"

	rows, err := h.DB.Psql.Query(context.Background(), sql)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Database connect fals",
			"Status":  "Error",
		})
	}
	defer rows.Close()

	var users []m.User
	for rows.Next() {
		var u m.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.Nationality); err == nil {
			users = append(users, u)
		}
	}

	c.JSON(http.StatusOK, users)
}

// PUT
func (h *UserHandler) PutUser(c *gin.Context) {
	id := c.Param("id")

	sql := "UPDATE id, name, surname, patronymic, age, gender, nationality, updated_at FROM users WHERE id=$1"

	var up m.User
	if err := c.ShouldBind(&up); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad request body",
			"Status":  "Error",
		})
	}

	up.UpdatedAt = time.Now()

	_, err := h.DB.Psql.Exec(context.Background(), sql,
		up.Name, up.Surname, up.Patronymic, up.Age, up.Gender, up.Nationality,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Update failed",
			"Status":  "Error",
		})
	}

	c.JSON(http.StatusAccepted, gin.H{
		"Message": "Data was updated",
		"Status":  "Successfule",
	})
}

// DELETE
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	sql := "DELETE FROM users WHERE id=$1 RETURNING id"

	var del_id int
	err := h.DB.Psql.QueryRow(context.Background(), sql, id).Scan(&del_id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "User not found",
			"Status":  "Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "User deleted",
		"Status":  "Successful",
	})
}
