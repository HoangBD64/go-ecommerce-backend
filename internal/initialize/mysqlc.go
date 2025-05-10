package initialize

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/anonystick/go-ecommerce-backend-api/global"
	"github.com/anonystick/go-ecommerce-backend-api/internal/errors"
	"gorm.io/gen"
)

func InitMysqlC() {
	m := global.Config.Mysql
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, m.Username, m.Password, m.Host, m.Port, m.Dbname)
	db, err := sql.Open("mysql", s)
	errors.Must(global.Logger, err, "InitMysql initialization error")
	global.Logger.Info("Initializing MySQL Successfully sql")
	global.Mdbc = db

	// set Pool
	SetPool()
	// genTableDAO()
	// migrateTables()
}

// InitMysql().SetPool()
func SetPoolC() {
	m := global.Config.Mysql
	sqlDb, err := global.Mdb.DB()
	if err != nil {
		fmt.Printf("mysql error: %s::", err)
	}
	sqlDb.SetConnMaxIdleTime(time.Duration(m.MaxIdleConns))
	sqlDb.SetMaxOpenConns(m.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))
}

func genTableDAOC() {
	// Initiate the tables
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/model",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	// // gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(global.Mdb) // reuse your gorm db
	// g.GenerateAllTable()
	g.GenerateModel("go_crm_user")
	// // Generate basic type-safe DAO API for struct `model.User` following conventions
	// g.ApplyBasic(model.User{})

	// // Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(Querier) {}, model.User{}, model.Company{})

	// Generate the code
	g.Execute()
}
func migrateTablesC() {
	err := global.Mdb.AutoMigrate(
	// &po.User{},
	// &po.Role{},
	// &model.GoCrmUserV2{},
	)
	if err != nil {
		fmt.Println("Migrating tables error:", err)
	}
}
