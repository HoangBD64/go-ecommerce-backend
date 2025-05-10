package user

import (
	"github.com/anonystick/go-ecommerce-backend-api/internal/controller/account"
	"github.com/anonystick/go-ecommerce-backend-api/internal/middlewares"
	"github.com/anonystick/go-ecommerce-backend-api/internal/wire"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// public router
	// this is non-dependency
	// ur := repo.NewUserRepository()
	// us := service.NewUserService(ur)
	// userHanlderNonDependency := controller.NewUserController(us)
	userController, _ := wire.InitUserRouterHanlder()
	// WIRE go
	// Dependency Injection (DI) DI java
	userRouterPublic := Router.Group("/user")
	{
		// userRouterPublic.POST("/register", userController.Register) // register -> YES -> No
		userRouterPublic.POST("/register", account.Login.Register)
		userRouterPublic.POST("/verify_account", account.Login.VerifyOTP)
		userRouterPublic.POST("/update_pass_register", account.Login.UpdatePasswordRegister)
		userRouterPublic.POST("/login", account.Login.Login) // login -> YES -> No
	}
	// private router
	userRouterPrivate := Router.Group("/user")
	userRouterPrivate.Use(middlewares.AuthenMiddleware())
	// userRouterPrivate.Use(limiter())
	// userRouterPrivate.Use(Authen())
	// userRouterPrivate.Use(Permission())
	{
		userRouterPrivate.GET("/get_info", userController.Register)
		userRouterPrivate.POST("/two-factor/setup", account.TwoFA.SetupTwoFactorAuth)
		userRouterPrivate.POST("/two-factor/verify", account.TwoFA.VerifyTwoFactorAuth)
	}
}
