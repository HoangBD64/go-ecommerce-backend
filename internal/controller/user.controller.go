package controller

import (
	"fmt"

	"github.com/anonystick/go-ecommerce-backend-api/internal/service"
	"github.com/anonystick/go-ecommerce-backend-api/internal/vo"
	"github.com/anonystick/go-ecommerce-backend-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(
	userService service.IUserService,
) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) Register(c *gin.Context) {
	var params vo.UserRegistratorRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamInvalid, err.Error())
		return
	}
	fmt.Printf("Email params: %s", params.Email)
	result := uc.userService.Register(params.Email, params.Purpose)
	response.SuccessResponse(c, result, nil)
}

// // uc user controller
// // us user service

// // controller -> service -> repo -> models -> dbs
// func (uc *UserController) GetUserByID(c *gin.Context) {

// 	// response.SuccessResponse(c, 20001, []string{"tipjs", "m10", "anonystick"})
// 	// response.ErrorResponse(c, 20003, "No need!!")
// }
