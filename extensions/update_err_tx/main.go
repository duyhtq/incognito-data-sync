package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/duyhtq/incognito-data-sync/config"
	pg "github.com/duyhtq/incognito-data-sync/databases/postgresql"
	"github.com/inc-backend/crypto-libs/incognito"
	_ "github.com/lib/pq"
)

func GetPrice(name, date string) (float64, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/history?date=%s", name, date)
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

func GetDataCoin() ([]InfoCoin, error) {
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

type Respone struct {
	Timestamp    int                    `json:"BeaconTimeStamp"`
	PDEPoolPairs map[string]interface{} `json:"PDEPoolPairs"`
}

type PairToken struct {
	Token1IDStr     string  `json:"Token1IDStr"`
	Token1PoolValue float64 `jsonL"Token1PoolValue"`
	Token2IDStr     string  `json:"Token2IDStr"`
	Token2PoolValue float64 `jsonL"Token2PoolValue"`
}

func GetPdexState(beaconHeght int32) (float64, error) {
	client := &http.Client{}
	//bc := incognito.NewBlockchain(client, "http://51.83.237.20:9944/", "", "", "https://mainnet.incognito.org/", "0000000000000000000000000000000000000000000000000000000000000004")
	conf := config.GetConfig()
	bc := incognito.NewBlockchain(client, conf.Incognito.ChainEndpoint, "", "", conf.Incognito.ChainEndpoint, "0000000000000000000000000000000000000000000000000000000000000004")

	pde, err := bc.GetPdeState(beaconHeght)
	if err != nil {
		return 0, err
	}
	resultRespStr, err := json.Marshal(pde)
	if err != nil {
		return 0, err
	}
	var result Respone
	err = json.Unmarshal(resultRespStr, &result)
	if err != nil {
		return 0, err
	}
	var value1 float64
	for _, k := range result.PDEPoolPairs {
		pair, err := json.Marshal(k)
		if err != nil {
			continue
		}
		var pairToken PairToken
		err = json.Unmarshal(pair, &pairToken)
		if err != nil {
			continue
		}

		if pairToken.Token1IDStr == "0000000000000000000000000000000000000000000000000000000000000004" && pairToken.Token2IDStr == "716fd1009e2a1669caacc36891e707bfdf02590f96ebd897548e8963c95ebac0" {
			x1 := math.Ceil((pairToken.Token1PoolValue * pairToken.Token2PoolValue) / (pairToken.Token1PoolValue + math.Pow10(int(9))))
			value1 = (pairToken.Token2PoolValue - x1) / math.Pow10(int(6))
		}
	}
	return value1, nil
}

const (
	IssuingResponseMeta    = 25 // shields ota
	IssuingETHResponseMeta = 81 // shields eta

	ContractingRequestMeta = 26  // centralized: btc, xmr...
	BurningRequestMeta     = 27  // eta
	BurningRequestMetaV2   = 240 // eta
)

func main() {
	conf := config.GetConfig()
	db, err := pg.Init(conf)
	if err != nil {
		return
	}

	txStore, err := pg.NewTransactionsStore(db)
	if err != nil {
		return
	}

	txs, err := txStore.GetTransac()
	fmt.Println("txs count: ", len(txs))
	for _, tx := range txs {
		fmt.Println("ID ID ID: ", tx.ID)
		var Price float64

		token, _ := txStore.GetToken(tx.TokenID)
		if token == nil {
			continue
		}
		dataCoin, _ := GetDataCoin()
		var coinID string
		for _, coinDetail := range dataCoin {
			if coinDetail.Symbol == strings.ToLower(token.Symbol) {
				coinID = coinDetail.ID
			}
		}
		if tx.TokenID == "4a790f603aa2e7afe8b354e63758bb187a4724293d6057a46859c81b7bd0e9fb" {
			coinID = "paxos-standard"
		}
		if tx.TokenID == "47129778a650d75a7d0e8432878f05e466f5ee64cde3cee91f343ec523e75b58" {
			coinID = "uniswap"
		}
		if tx.TokenID == "4077654cf585a99b448564d1ecc915baf7b8ac58693d9f0a6af6c12b18143044" || tx.TokenID == "0369ef074eee01fe42ce10bcddf4f411435598f451b31ad169f5aa8def9d940f" {
			coinID = "harmony"
		}

		price, _ := GetPrice(coinID, tx.CreatedTime.Format("02-01-2006"))
		Price = price
		if price != 0 {
			err = txStore.UpdatePrice(tx.ID, Price)
			if err == nil {
				fmt.Println("-----Update success!------", coinID, Price)
			}
		}

	}
}
