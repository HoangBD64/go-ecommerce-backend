package initialize

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"

	"github.com/anonystick/go-ecommerce-backend-api/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

var (
	redisRetryCount = 0        // Biến đếm số lần retry
	maxRetries      = 3        // Số lần retry tối đa
	redisMutex      sync.Mutex // Mutex để tránh race condition
)

// func handlePanic() {
// 	// if r := recover(); r != nil {
// 	// 	global.Logger.Error("Recovered from Redis panic", zap.Any("panic", r))
// 	// 	// Có thể thêm logic retry ở đây, nhưng không nên làm vậy nha Anh em, vì những lỗi mà dùng PANIC là đa số không được retry.

// 	// 	time.Sleep(3 * time.Second) //
// 	// }
// 	if r := recover(); r != nil {
// 		redisMutex.Lock()
// 		defer redisMutex.Unlock()

// 		// Log lỗi và thông tin retry
// 		global.Logger.Error("Recovered from Redis panic",
// 			zap.Any("error", r),
// 			zap.Int("retry_count", redisRetryCount),
// 			zap.String("stack", string(debug.Stack())),
// 		)

// 		// Kiểm tra điều kiện retry
// 		if redisRetryCount < maxRetries {
// 			redisRetryCount++
// 			backoff := time.Duration(redisRetryCount*redisRetryCount) * time.Second // Exponential backoff
// 			fmt.Println(">>>>>>backoff: ", backoff)
// 			global.Logger.Warn("Retrying Redis connection...",
// 				zap.Int("attempt", redisRetryCount),
// 				zap.Duration("backoff", backoff),
// 			)
// 			time.Sleep(backoff)
// 			InitRedis() // Gọi lại hàm khởi tạo
// 		} else {
// 			global.Logger.Fatal("Redis connection failed after max retries")
// 		}
// 	}
// }

func handlePanic() {
	// ✅ Quan trọng: recover phải nằm trong defer
	defer func() {
		if r := recover(); r != nil {
			redisMutex.Lock()
			defer redisMutex.Unlock()

			global.Logger.Error("Recovered from Redis panic",
				zap.Any("error", r),
				zap.Int("retry_count", redisRetryCount),
				zap.Int("maxRetries", maxRetries),
				zap.String("stack", string(debug.Stack())),
			)

			if redisRetryCount < maxRetries {
				redisRetryCount++
				backoff := time.Duration(redisRetryCount*redisRetryCount) * time.Second
				fmt.Println(">>>>>>backoff: ", backoff)
				global.Logger.Warn("Retrying Redis connection...",
					zap.Int("attempt", redisRetryCount),
					zap.Duration("backoff", backoff),
				)
				time.Sleep(backoff)
				InitRedis() // Gọi lại InitRedis() sau khi sleep
			} else {
				global.Logger.Fatal("Redis connection failed after max retries")
			}
		}
	}()
}

func InitRedis() {

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		redisMutex.Lock()
	// 		defer redisMutex.Unlock()

	// 		global.Logger.Error("Recovered from Redis panic",
	// 			zap.Any("error", r),
	// 			zap.Int("retry_count", redisRetryCount),
	// 			zap.Int("maxRetries", maxRetries),
	// 			zap.String("stack", string(debug.Stack())),
	// 		)

	// 		if redisRetryCount < maxRetries {
	// 			redisRetryCount++
	// 			backoff := time.Duration(redisRetryCount*redisRetryCount) * time.Second
	// 			fmt.Println(">>>>>>backoff: ", backoff)
	// 			global.Logger.Warn("Retrying Redis connection...",
	// 				zap.Int("attempt", redisRetryCount),
	// 				zap.Duration("backoff", backoff),
	// 			)
	// 			time.Sleep(backoff)
	// 			InitRedis()
	// 		} else {
	// 			global.Logger.Fatal("Redis connection failed after max retries")
	// 		}
	// 	}
	// }()

	// r := global.Config.Redis
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     fmt.Sprintf("%s:%v", r.Host, r.Port), // 55000
	// 	Password: r.Password,                           // no password set
	// 	DB:       r.Database,                           // use default DB
	// 	PoolSize: 10,                                   //
	// })

	// _, err := rdb.Ping(ctx).Result()
	// if err != nil {
	// 	errors.Must(global.Logger, err, "Redis initialization error")
	// 	// panic(err)
	// }

	// // fmt.Println("Initializing Redis Successfully")
	// global.Logger.Info("Initializing Redis Successfully")
	// redisRetryCount = 0 // Reset retry count khi thành công
	// global.Rdb = rdb
	// // redisExample()
	r := global.Config.Redis

	for redisRetryCount = 0; redisRetryCount <= maxRetries; redisRetryCount++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					global.Logger.Error("Recovered from Redis panic",
						zap.Any("error", r),
						zap.Int("retry_count", redisRetryCount),
						zap.Int("maxRetries", maxRetries),
						zap.String("stack", string(debug.Stack())),
					)
				}
			}()

			rdb := redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%v", r.Host, r.Port),
				Password: r.Password,
				DB:       r.Database,
				PoolSize: 10,
			})

			_, err := rdb.Ping(ctx).Result()
			if err != nil {
				panic(err)
			}

			global.Logger.Info("Initializing Redis Successfully")
			global.Rdb = rdb
			redisRetryCount = 0
		}()

		if global.Rdb != nil {
			// Khởi tạo thành công, break khỏi loop
			break
		}

		if redisRetryCount < maxRetries {
			backoff := time.Duration((redisRetryCount+1)*(redisRetryCount+1)) * time.Second
			fmt.Println(">>>>>>backoff: ", backoff)
			global.Logger.Warn("Retrying Redis connection...",
				zap.Int("attempt", redisRetryCount+1),
				zap.Duration("backoff", backoff),
			)
			time.Sleep(backoff)
		} else {
			global.Logger.Fatal("Redis connection failed after max retries")
		}
	}
}

// advanced
func InitRedisSentinel() {
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster", // Tên master do Sentinel quản lý
		SentinelAddrs: []string{"127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381"},
		DB:            0,        // Sử dụng database mặc định
		Password:      "123456", // Nếu Redis có mật khẩu, điền vào đây
	})

	// Check the connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis Sentinel: %v", err)
	}

	fmt.Println("Connected to Redis Sentinel successfully!")

	// Try setting and getting a value
	err = rdb.Set(ctx, "test_key", "Hello Redis Sentinel!", 0).Err()
	if err != nil {
		log.Fatalf("Error setting key: %v", err)
	}

	val, err := rdb.Get(ctx, "test_key").Result()
	if err != nil {
		log.Fatalf("Error getting key: %v", err)
	}

	fmt.Println("Value of test_key:", val)

	global.Logger.Info("Initializing RedisSentinel Successfully")
	global.Rdb = rdb
	// redisExample()
}

func redisExample() {
	err := global.Rdb.Set(ctx, "score", 100, 0).Err()
	if err != nil {
		fmt.Println("Error redis setting:", zap.Error(err))
		return
	}

	value, err := global.Rdb.Get(ctx, "score").Result()
	if err != nil {
		fmt.Println("Error redis setting:", zap.Error(err))
		return
	}

	global.Logger.Info("value score is::", zap.String("score", value))
}
