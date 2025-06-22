package notification

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/zaibon/surveilsense/proto"
)

// EmailNotifier implements Notifier for email notifications
// It uses SMTP to send emails

type EmailNotifier struct {
	SMTPServer string // e.g. smtp.gmail.com:587
	Username   string
	Password   string
	From       string
	To         []string
	UseTLS     bool
}

func (e *EmailNotifier) Notify(event *proto.DetectionEvent) error {
	subject := fmt.Sprintf("SurveilSense Alert: Detection on camera %s", event.CameraId)
	body := fmt.Sprintf("Detection event for camera %s at %d. Detections: %d", event.CameraId, event.Timestamp, len(event.Detections))
	msg := "From: " + e.From + "\r\n" +
		"To: " + strings.Join(e.To, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n"

	auth := smtp.PlainAuth("", e.Username, e.Password, strings.Split(e.SMTPServer, ":")[0])

	if e.UseTLS {
		// Connect using TLS
		c, err := smtp.Dial(e.SMTPServer)
		if err != nil {
			return fmt.Errorf("failed to dial SMTP server: %w", err)
		}
		defer c.Close()
		if err = c.StartTLS(&tls.Config{ServerName: strings.Split(e.SMTPServer, ":")[0]}); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
		if err = c.Mail(e.From); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}
		for _, addr := range e.To {
			if err = c.Rcpt(addr); err != nil {
				return fmt.Errorf("failed to set recipient: %w", err)
			}
		}
		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("failed to get data writer: %w", err)
		}
		_, err = w.Write([]byte(msg))
		if err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
		w.Close()
		return c.Quit()
	}

	// Non-TLS (STARTTLS is handled automatically by smtp.SendMail)
	return smtp.SendMail(e.SMTPServer, auth, e.From, e.To, []byte(msg))
}
