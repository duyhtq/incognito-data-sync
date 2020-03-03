package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/duyhtq/incognito-data-sync/models"
)

type Transaction struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) *Transaction {
	return &Transaction{db}
}

func (h *Transaction) Create(transaction *models.Transaction) error {
	return errors.Wrap(h.db.Create(transaction).Error, "h.db.Create")
}

func (r *Transaction) Update(h *models.Transaction) error {
	if err := r.db.Save(h).Error; err != nil {
		return errors.Wrap(err, "h.db.Save")
	}
	return nil
}
