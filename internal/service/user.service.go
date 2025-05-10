package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anonystick/go-ecommerce-backend-api/global"
	"github.com/anonystick/go-ecommerce-backend-api/internal/repo"
	"github.com/anonystick/go-ecommerce-backend-api/internal/utils/crypto"
	"github.com/anonystick/go-ecommerce-backend-api/internal/utils/random"
	"github.com/anonystick/go-ecommerce-backend-api/pkg/response"
	"github.com/segmentio/kafka-go"
)

type IUserService interface {
	Register(email string, purpose string) int
}

type userService struct {
	userRepo     repo.IUserRepository
	userAuthRepo repo.IUserAuthRepository
	//...
}

func NewUserService(
	userRepo repo.IUserRepository,
	userAuthRepo repo.IUserAuthRepository,
	//...
) IUserService {
	return &userService{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}

// Register implements IUserService.
func (us *userService) Register(email string, purpose string) int {
	// 0. hashEmail
	hashEmail := crypto.GetHash(email)
	fmt.Printf("hashEmail::%s", hashEmail)
	// 5. check OTP is available

	// 6. user spam ...

	// 1. check email exists in db
	if us.userRepo.GetUserByEmail(email) {
		return response.ErrCodeUserHasExists
	}
	// 2. new OTP -> ...
	otp := random.GenerateSixDigitOtp()
	if purpose == "TEST_USER" {
		otp = 123456
	}

	fmt.Printf("Otp is :::%d\n", otp)
	// 3. save OTP in Redis with expiration time
	err := us.userAuthRepo.AddOTP(hashEmail, otp, int64(10*time.Minute))
	// fmt.Printf("err is :::%d\n", err)
	if err != nil {
		return response.ErrInvalidOTP
	}
	// 4. send Email OTP
	// err = sendto.SendTemplateEmailOtp([]string{email}, "anonystick@gmail.com", "otp-auth.html", map[string]interface{}{
	// 	"otp": strconv.Itoa(otp),
	// })
	// fmt.Printf("err sendto :::%d\n", err)
	// if err != nil {
	// 	return response.ErrSendEmailOtp
	// }

	// send email OTP by JAVA
	// err = sendto.SendEmailToJavaByAPI(strconv.Itoa(otp), email, "otp-auth.html")
	// // fmt.Printf("err sendto :JAVA::%d\n", err)
	// if err != nil {
	// 	return response.ErrSendEmailOtp
	// }

	// send OTP via Kafak JAVA
	body := make(map[string]interface{})
	body["otp"] = otp
	body["email"] = email

	bodyRequest, _ := json.Marshal(body)

	message := kafka.Message{
		Key:   []byte("otp-auth"),
		Value: []byte(bodyRequest),
		Time:  time.Now(),
	}

	err = global.KafkaProducer.WriteMessages(context.Background(), message)
	if err != nil {
		fmt.Printf("err send to kafka::%v\n", err)
		return response.ErrSendEmailOtp
	}
	return response.ErrCodeSuccess
}
