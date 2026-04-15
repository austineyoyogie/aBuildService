package service

import (
	"aIBuildService/aPI/models"
)

type UserService interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	VerifyUserByEmail(string) (*models.User, error)
	IsEnabledUser(string) (*models.User, error)
	IsDisabledUser(string) (*models.User, error)
	GetUserById(string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUserById(user *models.User) error
	DeleteUserById(string) error
	AddToUserRole(*models.UserRole) (*models.UserRole, error)
	GetUserTokenEmailVerification(string, string) (*models.User, error)
	UpdateUserEmailVerification(user *models.User) error
	ClearRefreshToken(string) error
	StoreRefreshToken(token *models.RefreshToken) (*models.RefreshToken, error)
	GetRefreshToken(token *models.RefreshToken) (*models.RefreshToken, error)
	RevokeRefreshToken(string) error
	DeleteRefreshToken(UUID string, email string) error
	ResetUserPasswordToken(user *models.User) error
	VerifyUserPasswordToken(string) (*models.User, error)
	UpdateUserPasswordToken(user *models.User) error
	IncrementFailedLoginAttempts(user *models.User) error
	ResetFailedLoginAttempts(user *models.User) error
}
