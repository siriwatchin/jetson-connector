package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/siriwatchin/jetson-connector/handler"
	"github.com/siriwatchin/jetson-connector/middleware"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger())

	rawData := &handler.RawDataHandler{}
	v1 := r.Group("/api/v1")
	{
		v1.POST("/raw_data", rawData.Create)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server: %v", err)
	}
}
