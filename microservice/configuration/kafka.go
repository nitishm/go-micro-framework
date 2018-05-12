package configuration

type KafkaConfig struct {
	Brokers   []string `mapstructure:"brokers"`
	Group     int      `mapstructure:"group"`
	Offset    int64    `mapstructure:"offset"`
	Partition int32    `mapstructure:"partition"`
}

func MyKafkaConfig() (cfg *KafkaConfig) {
	return &KafkaConfig{}
}

func (cfg *KafkaConfig) LoadConfig(file string) (err error) {
	v, err := InitViper("kafka.json")
	if err != nil {
		return
	}

	cfg.Brokers = v.GetStringSlice("config.brokers")
	cfg.Group = v.GetInt("config.group")
	cfg.Offset = v.GetInt64("config.offset")
	cfg.Partition = int32(v.GetInt64("config.partition"))

	return
}
