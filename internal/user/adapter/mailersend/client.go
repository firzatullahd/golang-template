package mailersend

import (
	"context"
	"fmt"
	"time"

	"github.com/firzatullahd/golang-template/internal/user/model"
	mailersendSDK "github.com/mailersend/mailersend-go"
)

type Client struct {
	APIKey           string
	mailersendClient *mailersendSDK.Mailersend
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:           apiKey,
		mailersendClient: mailersendSDK.NewMailersend(apiKey),
	}
}

func (c *Client) SendEmail(ctx context.Context, input model.EmailPayload) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subject := "OTP"
	from := mailersendSDK.From{
		Name:  "Company Name",
		Email: "Company Email",
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
	message.SetTemplateID("template-id")
	message.SetPersonalization(variables)

	resp, err := c.mailersendClient.Email.Send(ctx, message)
	if err != nil {
		return err
	}

	fmt.Println("Email sent successfully", resp.Header.Get("X-Message-Id"))

	return nil
}
