package repo

import (
	"github.com/anonystick/go-ecommerce-backend-api/global"
	"github.com/anonystick/go-ecommerce-backend-api/internal/database"
)

// type UserRepo struct{}

// func NewUserRepo() *UserRepo {
// 	return &UserRepo{}
// }

// // user repo u
// func (ur *UserRepo) GetInfoUser() string {
// 	return "Tipjs"
// }

// INTERFACE_VERSION

type IUserRepository interface {
	GetUserByEmail(email string) bool
}

type userReppository struct {
	sqlc *database.Queries
}

// GetUserByEmail implements IUserRepository.
func (up *userReppository) GetUserByEmail(email string) bool {
	// SELECT * FROM user WHERE email = '??' ORDER BY email
	// row := global.Mdb.Table(TableNameGoCrmUser).Where("usr_email = ?", email).First(&model.GoCrmUser{}).RowsAffected
	// user, err := up.sqlc.GetUserByEmailSQLC(ctx, email)
	// if err != nil {
	// 	fmt.Printf("GetUserByEmail error: %v\n", err)
	// 	return false
	// }
	return 1 != 0
}

func NewUserRepository() IUserRepository {
	return &userReppository{
		sqlc: database.New(global.Mdbc),
	}
}
