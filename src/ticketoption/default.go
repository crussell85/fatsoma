package ticketoption

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"ticket-api/storage"
)

var _ Service = (*DefaultTicketOptionService)(nil)

type DefaultTicketOptionService struct {
	ticketOptionStorage storage.TicketOptionStorage
}

func NewDefaultTicketOptionService(ticketOptionStorage storage.TicketOptionStorage) *DefaultTicketOptionService {
	return &DefaultTicketOptionService{
		ticketOptionStorage: ticketOptionStorage,
	}
}

func (d *DefaultTicketOptionService) CreateTicketOption(ctx context.Context, input *CreateTicketOptionInput) (*TicketOption, error) {
	result, err := d.ticketOptionStorage.CreateTicketOption(ctx, &storage.CreateTicketOptionInput{
		Name:        input.Name,
		Description: input.Description,
		Allocation:  input.Allocation,
	})
	if err != nil {
		return nil, err
	}

	return &TicketOption{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Allocation:  result.Allocation,
	}, nil
}

func (d *DefaultTicketOptionService) GetTicketOption(ctx context.Context, ticketOptionId string) (*TicketOption, error) {
	result, err := d.ticketOptionStorage.GetTicketOption(ctx, ticketOptionId)
	if err != nil {
		if errors.Is(err, storage.ErrTicketOptionNotFound) {
			return nil, ErrTicketOptionNotFound
		}

		return nil, fmt.Errorf("unable to retrieve ticket from storage: %w", err)
	}

	return &TicketOption{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Allocation:  result.Allocation,
	}, nil
}

func (d *DefaultTicketOptionService) PurchaseTicketOption(ctx context.Context, input *PurchaseTicketOptionInput) error {
	result, err := d.ticketOptionStorage.GenerateTickets(ctx, &storage.GenerateTicketsInput{
		UserId:         input.UserID,
		TicketOptionId: input.TicketOptionId,
		Quantity:       input.Quantity,
	})
	if err != nil {
		if errors.Is(err, storage.ErrTicketOptionAllocationCheckFailed) {
			return ErrOverAllocatedTickets
		} else {
			return fmt.Errorf("unable to generate tickets: %w", err)
		}
	}

	if len(result.TicketIds) < input.Quantity {
		return ErrNotEnoughTicketsGenerated
	}

	log.Info().Str("purchase_id", result.PurchaseId).Msg("purchase generated")
	for _, ticketId := range result.TicketIds {
		log.Info().Str("ticket_id", ticketId).Msg("ticket generated")
	}

	return nil
}
