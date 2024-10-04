package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"ticket-api/ticketoption"
)

var _ TicketOption = (*TicketOptionController)(nil)

type TicketOptionController struct {
	ticketOptionService ticketoption.Service
}

func NewTicketOptionController(ticketOptionService ticketoption.Service) *TicketOptionController {
	return &TicketOptionController{
		ticketOptionService: ticketOptionService,
	}
}

func (t *TicketOptionController) HandleCreateTicketOption() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req CreateTicketOptionRequest
		if err := ctx.BindJSON(&req); err != nil {
			log.Error().Err(err).Msg("unable to parse create ticket option request")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request",
			})
			return
		}

		to, err := t.ticketOptionService.CreateTicketOption(ctx, &ticketoption.CreateTicketOptionInput{
			Name:        req.Name,
			Description: req.Description,
			Allocation:  req.Allocation,
		})
		if err != nil {
			log.Error().Err(err).Msg("unable to create ticket option")
			GenericErrorMessage(ctx)
			return
		}

		ctx.JSON(http.StatusOK, to)
	}
}

func (t *TicketOptionController) HandleGetTicketOption() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if len(id) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "id parameter is required",
			})
			return
		}

		result, err := t.ticketOptionService.GetTicketOption(ctx, id)
		if err != nil {
			if errors.Is(err, ticketoption.ErrTicketOptionNotFound) {
				ctx.JSON(http.StatusNotFound, nil)
			} else {
				log.Error().Err(err).Msg("unable to get ticket option")
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "something went wrong",
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, result)
	}
}

func (t *TicketOptionController) HandlePurchaseTicketOption() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if len(id) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "id parameter is required",
			})
			return
		}

		var purchaseTicketsRequest PurchaseTicketOptionRequest
		if err := ctx.BindJSON(&purchaseTicketsRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "unable to process request data",
			})
			return
		}

		if err := t.ticketOptionService.PurchaseTicketOption(ctx, &ticketoption.PurchaseTicketOptionInput{
			Quantity:       purchaseTicketsRequest.Quantity,
			UserID:         purchaseTicketsRequest.UserID,
			TicketOptionId: id,
		}); err != nil {
			if errors.Is(err, ticketoption.ErrOverAllocatedTickets) {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"quantity": purchaseTicketsRequest.Quantity,
					"message":  "quantity exceeds available allocation",
				})
			} else if errors.Is(err, ticketoption.ErrNotEnoughTicketsGenerated) {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "not enough tickets were generated",
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "something went wrong",
				})
				log.Error().Err(err).Msg("unable to purchase ticket option")
			}
		}

		ctx.Status(http.StatusOK)

		return
	}
}

func GenericErrorMessage(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": "something went wrong",
	})
}
