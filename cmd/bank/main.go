package main

import (
	"context"
	"log"
	"time"

	"bank-service/internal/adapter/db/postgres"
	"bank-service/internal/config"
	"bank-service/internal/logger"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()

	mylog := logger.SetupLogger(cfg.LogType)
	storage, err := postgres.New(ctx, cfg, mylog)
	if err != nil {
		log.Fatal("failed to coonect db")
	}
	defer storage.Close(ctx)

	printConfig(cfg)
}

func printConfig(cfg *config.Config) {
	time.Sleep(3 * time.Second) 

	log.Println("---------------------------------------")
	log.Println("Starting service:")
	log.Printf("Environment: %s \tLog level: %s\n", cfg.Env, cfg.LogType)
	log.Printf("Postgres: %s:%d \tDB: %s \tUser: %s\n", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName, cfg.Postgres.User)
	log.Printf("RabbitMQ: %s:%d \tUser: %s\n", cfg.Rabbitmq.Host, cfg.Rabbitmq.Port, cfg.Rabbitmq.User)
	log.Printf("REST API: Port: %d \tTimeouts: Read=%s, Write=%s, Idle=%s\n",
		cfg.REST.Port, cfg.REST.ReadTimeout, cfg.REST.WriteTimeout, cfg.REST.IdleTimeout)
	log.Printf("Mock DB: User Count=%d \tMessage Count=%d\n", cfg.MockDB.UserCount, cfg.MockDB.MsgCount)
	log.Println("---------------------------------------")
}
