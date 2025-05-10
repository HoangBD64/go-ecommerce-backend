package initialize

import (
	"github.com/anonystick/go-ecommerce-backend-api/global"
	"github.com/anonystick/go-ecommerce-backend-api/internal/database"
	"github.com/anonystick/go-ecommerce-backend-api/internal/service"
	"github.com/anonystick/go-ecommerce-backend-api/internal/service/impl"
)

func InitServiceInterface() {
	queries := database.New(global.Mdbc)
	// Kiểm tra Redis client có bị nil không
	if global.Rdb == nil {
		panic("global.Rdb is nil! Redis chưa được khởi tạo.")
	}

	// User Service Interface
	service.InitUserLogin(impl.NewUserLoginImpl(queries))
	// Ticker Service Interface
	// If this service use many services then pls use wire(Section wire)
	redisCache := impl.NewRedisCache(global.Rdb) // Khởi tạo IRedisCache implementation.
	localCache, err := impl.NewRistrettoCache()  // initialize ILocalCache implementation
	if err != nil {
		panic("failed to initialize local cache")
	}
	// ticketService, err := wire.InitializeTicketService()
	// if err != nil {
	// 	panic("failed to initialize services: " + err.Error())
	// }
	service.InitTicketItem(impl.NewTicketItemImpl(queries, redisCache, localCache))
}
