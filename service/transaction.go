package service

import (
	config "github.com/duyhtq/incognito-data-sync/config"
	"github.com/duyhtq/incognito-data-sync/databases/postgresql"
	"go.uber.org/zap"
)

type Transaction struct {
	conf        *config.Config
	logger      *zap.Logger
	transaction *postgresql.TransactionsStore
}

func NewTransactionService(conf *config.Config, transaction *postgresql.TransactionsStore, logger *zap.Logger) *Transaction {
	return &Transaction{
		conf:        conf,
		transaction: transaction,
		logger:      logger,
	}
}

func (t *Transaction) ListTransaction(publicKeys string, id int) ([]*postgresql.TransactionDb, error) {
	return t.transaction.GetTransactionByPublicKey(publicKeys, id)
}
