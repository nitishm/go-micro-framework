package configuration

type RPCConfig struct {
	Address string `mapstructure:"addr"`
}

func MyRPCConfig() (cfg *RPCConfig) {
	return &RPCConfig{}
}

func (cfg *RPCConfig) LoadConfig(file string) (err error) {
	v, err := InitViper(file)
	if err != nil {
		return
	}

	cfg.Address = v.GetString("config.addr")

	return
}
