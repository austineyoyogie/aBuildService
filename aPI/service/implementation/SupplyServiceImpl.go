package implementation

import (
	"aIBuildService/aPI/models/products"
	"aIBuildService/aPI/service"
	"aIBuildService/aPI/utils"
	"gorm.io/gorm"
	"time"
)

type SupplyServiceImpl struct {
	db *gorm.DB
}

func NewSupplyServiceImpl(db *gorm.DB) service.SupplyService {
	return &SupplyServiceImpl{db}
}

func (s *SupplyServiceImpl) Save(supply *products.Supply) (*products.Supply, error) {
	tx := s.db.Begin()

	supply.Name = utils.Escape(supply.Name)
	supply.Country = utils.Escape(supply.Country)
	SUID := utils.RandomUpperString(4)
	supply.SupplyID = "BV-" + SUID + "-CH"

	err := tx.Debug().Model(&products.Supply{}).Create(supply).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return supply, tx.Commit().Error
}

func (s *SupplyServiceImpl) Find(supplyId int64) (*products.Supply, error) {
	supply := &products.Supply{}
	err := s.db.Debug().Model(&products.Supply{}).Where("id = ?", supplyId).Find(&supply).Error
	if err != nil {
		return nil, err
	}
	return supply, err
}

func (s *SupplyServiceImpl) FindAll() ([]*products.Supply, error) {
	var roles []*products.Supply
	err := s.db.Debug().Model(&products.Supply{}).Find(&roles).Error
	return roles, err
}

func (s *SupplyServiceImpl) Exists(name string) (*products.Supply, error) {
	supply := &products.Supply{}
	err := s.db.Debug().Model(&products.Supply{}).Where("name = ?", name).Find(supply).Error
	if err != nil {
		return nil, err
	}
	return supply, err
}

func (s *SupplyServiceImpl) Update(supply *products.Supply) error {
	tx := s.db.Begin()
	columns := map[string]interface{}{
		"name":           supply.Name,
		"country":        supply.Country,
		"supply_type":    supply.SupplyType,
		"current_order":  supply.CurrentOrder,
		"order_received": supply.OrderReceived,
		"documents":      supply.Document,
		"updated_at":     time.Now(),
	}
	err := tx.Debug().Model(&products.Supply{}).Where("id = ?", supply.ID).UpdateColumns(columns).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *SupplyServiceImpl) Delete(supplyId int64) error {
	tx := s.db.Begin()
	err := tx.Debug().Model(&products.Supply{}).Where("id = ?", supplyId).Delete(&products.Supply{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
