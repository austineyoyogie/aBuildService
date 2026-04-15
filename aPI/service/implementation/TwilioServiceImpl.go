package implementation

import (
	"aIBuildService/aPI/models"
	"aIBuildService/aPI/service"
	"fmt"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"gorm.io/gorm"
	"time"
)

type TwilioServiceImpl struct {
	db     *gorm.DB
	client *twilio.RestClient
	sid    string
}

func NewTwilioServiceImpl(db *gorm.DB, accountSid, authToken, verifySid string) service.TwilioService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	return &TwilioServiceImpl{
		db:     db,
		client: client,
		sid:    verifySid,
	}
}
func (tw *TwilioServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	tx := tw.db.Begin()
	user := &models.User{}
	err := tx.Debug().Model(&models.User{}).Where("email = ?", email).Preload("Roles").Find(user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}
func (tw *TwilioServiceImpl) UpdateUser(user *models.User) error {
	tx := tw.db.Begin()
	columns := map[string]interface{}{
		"sms_enabled":      user.SMSEnabled,
		"totp_enabled":     user.TOTPEnabled,
		"is_authenticated": user.TOTPFactorSid,
		"updated_at":       time.Now(),
	}
	err := tx.Debug().Model(&models.User{}).Where("uuid = ?", user.ID).UpdateColumns(columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (tw *TwilioServiceImpl) SendSMSOTP(to string) error {
	params := &verify.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")
	_, err := tw.client.VerifyV2.CreateVerification(tw.sid, params)
	return err
}
func (tw *TwilioServiceImpl) VerifySMSOTP(to, code string) (bool, error) {
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)
	resp, err := tw.client.VerifyV2.CreateVerificationCheck(tw.sid, params)
	if err != nil {
		return false, err
	}
	return *resp.Status == "approved", nil
}
func (tw *TwilioServiceImpl) CreateTOTPFactor(identity, name string) (string, string, error) {
	params := &verify.CreateNewFactorParams{}
	params.SetFriendlyName(name + "'s totp")
	params.SetFactorType("totp")
	resp, err := tw.client.VerifyV2.CreateNewFactor(tw.sid, identity, params)
	if err != nil {
		return "", "", err
	}
	binding, ok := (*resp.Binding).(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("unexpected binding type")
	}
	uri, ok := binding["uri"].(string)
	if !ok {
		return "", "", fmt.Errorf("uri not found in binding or not a string")
	}
	return *resp.Sid, uri, nil
}
func (tw *TwilioServiceImpl) VerifyFactor(factorSid, code string, identity string) (bool, error) {
	params := &verify.UpdateFactorParams{}
	params.SetAuthPayload(code)
	resp, err := tw.client.VerifyV2.UpdateFactor(tw.sid, identity, factorSid, params)
	if err != nil {
		return false, err
	}
	return *resp.Status == "verified", nil
}
func (tw *TwilioServiceImpl) CreateTOTPChallenge(factorSid string, code string, identity string) (string, error) {
	params := &verify.CreateChallengeParams{}
	params.SetAuthPayload(code)
	params.SetFactorSid(factorSid)
	resp, err := tw.client.VerifyV2.CreateChallenge(tw.sid, identity, params)
	if err != nil {
		return "", err
	}
	return *resp.Sid, nil
}
