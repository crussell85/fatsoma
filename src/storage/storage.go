package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTicketOptionNotFound              = errors.New("ticket option not found")
	ErrTicketOptionAllocationCheckFailed = errors.New("ticket option allocation check failed")
)

// TicketOptionStorage an interface for interacting with ticket options storage implementations
type TicketOptionStorage interface {
	CreateTicketOption(ctx context.Context, input *CreateTicketOptionInput) (*TicketOptionResult, error)
	GetTicketOption(ctx context.Context, ticketOptionId string) (*TicketOptionResult, error)
	GenerateTickets(ctx context.Context, input *GenerateTicketsInput) (*GenerateTicketsResult, error)
	Close()
}

type CreateTicketOptionInput struct {
	Name        string
	Description string
	Allocation  int
}

type TicketOptionResult struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"desc"`
	Allocation  int       `db:"allocation"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreatePurchaseInput struct {
	Quantity       int
	UserId         string
	TicketOptionId string
}

type CreatePurchaseResult struct {
	ID string `db:"id"`
}

type CreateTicketsInput struct {
	TicketOptionId string
	PurchaseId     string
}

type CreateTicketsResult struct {
	TicketIds []string
}

type GenerateTicketsInput struct {
	UserId         string
	TicketOptionId string
	Quantity       int
}

type GenerateTicketsResult struct {
	PurchaseId string
	TicketIds  []string
}
