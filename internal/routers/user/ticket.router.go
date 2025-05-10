package user

import (
	"github.com/anonystick/go-ecommerce-backend-api/internal/controller/ticket"
	"github.com/gin-gonic/gin"
)

type TicketRouter struct{}

func (tr *TicketRouter) InitTicketRouter(Router *gin.RouterGroup) {
	// public router
	ticketRouterPublic := Router.Group("/ticket")
	{
		// ticketRouterPublic.GET("/search")
		ticketRouterPublic.GET("/item/:id", ticket.TicketItem.GetTicketItemById)

		ticketRouterPublic.POST("/order", ticket.TicketItem.PlaceOrderByUser)
	}
	// private router

}
