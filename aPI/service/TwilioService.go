package service

import "aIBuildService/aPI/models"

type TwilioService interface {
	GetUserByEmail(string) (*models.User, error)
	UpdateUser(user *models.User) error
	SendSMSOTP(to string) error
	VerifySMSOTP(to, code string) (bool, error)
	CreateTOTPFactor(identity, name string) (string, string, error)
	VerifyFactor(factorSid, code string, identity string) (bool, error)
	CreateTOTPChallenge(factorSid string, code string, identity string) (string, error)
}
