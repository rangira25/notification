package main

import (
    "github.com/gin-gonic/gin"
    "github.com/rangira25/notification/internal/kafka"
    "net/http"
)

func main() {
    r := gin.Default()

    brokers := []string{"185.239.209.252:9092"}
    topic := "notifications"

    producer, err := kafka.NewProducer(brokers)
    if err != nil {
        panic(err)
    }

    r.POST("/send", func(c *gin.Context) {
    var payload struct {
        Email   string `json:"email"`
        Subject string `json:"subject"`
        Body    string `json:"body"`
    }

    if err := c.BindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    msg := kafka.NotificationMessage{
        Email:   payload.Email,
        Subject: payload.Subject,
        Body:    payload.Body,
    }

    if err := producer.SendJSON(topic, msg); err != nil {
        c.JSON(500, gin.H{"status": "failed", "error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"status": "queued"})
})
}