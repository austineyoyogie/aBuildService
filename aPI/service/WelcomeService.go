package service

import "aIBuildService/aPI/models"

type WelcomeService interface {
	Find(uint64) (*models.User, error)
}
