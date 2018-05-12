package configuration

import "testing"

func TestKafkaConfig_LoadConfig(t *testing.T) {
	cfg := MyKafkaConfig()
	t.Run("KafkaTest", func(t *testing.T) {
		if err := cfg.LoadConfig("kafka.json"); err != nil {
			t.Errorf("KafkaConfig.LoadConfig() error = %v", err)
		}
	})
}
