package agents

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/duyhtq/incognito-data-sync/entities"
	"github.com/duyhtq/incognito-data-sync/models"
	"github.com/duyhtq/incognito-data-sync/utils"
	"github.com/inc-backend/crypto-libs/incognito"
)

const (
	BCHeightStartForExtractingInsts = 432622
)

type PDEInstructionsStore interface {
	GetLatestProcessedBCHeight() (uint64, error)
	StorePDETrade(*models.PDETrade) error
	GetPriceUsd(string) (float64, error)
	GetToken(string) (*models.PToken, error)
}

type PDEInstsExtractor struct {
	AgentAbs
	PDEInstructionsStore PDEInstructionsStore
}

// type InfoCoin struct {
// 	ID     string `json:"id"`
// 	Symbol string `json:"symbol"`
// }
// type MetaData struct {
// 	Type         int     `json:"Type"`
// 	BurnedAmount float64 `json:"BurnedAmount"`
// 	TokenID      string  `json:"TokenID"`
// }

// NewPDEInstsExtractor instantiates PDEInstsExtractor
func NewPDEInstsExtractor(
	name string,
	frequency int,
	rpcClient *utils.HttpClient,
	pdeInstructionsStore PDEInstructionsStore,
) *PDEInstsExtractor {
	pdeInstsExtractor := &PDEInstsExtractor{}
	pdeInstsExtractor.Name = name
	pdeInstsExtractor.Frequency = frequency
	pdeInstsExtractor.Quit = make(chan bool)
	pdeInstsExtractor.RPCClient = rpcClient
	pdeInstsExtractor.PDEInstructionsStore = pdeInstructionsStore
	return pdeInstsExtractor
}

func (pie *PDEInstsExtractor) extractPDEInstsFromBeaconBlk(beaconHeight uint64) (*entities.PDEInfoFromBeaconBlock, error) {
	params := []interface{}{
		map[string]interface{}{
			"BeaconHeight": beaconHeight,
		},
	}
	var pdeExtractedInstsRes entities.PDEExtractedInstsRes
	err := pie.RPCClient.RPCCall("extractpdeinstsfrombeaconblock", params, &pdeExtractedInstsRes)
	if err != nil {
		return nil, err
	}
	// bb, _ := json.Marshal(pdeExtractedInstsRes)
	// fmt.Println("hahaha: ", string(bb))
	if pdeExtractedInstsRes.RPCError != nil {
		return nil, errors.New(pdeExtractedInstsRes.RPCError.Message)
	}
	return pdeExtractedInstsRes.Result, nil
}

