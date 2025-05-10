package repo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/anonystick/go-ecommerce-backend-api/global"
)

type IUserAuthRepository interface {
	AddOTP(email string, otp int, expirationTime int64) error
}

type userAuthRepository struct{}

// AddOTP implements IUserAuthRepository.
func (u *userAuthRepository) AddOTP(email string, otp int, expirationTime int64) error {
	// panic("unimplemented")
	key := fmt.Sprintf("usr:%s:otp", email) // usr:email:otp
	// err := global.Rdb.SetEx(ctx, key, strconv.Itoa(otp), 10*time.Minute).Err()
	return global.Rdb.SetEx(ctx, key, strconv.Itoa(otp), 10*time.Minute).Err()
}

func NewUserAuthRepository() IUserAuthRepository {
	return &userAuthRepository{}
}
