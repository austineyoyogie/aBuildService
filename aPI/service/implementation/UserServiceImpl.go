package implementation

import (
	"aIBuildService/aPI/messages"
	"aIBuildService/aPI/models"
	"aIBuildService/aPI/service"
	"aIBuildService/aPI/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserServiceImpl struct {
	db *gorm.DB
}

func NewUserServiceImpl(db *gorm.DB) service.UserService {
	return &UserServiceImpl{db}
}

func (u *UserServiceImpl) AddToUserRole(userRole *models.UserRole) (*models.UserRole, error) {
	var user models.User
	tx := u.db.Begin()
	tx.Debug().Model(&models.User{}).Where("id = ?", userRole.UserId).First(&user)

	columns := map[string]interface{}{
		"user_id": user.ID,
		"role_id": true,
	}
	err := tx.Debug().Model(&models.UserRole{}).Create(&columns).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return userRole, tx.Commit().Error
}
func (u *UserServiceImpl) CreateUser(user *models.User) (*models.User, error) {
	tx := u.db.Begin()
	user.FirstName = utils.Escape(user.FirstName)
	user.LastName = utils.Escape(user.LastName)
	user.Email = utils.IsToLower(user.Email)
	PhoneNumber := user.PhoneNumber
	user.PhoneNumber = PhoneNumber
	hash, _ := utils.HashPassword(user.Password)
	user.Password = string(hash)
	token := utils.RandomString(255)
	user.VerifyToken = token
	user.PasswordLastChangedAt = time.Now().UTC()
	err := tx.Debug().Model(&models.User{}).Create(user).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	subject := "Account Email Verification"
	sendEmail := messages.Deliver([]string{user.Email}, subject)
	verification := fmt.Sprintf("http://localhost:4200/v3/verify/account?uuid=%s&token=%s", user.UUID, user.VerifyToken)
	//verification := fmt.Sprintf("http://localhost:4200/v3/verify/account?token=%s", user.UUID, user.VerifyToken)
	sendEmail.EmailTemplate("aPI/messages/verification.html", verification)
	return user, tx.Commit().Error
}
func (u *UserServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	tx := u.db.Begin()
	user := &models.User{}
	err := tx.Debug().Model(&models.User{}).Where("email = ?", email).Preload("Roles").Find(user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}
func (u *UserServiceImpl) VerifyUserByEmail(email string) (*models.User, error) {
	tx := u.db.Begin()
	user := &models.User{}
	err := tx.Debug().Model(&models.User{}).Where("email = ?", email).Find(user).Where("verified = ?", true).Find(user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}
func (u *UserServiceImpl) IsEnabledUser(email string) (*models.User, error) {
	tx := u.db.Begin()
	user := &models.User{}
	err := tx.Debug().Model(&models.User{}).Where("email = ?", email).Find(user).Where("enabled = ?", true).Find(user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}
func (u *UserServiceImpl) IsDisabledUser(email string) (*models.User, error) {
	tx := u.db.Begin()
	user := &models.User{}
	err := tx.Debug().Model(&models.User{}).Where("email = ?", email).Find(user).Where("disabled = ?", true).Find(user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}
func (u *UserServiceImpl) GetUserById(UUID string) (*models.User, error) {
	tx := u.db.Begin()
	user := &models.User{}
	err := tx.Debug().Model(&models.User{}).Where("uuid = ?", UUID).Preload("Roles").Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}
func (u *UserServiceImpl) GetAllUsers() ([]*models.User, error) {
	tx := u.db.Begin()
	var users []*models.User
	err := tx.Debug().Model(&models.User{}).Preload("Roles").Find(&users).Error
	return users, err
}
func (u *UserServiceImpl) UpdateUserById(user *models.User) error {
	tx := u.db.Begin()
	columns := map[string]interface{}{
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"phone_number": user.PhoneNumber,
		"updated_at":   time.Now(),
	}
	err := tx.Debug().Model(&models.User{}).Where("uuid = ?", user.UUID).UpdateColumns(columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (u *UserServiceImpl) DeleteUserById(UUID string) error {
	tx := u.db.Begin()
	columns := map[string]interface{}{
		"verified":    false,
		"enabled":     false,
		"disabled":    true,
		"disabled_at": time.Now(),
	}
	err := tx.Debug().Model(&models.User{}).Where("uuid = ?", UUID).UpdateColumns(&columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (u *UserServiceImpl) GetUserTokenEmailVerification(uuid string, token string) (*models.User, error) {
	tx := u.db.Begin()
	user := &models.User{}
	err := tx.Debug().Model(&models.User{}).Where("uuid = ?", uuid).Where("verify_token = ?", token).Find(user).Error
	return user, err
}
func (u *UserServiceImpl) UpdateUserEmailVerification(user *models.User) error {
	tx := u.db.Begin()
	columns := map[string]interface{}{
		"verify_token": "",
		"verified":     true,
		"enabled":      true,
	}
	err := tx.Debug().Model(&models.User{}).Where("email = ?", user.Email).UpdateColumns(&columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (u *UserServiceImpl) ClearRefreshToken(token string) error {
	tx := u.db.Begin()
	refreshToken := &models.RefreshToken{}
	err := tx.Debug().Model(&models.RefreshToken{}).Where("email = ?", token).Delete(&refreshToken).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (u *UserServiceImpl) StoreRefreshToken(token *models.RefreshToken) (*models.RefreshToken, error) {
	tx := u.db.Begin()
	err := tx.Debug().Model(&models.RefreshToken{}).Create(token).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return token, tx.Commit().Error
}
func (u *UserServiceImpl) GetRefreshToken(refresh *models.RefreshToken) (*models.RefreshToken, error) {
	tx := u.db.Begin()
	refreshToken := &models.RefreshToken{}
	err := tx.Debug().Model(&models.RefreshToken{}).Where("uuid = ?", refresh.UUID).Where("email = ?", refresh.UserEmail).Find(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return refreshToken, err
}
func (u *UserServiceImpl) RevokeRefreshToken(revokeId string) error {
	tx := u.db.Begin()
	columns := map[string]interface{}{
		"is_revoked": true,
	}
	err := tx.Debug().Model(&models.RefreshToken{}).Where("id = ?", revokeId).UpdateColumns(&columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (u *UserServiceImpl) DeleteRefreshToken(UUID string, refreshEmail string) error {
	tx := u.db.Begin()
	refreshToken := &models.RefreshToken{}
	err := tx.Debug().Model(&models.RefreshToken{}).Where("uuid = ?", UUID).Where("email = ?", refreshEmail).Delete(&refreshToken).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (u *UserServiceImpl) ResetUserPasswordToken(user *models.User) error {
	tx := u.db.Begin()
	user.Email = utils.IsToLower(user.Email)
	Token := utils.RandomString(255)
	user.VerifyToken = Token
	columns := map[string]interface{}{
		"verify_token": user.VerifyToken,
	}
	err := tx.Debug().Model(&models.User{}).Where("email = ?", user.Email).UpdateColumns(&columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	//subject := "Reset Password Email Verification"
	//msg := messages.Deliver([]string{user.Email}, subject)
	//verification := fmt.Sprintf("http://localhost:4200/v3/verify/password?token=%s", user.VerifyToken)
	//msg.EmailTemplate("aPI/messages/resetpassword.html", verification)
	return tx.Commit().Error
}
func (u *UserServiceImpl) VerifyUserPasswordToken(token string) (*models.User, error) {
	tx := u.db.Begin()
	user := &models.User{}
	err := tx.Debug().Model(&models.User{}).Where("verify_token = ?", token).Preload("Roles").Find(user).Error
	return user, err
}
func (u *UserServiceImpl) UpdateUserPasswordToken(user *models.User) error {
	tx := u.db.Begin()
	user.VerifyToken = utils.IsToLower(user.VerifyToken)
	hash, _ := utils.HashPassword(user.Password)
	user.Password = string(hash)
	columns := map[string]interface{}{
		"verify_token": "",
		"password":     user.Password,
	}
	err := tx.Debug().Model(&models.User{}).Where("verify_token = ?", user.VerifyToken).UpdateColumns(&columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	subject := "Account Reset Password Successfully"
	msg := messages.Deliver([]string{user.Email}, subject)
	msg.EmailTemplate("aPI/messages/success.html", subject)
	return tx.Commit().Error
}
func (u *UserServiceImpl) IncrementFailedLoginAttempts(user *models.User) error {
	tx := u.db.Begin()
	columns := map[string]interface{}{
		"failed_attempts":          user.FailedAttempts + 1,
		"last_failed_attempt_time": time.Now(),
	}
	err := tx.Debug().Model(&models.User{}).Where("email = ?", user.Email).UpdateColumns(&columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (u *UserServiceImpl) ResetFailedLoginAttempts(user *models.User) error {
	tx := u.db.Begin()
	columns := map[string]interface{}{
		"failed_attempts": 0,
	}
	err := tx.Debug().Model(&models.User{}).Where("email = ?", user.Email).UpdateColumns(&columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
