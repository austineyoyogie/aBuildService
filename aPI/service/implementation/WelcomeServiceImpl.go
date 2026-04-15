package implementation

import (
	"aIBuildService/aPI/models"
	"aIBuildService/aPI/service"
	"gorm.io/gorm"
)

type WelcomeServiceImpl struct {
	db *gorm.DB
}

func NewWelcomeServiceImpl(db *gorm.DB) service.WelcomeService {
	return &WelcomeServiceImpl{db}
}

func (u WelcomeServiceImpl) Find(Id uint64) (*models.User, error) {
	user := &models.User{}
	err := u.db.Debug().Model(&models.User{}).Where("id = ?", Id).Find(user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}
