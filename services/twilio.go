package services

import (
	"fmt"
	"sheleads-backend/config"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioService struct {
	client *twilio.RestClient
}

func NewTwilioService() *TwilioService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AppConfig.TwilioAccountSID,
		Password: config.AppConfig.TwilioAuthToken,
	})
	return &TwilioService{client: client}
}

func (t *TwilioService) SendBurnoutAlert(emergencyContact, userName string, hoursWorked float64) error {
	message := fmt.Sprintf(
		"ALERT: %s may be experiencing burnout. They have logged %.1f hours of unpaid care work in the past week, which exceeds healthy limits. Please check on their wellbeing. - SheLeads Care Platform",
		userName,
		hoursWorked,
	)

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(emergencyContact)
	params.SetMessagingServiceSid(config.AppConfig.TwilioMessagingServiceSID)
	params.SetBody(message)

	_, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}

func (t *TwilioService) SendMonthlyReport(phoneNumber, userName string, totalHours float64, totalPoints int) error {
	message := fmt.Sprintf(
		"Monthly Report for %s:\n\nTotal Hours: %.1f\nTotal Points: %d\n\nYour unpaid care work is valuable and recognized. Keep up the amazing work! - SheLeads",
		userName,
		totalHours,
		totalPoints,
	)

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetMessagingServiceSid(config.AppConfig.TwilioMessagingServiceSID)
	params.SetBody(message)

	_, err := t.client.Api.CreateMessage(params)
	return err
}

