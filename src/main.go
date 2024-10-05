package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"os"
	"ticket-api/controller"
	"ticket-api/storage"
	"ticket-api/ticketoption"
)

func main() {
	dynamoTable := os.Getenv("DYNAMO_TABLE")
	if len(dynamoTable) == 0 {
		log.Fatal().Msg("DYNAMO_TABLE env var must be set")
	}

	dynamoHost := os.Getenv("DYNAMO_HOST")
	if len(dynamoHost) == 0 {
		log.Fatal().Msg("DYNAMO_HOST env var must be set")
	}

	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"), config.WithCredentialsProvider(credentials.StaticCredentialsProvider{Value: aws.Credentials{
		AccessKeyID:     "123",
		SecretAccessKey: "123",
		SessionToken:    "dummy",
		Source:          "local dynamo",
	}}))
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load aws config")
	}

	dynamoCli := dynamodb.NewFromConfig(awsConfig, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(dynamoHost)
	})

	dynamoTicketStorage := storage.NewDynamoTicketOptionStorage(dynamoCli, dynamoTable)

	r := gin.Default()

	ticketOptionService := ticketoption.NewDefaultTicketOptionService(dynamoTicketStorage)
	ticketOptionController := controller.NewTicketOptionController(ticketOptionService)

	AddTicketOptionsRoutes(r, ticketOptionController)

	if err := r.Run(":3000"); err != nil {
		log.Fatal().Err(err).Msg("running router failed")
	}
}
