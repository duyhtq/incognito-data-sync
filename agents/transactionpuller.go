package agents

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	GetTransactionById(txID string) (*models.ShortTransaction, error)
	GetToken(token string) (*models.PToken, error)
}

type TransactionPuller struct {
	AgentAbs
	TransactionsStore TransactionsStore
	ShardID           int
}

const (
	IssuingResponseMeta    = 25 // shields ota
	IssuingETHResponseMeta = 81 // shields eta

	ContractingRequestMeta = 26  // centralized: btc, xmr...
	BurningRequestMeta     = 27  // eta
	BurningRequestMetaV2   = 240 // eta
)

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

func (puller *TransactionPuller) GetPrice(name, date string) (float64, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/history?date=%s", name, date)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	req.Header.Set("Accepts", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		return 0, err
	}
	if resp.Status != "200 OK" {
		return 0, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	type CurrentPrice struct {
		USD float64 `json:"usd"`
	}
	type MarketData struct {
		CurrentPrice CurrentPrice `json:"current_price"`
	}
	type CoingeckoResp struct {
		MarketData MarketData `json:"market_data"`
	}

	var exchangeRateResp CoingeckoResp
	if err := json.Unmarshal(respBody, &exchangeRateResp); err != nil {
		fmt.Println("error Unmarshal body ....", err)
		return 0, err
	}
	return exchangeRateResp.MarketData.CurrentPrice.USD, nil
}

type InfoCoin struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
}
type MetaData struct {
	Type          int     `json:"Type"`
	BurnedAmount  float64 `json:"BurnedAmount"`
	BurningAmount float64 `json:"BurningAmount"`
	TokenID       string  `json:"TokenID"`
}

