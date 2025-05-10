package ticket

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anonystick/go-ecommerce-backend-api/internal/model"
	"github.com/anonystick/go-ecommerce-backend-api/internal/service"
	"github.com/anonystick/go-ecommerce-backend-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// manager controller Ticket Item
var TicketItem = new(cTicketItem)

type cTicketItem struct{}

// NewTicketItem creates a new

func (p *cTicketItem) GetTicketItemById(ctx *gin.Context) {
	// get the ticket item
	ticket_item := ctx.Param("id")
	// Convert the string parameter to an integer.
	idInt, err := strconv.Atoi(ticket_item)
	if err != nil {
		// Handle the conversion error.  This is crucial!
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket item ID format"})
		return
	}
	// call implementation
	ticketItem, err := service.TicketItem().GetTicketItemById(ctx, idInt)
	if err != nil {
		if errors.Is(err, response.CouldNotGetTicketErr) {
			fmt.Println("4004???")
		}

		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())

	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, ticketItem)
}

// order ticket by user
// @Summary      order ticket by user
// @Description  Uorder ticket by user
// @Tags         vetautet-api service
// @Accept       json
// @Produce      json
// @Param        payload body model.OrderRequest true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /ticket/order  [post]
func (p *cTicketItem) PlaceOrderByUser(ctx *gin.Context) {

	//get context
	validation, exists := ctx.Get("validation")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Validator not found"})
		return
	}

	var orderDTO model.OrderRequest
	if err := ctx.ShouldBindJSON(&orderDTO); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	// check validation
	err := validation.(*validator.Validate).Struct(orderDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, orderDTO)
}
