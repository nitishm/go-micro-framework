package configuration

import "github.com/spf13/viper"

type Config interface {
	LoadConfig(file string) (err error)
}

func InitViper(file string) (cfg *viper.Viper, err error) {
	cfg = viper.GetViper()
	cfg.SetConfigFile(file)
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	return
}
