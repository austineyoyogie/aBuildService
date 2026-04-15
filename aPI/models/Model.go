package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type CustomTime time.Time

func (ct CustomTime) Value() (driver.Value, error) {
	return time.Time(ct), nil
}
func (ct *CustomTime) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan CustomTime")
	}
	*ct = CustomTime(t)
	return nil
}

type ChangedAt struct {
	PasswordLastChangedAt time.Time `gorm:json:"created_at"`
}

type Model struct {
	CreatedAt time.Time `gorm:json:"created_at"`
	UpdatedAt time.Time `gorm:json:"updated_at"`
}
