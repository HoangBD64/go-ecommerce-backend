package middlewares

import (
	"context"
	"log"

	"github.com/anonystick/go-ecommerce-backend-api/internal/utils/auth"
	"github.com/gin-gonic/gin"
)

func AuthenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the request url path
		uri := c.Request.URL.Path
		log.Println(" uri request: ", uri)
		// check headers authorization
		jwtToken, valid := auth.ExtractBearerToken(c)
		if !valid {
			c.AbortWithStatusJSON(401, gin.H{"code": 40001, "err": "Unauthorized", "description": ""})
			return
		}

		// validate jwt token by subject
		claims, err := auth.VerifyTokenSubject(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"code": 40001, "err": "invalid token", "description": ""})
			return
		}
		// update claims to context
		log.Println("claims::: UUID::", claims.Subject) // 11clitoken....
		ctx := context.WithValue(c.Request.Context(), "subjectUUID", claims.Subject)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
