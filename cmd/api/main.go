package main

import (
	"flag"
	"fmt"
	"log"
	"rdev-go-api/internal/config"
	"rdev-go-api/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	arg := flag.String("env", "local", "Config environment [local|docker]")
	flag.Parse()

	file := ""
	switch *arg {
	case "local":
		file = "config.yaml"
	case "docker":
		file = "config.docker.host.yaml"
	}

	cfg, err := config.LoadConfig(file)
	if err != nil {
		log.Println(fmt.Errorf("failed to load config: %w", err))
		return
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	db := config.ConnectDB(cfg.Database)
	redis := config.ConnectRedis(cfg.Redis)
	server.Container(router, db, redis, cfg)

	router.Run(cfg.Server.Port)
}
