package services

import (
	"errors"
	"fmt"
	"github.com/WildEgor/e-shop-auth/internal/adapters"
	"github.com/WildEgor/e-shop-auth/internal/configs"
	"github.com/WildEgor/e-shop-auth/internal/utils"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/notifier"
)

var (
	ErrSendSMS   = errors.New("failed to send SMS")
	ErrSendEmail = errors.New("failed to send e-mail")
)

type NotificationService struct {
	config    *configs.OTPConfig
	generator *utils.CodeGenerator
	notifier  *adapters.NotifierAdapter
}

func NewNotificationService(
	config *configs.OTPConfig,
	notifier *adapters.NotifierAdapter,
) *NotificationService {
	return &NotificationService{
		config:    config,
		notifier:  notifier,
		generator: utils.NewCodeGenerator(config.Length),
	}
}

func (ns *NotificationService) GenerateAndSendOTPSms(phone string) (string, error) {
	code := ns.generator.GenShortCode()
	msg := notifier.NewNotification(
		notifier.WithPhoneSettings(notifier.PhoneSettings{
			Phone: phone,
			Text:  fmt.Sprintf("Your code: %s", code),
		}),
	)

	if err := ns.notifier.Notify(&msg); err != nil {
		return "", ErrSendSMS
	}

	return code, nil
}

func (ns *NotificationService) GenerateAndSendEmailConfirm(email string) (string, error) {
	code := ns.generator.GenShortCode()

	m := notifier.NotificationPayload{
		Type: "email",
		EmailSettings: notifier.EmailSettings{
			Email:    email,
			Subject:  "Confirm email",
			Template: "email_confirm",
			Data: struct {
				Code string
			}{
				Code: code,
			},
		},
	}

	if err := ns.notifier.Notify(&m); err != nil {
		return "", ErrSendEmail
	}

	return code, nil
}

func (ns *NotificationService) GenerateAndSendPhoneConfirm(phone string) (string, error) {
	code := ns.generator.GenShortCode()
	msg := notifier.NewNotification(
		notifier.WithPhoneSettings(notifier.PhoneSettings{
			Phone: phone,
			Text:  fmt.Sprintf("Confirm code: %s", code),
		}),
	)

	if err := ns.notifier.Notify(&msg); err != nil {
		return "", ErrSendSMS
	}

	return code, nil
}
