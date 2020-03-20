package agents

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/duyhtq/incognito-data-sync/databases/postgresql"
	"github.com/duyhtq/incognito-data-sync/entities"
	"github.com/duyhtq/incognito-data-sync/models"
	"github.com/duyhtq/incognito-data-sync/utils"
)

type TransactionsStore interface {
	GetLatestProcessedShardHeight(shardID int) (uint64, error)
	StoreTransaction(txs *models.Transaction) error
	LatestProcessedTxByHeight(shardID int, blockHeight uint64) ([]string, error)
	ListProcessingTxByHeight(shardID int, blockHeight uint64) (*postgresql.ListProcessingTx, error)
	ListNeedProcessingTxByHeight(shardID int, blockHeight uint64) ([]*postgresql.ListProcessingTx, error)
	GetTransactionById(txID string) (*models.Transaction, error)
}

type TransactionPuller struct {
	AgentAbs
	TransactionsStore TransactionsStore
	ShardID           int
}

func NewTransactionPuller(name string,
	frequency int,
	rpcClient *utils.HttpClient,
	shardID int,
	transactionsStore TransactionsStore) *TransactionPuller {
	puller := &TransactionPuller{
		TransactionsStore: transactionsStore,
		ShardID:           shardID,
		AgentAbs: AgentAbs{
			Name:      name,
			RPCClient: rpcClient,
			Quit:      make(chan bool),
			Frequency: frequency,
		},
	}
	return puller
}

func (puller *TransactionPuller) getTransaction(txHash string) (*entities.TransactionDetail, error) {
	params := []interface{}{txHash}
	var shardBlockRes entities.TransactionDetailRes
	err := puller.RPCClient.RPCCall("gettransactionbyhash", params, &shardBlockRes)
	if err != nil {
		return nil, err
	}
	if shardBlockRes.RPCError != nil {
		return nil, errors.New(shardBlockRes.RPCError.Message)
	}
	return shardBlockRes.Result, nil
}

