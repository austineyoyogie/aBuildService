package service

import "aIBuildService/aPI/models"

type RoleService interface {
	Save(*models.Role) (*models.Role, error)
	Find(int64) (*models.Role, error)
	FindAll() (models.Role, error)
	Exists(name string) (*models.Role, error)
	Update(*models.Role) error
	Delete(int64) error
}
