package postgresql

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/duyhtq/incognito-data-sync/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// BeaconBlockStore is Beacon Block postgresql store
type TransactionsStore struct {
	PGStoreAbs
}

// NewBeaconBlockStore initialises BeaconBlockStore with pg connection
func NewTransactionsStore(db *sqlx.DB) (*TransactionsStore, error) {
	return &TransactionsStore{PGStoreAbs{DB: db}}, nil
}

func (st *TransactionsStore) GetLatestProcessedShardHeight(shardID int) (uint64, error) {
	sqlStr := `
		SELECT block_height FROM transactions
		WHERE shard_id=$1
		ORDER BY block_height DESC
		LIMIT 1
	`
	blockHeights := []uint64{}
	err := st.DB.Select(&blockHeights, sqlStr, shardID)
	if err != nil {
		return 0, err
	}
	if len(blockHeights) == 0 {
		return 0, nil
	}
	return blockHeights[0], nil
}

func (st *TransactionsStore) LatestProcessedTxByHeight(shardID int, blockHeight uint64) ([]string, error) {
	sql := `
		SELECT tx_id FROM transactions WHERE shard_id = $1 AND block_height = $2
	`
	result := []string{}
	err := st.DB.Select(&result, sql, shardID, blockHeight)
	if err != nil {
		return result, err
	}
	return result, err
}

type ListProcessingTx struct {
	BlockHeight   uint64
	BlockHash     string
	ShardID       int
	TxsHashString string
	TxsHash       []string ``
}

func (st *TransactionsStore) ListProcessingTxByHeight(shardID int, blockHeight uint64) (*ListProcessingTx, error) {
	sqlStr := `
		SELECT shardblock.block_height as BlockHeight, 
		shardblock.shard_id as ShardID, 
		shardblock.block_hash as BlockHash,  
		shardblock.list_hash_tx as TxsHash 

		
		FROM shard_blocks shardblock
		
		WHERE shard_id = $1 AND json_array_length(shardblock.list_hash_tx) > 0 AND block_height = $2
	`
	result := []*ListProcessingTx{}
	rows, err := st.DB.Query(sqlStr, shardID, blockHeight)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := &ListProcessingTx{
			TxsHash: []string{},
		}
		err1 := rows.Scan(&item.BlockHeight, &item.ShardID, &item.BlockHash, &item.TxsHashString)
		if err1 != nil {
			return nil, err1
		}
		err1 = json.Unmarshal([]byte(item.TxsHashString), &item.TxsHash)
		if err1 != nil {
			return nil, err1
		}
		result = append(result, item)
	}
	if len(result) > 0 {
		return result[0], err
	} else {
		return nil, nil
	}
}

func (st *TransactionsStore) ListNeedProcessingTxByHeight(shardID int, blockHeight uint64) ([]*ListProcessingTx, error) {
	sqlStr := `
		SELECT shardblock.block_height as BlockHeight, 
		shardblock.shard_id as ShardID, 
		shardblock.block_hash as BlockHash,  
		shardblock.list_hash_tx as TxsHash 
		
		FROM shard_blocks shardblock
		
		WHERE shard_id = $1 AND json_array_length(shardblock.list_hash_tx) > 0 AND block_height >= $2
		
		ORDER BY block_height ASC
	`
	result := []*ListProcessingTx{}
	rows, err := st.DB.Query(sqlStr, shardID, blockHeight)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := &ListProcessingTx{
			TxsHash: []string{},
		}
		err1 := rows.Scan(&item.BlockHeight, &item.ShardID, &item.BlockHash, &item.TxsHashString)
		if err1 != nil {
			return nil, err1
		}
		err1 = json.Unmarshal([]byte(item.TxsHashString), &item.TxsHash)
		if err1 != nil {
			return nil, err1
		}
		result = append(result, item)
	}
	if len(result) > 0 {
		return result, err
	} else {
		return nil, nil
	}
}

func (st *TransactionsStore) GetTransactionById(txID string) (*models.Transaction, error) {
	sql := `SELECT * FROM transactions WHERE tx_id=$1`
	result := []*models.Transaction{}
	err := st.DB.Select(&result, sql, txID)
	if err != nil {
		return nil, err
	} else {
		if len(result) == 0 {
			return nil, nil
		}
		return result[0], nil
	}
}

