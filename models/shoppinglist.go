package models

import (
	"gorm.io/gorm"
)

type ShoppingList struct {
	// Standard fields: ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model

	// ShoppingList specific fields
	Name    string `gorm:"not null;check:name <> ''"`
	Creator string `gorm:"not null;check:creator <> ''"`

	// 1:n collection of Item objects
	Items []Item
}
