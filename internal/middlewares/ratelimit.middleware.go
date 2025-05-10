package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anonystick/go-ecommerce-backend-api/global"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	redisStore "github.com/ulule/limiter/v3/drivers/store/redis"
)

type RateLimiter struct {
	// Implement rate limiter logic here
	globalRateLimiter         *limiter.Limiter
	publicAPIRateLimiter      *limiter.Limiter
	userPrivateAPIRateLimiter *limiter.Limiter
}

func NewRateLimiter() *RateLimiter {
	rateLimit := &RateLimiter{
		globalRateLimiter:         rateLimiter("10000-S"),
		publicAPIRateLimiter:      rateLimiter("80-S"),
		userPrivateAPIRateLimiter: rateLimiter("50-S"),
	}

	return rateLimit
}

func rateLimiter(interval string) *limiter.Limiter {
	store, err := redisStore.NewStoreWithOptions(global.Rdb, limiter.StoreOptions{
		Prefix:          "rate-limiter", // u:uuid -> u:1001
		MaxRetry:        3,
		CleanUpInterval: time.Hour,
	})
	if err != nil {
		return nil
	}
	rate, err := limiter.NewRateFromFormatted(interval) // 5-S, 10-M
	if err != nil {
		panic(err)
	}
	instance := limiter.New(store, rate)
	return instance
}

// GLOBAL LIMITER
func (rl *RateLimiter) GlobalRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "global" // unit
		log.Println("global--->")
		limitContext, err := rl.globalRateLimiter.Get(c, key)
		if err != nil {
			fmt.Println("Failed to check rate limit GLOBAL", err)
			c.Next()
			return
		}
		if limitContext.Reached {
			log.Printf("Rate limit breached GLOBAL %s", key)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit breached GLOBAL, try later"})
			return
		}
		c.Next()
	}

}

// PUBLIC API LIMITER
func (rl *RateLimiter) PublicAPIRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		urlPath := c.Request.URL.Path // urlPath: /ping/80 80 >
		rateLimitPath := rl.filterLimitUrlPath(urlPath)
		if rateLimitPath != nil {
			log.Println("Client Ip--->", c.ClientIP())

			key := fmt.Sprintf("%s-%s", "111-222-333-44", urlPath)
			limitContext, err := rateLimitPath.Get(c, key)
			if err != nil {
				fmt.Println("Failed to check rate limit", err)
				c.Next()
				return
			}
			if limitContext.Reached {
				log.Printf("Rate limit breached %s", key)
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit breached, try later"})
				return
			}
		}

		c.Next()
	}
}

// PRIVATE API LIMITER
func (rl *RateLimiter) UserAndPrivateRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		urlPath := c.Request.URL.Path
		rateLimitPath := rl.filterLimitUrlPath(urlPath)
		if rateLimitPath != nil {
			userId := 1001
			key := fmt.Sprintf("%d-%s", userId, urlPath)
			limitContext, err := rateLimitPath.Get(c, key)
			if err != nil {
				fmt.Println("Failed to check rate limit", err)
				c.Next()
				return
			}
			if limitContext.Reached {
				log.Printf("Rate limit breached %s", key)
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit breached, try later"})
				return
			}
		}

		c.Next()

	}
}

func (rl *RateLimiter) filterLimitUrlPath(urlPath string) *limiter.Limiter {
	if urlPath == "/v1/2024/user/login" || urlPath == "/ping/80" {
		return rl.publicAPIRateLimiter
	} else if urlPath == "/v1/2024/user/info" || urlPath == "/ping/50" {
		return rl.userPrivateAPIRateLimiter
	} else {
		return rl.globalRateLimiter
	}
}
