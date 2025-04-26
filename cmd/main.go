package main

import (
	"effective-mobile-task/internal/handler"
	psql "effective-mobile-task/internal/storage"
	"effective-mobile-task/pkg/config"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()
	cfg := config.Config()

	url := cfg.DBaseURL()
	conn, err := psql.Open(url)
	if err != nil {
		log.Fatalf("Unable to connect to db: %v\n", err)
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
		log.Printf("Port not found")
	}

	log.Println("Starting server on port", port)
	if err := g.Run(":" + port); err != nil {
		log.Fatal("Server startup error: " + err.Error())
	}
}
