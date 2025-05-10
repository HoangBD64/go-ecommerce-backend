package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/anonystick/go-ecommerce-backend-api/global"
	consts "github.com/anonystick/go-ecommerce-backend-api/internal/const"
	"github.com/anonystick/go-ecommerce-backend-api/internal/database"
	"github.com/anonystick/go-ecommerce-backend-api/internal/model"
	"github.com/anonystick/go-ecommerce-backend-api/internal/utils"
	"github.com/anonystick/go-ecommerce-backend-api/internal/utils/auth"
	"github.com/anonystick/go-ecommerce-backend-api/internal/utils/crypto"
	"github.com/anonystick/go-ecommerce-backend-api/internal/utils/random"
	"github.com/anonystick/go-ecommerce-backend-api/internal/utils/sendto"
	"github.com/anonystick/go-ecommerce-backend-api/pkg/response"
	"github.com/redis/go-redis/v9"
)

type sUserLogin struct {
	// Implement the IUserLogin interface here
	r *database.Queries
}

func NewUserLoginImpl(r *database.Queries) *sUserLogin {
	return &sUserLogin{
		r: r,
	}
}

// ---- TWO FACTOR AUTHEN -----

// two-factor authentication
func (s *sUserLogin) IsTwoFactorEnabled(ctx context.Context, userId int) (codeResult int, rs bool, err error) {
	return 200, true, nil
}

// setup authentication
func (s *sUserLogin) SetupTwoFactorAuth(ctx context.Context, in *model.SetupTwoFactorAuthInput) (codeResult int, err error) {
	// Logic
	// 1. Check isTwoFactorEnabled -> true return
	isTwoFactorAuth, err := s.r.IsTwoFactorEnabled(ctx, in.UserId)
	if err != nil {
		return response.ErrCodeTwoFactorAuthSetupFailed, err
	}
	if isTwoFactorAuth > 0 {
		return response.ErrCodeTwoFactorAuthSetupFailed, fmt.Errorf("Two-factor authentication is already enabled")
	}
	// 2. crate new type Authe
	err = s.r.EnableTwoFactorTypeEmail(ctx, database.EnableTwoFactorTypeEmailParams{
		UserID:            in.UserId,
		TwoFactorAuthType: database.PreGoAccUserTwoFactor9999TwoFactorAuthTypeEMAIL,
		TwoFactorEmail:    sql.NullString{String: in.TwoFactorEmail, Valid: true},
	})
	if err != nil {
		return response.ErrCodeTwoFactorAuthSetupFailed, err
	}

	// 3. send otp to in.TwoFactorEmail
	keyUserTwoFator := crypto.GetHash("2fa:" + strconv.Itoa(int(in.UserId)))
	go global.Rdb.Set(ctx, keyUserTwoFator, "123456", time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
	// if err != nil {
	// 	return response.ErrCodeTwoFactorAuthSetupFailed, err
	// }
	return response.ErrCodeSuccess, nil
}

// Verify Two Factor Authentication
func (s *sUserLogin) VerifyTwoFactorAuth(ctx context.Context, in *model.TwoFactorVerificationInput) (codeResult int, err error) {
	// 1. Check isTwoFactorEnabled
	isTwoFatorAuth, err := s.r.IsTwoFactorEnabled(ctx, in.UserId)
	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}
	if isTwoFatorAuth > 0 {
		return response.ErrCodeTwoFactorAuthVerifyFailed, fmt.Errorf("Two-factor authentication is not enabled")
	}

	// 2. Check Otp in redis avaible
	keyUserTwoFator := crypto.GetHash("2fa:" + strconv.Itoa(int(in.UserId)))
	otpVerifyAuth, err := global.Rdb.Get(ctx, keyUserTwoFator).Result()
	if err == redis.Nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, fmt.Errorf("Key %s does not exists", keyUserTwoFator)
	} else if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}
	// 3. check otp
	if otpVerifyAuth != in.TwoFactorCode {
		return response.ErrCodeTwoFactorAuthVerifyFailed, fmt.Errorf("OTP does not match")
	}

	// 4. udpoate status
	err = s.r.UpdateTwoFactorStatus(ctx, database.UpdateTwoFactorStatusParams{
		UserID:            in.UserId,
		TwoFactorAuthType: database.PreGoAccUserTwoFactor9999TwoFactorAuthTypeEMAIL,
	})
	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}
	// 5. remove otp
	_, err = global.Rdb.Del(ctx, keyUserTwoFator).Result()
	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}
	return 200, nil
}

