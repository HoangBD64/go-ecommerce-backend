package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/anonystick/go-ecommerce-backend-api/global"
	consts "github.com/anonystick/go-ecommerce-backend-api/internal/const"
	"github.com/anonystick/go-ecommerce-backend-api/internal/database"
	"github.com/anonystick/go-ecommerce-backend-api/internal/model"
	"github.com/anonystick/go-ecommerce-backend-api/internal/model/mapper"
	"github.com/anonystick/go-ecommerce-backend-api/internal/service"
	"github.com/anonystick/go-ecommerce-backend-api/pkg/response"
	"go.uber.org/zap"
)

type sTicketItem struct {
	// implementation interface here
	r                *database.Queries
	distributedCache service.IRedisCache
	localCache       service.ILocalCache // Video Go 39: Add local cache
}

func NewTicketItemImpl(r *database.Queries, redisCache service.IRedisCache, localCache service.ILocalCache) *sTicketItem {
	return &sTicketItem{
		r:                r,
		distributedCache: redisCache,
		localCache:       localCache,
	}
}

func (s *sTicketItem) GetTicketItemById(ctx context.Context, ticketId int) (out model.TicketItemsOutput, err error) {

	// 1. get ticket item from local cache
	fmt.Println("START GET TICKET >>>>>> WITH TICKETID -> | ", ticketId)

	out, err = s.getTicketItemFromLocalCache(ctx, ticketId, "v1") // version...
	if err != nil {
		return out, fmt.Errorf("%w with id = %d -> err: %w", response.CouldNotGetTicketErr, ticketId, err)
	}

	if (out != model.TicketItemsOutput{}) {
		// fmt.Println("12 - RESPONSE TICKET ITEM LOCAL CACHE -> CHECK DATA TICKET WITH ID -> ", ticketId)
		// global.Logger.Info("12 - RESPONSE TICKET ITEM LOCAL CACHE -> CHECK DATA TICKET WITH ID -> ", zap.Any("ticketId", ticketId))
		return out, nil
	}
	// 1 get cache from distributed cache
	out, err = s.getTicketItemFromDistributedCache(ctx, ticketId)
	if err != nil {
		return out, fmt.Errorf("%w with id = %d -> err: %w", response.CouldNotGetTicketErr, ticketId, err)
	}

	if (out != model.TicketItemsOutput{}) {
		//fmt.Println("13 - RESPONSE TICKET ITEM DISTRIBUTED CACHE -> CHECK DATA TICKET WITH ID -> ", ticketId)
		//global.Logger.Info("13 - RESPONSE TICKET ITEM DISTRIBUTED CACHE -> CHECK DATA TICKET WITH ID  -> ", zap.Any("ticketId", ticketId))
		return out, nil
	}

	// out, err = s.getTicketItemFromDatabase(ctx, ticketId)
	out, err = s.getTicketItemFromDatabaseLock(ctx, ticketId)
	if err != nil {
		return out, fmt.Errorf("%w with id = %d -> err: %w", response.CouldNotGetTicketErr, ticketId, err)
	}
	// fmt.Println("11 -RESPONSE TICKET ITEM MYSQL -> CHECK DATA TICKET WITH ID -> ", ticketId)
	// global.Logger.Info("11 -RESPONSE TICKET ITEM MYSQL+++")
	return out, nil
}

func (s *sTicketItem) getTicketItemFromDatabaseLock(ctx context.Context, ticketId int) (out model.TicketItemsOutput, err error) {
	lockKey := "lock:ticketItem:" + strconv.Itoa(ticketId)
	// 1. lock
	err = s.distributedCache.WithDistributedLock(ctx, lockKey, 5, func(ctx context.Context) error {
		global.Logger.Info("LOCK ACQUIRED → QUERY DATABASE →", zap.Any("ticketId", ticketId))

		// 2. get data from mysql
		ticketItem, err := s.r.GetTicketItemById(ctx, int64(ticketId))
		if err != nil {
			return err
		}

		ticketItemCacheJSON, err := json.Marshal(ticketItem)
		if err != nil {
			return fmt.Errorf("marshal json failed: %w", err)
		}

		err = global.Rdb.Set(
			ctx, s.getKeyTicketItemCache(ticketId),
			ticketItemCacheJSON,
			time.Duration(consts.TIME_2FA_OTP_REGISTER)*time.Minute,
		).Err()
		if err != nil {
			return fmt.Errorf("save redis failed: %w", err)
		}

		isSuccess := s.localCache.SetWithTTL(ctx, s.getKeyTicketItemCache(ticketId), ticketItem)
		if !isSuccess {
			return fmt.Errorf("save local cache failed")
		}
		// if err := s.updateCaches(ctx, ticketId, ticketItemCacheJSON, ticketItem); err != nil {
		// 	return err
		// }

		out = mapper.ToTicketItemDTO(ticketItem)
		return nil
	})

	return
}

