package errors

import (
	"github.com/anonystick/go-ecommerce-backend-api/global"
	"github.com/anonystick/go-ecommerce-backend-api/pkg/logger"
	"go.uber.org/zap"
)

func MustCheck(err error, msg string) {
	if err != nil {
		// Logger được inject hoặc khởi tạo trong package này
		global.Logger.Error(msg, zap.Error(err))
		panic(err)
	}
}

// not dependentcy looger
func Must(logger *logger.LoggerZap, err error, component string) {
	if err != nil {
		// logger.Error("Failed to initialize "+component, zap.Error(err))
		panic(err)
	}
}