// ---- END TWO FACTOR AUTHEN ----

// Implement the IUserLogin interface here
func (s *sUserLogin) Login(ctx context.Context, in *model.LoginInput) (codeResult int, out model.LoginOutput, err error) {
	// 1. logic login
	userBase, err := s.r.GetOneUserInfo(ctx, in.UserAccount)
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}
	// 2. check password?
	if !crypto.MatchingPassword(userBase.UserPassword, in.UserPassword, userBase.UserSalt) {
		return response.ErrCodeAuthFailed, out, fmt.Errorf("does not match password")
	}
	// 3. check two-factor authentication
	isTwoFactorEnable, err := s.r.IsTwoFactorEnabled(ctx, uint32(userBase.UserID))
	if err != nil {
		return response.ErrCodeAuthFailed, out, fmt.Errorf("does not match password")
	}
	if isTwoFactorEnable > 0 {
		// sen otp to in.TwoFactorEmail
		keyUserLoginTwoFactor := crypto.GetHash("2fa:otp:" + strconv.Itoa(int(userBase.UserID)))
		err = global.Rdb.SetEx(ctx, keyUserLoginTwoFactor, "111111", time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
		if err != nil {
			return response.ErrCodeAuthFailed, out, fmt.Errorf("set otp redis faiuled")
		}
		// send otp via twofactorEmail
		// get email 2fa
		infoUserTwoFactor, err := s.r.GetTwoFactorMethodByIDAndType(ctx, database.GetTwoFactorMethodByIDAndTypeParams{
			UserID:            uint32(userBase.UserID),
			TwoFactorAuthType: database.PreGoAccUserTwoFactor9999TwoFactorAuthTypeEMAIL,
		})
		if err != nil {
			return response.ErrCodeAuthFailed, out, fmt.Errorf("get two factor method failed")
		}
		// go sendto.SendEmailToJavaByAPI()
		log.Println("send OTP 2FA to Email::", infoUserTwoFactor.TwoFactorEmail)
		go sendto.SendTextEmailOtp([]string{infoUserTwoFactor.TwoFactorEmail.String}, consts.HOST_EMAIL, "111111")

		out.Message = "send OTP 2FA to Email, pls het OTP by Email.."
		return response.ErrCodeSuccess, out, nil
	}
	// 4. update password time
	go s.r.LoginUserBase(ctx, database.LoginUserBaseParams{
		UserLoginIp:  sql.NullString{String: "127.0.0.1", Valid: true},
		UserAccount:  in.UserAccount,
		UserPassword: in.UserPassword, // khong can
	})

	// 5. Create UUID User
	subToken := utils.GenerateCliTokenUUID(int(userBase.UserID))
	log.Println("subtoken:", subToken)
	// 6. get user_info table
	infoUser, err := s.r.GetUser(ctx, uint64(userBase.UserID))
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}
	// convert to json
	infoUserJson, err := json.Marshal(infoUser)
	if err != nil {
		return response.ErrCodeAuthFailed, out, fmt.Errorf("convert to json failed: %v", err)
	}
	// 7. give infoUserJson to redis with key = subToken
	err = global.Rdb.Set(ctx, subToken, infoUserJson, time.Duration(consts.TIME_2FA_OTP_REGISTER)*time.Minute).Err()
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}
	// 8. create token
	out.Token, err = auth.CreateToken(subToken)
	if err != nil {
		return
	}
	return 200, out, nil
}

