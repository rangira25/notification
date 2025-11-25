package kafka

import "github.com/IBM/sarama"

func NewKafkaConfig() *sarama.Config {
    config := sarama.NewConfig()
    config.Version = sarama.V3_4_0_0

    // Producer settings
    config.Producer.Return.Successes = true
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5

    // Consumer settings
    config.Consumer.Return.Errors = true

    return config
}
