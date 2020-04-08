package postgresql

import (
	"github.com/duyhtq/incognito-data-sync/models"
	"github.com/jmoiron/sqlx"
)

type PTokensStore struct {
	PGStoreAbs
}

func NewPTokensStore(db *sqlx.DB) (*PTokensStore, error) {
	return &PTokensStore{PGStoreAbs{DB: db}}, nil
}

func (st *PTokensStore) StorePToken(token *models.PToken) error {
	sqlStr := `
		INSERT INTO p_tokens (token_id, name, symbol, decimal, price)
		VALUES (:token_id, :name, :symbol, :decimal, :price)
		RETURNING id
	`
	tx := st.DB.MustBegin()
	defer tx.Commit()
	_, err := tx.NamedQuery(sqlStr, token)
	return err
}
