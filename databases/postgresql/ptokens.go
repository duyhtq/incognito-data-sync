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
func (st *PTokensStore) UpdatePToken(id int, price float64) error {
	sqlStr := `
		UPDATE p_tokens SET price=$1 WHERE id=$2
		RETURNING id
	`
	tx := st.DB.MustBegin()
	defer tx.Commit()
	_, err := tx.Exec(sqlStr, price, id)
	return err
}
func (st *PTokensStore) GetPtoken(PtokenID string) (int, error) {
	sqlStr := `
		SELECT id FROM p_tokens
		WHERE token_id=$1
		LIMIT 1
	`
	result := []int{}
	err := st.DB.Select(&result, sqlStr, PtokenID)
	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0], nil
}