func (st *TransactionsStore) StoreTransaction(txs *models.Transaction) error {

	// fmt.Println(txs.Data, txs.TxID, txs.TxVersion, txs.ShardID, txs.TxType, txs.PRVFee, txs.Info, txs.Proof, txs.ProofDetail, txs.Metadata, txs.TransactedPrivacyCoin, txs.TransactedPrivacyCoinProofDetail, txs.TransactedPrivacyCoinFee, txs.CreatedTime, txs.BlockHeight, txs.BlockHash, pq.Array(txs.SerialNumberList), pq.Array(txs.PublickeyList), pq.Array(txs.CoinCommitmentList), txs.MetaDataType)
	fmt.Println("========================================")

	// sqlStr := `
	// 	INSERT INTO transactions (data, tx_id, tx_version, shard_id, tx_type, prv_fee, info, proof, proof_detail, metadata, transacted_privacy_coin, transacted_privacy_coin_proof_detail, transacted_privacy_coin_fee, created_time, block_height, block_hash, serial_number_list, public_key_list, coin_commitment_list)
	// 	VALUES (:data, :tx_id, :tx_version, :shard_id, :tx_type, :prv_fee, :info, :proof, :proof_detail, :metadata, :transacted_privacy_coin, :transacted_privacy_coin_proof_detail, :transacted_privacy_coin_fee, :created_time, :block_height, :block_hash, :serial_number_list, :public_key_list, :coin_commitment_list)
	// 	RETURNING id
	// `

	// db.Exec("INSERT INTO temp VALUES ($1)", []int32{1, 2, 3})

	// or:
	// values := []string{"foo", "bar"}
	// err := db.Query(`SELECT * FROM foo WHERE bar = ANY($1)`, pq.Array(values))

	tx := st.DB.MustBegin()
	// defer tx.Commit()

	tx.MustExec("INSERT INTO transactions (data, tx_id, tx_version, shard_id, tx_type, prv_fee, info, proof, proof_detail, metadata, transacted_privacy_coin, transacted_privacy_coin_proof_detail, transacted_privacy_coin_fee, created_time, block_height, block_hash, serial_number_list, public_key_list, coin_commitment_list, meta_data_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17 , $18, $19, $20) ",
		txs.Data, txs.TxID, txs.TxVersion, txs.ShardID, txs.TxType, txs.PRVFee, txs.Info, txs.Proof, txs.ProofDetail, txs.Metadata, txs.TransactedPrivacyCoin, txs.TransactedPrivacyCoinProofDetail, txs.TransactedPrivacyCoinFee, txs.CreatedTime, txs.BlockHeight, txs.BlockHash, pq.Array(txs.SerialNumberList), pq.Array(txs.PublickeyList), pq.Array(txs.CoinCommitmentList), txs.MetaDataType)

	// _, err := tx.NamedQuery(sqlStr, txs)

	err := tx.Commit()

	return err
}

type TransactionDb struct {
	ID        string `db:"id"`
	TxID      string `db:"tx_id"`
	TxVersion int8   `db:"tx_version"`
	TxType    string `db:"tx_type"`

	Data                             string    `db:"data"`
	ShardID                          int       `db:"shard_id"`
	PRVFee                           uint64    `db:"prv_fee"`
	Info                             string    `db:"info"`
	Proof                            *string   `db:"proof"`
	ProofDetail                      *string   `db:"proof_detail"`
	Metadata                         *string   `db:"metadata"`
	TransactedPrivacyCoin            *string   `db:"transacted_privacy_coin"`
	TransactedPrivacyCoinProofDetail *string   `db:"transacted_privacy_coin_proof_detail"`
	TransactedPrivacyCoinFee         uint64    `db:"transacted_privacy_coin_fee"`
	CreatedTime                      time.Time `db:"created_time"`
	BlockHeight                      uint64    `db:"block_height"`
	BlockHash                        string    `db:"block_hash"`
	MetaDataType                     int       `db:"meta_data_type"`
	SerialNumberList                 string    `db:"serial_number_list"`
	PublickeyList                    string    `db:"public_key_list"`
	CoinCommitmentList               string    `db:"coin_commitment_list"`
}

func (st *TransactionsStore) GetTransactionByPublicKey(publicKeys string, id int) ([]*TransactionDb, error) {
	var sql string
	var err error
	result := []*TransactionDb{}

	PublicKeys := strings.Split(publicKeys, ",")

	sql = `SELECT * FROM transactions WHERE public_key_list && $1 AND id > $2`
	err = st.DB.Select(&result, sql, pq.Array(PublicKeys), id)

	if err != nil {
		return nil, err
	} else {
		if len(result) == 0 {
			return nil, nil
		}
		return result, nil
	}
}