func (puller *TransactionPuller) Execute() {
	fmt.Println("time1: ", time.Now())
	fmt.Println("[Transaction puller] Agent is executing...")

	processingTxs := []string{}
	latestBlockHeight, err := puller.TransactionsStore.GetLatestProcessedShardHeight(puller.ShardID)
	if err != nil {
		log.Printf("[Transaction puller] An error occured while GetLatestProcessedShardHeight the latest processed shard %d: %+v \n", puller.ShardID, err)
		return
	}
	if latestBlockHeight > 0 {
		latestTxs, err := puller.TransactionsStore.LatestProcessedTxByHeight(puller.ShardID, latestBlockHeight)
		if err != nil {
			log.Printf("[Transaction puller] An error occured while LatestProcessedTxByHeight the latest processed shard %d block height %d: %+v \n", puller.ShardID, latestBlockHeight, err)
			return
		}

		temp, err := puller.TransactionsStore.ListProcessingTxByHeight(puller.ShardID, latestBlockHeight)
		if err != nil {
			log.Printf("[Transaction puller] An error occured while ListProcessingTxByHeight shard %d block height %d: %+v \n", puller.ShardID, latestBlockHeight, err)
			return
		}

		if len(temp.TxsHash) > len(latestTxs) {
			fmt.Println("time2: ", time.Now())
			for _, a := range temp.TxsHash {
				if !utils.StringInSlice(a, latestTxs) {
					processingTxs = append(processingTxs, a)
				}
			}
			fmt.Println("time3", time.Now())
		}
	} else {
		latestBlockHeight = 0
	}

	latestBlockHeight += 1
	for {
		fmt.Println("for.....")
		fmt.Println("time for ", time.Now())
		temp, err := puller.TransactionsStore.ListNeedProcessingTxByHeight(puller.ShardID, latestBlockHeight)
		if len(temp) > 0 && err == nil {
			for _, t := range temp {
				processingTxs = append(processingTxs, t.TxsHash...)
			}
			latestBlockHeight = temp[len(temp)-1].BlockHeight
		} else {
			log.Printf("[Transaction puller] No more tx to process\n")
			continue
		}

		if len(processingTxs) > 0 {
			for _, t := range processingTxs {
				fmt.Println("time for for ", time.Now())
				txByID, err := puller.TransactionsStore.GetTransactionById(t)
				if err != nil || txByID != nil {
					fmt.Println("err get tx id: ", err)
					continue
				}
				time.Sleep(5 * time.Millisecond)
				tx, e := puller.getTransaction(t)
				if e != nil {
					log.Printf("[Transaction puller] An error occured while getting transaction %s : %+v\n", t, e)
					continue
				}

				fmt.Println("execute tx to db begin: ", time.Now())
				txModel := models.Transaction{
					ShardID:     puller.ShardID,
					BlockHeight: tx.BlockHeight,
					BlockHash:   tx.BlockHash,
					Info:        tx.Info,
					TxID:        tx.Hash,
					TxType:      tx.Type,
					TxVersion:   tx.Version,
					PRVFee:      tx.Fee,
					//Data:
					Proof: &tx.Proof,
					//ProofDetail: tx.ProofDetail
					TransactedPrivacyCoinFee: tx.PrivacyCustomTokenFee,
					//TransactedPrivacyCoinProofDetail: tx.PrivacyCustomTokenProofDetail,
				}
				dataJson, err := json.Marshal(tx)
				if err == nil {
					txModel.Data = string(dataJson)
				}

				if len(tx.Metadata) > 0 {
					txModel.Metadata = &tx.Metadata

					// get meta data type:
					type MetaData struct {
						Type int `json:"Type"`
					}
					data := &MetaData{Type: 0}
					err := json.Unmarshal([]byte(tx.Metadata), data)
					if err == nil {
						txModel.MetaDataType = data.Type
					}
				}

				// Begin custom data =========================================
				var listSerialNumber, listCoinCommitment, listPublicKey []string

				// 1. for ProofDetail:
				// 1.a. Get InputCoins in SerialNumber:
				if tx.ProofDetail.InputCoins != nil {
					for _, coinDetail := range tx.ProofDetail.InputCoins {
						if coinDetail.CoinDetails.SerialNumber != "" {
							listSerialNumber = append(listSerialNumber, coinDetail.CoinDetails.SerialNumber)
						}
					}
				}
				// 1.b. Get CoinCommitment, PublicKey in OutputCoins:
				if tx.ProofDetail.OutputCoins != nil {
					for _, coinDetail := range tx.ProofDetail.OutputCoins {
						if coinDetail.CoinDetails.PublicKey != "" {
							listPublicKey = append(listPublicKey, coinDetail.CoinDetails.PublicKey)
						}
						if coinDetail.CoinDetails.CoinCommitment != "" {
							listCoinCommitment = append(listCoinCommitment, coinDetail.CoinDetails.CoinCommitment)
						}
					}
				}
				// 2. for PrivacyCustomTokenProofDetail:
				// 2.a. Get SerialNumber in InputCoins:
				if tx.PrivacyCustomTokenProofDetail.InputCoins != nil {
					for _, coinDetail := range tx.PrivacyCustomTokenProofDetail.InputCoins {
						if coinDetail.CoinDetails.SerialNumber != "" {
							listSerialNumber = append(listSerialNumber, coinDetail.CoinDetails.SerialNumber)
						}
					}
				}
				// 1.b. Get CoinCommitment, PublicKey in OutputCoins:
				if tx.PrivacyCustomTokenProofDetail.OutputCoins != nil {
					for _, coinDetail := range tx.PrivacyCustomTokenProofDetail.OutputCoins {
						if coinDetail.CoinDetails.PublicKey != "" {
							listPublicKey = append(listPublicKey, coinDetail.CoinDetails.PublicKey)
						}
						if coinDetail.CoinDetails.CoinCommitment != "" {
							listCoinCommitment = append(listCoinCommitment, coinDetail.CoinDetails.CoinCommitment)
						}
					}
				}
				//3. set to db:
				txModel.SerialNumberList = listSerialNumber
				txModel.PublickeyList = listPublicKey
				txModel.CoinCommitmentList = listCoinCommitment

				//=============== end struct custom filed in db ===========

				// continue ....
				proofDetailJson, err := json.Marshal(tx.ProofDetail)
				if err == nil {
					temp := string(proofDetailJson)
					txModel.ProofDetail = &temp
				}
				privacyCustomTokenProofDetailJson, err := json.Marshal(tx.PrivacyCustomTokenProofDetail)
				if err == nil {
					temp := string(privacyCustomTokenProofDetailJson)
					txModel.TransactedPrivacyCoinProofDetail = &temp
				}
				if len(tx.PrivacyCustomTokenData) > 0 {
					txModel.TransactedPrivacyCoin = &tx.PrivacyCustomTokenData
				}
				txModel.CreatedTime, _ = time.Parse("2006-01-02T15:04:05.999999", tx.LockTime)
				err = puller.TransactionsStore.StoreTransaction(&txModel)

				fmt.Println("execute tx to db end: ", time.Now())

				if err != nil {
					log.Printf("[Transaction puller err] An error occured while storing tx %s, shard %d err: %+v\n", t, puller.ShardID, err)
					continue
				} else {
					log.Printf("[Transaction puller success] Success storing tx %s, shard %d\n", t, puller.ShardID)
				}
			}
		}
		latestBlockHeight++
	}

	fmt.Println("[Transaction puller] Agent is finished...")
}

/*
- List input coin in ProofDetail (SerialNumber)
- List output coin in ProofDetail (CoinCommitment, PublicKey)

- List input coin in Privacy coin proof detail (SerialNumber)
- List output coin in Privacy coin proof detail (CoinCommitment, PublicKey)


*/
