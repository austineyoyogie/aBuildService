package service

import "aIBuildService/aPI/models/products"

type CategoryService interface {
	Save(category *products.Category) (*products.Category, error)
	Exists(name string) (*products.Category, error)
	Find(int64) (*products.Category, error)
	FindAll() ([]*products.Category, error)
	Update(category *products.Category) error
	Delete(int64) error
}
