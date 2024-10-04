package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"strconv"
	"time"
)

var _ TicketOptionStorage = (*DynamoTicketOptionStorage)(nil)

type DynamoTicketOptionStorage struct {
	cli       *dynamodb.Client
	tableName string
}

func NewDynamoTicketOptionStorage(cli *dynamodb.Client, tableName string) *DynamoTicketOptionStorage {
	return &DynamoTicketOptionStorage{
		cli:       cli,
		tableName: tableName,
	}
}

func (d *DynamoTicketOptionStorage) CreateTicketOption(ctx context.Context, input *CreateTicketOptionInput) (*TicketOptionResult, error) {
	toId := uuid.New()

	allocation := strconv.FormatInt(int64(input.Allocation), 10)

	_, err := d.cli.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"pk":          &types.AttributeValueMemberS{Value: toId.String()},
			"sk":          &types.AttributeValueMemberS{Value: "TICKETOPTION"},
			"name":        &types.AttributeValueMemberS{Value: input.Name},
			"description": &types.AttributeValueMemberS{Value: input.Description},
			"allocation":  &types.AttributeValueMemberN{Value: allocation},
		},
		TableName: aws.String(d.tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("uanble to create ticket option: %w", err)
	}

	return &TicketOptionResult{
		ID:          toId.String(),
		Name:        input.Name,
		Description: input.Description,
		Allocation:  input.Allocation,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}, nil
}

func (d *DynamoTicketOptionStorage) GetTicketOption(ctx context.Context, ticketOptionId string) (*TicketOptionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DynamoTicketOptionStorage) GenerateTickets(ctx context.Context, input *GenerateTicketsInput) (*GenerateTicketsResult, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DynamoTicketOptionStorage) Close() {
	// Not implemented
}
