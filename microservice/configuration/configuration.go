package configuration

import "github.com/spf13/viper"

// Config provides a generic configuration interface
type Config interface {
	// LoadConfig is used to load a config from any kind of source
	LoadConfig(file string) (err error)
}

// InitViper is a helper function that uses spf13/Viper to load a configuration
// from the provided file.
func InitViper(file string) (cfg *viper.Viper, err error) {
	cfg = viper.GetViper()
	cfg.SetConfigFile(file)
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	return
}
