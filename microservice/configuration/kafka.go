package configuration

// KafkaConfig encapsulates all the configuration parameters used
// to setup a kafka consumer.
type KafkaConfig struct {
	// A list of broker IPPort strings
	Brokers []string `mapstructure:"brokers"`
	// The kafka consumer group number
	Group int `mapstructure:"group"`
	// The topic Offset (-1 => newest and 0 => from the beginning).
	Offset int64 `mapstructure:"offset"`
	// The partition to consume the messages from registed on the kafka broker.
	Partition int32 `mapstructure:"partition"`
}

// MyKafkaConfig returns a KafkaConfig instance
func MyKafkaConfig() (cfg *KafkaConfig) {
	return &KafkaConfig{}
}

// LoadConfig implements the Config interface method LoadConfig
// Kafka LoadConfig loads the config from a file on the system using
// the Viper helper function.
func (cfg *KafkaConfig) LoadConfig(file string) (err error) {
	v, err := InitViper(file)
	if err != nil {
		return
	}

	cfg.Brokers = v.GetStringSlice("config.brokers")
	cfg.Group = v.GetInt("config.group")
	cfg.Offset = v.GetInt64("config.offset")
	cfg.Partition = int32(v.GetInt64("config.partition"))

	return
}
