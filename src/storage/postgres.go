package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	_insertSql                 = `INSERT INTO ticket_options(name, "desc", allocation, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	_getTicketSql              = `SELECT id, name, "desc", allocation, created_at, updated_at FROM ticket_options WHERE id = $1`
	_updateTicketAllocationSql = `UPDATE ticket_options SET allocation = allocation - $1, updated_at = now() AT TIME ZONE 'UTC' WHERE id = $2`
	_insertPurchaseSql         = `INSERT INTO purchases(quantity, user_id, ticket_option_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	_insertTicketSql           = `INSERT INTO tickets(ticket_option_id, purchase_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`
)

var _ TicketOptionStorage = (*PostgresTicketOptionStorage)(nil)

type PostgresStorageConfig struct {
	Host     string
	Username string
	Password string
	DbName   string
}

type PostgresTicketOptionStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresTicketOptionStorage(config *PostgresStorageConfig) *PostgresTicketOptionStorage {
	pgUrl := fmt.Sprint("postgres://", config.Username, ":", config.Password, "@", config.Host, "/", config.DbName)

	pgPool, err := pgxpool.New(context.TODO(), pgUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to postgres")
	}

	return &PostgresTicketOptionStorage{
		pool: pgPool,
	}
}

func (p *PostgresTicketOptionStorage) CreateTicketOption(ctx context.Context, input *CreateTicketOptionInput) (*TicketOptionResult, error) {
	result := p.pool.QueryRow(ctx, _insertSql, input.Name, input.Description, input.Allocation, time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))

	var id string
	if err := result.Scan(&id); err != nil {
		return nil, fmt.Errorf("unable to insert ticket option: %w", err)
	}

	if len(id) == 0 {
		return nil, fmt.Errorf("invalid id returned")
	}

	return &TicketOptionResult{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Allocation:  input.Allocation,
	}, nil
}

func (p *PostgresTicketOptionStorage) GetTicketOption(ctx context.Context, ticketOptionId string) (*TicketOptionResult, error) {
	result, err := p.pool.Query(ctx, _getTicketSql, ticketOptionId)
	if err != nil {
		return nil, fmt.Errorf("unable to query ticket options: %w", err)
	}

	ticketOptions, err := pgx.CollectRows(result, pgx.RowToStructByName[TicketOptionResult])
	if err != nil {
		return nil, fmt.Errorf("unable to map results to ticket option: %w", err)
	}

	if len(ticketOptions) == 0 {
		return nil, ErrTicketOptionNotFound
	}

	return &ticketOptions[0], nil
}

func (p *PostgresTicketOptionStorage) GenerateTickets(ctx context.Context, input *GenerateTicketsInput) (*GenerateTicketsResult, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}

	// allocate tickets
	_, err = tx.Exec(ctx, _updateTicketAllocationSql, input.Quantity, input.TicketOptionId)
	if err != nil {
		tx.Rollback(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "ticket_options_allocation_check" {
				return nil, ErrTicketOptionAllocationCheckFailed
			}
		}

		return nil, fmt.Errorf("unable to allocate tickets %w", err)
	}

	// create purchase
	purchaseResult := tx.QueryRow(ctx, _insertPurchaseSql, input.Quantity, input.UserId, input.TicketOptionId, time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))

	var purchaseId string
	if err := purchaseResult.Scan(&purchaseId); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("unable to insert purchase: %w", err)
	}

	if len(purchaseId) == 0 {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("invalid purchase id returned")
	}

	// generate tickets
	ticketBatch := pgx.Batch{}
	for i := 1; i <= input.Quantity; i++ {
		ticketBatch.Queue(_insertTicketSql, input.TicketOptionId, purchaseId, time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))
	}

	ticketBatchResults := tx.SendBatch(ctx, &ticketBatch)

	// get the ticket ids
	var ticketIds []string
	for i := 1; i <= input.Quantity; i++ {
		var ticketId string
		insertedRow := ticketBatchResults.QueryRow()
		if err := insertedRow.Scan(&ticketId); err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("unable to get ticket id from ticket insert: %w", err)
		}

		ticketIds = append(ticketIds, ticketId)
	}

	if err := ticketBatchResults.Close(); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("unable to close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("unable to commit generate tickets transaction: %w", err)
	}

	return &GenerateTicketsResult{
		PurchaseId: purchaseId,
		TicketIds:  ticketIds,
	}, nil
}

func (p *PostgresTicketOptionStorage) Close() {
	p.pool.Close()
}
