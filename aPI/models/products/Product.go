package products

import (
	"database/sql"
	"gorm.io/gorm"
)

type Product struct {
	ID           int64        `gorm:"column:id;primary_key;auto_increment" json:"id"`
	UserID       uint         `gorm:"column:user_id;size:45;not null;" json:"user_id"`
	SkuID        string       `gorm:"column:sku_id;size:45;not null;unique" json:"sku_id"`
	SupplierID   uint         `gorm:"column:supplier_id;size:45;not null;" json:"supplier_id"`
	CategoryID   uint         `gorm:"column:category_id;size:45;not null;" json:"category_id"`
	Name         string       `gorm:"column:name;size:255;not null;unique" validate:"required,min=2,max=255" json:"name"`
	Description  string       `gorm:"column:description;size:255;not null;" validate:"required,min=2,max=255" json:"description"`
	Picture      string       `gorm:"column:picture;size:255;not null;" validate:"required,min=2,max=255" json:"picture"`
	Quantity     int64        `gorm:"column:quantities;size:45;not null;" validate:"required" json:"quantities"`
	UnitsInStock int64        `gorm:"column:units_in_stock;size:45;not null;" validate:"required" json:"units_in_stock"`
	Active       sql.NullBool `gorm:"column:active;default:true" json:"active"`
	Enabled      sql.NullBool `gorm:"column:enabled;default:false" json:"enabled"`
	Deleted      sql.NullBool `gorm:"column:deleted;default:false" json:"deleted"`
	Categories   *[]Category  `gorm:"many2many:product_categories" json:"category"`
	Model
}

type Category struct {
	ID          int64        `gorm:"column:id;primary_key;auto_increment" json:"id"`
	UserID      uint         `gorm:"column:user_id;size:45;not null;" json:"user_id"`
	Name        string       `gorm:"column:name;size:255;not null;unique" validate:"required,min=2,max=255" json:"name"`
	Description string       `gorm:"column:description;size:255;not null;" validate:"required,min=2,max=255" json:"description"`
	Picture     string       `gorm:"column:picture;size:255;not null;" validate:"required,min=2,max=255" json:"picture"`
	Active      sql.NullBool `gorm:"column:active;default:true" json:"active"`
	Enabled     sql.NullBool `gorm:"column:enabled;default:false" json:"enabled"`
	Deleted     sql.NullBool `gorm:"column:deleted;default:false" json:"deleted"`
	Products    *[]Product   `gorm:"many2many:product_categories" json:"product"`
	Model
}

type ProductCategory struct {
	ID         int64 `gorm:"column:id;primary_key;auto_increment" json:"id"`
	ProductId  uint
	CategoryId uint
	Model
}

type Supply struct {
	ID            int64  `gorm:"column:id;primary_key;auto_increment" json:"id"`
	UserID        uint   `gorm:"column:user_id;size:45;not null;" json:"user_id"`
	SupplyID      string `gorm:"column:supply_id;size:45;not null;unique" json:"supply_id"`
	Name          string `gorm:"column:name;size:255;not null;" validate:"required,min=2,max=255" json:"name"`
	Country       string `gorm:"column:country;size:255;not null;" validate:"required,min=2,max=255" json:"country"`
	SupplyType    string `gorm:"column:supply_type;size:255;not null;" validate:"required,min=2,max=255" json:"supply_type"`
	CurrentOrder  int64  `gorm:"column:current_order;not null;" validate:"required" json:"current_order"`
	OrderReceived int64  `gorm:"column:order_received;not null;" validate:"required" json:"order_received"`
	Document      string `gorm:"column:documents;size:255;not null;" validate:"required,min=2,max=255" json:"documents"`
	Model
}

func ProductMigration(db *gorm.DB) {
	err := db.Debug().Migrator().AutoMigrate(&Product{}, &Category{}, &ProductCategory{}, &Supply{})
	//err := db.Debug().Migrator().DropTable(&Product{}, &Category{}, &ProductCategory{}, &Supply{})
	if err != nil {
		return
	}
}
