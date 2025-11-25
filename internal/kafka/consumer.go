package kafka

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func StartConsumerWithHandler(brokers []string, topic string, handler func([]byte)) {
	config := NewKafkaConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range partitions {
		go func(partition int32) {
			pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Println("Error:", err)
				return
			}

			fmt.Println("Listening on partition", partition)

			for msg := range pc.Messages() {
				handler(msg.Value)
			}
		}(p)
	}
}