func (puller *TransactionPuller) GetDataCoin() ([]InfoCoin, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/list")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	req.Header.Set("Accepts", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		return nil, err
	}
	if resp.Status != "200 OK" {
		return nil, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	var respone []InfoCoin
	if err := json.Unmarshal(respBody, &respone); err != nil {
		fmt.Println("error Unmarshal body ....", err)
		return nil, err
	}
	return respone, nil
}
func (puller *TransactionPuller) Price(symbol, tokenId, timestamp string) float64 {
	dataCoin, _ := puller.GetDataCoin()
	var coinID string
	for _, coinDetail := range dataCoin {
		if coinDetail.Symbol == strings.ToLower(symbol) {
			coinID = coinDetail.ID
		}
	}
	if tokenId == "4a790f603aa2e7afe8b354e63758bb187a4724293d6057a46859c81b7bd0e9fb" {
		coinID = "paxos-standard"
	}
	if tokenId == "47129778a650d75a7d0e8432878f05e466f5ee64cde3cee91f343ec523e75b58" {
		coinID = "uniswap"
	}
	if tokenId == "4077654cf585a99b448564d1ecc915baf7b8ac58693d9f0a6af6c12b18143044" || tokenId == "0369ef074eee01fe42ce10bcddf4f411435598f451b31ad169f5aa8def9d940f" {
		coinID = "harmony"
	}
	timeDate, _ := time.Parse("2006-01-02T15:04:05.999999", timestamp)
	price, _ := puller.GetPrice(coinID, timeDate.Format("02-01-2006"))
	return price
}
func (puller *TransactionPuller) Execute() {
	// fmt.Println("time1: ", time.Now())
	fmt.Println("[Transaction puller] Agent is executing...")

	processingTxs := []string{}
	latestBlockHeight, err := puller.TransactionsStore.GetLatestProcessedShardHeight(puller.ShardID)
	if err != nil {
		log.Printf("[Transaction puller] An error occured while GetLatestProcessedShardHeight the latest processed shard %d: %+v \n", puller.ShardID, err)
		return
	}
	if latestBlockHeight > 0 {
		latestTxs, err := puller.TransactionsStore.LatestProcessedTxByHeight(puller.ShardID, latestBlockHeight)
		if err != nil || latestTxs == nil {
			log.Printf("[Transaction puller] An error occured while LatestProcessedTxByHeight the latest processed shard %d block height %d: %+v \n", puller.ShardID, latestBlockHeight, err)
			return
		}

		// long time:
		start := time.Now()
		temp, err := puller.TransactionsStore.ListProcessingTxByHeight(puller.ShardID, latestBlockHeight)
		elapsed := time.Since(start)
		log.Printf("ListProcessingTxByHeight took %s", elapsed)

		if err != nil || temp == nil {
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
		fmt.Println("for.....", latestBlockHeight)
		fmt.Println("time for ", time.Now())
		start := time.Now()
		temp, err := puller.TransactionsStore.ListNeedProcessingTxByHeight(puller.ShardID, latestBlockHeight)
		elapsed := time.Since(start)
		log.Printf("ListNeedProcessingTxByHeight took %s", elapsed)

		if len(temp) > 0 && err == nil {
			for _, t := range temp {
				processingTxs = append(processingTxs, t.TxsHash...)
			}
			latestBlockHeight = temp[len(temp)-1].BlockHeight
		} else {
			fmt.Println("[Transaction puller] No more tx to process ShardID, blockHeight, err", puller.ShardID, latestBlockHeight, err)
			continue
		}

		if len(processingTxs) > 0 {
			for _, t := range processingTxs {
				fmt.Println("time for for ", time.Now())
				fmt.Println("Transaction ID: ", t)
				txByID, err := puller.TransactionsStore.GetTransactionById(t)
				if err != nil || txByID != nil {
					fmt.Println("err get tx id: ", err)
					continue
				}
				time.Sleep(500 * time.Millisecond)
				tx, e := puller.getTransaction(t)
				if e != nil {
					log.Printf("[Transaction puller] An error occured while getting transaction %s : %+v\n", t, e)
					continue
				}

				fmt.Println("execute tx to db begin: ", time.Now())
				txModel := models.Transaction{
					ShardID:                  puller.ShardID,
					BlockHeight:              tx.BlockHeight,
					BlockHash:                tx.BlockHash,
					Info:                     tx.Info,
					TxID:                     tx.Hash,
					TxType:                   tx.Type,
					TxVersion:                tx.Version,
					PRVFee:                   tx.Fee,
					Proof:                    &tx.Proof,
					TransactedPrivacyCoinFee: tx.PrivacyCustomTokenFee,
					//TransactedPrivacyCoinProofDetail: tx.PrivacyCustomTokenProofDetail,
					//ProofDetail: tx.ProofDetail
				}
				dataJson, err := json.Marshal(tx)
				if err == nil {
					txModel.Data = string(dataJson)
				}

				if len(tx.Metadata) > 0 {
					txModel.Metadata = &tx.Metadata
					data := &MetaData{Type: 0}
					err := json.Unmarshal([]byte(tx.Metadata), data)
					if err == nil {
						txModel.MetaDataType = data.Type
					}
				}
				if txModel.MetaDataType == IssuingResponseMeta || txModel.MetaDataType == IssuingETHResponseMeta {
					type PrivacyCustomTokenData struct {
						Amount       float64 `json:"Amount"`
						PropertyID   string  `json:"PropertyID"`
						PropertyName string  `json:"PropertyName"`
					}
					data := &PrivacyCustomTokenData{}
					_ = json.Unmarshal([]byte(tx.PrivacyCustomTokenData), data)
					token, _ := puller.TransactionsStore.GetToken(data.PropertyID)

					if token != nil {
						fmt.Printf("=======tx: %+v \n", tx)
						fmt.Printf("=======token: %+v \n", token)
						fmt.Printf("=======data: %+v \n", data)
						txModel.AmountSheld = data.Amount / token.Decimal
						txModel.ShieldType = 1
						txModel.TokenName = token.Name
						txModel.TokenID = data.PropertyID
						price := puller.Price(token.Symbol, data.PropertyID, tx.LockTime)
						txModel.Price = price
					} else {
						log.Printf("[Transaction puller] An error occured while getting token err %+v\n", token)

					}
				}

				if txModel.MetaDataType == ContractingRequestMeta || txModel.MetaDataType == BurningRequestMeta || txModel.MetaDataType == BurningRequestMetaV2 {
					data := &MetaData{}
					err := json.Unmarshal([]byte(tx.Metadata), data)
					if err != nil {
						log.Printf("[Transaction puller] An error occured while getting transaction %s : %+v\n", t, err)
						continue
					}
					token, _ := puller.TransactionsStore.GetToken(data.TokenID)
					if token != nil {
						fmt.Printf("=======tx: %+v \n", tx)
						fmt.Printf("=======token: %+v \n", token)
						fmt.Printf("=======data: %+v \n", data)

						if txModel.MetaDataType == ContractingRequestMeta {
							txModel.AmountSheld = data.BurnedAmount / token.Decimal
						} else {
							txModel.AmountSheld = data.BurningAmount / token.Decimal
						}
						txModel.ShieldType = 2
						txModel.TokenName = token.Name
						txModel.TokenID = data.TokenID
						price := puller.Price(token.Symbol, data.TokenID, tx.LockTime)
						txModel.Price = price
					} else {
						log.Printf("[Transaction puller] An error occured while getting token err %+v\n", token)

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
}

/*
- List input coin in ProofDetail (SerialNumber)
- List output coin in ProofDetail (CoinCommitment, PublicKey)

- List input coin in Privacy coin proof detail (SerialNumber)
- List output coin in Privacy coin proof detail (CoinCommitment, PublicKey)


*/
