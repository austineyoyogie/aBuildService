package implementation

import (
	"aIBuildService/aPI/models/products"
	"aIBuildService/aPI/service"
	"aIBuildService/aPI/utils"
	"gorm.io/gorm"
	"time"
)

type ProductServiceImpl struct {
	db *gorm.DB
}

func NewProductServiceImpl(db *gorm.DB) service.ProductService {
	return &ProductServiceImpl{db}
}

func (p *ProductServiceImpl) AddToProductCategory(productCategory *products.ProductCategory) (*products.ProductCategory, error) {
	var product products.Product
	tx := p.db.Begin()
	tx.Debug().Model(&products.Product{}).Where("id = ?", productCategory.ProductId).First(&product)

	columns := map[string]interface{}{
		"product_id":  product.ID,
		"category_id": true,
	}
	err := tx.Debug().Model(&products.ProductCategory{}).Create(&columns).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return productCategory, tx.Commit().Error
}

func (p *ProductServiceImpl) Save(product *products.Product) (*products.Product, error) {
	tx := p.db.Begin()
	product.Name = utils.Escape(product.Name)
	product.Description = utils.Escape(product.Description)

	SUKID := utils.RandomUpperString(4)
	product.SkuID = "SUK-" + SUKID + "-MV"
	err := tx.Debug().Model(&products.Product{}).Create(product).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return product, tx.Commit().Error
}

func (p *ProductServiceImpl) Find(productId int64) (*products.Product, error) {
	product := &products.Product{}
	err := p.db.Debug().Model(&products.Product{}).Where("id = ?", productId).Preload("Categories").Find(&product).Error
	if err != nil {
		return nil, err
	}
	return product, err
}

func (p *ProductServiceImpl) FindAll() ([]*products.Product, error) {
	var product []*products.Product
	err := p.db.Debug().Model(&products.Product{}).Find(&product).Error
	return product, err
}

func (p *ProductServiceImpl) Exists(name string) (*products.Product, error) {
	product := &products.Product{}
	err := p.db.Debug().Model(&products.Product{}).Where("name = ?", name).Preload("Categories").Find(product).Error
	if err != nil {
		return nil, err
	}
	return product, err
}

func (p *ProductServiceImpl) Update(product *products.Product) error {
	tx := p.db.Begin()
	columns := map[string]interface{}{
		"name":        product.Name,
		"description": product.Description,
		"picture":     product.Picture,
		"active":      product.Active,
		"enabled":     product.Enabled,
		"deleted":     product.Deleted,
		"updated_at":  time.Now(),
	}
	err := tx.Debug().Model(&products.Product{}).Where("id = ?", product.ID).UpdateColumns(columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (p *ProductServiceImpl) Delete(productId int64) error {
	tx := p.db.Begin()
	columns := map[string]interface{}{
		"active":     false,
		"enabled":    false,
		"deleted":    true,
		"updated_at": time.Now(),
	}
	err := tx.Debug().Model(&products.Product{}).Where("id = ?", productId).UpdateColumns(&columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
