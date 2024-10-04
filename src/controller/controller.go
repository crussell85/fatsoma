package controller

import (
	"github.com/gin-gonic/gin"
)

type TicketOption interface {
	HandleCreateTicketOption() func(ctx *gin.Context)
	HandleGetTicketOption() func(ctx *gin.Context)
	HandlePurchaseTicketOption() func(ctx *gin.Context)
}

type CreateTicketOptionRequest struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	Allocation  int    `json:"allocation"`
}

type PurchaseTicketOptionRequest struct {
	Quantity int    `json:"quantity"`
	UserID   string `json:"user_id"`
}
