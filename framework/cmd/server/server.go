package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Twsouza/codeflix-encoder/application/services"
	"github.com/Twsouza/codeflix-encoder/framework/database"
	"github.com/Twsouza/codeflix-encoder/framework/queue"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var db database.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("error loading AUTO_MIGRATE_DB env var: %s", os.Getenv("AUTO_MIGRATE_DB"))
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("error loading DEBUG env var: %s", os.Getenv("DEBUG"))
	}

	db.Env = os.Getenv("ENV")
	db.Debug = debug
	db.AutoMigrateDb = autoMigrateDb

	db.Dsn = os.Getenv("DSN")
	db.DsnTest = os.Getenv("DSN_TEST")

	db.DbType = os.Getenv("DB_TYPE")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
}

func main() {
	messageChannel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan services.JobWorkerResult)
	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalf("error connecting to database")
	}
	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := services.NewJobManager(dbConnection, rabbitMQ, messageChannel, jobReturnChannel)
	jobManager.Start(ch)
}
