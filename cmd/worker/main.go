package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    "github.com/rangira25/notification/internal/kafka"
    "github.com/rangira25/notification/internal/services"
)

func main() {

    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  .env file not found, using system environment variables.")
    }

    // Read Email configs
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
    smtpUser := os.Getenv("SMTP_USER")
    smtpPass := os.Getenv("SMTP_PASS")

    if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
        log.Fatal("❌ Missing SMTP_* environment variables in .env")
    }

    // Kafka configs
    brokers := []string{os.Getenv("KAFKA_BROKER")}
    if brokers[0] == "" {
        brokers = []string{"185.***.***.***:9092"} // fallback
    }
    topic := "notifications"

    // Initialize Email Service
    emailSvc := services.NewEmailService(smtpHost, smtpPort, smtpUser, smtpPass)

    kafka.StartConsumerWithHandler(brokers, topic, func(msgBytes []byte) {
        var msg kafka.NotificationMessage

        // Parse the structured JSON message
        if err := json.Unmarshal(msgBytes, &msg); err != nil {
            fmt.Println("❌ Invalid JSON message:", err)
            return
        }

        fmt.Printf("📨 Sending email to %s (subject: %s)\n", msg.Email, msg.Subject)

        // Context with timeout for SMTP
        ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
        defer cancel()

        // Updated call — now correct
        err := emailSvc.SendEmail(ctx, msg.Email, msg.Subject, msg.Body)
        if err != nil {
            log.Println("❌ Email failed:", err)
        } else {
            fmt.Println("✅ Email sent!")
        }
    })

    select {} // keeps worker alive
}
