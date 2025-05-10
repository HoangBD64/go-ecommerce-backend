package initialize

import (
	"fmt"

	"github.com/anonystick/go-ecommerce-backend-api/global"
	"github.com/anonystick/go-ecommerce-backend-api/internal/errors"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper := viper.New()
	viper.AddConfigPath("./config/") // path to config
	viper.SetConfigName("local")     // ten file
	viper.SetConfigType("yaml")

	// read configuration
	err := viper.ReadInConfig()
	if err != nil {
		errors.Must(global.Logger, err, "Error loading configuration")
	}
	// read server configuration
	fmt.Println("Server Port::", viper.GetInt("server.port"))
	fmt.Println("Server Port::", viper.GetString("security.jwt.key"))

	// configure structur
	if err := viper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("Unable to decode configuration %v", err)
	}
}