func (s *sUserLogin) Register(ctx context.Context, in *model.RegisterInput) (codeResult int, err error) {
	// logic
	// 1. hash email
	fmt.Printf("VerifyKey: %s\n", in.VerifyKey)
	fmt.Printf("VerifyType: %d\n", in.VerifyType)
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))
	fmt.Printf("hashKey: %s\n", hashKey)

	// 2. check user exists in user base

	userFound, err := s.r.CheckUserBaseExists(ctx, in.VerifyKey)
	/*
		SELECT count(email) FROM user WHERE email = '%s' --> 2
		// user111@gmail.com

	*/
	if err != nil {
		return response.ErrCodeUserHasExists, err
	}

	if userFound > 0 {
		return response.ErrCodeUserHasExists, fmt.Errorf("user has already registered")
	}

	// 3. Create OTP
	userKey := utils.GetUserKey(hashKey) //fmt.Sprintf("u:%s:otp", hashKey)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()

	// util..
	switch {
	case err == redis.Nil:
		fmt.Println("Key does not exist")
	case err != nil:
		fmt.Println("Get failed::", err)
		return response.ErrInvalidOTP, err
	case otpFound != "":
		return response.ErrCodeOtpNotExists, fmt.Errorf("")
	}

	// 4. Generate OTP
	otpNew := random.GenerateSixDigitOtp()
	if in.VerifyPurpose == "TEST_USER" {
		otpNew = 123456
	}
	fmt.Printf("Otp is :::%d\n", otpNew)
	// 5. save OTP in Redis with expiration time
	err = global.Rdb.SetEx(ctx, userKey, strconv.Itoa(otpNew), time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()

	if err != nil {
		return response.ErrInvalidOTP, err
	}
	// 6/ Sen OTP
	switch in.VerifyType {
	case consts.EMAIL:
		err = sendto.SendTextEmailOtp([]string{in.VerifyKey}, consts.HOST_EMAIL, strconv.Itoa(otpNew))
		if err != nil {
			return response.ErrSendEmailOtp, err
		}
		// 7. save OTP to MYSQL
		result, err := s.r.InsertOTPVerify(ctx, database.InsertOTPVerifyParams{
			VerifyOtp:     strconv.Itoa(otpNew),
			VerifyType:    sql.NullInt32{Int32: 1, Valid: true},
			VerifyKey:     in.VerifyKey,
			VerifyKeyHash: hashKey,
		})

		if err != nil {
			return response.ErrSendEmailOtp, err
		}

		// 8. getlasId
		lastIdVerifyUser, err := result.LastInsertId()
		if err != nil {
			return response.ErrSendEmailOtp, err
		}
		log.Println("lastIdVerifyUser", lastIdVerifyUser)
		return response.ErrCodeSuccess, nil
	case consts.MOBILE:
		return response.ErrCodeSuccess, nil
	}

	return response.ErrCodeSuccess, nil
}

func (s *sUserLogin) VerifyOTP(ctx context.Context, in *model.VerifyInput) (out model.VerifyOTPOutput, err error) {
	// logic
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))

	// get otp
	otpFound, err := global.Rdb.Get(ctx, utils.GetUserKey(hashKey)).Result()
	if err != nil {
		return out, err
	}
	if in.VerifyCode != otpFound {
		// Neu nhu ma sai 3 lan trong vong 1 phut??

		return out, fmt.Errorf("OTP not match")
	}
	infoOTP, err := s.r.GetInfoOTP(ctx, hashKey)
	if err != nil {
		return out, err
	}
	// uopdate status verified
	err = s.r.UpdateUserVerificationStatus(ctx, hashKey)
	if err != nil {
		return out, err
	}

	// output
	out.Token = infoOTP.VerifyKeyHash // token temp
	out.Message = "success"

	return out, err
}

func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context, token string, password string) (userId int, err error) {
	// 1. token is already verified : user_verify table
	infoOTP, err := s.r.GetInfoOTP(ctx, token)
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	// 1 check isVerified OK
	if infoOTP.IsVerified.Int32 == 0 {
		return response.ErrCodeUserOtpNotExists, fmt.Errorf("user OTP not verified")
	}
	// 2. check token is exists in user_base
	//update user_base table
	log.Println("infoOTP::", infoOTP)
	userBase := database.AddUserBaseParams{}
	userBase.UserAccount = infoOTP.VerifyKey
	userSalt, err := crypto.GenerateSalt(16)
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	userBase.UserSalt = userSalt
	userBase.UserPassword = crypto.HashPassword(password, userSalt)
	// add userBase to user_base table
	newUserBase, err := s.r.AddUserBase(ctx, userBase)
	log.Println("newUserBase::", newUserBase, userBase)
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	user_id, err := newUserBase.LastInsertId()
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	// add user_id to user_info table
	newUserInfo, err := s.r.AddUserHaveUserId(ctx, database.AddUserHaveUserIdParams{
		UserID:               uint64(user_id),
		UserAccount:          infoOTP.VerifyKey,
		UserNickname:         sql.NullString{String: infoOTP.VerifyKey, Valid: true},
		UserAvatar:           sql.NullString{String: "", Valid: true},
		UserState:            1,
		UserMobile:           sql.NullString{String: "", Valid: true},
		UserGender:           sql.NullInt16{Int16: 0, Valid: true},
		UserBirthday:         sql.NullTime{Time: time.Time{}, Valid: false},
		UserEmail:            sql.NullString{String: infoOTP.VerifyKey, Valid: true},
		UserIsAuthentication: 1,
	})
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	user_id, err = newUserInfo.LastInsertId()
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	return int(user_id), nil
}
