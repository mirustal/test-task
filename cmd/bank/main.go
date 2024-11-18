package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"bank-service/internal/adapter/db/postgres"
	"bank-service/internal/api"
	"bank-service/internal/modules/client"
	"bank-service/internal/modules/transaction"
	mockdb "bank-service/internal/utils/mockDB"
	"bank-service/pkg/config"
	"bank-service/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	printConfig(cfg)

	ctx := context.Background()

	mylog := logger.SetupLogger(cfg.LogType)

	storage, err := postgres.New(ctx, cfg, mylog)
	if err != nil {
		log.Fatal("failed to coonect db")
	}
	defer storage.Close(ctx)
	if err := mockdb.SeedDataBase(ctx, storage); err != nil {
		mylog.Error("can't create mock for DB")
	}
	

	clService := client.New(mylog, storage, storage)
	trService := transaction.New(mylog, storage, storage)

	app := api.NewRouter(cfg, mylog, clService, trService)
	app.Listen(fmt.Sprintf(":%d", cfg.REST.Port))

}

func printConfig(cfg *config.Config) {
	time.Sleep(3 * time.Second)

	log.Println("-------------------------------------------------------")
	log.Println("Starting service:")
	log.Printf("Environment:\t%s\tLog_level: %s\n", cfg.Env, cfg.LogType)
	log.Printf("Postgres:\t%s:%d\tDB: %s\tUser: %s\n", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName, cfg.Postgres.User)
	log.Printf("RabbitMQ:\t%s:%d\tUser: %s\n", cfg.Rabbitmq.Host, cfg.Rabbitmq.Port, cfg.Rabbitmq.User)
	log.Printf("REST API:\tPort: %d\tTimeouts: Read=%s, Write=%s, Idle=%s\n",
		cfg.REST.Port, cfg.REST.ReadTimeout, cfg.REST.WriteTimeout, cfg.REST.IdleTimeout)
	log.Printf("Mock DB:\tUser Count=%d\tMessage_Count=%d\n", cfg.MockDB.UserCount, cfg.MockDB.MsgCount)
	log.Println("-------------------------------------------------------")
}
