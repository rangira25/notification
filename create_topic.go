package main

import (
    "fmt"
    "github.com/IBM/sarama"
)

func main() {
    brokers := []string{"185.239.209.252:9092"}

    config := sarama.NewConfig()
    config.Version = sarama.V3_4_0_0

    admin, err := sarama.NewClusterAdmin(brokers, config)
    if err != nil {
        panic(err)
    }
    defer admin.Close()

    topicName := "notifications"

    topicDetail := &sarama.TopicDetail{
        NumPartitions:     1,
        ReplicationFactor: 1,
    }

    err = admin.CreateTopic(topicName, topicDetail, false)
    if err != nil {
        panic(err)
    }

    fmt.Println("Topic created successfully:", topicName)
}
