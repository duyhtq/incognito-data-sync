package postgresql

import (
	"encoding/json"
	"fmt"
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

func (st *TransactionsStore) GetTransactionById(txID string) (*models.ShortTransaction, error) {
	sql := `SELECT id, tx_id, data, proof, proof_detail, metadata, transacted_privacy_coin, transacted_privacy_coin_proof_detail FROM transactions WHERE tx_id=$1`
	fmt.Println("sql", sql)
	result := []*models.ShortTransaction{}
	err := st.DB.Select(&result, sql, txID)

	fmt.Println("result", result)

	if err != nil {
		return nil, err
	} else {
		if len(result) == 0 {
			return nil, nil
		}
		return result[0], nil
	}
}

func (st *TransactionsStore) GetTransactions() ([]*models.ShortTransaction, error) {
	sql := "SELECT id, tx_id, data, proof, proof_detail, metadata, transacted_privacy_coin, transacted_privacy_coin_proof_detail FROM transactions where id > 153418"
	// fmt.Println("sql", sql)
	result := []*models.ShortTransaction{}

	// result1 := models.Transaction2{}

	// rows, err := st.DB.Query(sql)

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// for rows.Next() {
	// 	err = rows.Scan(&result1)
	// 	fmt.Println("result1-->", result1)
	// 	result = append(result, result1)
	// }

	err := st.DB.Select(&result, sql)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// fmt.Println("results", result)

	return result, nil

	// if err != nil {
	// 	return nil, err
	// } else {
	// 	if len(result) == 0 {
	// 		return nil, nil
	// 	}
	// 	return result, nil
	// }
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

	// fmt.Println("txs.Data, txs.TxID, txs.TxVersion, txs.ShardID, txs.TxType, txs.PRVFee, txs.Info, txs.Proof, txs.ProofDetail, txs.Metadata, txs.TransactedPrivacyCoin, txs.TransactedPrivacyCoinProofDetail, txs.TransactedPrivacyCoinFee, txs.CreatedTime, txs.BlockHeight, txs.BlockHash, pq.Array(txs.SerialNumberList), pq.Array(txs.PublickeyList), pq.Array(txs.CoinCommitmentList), txs.MetaDataType", txs.Data, txs.TxID, txs.TxVersion, txs.ShardID, txs.TxType, txs.PRVFee, txs.Info, txs.Proof, txs.ProofDetail, txs.Metadata, txs.TransactedPrivacyCoin, txs.TransactedPrivacyCoinProofDetail, txs.TransactedPrivacyCoinFee, txs.CreatedTime, txs.BlockHeight, txs.BlockHash, pq.Array(txs.SerialNumberList), pq.Array(txs.PublickeyList), pq.Array(txs.CoinCommitmentList), txs.MetaDataType)

	tx.MustExec("INSERT INTO transactions (data, tx_id, tx_version, shard_id, tx_type, prv_fee, info, proof, proof_detail, metadata, transacted_privacy_coin, transacted_privacy_coin_proof_detail, transacted_privacy_coin_fee, created_time, block_height, block_hash, serial_number_list, public_key_list, coin_commitment_list, meta_data_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17 , $18, $19, $20) ",
		txs.Data, txs.TxID, txs.TxVersion, txs.ShardID, txs.TxType, txs.PRVFee, txs.Info, txs.Proof, txs.ProofDetail, txs.Metadata, txs.TransactedPrivacyCoin, txs.TransactedPrivacyCoinProofDetail, txs.TransactedPrivacyCoinFee, txs.CreatedTime, txs.BlockHeight, txs.BlockHash, pq.Array(txs.SerialNumberList), pq.Array(txs.PublickeyList), pq.Array(txs.CoinCommitmentList), txs.MetaDataType)

	// _, err := tx.NamedQuery(sqlStr, txs)

	err := tx.Commit()

	return err
}

type ReportData struct {
	Day         time.Time `db:"day"`
	Total       int       `db:"total"`
	TotalVolume float64   `db:"total_volume"`
}

func (st *TransactionsStore) ReportPdexTrading() ([]*ReportData, error) {
	var sql string
	var err error
	result := []*ReportData{}

	sql = `
	SELECT
		date_trunc('day', b.created_date) "day",
		COUNT(CAST(b.created_date AS DATE)) AS total,
		Round(SUM(b.usd_value)::NUMERIC, 2) AS total_volume
	FROM (
		SELECT
			CAST(o.beacon_time_stamp AS DATE) AS created_date,
			(((o.receive_amount + 0)::decimal / p.decimal::decimal) * p.price) AS usd_value
		FROM
			pde_trades AS o,
			p_tokens AS p
		WHERE
			o.receiving_tokenid_str = p.token_id
			AND o.status = 'accepted'
		ORDER BY
			created_date) AS b
	GROUP BY
		CAST(b.created_date AS DATE);
		`
	err = st.DB.Select(&result, sql)

	if err != nil {
		return nil, err
	} else {
		if len(result) == 0 {
			return nil, nil
		}
		return result, nil
	}
}
func (st *TransactionsStore) UpdateCustomFiledTransaction(ID string, serialNumberList, publickeyList, coinCommitmentList []string, metaDataType int) error {

	tx := st.DB.MustBegin()

	tx.MustExec("Update transactions set serial_number_list= $1, public_key_list = $2, coin_commitment_list=$3, meta_data_type=$4 WHERE id = $5",
		pq.Array(serialNumberList), pq.Array(publickeyList), pq.Array(coinCommitmentList), metaDataType, ID)

	err := tx.Commit()

	return err
}
