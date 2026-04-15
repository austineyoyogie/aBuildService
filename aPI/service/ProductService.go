package service

import (
	"aIBuildService/aPI/models/products"
)

type ProductService interface {
	Save(product *products.Product) (*products.Product, error)
	Exists(name string) (*products.Product, error)
	Find(int64) (*products.Product, error)
	FindAll() ([]*products.Product, error)
	Update(product *products.Product) error
	Delete(int64) error
	AddToProductCategory(*products.ProductCategory) (*products.ProductCategory, error)
}
