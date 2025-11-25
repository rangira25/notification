package kafka

import (
    "encoding/json"

    "github.com/IBM/sarama"
)

type Producer struct {
    producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
    config := NewKafkaConfig()

    p, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }

    return &Producer{producer: p}, nil
}

func (p *Producer) SendJSON(topic string, payload interface{}) error {
    jsonBytes, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    msg := &sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.ByteEncoder(jsonBytes),
    }

    _, _, err = p.producer.SendMessage(msg)
    return err
}