// get data from database
func (s *sTicketItem) getTicketItemFromDatabase(ctx context.Context, ticketId int) (out model.TicketItemsOutput, err error) {

	fmt.Println("07 - QUERY DATABASE -> CHECK DATA TICKET WITH ID -> ", ticketId)

	// lock here
	ticketItem, err := s.r.GetTicketItemById(ctx, int64(ticketId))
	if err != nil {
		return out, err
	}
	// add to redis cache -> not LOCK ...
	fmt.Println("08 - QUERY DATABASE: FOUND -> CHECK DATA TICKET WITH ID -> ", ticketId)
	ticketItemCacheJSON, err := json.Marshal(ticketItem)

	if err != nil {
		return out, fmt.Errorf("convert to json failed: %v", err)
	}

	err = global.Rdb.Set(
		ctx, s.getKeyTicketItemCache(ticketId),
		ticketItemCacheJSON,
		time.Duration(consts.TIME_2FA_OTP_REGISTER)*time.Minute,
	).Err()

	if err != nil {
		return out, fmt.Errorf("save redis failed: %v", err)
	}
	fmt.Println("09 - QUERY DATABASE: SET DATA TO DISTRIBUTED CACHE -> SUCCESS -> ", ticketId)

	// ticketItemCache, err := json.Marshal(ticketItemCacheJSON) // conver to byte
	// fmt.Println("ticketItemticketItem", ticketItemCacheJSON)

	isSuccess := s.localCache.SetWithTTL(ctx, s.getKeyTicketItemCache(ticketId), ticketItem)

	if !isSuccess {
		return out, fmt.Errorf("save localcache failed: %w", err)
	}

	fmt.Println("10 - QUERY DATABASE: SET DATA TO LOCAL CACHE -> SUCCESS -> ", ticketId)
	// reponse to client request
	out = mapper.ToTicketItemDTO(ticketItem) // Ở đây có mùi của java kakak...

	return out, nil
}

// get data from redis distributed

func (s *sTicketItem) getTicketItemFromDistributedCache(ctx context.Context, ticketId int) (out model.TicketItemsOutput, err error) {

	// ticketItemCache, err := global.Rdb.Get(ctx, s.getKeyTicketItemCache(ticketId)).Result()
	// if err != nil {
	// 	if errors.Is(err, redis.Nil) {
	// 		// Trả về lỗi riêng khi key không tồn tại
	// 		return out, nil
	// 	}
	// 	return out, fmt.Errorf("failed to get ticket item cache: %v", err)
	// }
	fmt.Println("04 - DISTRIBUTED CACHE -> CHECK DATA TICKET WITH ID -> ", ticketId)

	ticketItemCache, err := s.distributedCache.Get(ctx, s.getKeyTicketItemCache(ticketId))
	if err != nil {
		return out, fmt.Errorf("failed to get ticket item cache: %v", err)
	}
	if ticketItemCache == "" {
		fmt.Println("05 - DISTRIBUTED CACHE: NOT FOUND -> CHECK DATA TICKET WITH ID -> ", ticketId)
		return out, nil
	}

	if err := json.Unmarshal([]byte(ticketItemCache), &out); err != nil {
		return out, fmt.Errorf("parse redis data failed: %v", err)
	}
	fmt.Println("06 - DISTRIBUTED CACHE: FOUND -> CHECK DATA TICKET WITH ID -> ", ticketId, ticketItemCache)

	// put to local cache
	s.localCache.SetWithTTL(ctx, s.getKeyTicketItemCache(ticketId), out)
	return out, nil
}

// get data from local cache
func (s *sTicketItem) getTicketItemFromLocalCache(ctx context.Context, ticketId int, version string) (out model.TicketItemsOutput, err error) {
	// global.Logger.Info("getTicketItemFromLocalCache with ticketId: ", zap.Int("TicketId", ticketId))
	// var out model.TicketItemsOutput
	fmt.Println("01 - LOCAL CACHE -> CHECK DATA TICKET WITH ID -> ", ticketId)

	ticketItemLocalCache, isFound := s.localCache.Get(ctx, s.getKeyTicketItemCache(ticketId))
	if !isFound {
		// global.Logger.Info("getTicketItemFromLocalCache with ticketId is notfound: ", zap.Int("TicketId", ticketId))
		// fmt.Println("getTicketItemFromLocalCache with ticketId is notfound: ", ticketId)
		fmt.Println("02 - LOCAL CACHE: NOT FOUND -> CHECK DATA TICKET WITH ID -> ", ticketId)
		return out, nil
	}
	// fmt.Println(">>>", ticketItemLocalCache)
	// fmt.Printf(">>> Value: %+v, Type: %s\n", ticketItemLocalCache, reflect.TypeOf(ticketItemLocalCache))

	fmt.Println("03 - LOCAL CACHE: FOUND -> CHECK DATA TICKET WITH ID -> ", ticketId)

	// // Type Assertion to string
	jsonTicketString, ok := ticketItemLocalCache.(string)
	if !ok {
		fmt.Printf("ERROR: Local cache item with key %d is not a string\n", ticketId)
	}

	if err := json.Unmarshal([]byte(jsonTicketString), &out); err != nil {
		return out, fmt.Errorf("parse redis data failed: %v", err)
	}

	return out, nil
}

// generate key cache
func (s *sTicketItem) getKeyTicketItemCache(ticketId int) string {
	return "PRO_TICKET:ITEM:" + strconv.Itoa(ticketId)
}

func (s *sTicketItem) updateCaches(ctx context.Context, ticketId int, ticketItemCacheJSON []byte, rawModel any) error {
	// Save to Redis
	err := global.Rdb.Set(
		ctx,
		s.getKeyTicketItemCache(ticketId),
		ticketItemCacheJSON,
		time.Duration(consts.TIME_2FA_OTP_REGISTER)*time.Minute,
	).Err()
	if err != nil {
		return fmt.Errorf("save redis failed: %w", err)
	}

	// Save to local cache
	isSuccess := s.localCache.SetWithTTL(ctx, s.getKeyTicketItemCache(ticketId), rawModel)
	if !isSuccess {
		return fmt.Errorf("save local cache failed")
	}

	return nil
}
