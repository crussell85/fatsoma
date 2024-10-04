package ticketoption

import (
	"context"
	"errors"
)

var (
	ErrTicketOptionNotFound      = errors.New("ticket option not found")
	ErrOverAllocatedTickets      = errors.New("no more ticket allocation available")
	ErrNotEnoughTicketsGenerated = errors.New("not enough tickets generated")
)

type Service interface {
	CreateTicketOption(ctx context.Context, input *CreateTicketOptionInput) (*TicketOption, error)
	GetTicketOption(ctx context.Context, ticketOptionId string) (*TicketOption, error)
	PurchaseTicketOption(ctx context.Context, input *PurchaseTicketOptionInput) error
}

// CreateTicketOptionInput has the required fields to create a ticket option
type CreateTicketOptionInput struct {
	Name        string
	Description string
	Allocation  int
}

// PurchaseTicketOptionInput has the required fields for purchasing ticket options
type PurchaseTicketOptionInput struct {
	Quantity       int
	UserID         string
	TicketOptionId string
}

// TicketOption represents data for a ticket option
type TicketOption struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	Allocation  int    `json:"allocation"`
}
