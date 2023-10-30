package models

import (
	"gorm.io/gorm"
)

type Item struct {
	// Standard fields: ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model

	// Item specific fields
	Name     string  `gorm:"not null;check:name <> ''"`
	Quantity float32 `gorm:"not null;check:quantity >= 0.0"`
	Uom      string

	// 1:1 association with Store object
	Store Store

	// 1:n association with ShoppingList
	ShoppingListID uint
}
