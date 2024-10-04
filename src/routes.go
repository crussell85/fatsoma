package main

import (
	"github.com/gin-gonic/gin"
	"ticket-api/controller"
)

func AddTicketOptionsRoutes(engine *gin.Engine, controller controller.TicketOption) {
	grp := engine.Group("/ticket_options")

	grp.POST("/", controller.HandleCreateTicketOption())
	grp.GET("/:id", controller.HandleGetTicketOption())
	grp.POST("/:id/purchases", controller.HandlePurchaseTicketOption())
}
