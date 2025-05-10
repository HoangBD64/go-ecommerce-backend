package initialize

import (
	"github.com/gin-gonic/gin"
)

func Run() *gin.Engine {
	// load configuration
	LoadConfig()
	// m := global.Config.Mysql
	// fmt.Println("Loading configuration nysql", m.Username, m.Password)
	InitLogger()
	// InitRedis()
	// global.Logger.Info("Config Log ok!!", zap.String("ok", "success"))
	InitMysql()
	InitMysqlC()
	// InitRedis()
	InitRedisSentinel()
	InitKafka()
	InitServiceInterface()
	r := InitRouter()
	return r
	// r.Run(":8002")
}
