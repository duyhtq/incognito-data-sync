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

func (t *Transaction) ReportPdexTrading(rangeFilter, token string) ([]*postgresql.ReportData, error) {
	return t.transaction.ReportPdexTrading(rangeFilter, token)
}
func (t *Transaction) Report24h() ([]*postgresql.ReportData, error) {
	return t.transaction.Report24h()
}

func (t *Transaction) PdexVolume(token1str, token2str string) (float64, error) {
	return t.transaction.PdexVolume(token1str, token2str)
}

func (t *Transaction) PdexTradingV2(rangeFilter, token string) ([]*postgresql.ReportData, error) {
	return t.transaction.PdexTradingV2(rangeFilter, token)
}
func (t *Transaction) Report24hV2() ([]*postgresql.ReportData, error) {
	return t.transaction.Report24hV2()
}
func (t *Transaction) Shield() ([]*postgresql.ReportData, error) {
	return t.transaction.Shield()
}
func (t *Transaction) Unshield() ([]*postgresql.ReportData, error) {
	return t.transaction.Unshield()
}
