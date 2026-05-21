package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/siriwatchin/jetson-connector/config"
	appdb "github.com/siriwatchin/jetson-connector/db"
	"github.com/siriwatchin/jetson-connector/handler"
	"github.com/siriwatchin/jetson-connector/middleware"
)

func main() {
	cfg := config.Load()

	db, err := appdb.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger())

	rawData := &handler.RawDataHandler{DB: db, EnableWrite: cfg.EnableWrite}
	v1 := r.Group("/api/v1")
	{
		v1.POST("/raw_data", rawData.Create)
		v1.POST("/raw_data/batch", rawData.CreateBatch)
	}

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}
