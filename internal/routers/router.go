package routers

// import (
// 	"fmt"

// 	c "github.com/anonystick/go-ecommerce-backend-api/internal/controller"
// 	"github.com/anonystick/go-ecommerce-backend-api/internal/middlewares"
// 	"github.com/gin-gonic/gin"
// )

// func AA() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Println("Before -->AA")
// 		c.Next()
// 		fmt.Println("Alter -->AA")

// 	}
// }

// func BB() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Println("Before -->BB")
// 		c.Next()
// 		fmt.Println("Alter -->BB")

// 	}
// }

// func CC(c *gin.Context) {
// 	fmt.Println("Before -->CC")
// 	c.Next()
// 	fmt.Println("Alter -->CC")
// }

// func NewRouter() *gin.Engine {
// 	r := gin.Default()
// 	// use the middleware
// 	r.Use(middlewares.AuthenMiddleware(), BB(), CC)

// 	v1 := r.Group("/v1/2024")
// 	{
// 		v1.GET("/ping", c.NewPongController().Pong)          // /v1/2024/ping
// 		v1.GET("/user/1", c.NewUserController().GetUserByID) // /v1/2024/ping
// 		// v1.PUT("/ping", Pong)
// 		// v1.PATCH("/ping", Pong)
// 		// v1.DELETE("/ping", Pong)
// 		// v1.HEAD("/ping", Pong)
// 		// v1.OPTIONS("/ping", Pong)
// 	}

// 	// v2 := r.Group("/v2/2024")
// 	// {
// 	// 	v2.GET("/ping", Pong) // /v2/2024/ping
// 	// 	v2.PUT("/ping", Pong)
// 	// 	v2.PATCH("/ping", Pong)
// 	// 	v2.DELETE("/ping", Pong)
// 	// 	v2.HEAD("/ping", Pong)
// 	// 	v2.OPTIONS("/ping", Pong)
// 	// }
// 	return r
// }
