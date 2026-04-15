package models

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

var (
	ErrRoleEmptyName = errors.New("role.permission name can't be empty")
)

type User struct {
	ID                    int64        `gorm:"column:id;primary_key;" json:"id"`
	UUID                  string       `gorm:"column:uuid;size:50;not null" json:"uuid"`
	FirstName             string       `gorm:"column:first_name;size:45;not null;"`
	LastName              string       `gorm:"column:last_name;size:45;not null;"`
	Email                 string       `gorm:"column:email;size:45;not null;unique" `
	Password              string       `gorm:"column:password;size:255;not null;"`
	PhoneNumber           string       `gorm:"column:phone_number;size:45;not null;"`
	VerifyToken           string       `gorm:"column:verify_token;size:255;not null;" json:"verify_token"`
	TOTPFactorSid         string       `gorm:"column:totp_factor_sid;size:255;not null;" json:"totp_factor_sid"`
	SMSEnabled            bool         `gorm:"column:sms_enabled;default:false" json:"sms_enabled"`
	TOTPEnabled           bool         `gorm:"column:totp_enabled;default:false" json:"totp_enabled"`
	IsAuthenticated       bool         `gorm:"column:is_authenticated;default:false" json:"is_authenticated"`
	FailedAttempts        int          `gorm:"column:failed_attempts;default:0" json:"failed_attempts"`
	LastFailedAttemptTime sql.NullTime `gorm:"column:last_failed_attempt_time;" json:"last_failed_attempt_time"`
	PasswordLastChangedAt time.Time    `gorm:"column:password_last_changed_at" json:"password_last_changed_at"`
	Verified              sql.NullBool `gorm:"column:verified;default:false" json:"verified"`
	Enabled               sql.NullBool `gorm:"column:enabled;default:false" json:"enabled"`
	Disabled              sql.NullBool `gorm:"column:disabled;default:false" json:"disabled"`
	DisabledAt            string       `gorm:"column:disabled_at;size:45;not null;" json:"disabled_at"`
	Roles                 *[]Role      `gorm:"many2many:user_roles" json:"roles"`
	Model
}
type Role struct {
	ID          int64   `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Name        string  `gorm:"column:name;size:100;not null;" validate:"required,min=2,max=255" json:"name"`
	Permissions string  `gorm:"column:permissions;size:255;not null;" validate:"required,min=2,max=255" json:"permissions"`
	Users       *[]User `gorm:"many2many:user_roles" json:"user"`
}
type UserRole struct {
	ID     int64 `gorm:"column:id;primary_key;auto_increment" json:"id"`
	UserId int64
	RoleId int64
}

func (user *User) BeforeCreate(*gorm.DB) error {
	tokenID, _ := uuid.NewV7()
	user.UUID = tokenID.String()
	return nil
}

type CreateUserRequest struct {
	FirstName   string `binding:"required" validate:"min=2,max=45" json:"first_name"`
	LastName    string `binding:"required" validate:"min=2,max=45" json:"last_name" `
	Email       string `binding:"required" validate:"email" json:"email"`
	Password    string `binding:"required" validate:"min=8,max=255" json:"password"`
	PhoneNumber string `binding:"required" json:"phone_number"`
}
type UpdateUserRequest struct {
	FirstName   string `binding:"required" validate:"min=2,max=45" json:"first_name"`
	LastName    string `binding:"required" validate:"min=2,max=45" json:"last_name" `
	PhoneNumber string `binding:"required" json:"phone_number"`
}
type LoginCredentials struct {
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}
type LoginUserResponse struct {
	FirstName string  `binding:"required" validate:"min=2,max=45" json:"first_name"`
	LastName  string  `binding:"required" validate:"min=2,max=45" json:"last_name" `
	Roles     []*Role `json:"roles"`
}
type LoginCredentialsResponse struct {
	UUID                  string    `json:"uuid"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  *User     `json:"user"`
}
type ForgotPasswordRequest struct {
	Email string `binding:"required" validate:"email" json:"email"`
}
type ChangePasswordWithTokenRequest struct {
	Token    string `binding:"required" json:"token"`
	Password string `binding:"required" json:"password"`
}
type RefreshToken struct {
	ID           string     `gorm:"column:id;primary_key;size:255" json:"id"`
	UUID         string     `gorm:"column:uuid;size:50;not null" json:"uuid"`
	UserEmail    string     `gorm:"column:email;size:255;not null;" json:"email"`
	RefreshToken string     `gorm:"column:refresh_token;size:512;not null;" json:"refresh_token"`
	IsRevoked    bool       `gorm:"column:is_revoked;default:false" json:"is_revoked"`
	ExpiresAt    CustomTime `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
}
type RefreshTokenRequest struct {
	RefreshToken string `gorm:"refresh_token" json:"refresh_token"`
}
type RefreshTokenResponse struct {
	AccessToken          string    `gorm:"access_token" json:"access_token"`
	AccessTokenExpiresAt time.Time `gorm:"access_token_expires_at" json:"access_token_expires_at"`
}

func AutoMigration(db *gorm.DB) {
	err := db.Debug().Migrator().AutoMigrate(&Role{}, &User{}, &UserRole{}, RefreshToken{})
	//err := db.Debug().Migrator().DropTable(RefreshToken{})
	if err != nil {
		return
	}
}
func (r *Role) Validate() error {
	if r.Name == "" {
		return ErrRoleEmptyName
	}
	return nil
}

// OAuth Golang Tutorial for Authentication
// https://www.youtube.com/watch?v=iHFQyd__2A0
// https://www.youtube.com/watch?v=WvEzw7wTuzE

// https://www.twilio.com/en-us/blog/developers/community/multi-factor-authentication-go-twilio-verify
// https://codevoweb.com/two-factor-authentication-2fa-in-golang/
// Rename the method struct
// https://abrialstha.medium.com/going-with-go-authentication-route-protection-02c52fc11382
