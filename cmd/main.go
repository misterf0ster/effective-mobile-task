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
	cfg := config.Config()

	url := cfg.DBaseURL()
	conn, err := psql.Open(url)
	if err != nil {
		log.Fatalf("Unable to connect to db: %v\n", err)
	}
	defer conn.Close()

	g := gin.Default()

	m := g.Group("/users")
	{
		m.POST("/", handler.PostUser)
		m.GET("/", handler.GetUsers)
		m.PUT("/:id", handler.PutUser)
		m.DELETE("/:id", handler.DeleteUser)
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
