package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GenerateISOString() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")
}

// Base contains common columns for all tables
type Base struct {
	ID        uint      `gorm:"primaryKey; autoIncrement" json:"id"`
	UUID      uuid.UUID `json:"_id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt string    `json:"-"`
	UpdatedAt string    `json:"-"`
}

// BeforeCreate will set Base struct before every insert
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.UUID = uuid.New()

	// generate timestamps
	t := GenerateISOString()
	base.CreatedAt, base.UpdatedAt = t, t
	return nil
}

// AfterUpdate will update the Base struct after every update
func (base *Base) AfterUpdate(tx *gorm.DB) error {

	// update timestamps
	base.UpdatedAt = GenerateISOString()
	return nil
}
