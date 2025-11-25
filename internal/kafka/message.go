package kafka

type NotificationMessage struct {
    Email   string `json:"email"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}