// Execute is to pull pde state from incognito chain and put into db
func (pie *PDEInstsExtractor) Execute() {
	fmt.Println("PDE InstsExtractor agent is executing...")
	bcHeight, err := pie.PDEInstructionsStore.GetLatestProcessedBCHeight()
	if err != nil {
		log.Printf("[Instructions Extractor] An error occured while getting the latest processed beacon height: %+v \n", err)
		return
	}
	if bcHeight == 0 {
		bcHeight = uint64(BCHeightStartForExtractingInsts)
	} else {
		bcHeight++
	}

	for {
		time.Sleep(500 * time.Millisecond)
		// time.Sleep(5 * time.Second)
		log.Printf("[Instructions Extractor] Proccessing for beacon height: %d\n", bcHeight)
		insts, err := pie.extractPDEInstsFromBeaconBlk(bcHeight)
		if err != nil {
			fmt.Println("An error occured while extracting pde instruction from chain: ", err)
			return
		}
		if insts == nil {
			break
		}
		for _, trade := range insts.PDETrades {
			fmt.Println("PDETrades PDETrades PDETrades")
			time.Sleep(1 * time.Second)
			price, amount := pie.CheckPrice(trade.ReceivingTokenIDStr, trade.ReceiveAmount, trade.BeaconHeight, time.Unix(insts.BeaconTimeStamp, 0))
			tradeModel := models.PDETrade{
				TraderAddressStr:    trade.TraderAddressStr,
				ReceivingTokenIDStr: trade.ReceivingTokenIDStr,
				ReceiveAmount:       trade.ReceiveAmount,
				Token1IDStr:         trade.Token1IDStr,
				Token2IDStr:         trade.Token2IDStr,
				ShardID:             trade.ShardID,
				RequestedTxID:       trade.RequestedTxID,
				Status:              trade.Status,
				BeaconHeight:        trade.BeaconHeight,
				BeaconTimeStamp:     time.Unix(insts.BeaconTimeStamp, 0),
				Price:               price,
				Amount:              amount,
			}
			err = pie.PDEInstructionsStore.StorePDETrade(&tradeModel)
			if err != nil {
				fmt.Println("An error occured while storing pde trade")
				continue
			}
		}
		for _, trade := range insts.PDERefundedTradesV2 {
			fmt.Println("PDERefundedTradesV2 PDERefundedTradesV2 PDERefundedTradesV2")

			time.Sleep(1 * time.Second)
			price, amount := pie.CheckPrice(trade.TokenIDStr, trade.ReceiveAmount, trade.BeaconHeight, time.Unix(insts.BeaconTimeStamp, 0))
			tradeModel := models.PDETrade{
				TraderAddressStr:    trade.TraderAddressStr,
				ReceivingTokenIDStr: trade.TokenIDStr,
				ReceiveAmount:       trade.ReceiveAmount,
				Token1IDStr:         trade.TokenIDStr,
				Token2IDStr:         trade.TokenIDStr,
				ShardID:             trade.ShardID,
				RequestedTxID:       trade.RequestedTxID,
				Status:              trade.Status,
				BeaconHeight:        trade.BeaconHeight,
				BeaconTimeStamp:     time.Unix(insts.BeaconTimeStamp, 0),
				Price:               price,
				Amount:              amount,
			}
			err = pie.PDEInstructionsStore.StorePDETrade(&tradeModel)
			if err != nil {
				fmt.Println("An error occured while storing pde trade")
				continue
			}
		}
		for _, trade := range insts.PDEAcceptedTradesV2 {
			fmt.Println("PDEAcceptedTradesV2 PDEAcceptedTradesV2 PDEAcceptedTradesV2")

			time.Sleep(1 * time.Second)
			var tradeModel models.PDETrade
			if len(trade.TradePaths) == 1 {
				price, amount := pie.CheckPrice(trade.TradePaths[0].TokenIDToBuyStr, trade.TradePaths[0].ReceiveAmount, trade.BeaconHeight, time.Unix(insts.BeaconTimeStamp, 0))
				tradeModel = models.PDETrade{
					TraderAddressStr:    trade.TraderAddressStr,
					ReceivingTokenIDStr: trade.TradePaths[0].TokenIDToBuyStr,
					ReceiveAmount:       trade.TradePaths[0].ReceiveAmount,
					Token1IDStr:         trade.TradePaths[0].Token1IDStr,
					Token2IDStr:         trade.TradePaths[0].Token2IDStr,
					ShardID:             trade.ShardID,
					RequestedTxID:       trade.RequestedTxID,
					Status:              trade.Status,
					BeaconHeight:        trade.BeaconHeight,
					BeaconTimeStamp:     time.Unix(insts.BeaconTimeStamp, 0),
					Price:               price,
					Amount:              amount,
				}
			} else {
				var Token1IDStr, Token2IDStr, ReceivingTokenIDStr string
				var ReceiveAmount uint64
				if trade.TradePaths[0].TokenIDToBuyStr == "0000000000000000000000000000000000000000000000000000000000000004" {
					Token1IDStr = trade.TradePaths[0].Token2IDStr
					Token2IDStr = trade.TradePaths[1].Token2IDStr
					ReceivingTokenIDStr = trade.TradePaths[1].TokenIDToBuyStr
					ReceiveAmount = trade.TradePaths[1].ReceiveAmount

				} else {
					Token1IDStr = trade.TradePaths[1].Token2IDStr
					Token2IDStr = trade.TradePaths[0].Token2IDStr
					ReceivingTokenIDStr = trade.TradePaths[0].TokenIDToBuyStr
					ReceiveAmount = trade.TradePaths[0].ReceiveAmount
				}

				price, amount := pie.CheckPrice(ReceivingTokenIDStr, ReceiveAmount, trade.BeaconHeight, time.Unix(insts.BeaconTimeStamp, 0))
				tradeModel = models.PDETrade{
					TraderAddressStr:    trade.TraderAddressStr,
					ReceivingTokenIDStr: ReceivingTokenIDStr,
					ReceiveAmount:       ReceiveAmount,
					Token1IDStr:         Token1IDStr,
					Token2IDStr:         Token2IDStr,
					ShardID:             trade.ShardID,
					RequestedTxID:       trade.RequestedTxID,
					Status:              trade.Status,
					BeaconHeight:        trade.BeaconHeight,
					BeaconTimeStamp:     time.Unix(insts.BeaconTimeStamp, 0),
					Price:               price,
					Amount:              amount,
				}
			}

			err := pie.PDEInstructionsStore.StorePDETrade(&tradeModel)
			if err != nil {
				fmt.Println("An error occured while storing pde trade")
				continue
			}
		}
		bcHeight++
	}
	fmt.Println("PDE InstsExtractor agent is finished...")
}

func (pie *PDEInstsExtractor) CheckPrice(tokenID string, amount, beaconHeght uint64, time time.Time) (float64, float64) {
	var Price float64
	var Amount float64
	if tokenID == "0000000000000000000000000000000000000000000000000000000000000004" {
		pricePrv, _ := GetPdexState(int32(beaconHeght))
		Price = pricePrv
		Amount = float64(amount) / math.Pow10(int(9))
	} else {
		token, _ := pie.PDEInstructionsStore.GetToken(tokenID)
		if token == nil {
			return 0, 0
		}
		dataCoin, _ := GetDataCoin()
		var coinID string
		for _, coinDetail := range dataCoin {
			if coinDetail.Symbol == strings.ToLower(token.Symbol) {
				coinID = coinDetail.ID
			}
		}
		price, err := GetPrice(coinID, time.Format("02-01-2006"))
		fmt.Println("error GetPrice ....", err)
		Price = price
		Amount = float64(amount) / token.Decimal
	}
	return Price, Amount
}

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
	bc := incognito.NewBlockchain(client, "http://51.83.237.20:9944/", "", "", "https://mainnet.incognito.org/", "0000000000000000000000000000000000000000000000000000000000000004")
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
