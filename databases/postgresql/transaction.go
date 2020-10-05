package postgresql

import (
	"encoding/json"
	"fmt"
	"strings"

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
func (st *TransactionsStore) GetToken(token string) (*models.PToken, error) {
	sql := `SELECT token_id, name, decimal,symbol, price FROM p_tokens WHERE token_id=$1`
	result := []*models.PToken{}
	err := st.DB.Select(&result, sql, token)
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

func (st *TransactionsStore) GetTransactionsNotFix() ([]*models.ShortTransaction, error) {
	sql := "SELECT id, tx_id, data, meta_data_type, created_time, proof, proof_detail, metadata, transacted_privacy_coin, transacted_privacy_coin_proof_detail FROM transactions where meta_data_type in (25,81, 26, 27, 240)  and shield_type IS NULL"
	result := []*models.ShortTransaction{}
	err := st.DB.Select(&result, sql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
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

	tx.MustExec("INSERT INTO transactions (data, tx_id, tx_version, shard_id, tx_type, prv_fee, info, proof, proof_detail, metadata, transacted_privacy_coin, transacted_privacy_coin_proof_detail, transacted_privacy_coin_fee, created_time, block_height, block_hash, serial_number_list, public_key_list, coin_commitment_list, meta_data_type, shield_type, amount_shield, price, token_id, token_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17 , $18, $19, $20, $21, $22, $23, $24, $25) ",
		txs.Data, txs.TxID, txs.TxVersion, txs.ShardID, txs.TxType, txs.PRVFee, txs.Info, txs.Proof, txs.ProofDetail, txs.Metadata, txs.TransactedPrivacyCoin, txs.TransactedPrivacyCoinProofDetail, txs.TransactedPrivacyCoinFee, txs.CreatedTime, txs.BlockHeight, txs.BlockHash, pq.Array(txs.SerialNumberList), pq.Array(txs.PublickeyList), pq.Array(txs.CoinCommitmentList), txs.MetaDataType, txs.ShieldType, txs.AmountSheld, txs.Price, txs.TokenID, txs.TokenName)

	// _, err := tx.NamedQuery(sqlStr, txs)

	err := tx.Commit()

	return err
}

type ReportData struct {
	Day         string  `db:"day"`
	Total       int     `db:"total"`
	TotalVolume float64 `db:"total_volume"`
	Month       int     `db:"month"`
	Year        int     `db:"year"`
}

func (st *TransactionsStore) ReportPdexTrading(rangeFilter, token string) ([]*ReportData, error) {
	var sql string
	var err error
	var depositQuery string
	var filter string
	result := []*ReportData{}

	sql = `
	SELECT
		%s,
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
			AND o.status = 'accepted' %s
		ORDER BY
			created_date) AS b
	GROUP BY
		%s;
		`

	if token != "" {
		p := strings.Split(token, ",")
		filter = fmt.Sprintf("and token1_id_str in ('%s', '%s') and token2_id_str in ('%s', '%s')", p[0], p[1], p[0], p[1])
	} else {
		filter = ""
	}

	switch rangeFilter {
	case "week":
		depositQuery = fmt.Sprintf(sql, "TO_CHAR(DATE_TRUNC('week', b.created_date), 'YYYYWW') AS day", filter, "DATE_TRUNC('week', b.created_date) ORDER BY day")
	case "month":
		depositQuery = fmt.Sprintf(sql, "extract(year FROM b.created_date) AS year, extract(MONTH FROM b.created_date) AS month", filter, "extract(MONTH FROM b.created_date), extract(year FROM b.created_date) order by year, month")
	case "year":
		depositQuery = fmt.Sprintf(sql, "extract(year FROM b.created_date) AS year", filter, "extract(YEAR FROM b.created_date)")
	default:
		depositQuery = fmt.Sprintf(sql, "CAST(b.created_date AS date) AS day", filter, "CAST(b.created_date AS DATE)")
	}
	err = st.DB.Select(&result, depositQuery)

	if err != nil {
		return nil, err
	} else {
		if len(result) == 0 {
			return nil, nil
		}
		return result, nil
	}
}

func (st *TransactionsStore) Report24h() ([]*ReportData, error) {
	var sql string
	var err error
	result := []*ReportData{}

	sql = `
			SELECT
			COUNT(CAST(b.created_date AS DATE)) AS total,
			coalesce(Round(SUM(b.usd_value)::NUMERIC, 2), 0) AS total_volume
		FROM (
			SELECT
				CAST(o.beacon_time_stamp AS DATE) AS created_date,
				(((o.receive_amount + 0)::decimal / p.decimal::decimal) * p.price) AS usd_value
			FROM
				pde_trades AS o,
				p_tokens AS p
			WHERE
				o.receiving_tokenid_str = p.token_id
				AND o.status = 'accepted' and o.beacon_time_stamp >= NOW() - INTERVAL '24 HOURS'
			ORDER BY
				created_date) AS b
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
func (st *TransactionsStore) UpdateTransaction(ID string, TokenName, TokenID string, AmountSheld, Price float64, ShieldType int) error {
	tx := st.DB.MustBegin()
	tx.MustExec("Update transactions set token_name= $1, token_id = $2, amount_shield=$3, price=$4, shield_type=$5 WHERE id = $6",
		TokenName, TokenID, AmountSheld, Price, ShieldType, ID)
	err := tx.Commit()
	return err
}

type PdexVolume struct {
	Sum   float64 `db:"total_usd"`
	Count int     `db:"total"`
}

func (st *TransactionsStore) PdexVolume(token1str, token2str string) (float64, error) {

	result := []*PdexVolume{}

	var sql = `
		SELECT
			COALESCE(sum((((o.receive_amount + 0)::decimal / p.decimal::decimal) * p.price)), 0) AS total_usd,
			count (*)  AS total
		FROM
			pde_trades AS o,
			p_tokens AS p
		WHERE
			o.receiving_tokenid_str = p.token_id
			AND o.status = 'accepted'
			AND o.token1_id_str = $1
			AND o.token2_id_str = $2
			AND o.beacon_time_stamp >= NOW() - INTERVAL '24 HOURS'`

	err := st.DB.Select(&result, sql, token1str, token2str)

	if err != nil {
		return 0, err
	} else {
		return result[0].Sum, nil
	}
}

func (st *TransactionsStore) PdexTradingV2(rangeFilter, token string) ([]*ReportData, error) {
	var sql string
	var err error
	var depositQuery string
	var filter string
	result := []*ReportData{}

	sql = `
	SELECT
		%s,
		COUNT(CAST(b.created_date AS DATE)) AS total,
		Round(SUM(b.usd_value)::NUMERIC, 2) AS total_volume
	FROM (
		SELECT
			CAST(o.beacon_time_stamp AS DATE) AS created_date,
			((o.amount+ 0) * (o.price+ 0)) AS usd_value
		FROM
			pde_trades AS o
		WHERE
		
			o.status = 'accepted' %s
		ORDER BY
			created_date) AS b
	GROUP BY
		%s;
		`

	if token != "" {
		p := strings.Split(token, ",")
		filter = fmt.Sprintf("and token1_id_str in ('%s', '%s') and token2_id_str in ('%s', '%s')", p[0], p[1], p[0], p[1])
	} else {
		filter = ""
	}

	switch rangeFilter {
	case "week":
		depositQuery = fmt.Sprintf(sql, "TO_CHAR(DATE_TRUNC('week', b.created_date), 'YYYYWW') AS day", filter, "DATE_TRUNC('week', b.created_date) ORDER BY day")
	case "month":
		depositQuery = fmt.Sprintf(sql, "extract(year FROM b.created_date) AS year, extract(MONTH FROM b.created_date) AS month", filter, "extract(MONTH FROM b.created_date), extract(year FROM b.created_date) order by year, month")
	case "year":
		depositQuery = fmt.Sprintf(sql, "extract(year FROM b.created_date) AS year", filter, "extract(YEAR FROM b.created_date)")
	default:
		depositQuery = fmt.Sprintf(sql, "CAST(b.created_date AS date) AS day", filter, "CAST(b.created_date AS DATE)")
	}
	err = st.DB.Select(&result, depositQuery)

	if err != nil {
		return nil, err
	} else {
		if len(result) == 0 {
			return nil, nil
		}
		return result, nil
	}
}

func (st *TransactionsStore) Report24hV2() ([]*ReportData, error) {
	var sql string
	var err error
	result := []*ReportData{}
	sql = `
			SELECT
			COUNT(CAST(b.created_date AS DATE)) AS total,
			coalesce(Round(SUM(b.usd_value)::NUMERIC, 2), 0) AS total_volume
		FROM (
			SELECT
				CAST(o.beacon_time_stamp AS DATE) AS created_date,
				((o.amount+ 0) * (o.price+ 0)) AS usd_value
			FROM
				pde_trades AS o
			WHERE
				o.status = 'accepted' and o.beacon_time_stamp >= NOW() - INTERVAL '24 HOURS'
					ORDER BY
				created_date) AS b
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

func (st *TransactionsStore) Shield() ([]*ReportData, error) {
	var sql string
	var err error
	result := []*ReportData{}
	sql = `
			SELECT
			COUNT(CAST(b.created_date AS DATE)) AS total,
			coalesce(Round(SUM(b.usd_value)::NUMERIC, 2), 0) AS total_volume
			FROM (
			SELECT
				CAST(o.created_time AS DATE) AS created_date,
				((o.amount_shield + 0) * (o.price + 0)) AS usd_value
			FROM
				transactions AS o
			WHERE
				o.shield_type = 1
			ORDER BY
				created_date) AS b
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

func (st *TransactionsStore) Unshield() ([]*ReportData, error) {
	var sql string
	var err error
	result := []*ReportData{}
	sql = `		
			SELECT
				COUNT(CAST(b.created_date AS DATE)) AS total,
				coalesce(Round(SUM(b.usd_value)::NUMERIC, 2), 0) AS total_volume
			FROM (
				SELECT
					CAST(o.created_time AS DATE) AS created_date,
					((o.amount_shield + 0) * (o.price + 0)) AS usd_value
				FROM
					transactions AS o
				WHERE
					o.shield_type = 2
				ORDER BY
					created_date) AS b
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
