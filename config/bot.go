package config

import (
	"strings"

	"github.com/spf13/viper"
)

type BotConfiguration struct {
	OwnerNumbers []string
	Prefix       string
}

func SetupBotConfiguration() BotConfiguration {

	prefix := viper.GetString("WABOT_PREFIX")

	numbers := viper.GetString("WABOT_OWNER_NUMBERS")
	slicedNumbers := strings.Split(numbers, ",")

	var trimmedValues []string

	for _, value := range slicedNumbers {
		trimmedValues = append(trimmedValues, strings.TrimSpace(value))
	}

	return BotConfiguration{
		OwnerNumbers: trimmedValues,
		Prefix:       prefix,
	}
}
