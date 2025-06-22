package notification

import (
	"fmt"

	"github.com/zaibon/surveilsense/proto"
)

// SMSNotifier implements Notifier for SMS notifications
// Integrate with your SMS provider (e.g., Twilio) in the SendSMS method

type SMSNotifier struct {
	AccountSID string // For Twilio or similar
	AuthToken  string
	From       string
	To         []string
	// Add other provider-specific fields as needed
}

func (s *SMSNotifier) Notify(event *proto.DetectionEvent) error {
	body := fmt.Sprintf("SurveilSense Alert: Detection on camera %s at %d. Detections: %d", event.CameraId, event.Timestamp, len(event.Detections))
	for _, recipient := range s.To {
		if err := s.SendSMS(recipient, body); err != nil {
			return fmt.Errorf("failed to send SMS to %s: %w", recipient, err)
		}
	}
	return nil
}

// SendSMS is a placeholder for actual SMS sending logic (e.g., via Twilio API)
func (s *SMSNotifier) SendSMS(to, body string) error {
	// TODO: Integrate with SMS provider API
	fmt.Printf("[SMS] Would send to %s: %s\n", to, body)
	return nil
}
