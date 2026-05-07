package main

import (
	"fmt"
	"log"
	"rdev-go-api/internal/config"
	"rdev-go-api/internal/data"
	"rdev-go-api/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Println(fmt.Errorf("failed to load config: %w", err))
		return
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	db := data.ConnectDB(cfg.Database.URL, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns)
	server.RegisterRouters(router, db, cfg)

	router.Run(cfg.Server.Port)
}
