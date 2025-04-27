package handler

import (
	"context"
	m "effective-mobile-task/internal/model"
	s "effective-mobile-task/internal/service"
	psql "effective-mobile-task/internal/storage"
	"effective-mobile-task/pkg/logger"
	"fmt"
	"net/http"
	"strings"
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
		logger.LogError("Bad Request", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
			"Status":  "Error",
		})
		return
	}

	logger.Log.Printf("Debug: Received POST with name=%s, surname=%s, patronymic=%s", req.Name, req.Surname, req.Patronymic)

	// Данные от 3 API
	age, gender, nationality, err := s.APIRespData(req.Name)
	if err != nil {
		logger.LogError("Failed to fetch data: %v", err)
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

	logger.LogInfo(fmt.Sprintf("User created with id=%d", personID))
	person.ID = personID

	c.JSON(http.StatusCreated, person)
}

// GET
func (h *UserHandler) GetUsers(c *gin.Context) {
	sql := "SELECT id, name, surname, patronymic, age, gender, nationality FROM users"

	rows, err := h.DB.Psql.Query(context.Background(), sql)
	if err != nil {
		logger.LogError("Database connect fals", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Database connect fals",
			"Status":  "Error",
		})
		return
	}
	defer rows.Close()

	var users []m.User
	for rows.Next() {
		var u m.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.Nationality); err == nil {
			users = append(users, u)
		}
		return
	}

	logger.LogInfo("Users get successfully")
	c.JSON(http.StatusOK, users)
}

// PUT
func (h *UserHandler) PutUser(c *gin.Context) {
	id := c.Param("id")

	var upUser map[string]interface{}
	if err := c.ShouldBind(&upUser); err != nil {
		logger.LogError("Bad request body updated", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad request body",
			"Status":  "Error",
		})
		return
	}

	parametrs := []string{}
	data := []interface{}{}
	i := 1

	for key, value := range upUser {
		parametrs = append(parametrs, fmt.Sprintf("%s=%d", key, i))
		data = append(data, value)
	}

	parametrs = append(parametrs, fmt.Sprintf("updated_at=$%d", i))
	data = append(data, time.Now())
	i++
	data = append(data, id)

	sql := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d", strings.Join(parametrs, ", "), i)

	_, err := h.DB.Psql.Exec(context.Background(), sql, data...)
	if err != nil {
		logger.LogError("Update error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Update failed",
			"Status":  "Error",
		})
		return
	}

	logger.LogInfo(fmt.Sprintf("User with id=%s updated successfully", id))
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
		logger.LogError("Failed to delete user", err)
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "User not found",
			"Status":  "Error",
		})
		return
	}

	logger.LogInfo(fmt.Sprintf("User with id=%s deleted successfully", id))
	c.JSON(http.StatusOK, gin.H{
		"Message": "User deleted",
		"Status":  "Successful",
	})
}
