package config

import "github.com/spf13/viper"

var Cfg *Config

func LoadConfig() (*Config, error){
	cfg := &Config{}
	viper.AddConfigPath(".")
	
	// Try to load local config first for development
	viper.SetConfigName("config.local")
	viper.SetConfigType("yaml")
	
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		// If local config fails, try the default config
		viper.SetConfigName("config")
		err = viper.ReadInConfig()
		if err != nil {
			return cfg, err
		}
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return	cfg, err
	}
	return cfg, err
}