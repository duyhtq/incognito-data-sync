package models

import "time"

type Transaction struct {
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
	PublickeyList                    []string  `db:"public_key_list"`
	SerialNumberList                 []string  `db:"serial_number_list"`
	CoinCommitmentList               []string  `db:"coin_commitment_list"`

	ShieldType  int     `db:"shield_type"`
	AmountSheld float64 `db:"amount_shield"`
	Price       float64 `db:"price"`
	TokenID     string  `db:"token_id"`
	TokenName   string  `db:"token_name"`
}

type ShortTransaction struct {
	ID   string `db:"id"`
	TxID string `db:"tx_id"`

	Data                             string    `db:"data"`
	Proof                            *string   `db:"proof"`
	ProofDetail                      *string   `db:"proof_detail"`
	Metadata                         *string   `db:"metadata"`
	TransactedPrivacyCoin            *string   `db:"transacted_privacy_coin"`
	TransactedPrivacyCoinProofDetail *string   `db:"transacted_privacy_coin_proof_detail"`
	MetaDataType                     int       `db:"meta_data_type"`
	CreatedTime                      time.Time `db:"created_time"`
}
