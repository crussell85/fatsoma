package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"os"
	"ticket-api/controller"
	"ticket-api/storage"
	"ticket-api/ticketoption"
)

func main() {
	pgHost := os.Getenv("PG_HOST")
	if len(pgHost) == 0 {
		log.Fatal().Msg("PG_HOST env var must be set")
	}

	username := os.Getenv("PG_USER")
	if len(username) == 0 {
		log.Fatal().Msg("PG_USER env var must be set")
	}

	password := os.Getenv("PG_PASS")
	if len(password) == 0 {
		log.Fatal().Msg("PG_PASS env var must be set")
	}

	dbName := os.Getenv("PG_NAME")
	if len(dbName) == 0 {
		log.Fatal().Msg("PG_NAME env var must be set")
	}

	pgTicketOptionStorage := storage.NewPostgresTicketOptionStorage(&storage.PostgresStorageConfig{
		Host:     pgHost,
		Username: username,
		Password: password,
		DbName:   dbName,
	})

	defer pgTicketOptionStorage.Close()

	r := gin.Default()

	ticketOptionService := ticketoption.NewDefaultTicketOptionService(pgTicketOptionStorage)
	ticketOptionController := controller.NewTicketOptionController(ticketOptionService)

	AddTicketOptionsRoutes(r, ticketOptionController)

	if err := r.Run(":3000"); err != nil {
		log.Fatal().Err(err).Msg("running router failed")
	}
}
