package util

import "github.com/spf13/viper"

type Config struct {
	// viper use mapstructure package under the hood
	// to unmarshall value
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig will read configuration from config file
// inside a path if it exists
// or override the environmental variables if it's provided
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	// tell viper to look for a file with this name
	viper.SetConfigName("app")
	// type of the app file can use xml, json
	viper.SetConfigType("env")
	// check if the environment var in the file
	// match the keys or not
	// then override, load into viper
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}
	// unmarshall config into a struct
	err = viper.Unmarshal(&config)
	return
}
