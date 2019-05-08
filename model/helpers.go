package model

import (
	"regexp"

	"github.com/jinzhu/gorm"
)

// AutoMigrate ...
func AutoMigrate(db *gorm.DB) *gorm.DB {
	return db.AutoMigrate(&Player{}, &Game{}, &EventContainer{})
}

// IsDuplicateError check for a "Duplicate entry" error.
func IsDuplicateError(err error) bool {
	m, _ := regexp.MatchString("duplicate key value violates unique constraint", err.Error())

	return m
}
