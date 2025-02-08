package routers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine{
	r := gin.Default()
	v1 := r.Group("/v1/2024") 
	{
		v1.GET("/ping", pong)
		v1.POST("/ping", pong)
		v1.PUT("/ping", pong)
		v1.PATCH("/ping", pong)
		v1.DELETE("/ping", pong)
		v1.HEAD("/ping", pong)
		v1.OPTIONS("/ping", pong)
	}
	v2 := r.Group("/v2/2024") 
	{
		v2.GET("/ping", pong)
		v2.POST("/ping", pong)
		v2.PUT("/ping", pong)
		v2.PATCH("/ping", pong)
		v2.DELETE("/ping", pong)
		v2.HEAD("/ping", pong)
		v2.OPTIONS("/ping", pong)
	}
	return r
}

func pong(c *gin.Context) {
	name := c.DefaultQuery("name", "hoangbd")

	uid := c.Query("uid")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong" + name,
		"uid":     uid,
		"users":    []string{"nguyen", "dinh", "hoang"},
	})
}