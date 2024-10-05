package storage

import (
	"context"
	"errors"
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
	purchaseId := uuid.New().String()

	_, err := d.cli.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: input.TicketOptionId},
			"sk": &types.AttributeValueMemberS{Value: "TICKETOPTION"},
		},
		TableName:           aws.String(d.tableName),
		ConditionExpression: aws.String("allocation >= :zeroValue"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":zeroValue": &types.AttributeValueMemberN{Value: "0"},
			":q":         &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(input.Quantity*-1), 10)},
		},
		UpdateExpression: aws.String("ADD allocation :q"),
	})
	if err != nil {
		var cExpErr *types.ConditionalCheckFailedException
		if errors.As(err, &cExpErr) {
			return nil, ErrTicketOptionAllocationCheckFailed
		} else {
			return nil, fmt.Errorf("unable to allocate tickets: %w", err)
		}
	}

	_, err = d.cli.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"pk":       &types.AttributeValueMemberS{Value: input.TicketOptionId},
			"sk":       &types.AttributeValueMemberS{Value: "PURCHASE"},
			"user_id":  &types.AttributeValueMemberS{Value: input.UserId},
			"quantity": &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(input.Quantity), 10)},
			"id":       &types.AttributeValueMemberS{Value: purchaseId},
		},
		TableName: aws.String(d.tableName),
	})
	if err != nil {
		// TODO roll back allocation change
		return nil, fmt.Errorf("unable to create purchase: %w", err)
	}

	var batchTicketWriteRequests []types.WriteRequest

	var ticketIds []string
	for i := 1; i <= input.Quantity; i++ {
		ticketId := uuid.New().String()
		batchTicketWriteRequests = append(batchTicketWriteRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"pk":          &types.AttributeValueMemberS{Value: input.TicketOptionId},
					"sk":          &types.AttributeValueMemberS{Value: fmt.Sprint("TICKET#", ticketId)},
					"purchase_id": &types.AttributeValueMemberS{Value: purchaseId},
					"ticket_id":   &types.AttributeValueMemberS{Value: ticketId},
				},
			},
		})

		ticketIds = append(ticketIds, ticketId)
	}

	_, err = d.cli.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			d.tableName: batchTicketWriteRequests,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to write tickets: %w", err)
	}

	return &GenerateTicketsResult{
		PurchaseId: purchaseId,
		TicketIds:  ticketIds,
	}, nil

}

func (d *DynamoTicketOptionStorage) Close() {
	// Not implemented
}
