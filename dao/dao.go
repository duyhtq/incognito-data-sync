package dao

import (
	"github.com/duyhtq/incognito-data-sync/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	transactionTable = []interface{}{(*models.Transaction)(nil)}
)

func AutoMigrate(db *gorm.DB) error {
	allTables := make([]interface{}, 0, len(transactionTable))

	allTables = append(allTables, transactionTable...)

	if err := db.AutoMigrate(allTables...).Error; err != nil {
		return errors.Wrap(err, "db.AutoMigrate")
	}

	return nil
}
