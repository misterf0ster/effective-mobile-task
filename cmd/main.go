package main

import (
	"effective-mobile-task/internal/handler"
	psql "effective-mobile-task/internal/storage"
	"effective-mobile-task/pkg/config"
	"effective-mobile-task/pkg/logger"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	logger.InitLogger()

	logger.Log.Info("Starting the service...")

	config.LoadEnv()
	cfg := config.Config()

	url := cfg.DBaseURL()
	conn, err := psql.Open(url)
	if err != nil {
		logger.Log.Fatalf("Unable to connect to db: %v\n", err)
	}
	defer conn.Close()

	h := handler.CreateUserHandler(conn)
	g := gin.Default()

	{
		m := g.Group("/users")
		m.POST("/", h.CreateUser)
		m.GET("/", h.GetUsers)
		m.PUT("/:id", h.PutUser)
		m.DELETE("/:id", h.DeleteUser)
	}

	port := os.Getenv("PORT")
	if port == "" {
		logger.Log.Printf("Port not found")
	}

	logger.Log.Println("Starting server on port", port)
	if err := g.Run(":" + port); err != nil {
		logger.Log.Fatal("Server startup error: " + err.Error())
	}
}
