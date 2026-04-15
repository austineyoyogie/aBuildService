package service

import (
	"aIBuildService/aPI/models/products"
)

type SupplyService interface {
	Save(supply *products.Supply) (*products.Supply, error)
	Exists(name string) (*products.Supply, error)
	Find(int64) (*products.Supply, error)
	FindAll() ([]*products.Supply, error)
	Update(supply *products.Supply) error
	Delete(int64) error
}
