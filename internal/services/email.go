package services

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"time"
)

// EmailService sends emails. Replace smtp usage with a provider if desired.
type EmailService struct {
	host string
	port int
	user string
	pass string
}

func NewEmailService(host string, portStr string, user, pass string) *EmailService {
	p, _ := strconv.Atoi(portStr)
	return &EmailService{
		host: host,
		port: p,
		user: user,
		pass: pass,
	}
}

// Generic email sender for queue workers or other callers.
// This keeps backward compatibility by reusing the same SMTP logic.
func (s *EmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	return s.send(ctx, to, subject, body)
}

// Specific welcome-email function (still works exactly the same).
func (s *EmailService) SendWelcome(ctx context.Context, to, subject, body string) error {
	return s.send(ctx, to, subject, body)
}

// Internal reusable SMTP sending logic.
// Both SendEmail and SendWelcome call this.
func (s *EmailService) send(ctx context.Context, to, subject, body string) error {
	if s.host == "" || s.port == 0 {
		log.Println("EmailService is in NOOP mode, skipping email send.")
		return nil
	}

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	auth := smtp.PlainAuth("", s.user, s.pass, s.host)

	done := make(chan error, 1)
	go func() {
		err := smtp.SendMail(addr, auth, s.user, []string{to}, msg)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		if err != nil {
			log.Printf("smtp send error: %v\n", err)
			return err
		}
		return nil
	case <-time.After(30 * time.Second):
		return fmt.Errorf("smtp send timeout")
	}
}

// For local dev/testing
func NewNoopEmailService() *EmailService {
	return &EmailService{}
}
