package mailersend

import (
	"context"
	"fmt"
	"time"

	"github.com/firzatullahd/golang-template/internal/user/model"
	"github.com/firzatullahd/golang-template/utils/logger"
	mailersendSDK "github.com/mailersend/mailersend-go"
)

type Client struct {
	templateOTP      string
	emailfrom        string
	mailersendClient *mailersendSDK.Mailersend
}

func NewClient(apiKey, emailfrom, templateOtp string) *Client {
	return &Client{
		templateOTP:      templateOtp,
		emailfrom:        emailfrom,
		mailersendClient: mailersendSDK.NewMailersend(apiKey),
	}
}

func (c *Client) SendEmail(ctx context.Context, input model.EmailPayload) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subject := "OTP"
	from := mailersendSDK.From{
		Name:  "Firza Playground",
		Email: c.emailfrom,
	}

	recipients := []mailersendSDK.Recipient{
		{
			Name:  input.Name,
			Email: input.Email,
		},
	}

	variables := []mailersendSDK.Personalization{
		{
			Email: input.Email,
			Data: map[string]any{
				"verification_code": input.VerificationCode,
			},
		},
	}

	message := c.mailersendClient.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetTemplateID(c.templateOTP)
	message.SetPersonalization(variables)

	resp, err := c.mailersendClient.Email.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("Send Email failed %w", err)
	}

	logger.Log.Info("Email sent successfully ", resp.Header.Get("X-Message-Id"))

	return nil
}
