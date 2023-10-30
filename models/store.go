package models

import "gorm.io/gorm"

type Store struct {
	// Standard fields: ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model

	// Store specific fields
	Name    string `gorm:"not null"`
	Address string

	// Association with Item
	ItemID uint
}
