package mapper

import (
	"github.com/anonystick/go-ecommerce-backend-api/internal/database"
	"github.com/anonystick/go-ecommerce-backend-api/internal/model"
)

func ToTicketItemDTO(ticketItem database.GetTicketItemByIdRow) model.TicketItemsOutput {
	return model.TicketItemsOutput{
		TicketId:       int(ticketItem.ID),
		TicketName:     ticketItem.Name,
		StockInitial:   int(ticketItem.StockInitial),
		StockAvailable: int(ticketItem.StockAvailable),
	}
}
