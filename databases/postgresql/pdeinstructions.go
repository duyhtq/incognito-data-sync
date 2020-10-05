package postgresql

import (
	"fmt"

	"github.com/duyhtq/incognito-data-sync/models"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// PDEInstructionsPGStore is PDE instructions postgresql store
type PDEInstructionsPGStore struct {
	PGStoreAbs
}

// NewPDEInstructionsPGStore initialises PDEStatePGStore with pg connection
func NewPDEInstructionsPGStore(db *sqlx.DB) (*PDEInstructionsPGStore, error) {
	return &PDEInstructionsPGStore{PGStoreAbs{DB: db}}, nil
}

func (is *PDEInstructionsPGStore) EnsurePDEInstructionsPGStore() {
	// TODO: implement it
}

func (is *PDEInstructionsPGStore) GetLatestProcessedBCHeight() (uint64, error) {
	sqlStr := `
		SELECT beacon_height FROM pde_trades
		ORDER BY beacon_height DESC
		LIMIT 1
	`
	bcHeights := []uint64{}
	err := is.DB.Select(&bcHeights, sqlStr)
	if err != nil {
		return 0, err
	}
	if len(bcHeights) == 0 {
		return 0, nil
	}
	return bcHeights[0], nil
}

func (is *PDEInstructionsPGStore) GetPdex() ([]*models.PDETrade, error) {
	sql := "SELECT requested_tx_id, receiving_tokenid_str, receive_amount, beacon_height, beacon_time_stamp FROM pde_trades where price is null or amount is null"
	result := []*models.PDETrade{}
	err := is.DB.Select(&result, sql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}
func (is *PDEInstructionsPGStore) GetToken(token string) (*models.PToken, error) {
	sql := `SELECT token_id, name, symbol, decimal, price FROM p_tokens WHERE token_id=$1`
	result := []*models.PToken{}
	err := is.DB.Select(&result, sql, token)
	if err != nil {
		return nil, err
	} else {
		if len(result) == 0 {
			return nil, nil
		}
		return result[0], nil
	}
}

func (is *PDEInstructionsPGStore) UpdateCustomFiled(txid string, price, amount float64) error {
	tx := is.DB.MustBegin()
	tx.MustExec("Update pde_trades set price= $1, amount = $2 WHERE requested_tx_id = $3",
		price, amount, txid)
	err := tx.Commit()
	return err
}

func (is *PDEInstructionsPGStore) GetPriceUsd(token string) (float64, error) {
	sqlStr := `
		SELECT price FROM p_tokens
		WHERE token_id=$1
		LIMIT 1
	`
	result := []float64{}
	err := is.DB.Select(&result, sqlStr, token)
	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0], nil
}
func (is *PDEInstructionsPGStore) StorePDETrade(tradeModel *models.PDETrade) error {
	fmt.Println("StorePDETrade StorePDETrade StorePDETrade: ", tradeModel.RequestedTxID)

	sqlStr := `
		INSERT INTO pde_trades (
			trader_address_str,
			receiving_tokenid_str,
			receive_amount,
			token1_id_str,
			token2_id_str,
			shard_id,
			requested_tx_id,
			status,
			beacon_height,
			beacon_time_stamp,
			price,
			amount
		)
		VALUES (
			:trader_address_str,
			:receiving_tokenid_str,
			:receive_amount,
			:token1_id_str,
			:token2_id_str,
			:shard_id,
			:requested_tx_id,
			:status,
			:beacon_height,
			:beacon_time_stamp,
			:price,
			:amount
		)
		RETURNING id
	`
	tx := is.DB.MustBegin()
	defer tx.Commit()
	_, err := tx.NamedQuery(sqlStr, tradeModel)
	return err
}
