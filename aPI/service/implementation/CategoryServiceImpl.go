package implementation

import (
	"aIBuildService/aPI/models/products"
	"aIBuildService/aPI/service"
	"aIBuildService/aPI/utils"
	"gorm.io/gorm"
	"time"
)

type CategoryServiceImpl struct {
	db *gorm.DB
}

func NewCategoryServiceImpl(db *gorm.DB) service.CategoryService {
	return &CategoryServiceImpl{db}
}

func (c *CategoryServiceImpl) Save(category *products.Category) (*products.Category, error) {
	tx := c.db.Begin()
	category.Name = utils.Escape(category.Name)
	category.Description = utils.Escape(category.Description)
	err := tx.Debug().Model(&products.Category{}).Create(category).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return category, tx.Commit().Error
}

func (c *CategoryServiceImpl) Find(categoryId int64) (*products.Category, error) {
	category := &products.Category{}
	err := c.db.Debug().Model(&products.Category{}).Where("id = ?", categoryId).Preload("Products").Find(&category).Error
	if err != nil {
		return nil, err
	}
	return category, err
}

func (c *CategoryServiceImpl) FindAll() ([]*products.Category, error) {
	var categories []*products.Category
	err := c.db.Debug().Model(&products.Category{}).Find(&categories).Error
	return categories, err
}

func (c *CategoryServiceImpl) Exists(name string) (*products.Category, error) {
	category := &products.Category{}
	err := c.db.Debug().Model(&products.Category{}).Where("name = ?", name).Preload("Products").Find(category).Error
	if err != nil {
		return nil, err
	}
	return category, err
}

func (c *CategoryServiceImpl) Update(category *products.Category) error {
	tx := c.db.Begin()
	columns := map[string]interface{}{
		"name":        category.Name,
		"description": category.Description,
		"picture":     category.Picture,
		"active":      category.Active,
		"enabled":     category.Enabled,
		"deleted":     category.Deleted,
		"updated_at":  time.Now(),
	}
	err := tx.Debug().Model(&products.Supply{}).Where("id = ?", category.ID).UpdateColumns(columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (c *CategoryServiceImpl) Delete(categoryId int64) error {
	tx := c.db.Begin()
	columns := map[string]interface{}{
		"active":     false,
		"enabled":    false,
		"deleted":    true,
		"updated_at": time.Now(),
	}
	err := tx.Debug().Model(&products.Category{}).Where("id = ?", categoryId).UpdateColumns(&columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
