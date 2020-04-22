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

func (t *Transaction) ReportPdexTrading() ([]*postgresql.ReportData, error) {
	return t.transaction.ReportPdexTrading()
}

func (t *Transaction) PdexVolume(token1str, token2str string) (float64, error) {
	return t.transaction.PdexVolume(token1str, token2str)
}
