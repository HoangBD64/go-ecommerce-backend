package model

// VO: Get ticketItems returns
type TicketItemsOutput struct {
	TicketId       int    `json:"ID"`             // Sửa tag và thêm field
	TicketName     string `json:"Name"`           // Sửa tag
	StockAvailable int    `json:"StockAvailable"` // Sửa tag
	StockInitial   int    `json:"StockInitial"`   // Sửa tag
}

// DTO
type TicketItemRequest struct {
	TicketId string `json:"ticket_Id"`
}

// DTO Request Order

type OrderRequest struct {
	TicketId int    `json:"ticket_Id" validate:"required"`
	UserId   int    `json:"user_Id"`
	Quantity int    `json:"quantity" validate:"gte=1"` // validate:"gte=0,lte=150" Age   int    `json:"age" validate:"gte=0,lte=150"`
	Notes    string `json:"notes"`
}
