package config

import (
	"fmt"

	db "github.com/drizion/wabot-go/database"
	"github.com/spf13/viper"
)

type Configuration struct {
	Database DatabaseConfiguration
	Bot      BotConfiguration
}

var (
	Bot *BotConfiguration
	DB  *DatabaseConfiguration
)

// SetupConfig configuration
func SetupConfig() error {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error to reading config file, %s", err)
		// logger.Errorf("Error to reading config file, %s", err)
		return err
	}

	err := viper.Unmarshal(&Configuration{})
	if err != nil {
		// logger.Errorf("error to decode, %v", err)
		return err
	}

	DB = &DatabaseConfiguration{
		Driver:   "postgres",
		Dbname:   viper.GetString("DB_NAME"),
		Username: viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		LogMode:  viper.GetBool("DB_DEBUG"),
	}

	botConfig := SetupBotConfiguration()

	Bot = &botConfig

	err = db.DbConnection(GetDBDSN())
	if err != nil {
		// logger.Errorf("error to connect database, %v", err)
		return err
	}

	return nil
}
