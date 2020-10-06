package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	config "github.com/duyhtq/incognito-data-sync/config"
	pg "github.com/duyhtq/incognito-data-sync/databases/postgresql"
	"github.com/duyhtq/incognito-data-sync/models"
	"github.com/robfig/cron"
)

type ExchangeRateResp struct {
	Data struct {
		Quote struct {
			USD struct {
				Price float64 `json:"price"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

type PToken struct {
	TokenID   string  `json:"TokenID"`
	Symbol    string  `json:"Symbol"`
	Name      string  `json:"Name"`
	Decimals  int     `json:"Decimals"`
	PDecimals int     `json:"PDecimals"`
	PriceUsd  float64 `json:"PriceUsd"`
}
type PTokenResp struct {
	Result []PToken `json:"Result"`
}
type CronUpdatePrice struct {
	running bool
}

// NewServer is to new server instance
func NewServer() (*CronUpdatePrice, error) {
	return &CronUpdatePrice{}, nil
}

func main() {
	c := cron.New()
	if err := c.AddFunc("@every 30m", func() {
		run()
	}); err != nil {
		fmt.Println("err", err)
	}
	c.Start()
	select {}
}

func run() {
	fmt.Println("----------- Start update price -----------")
	conf := config.GetConfig()
	db, err := pg.Init(conf)
	if err != nil {
		return
	}

	pTokenStore, err := pg.NewPTokensStore(db)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	s, err := NewServer()
	if err != nil {
		fmt.Println("err", err)
		return
	}
	tokens, err := s.getTokens()
	for _, ptoken := range tokens.Result {
		id, err := pTokenStore.GetPtoken(ptoken.TokenID)
		if err != nil {
			fmt.Println("err", err)
			return
		}
		if id == 0 {
			fmt.Println("create", ptoken.Name, ptoken.TokenID)
			pTokenModel := models.PToken{
				TokenID: ptoken.TokenID,
				Name:    ptoken.Symbol,
				Symbol:  ptoken.Symbol,
				Decimal: math.Pow10(int(ptoken.PDecimals)),
				Price:   ptoken.PriceUsd,
			}
			err = pTokenStore.StorePToken(&pTokenModel)
			if err != nil {
				fmt.Println("StorePToken err:", err)
			}
		} else {
			fmt.Println("update", ptoken.Name)
			err = pTokenStore.UpdatePToken(id, ptoken.PriceUsd, math.Pow10(int(ptoken.PDecimals)))
			if err != nil {
				continue
			} else {
				fmt.Println("UpdatePToken err:", err)
			}
		}
	}
}
func (c *CronUpdatePrice) getTokens() (*PTokenResp, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.incognito.org/ptoken/list", nil)
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
	var exchangeRateResp PTokenResp
	if err := json.Unmarshal(respBody, &exchangeRateResp); err != nil {
		fmt.Println("error Unmarshal body ....", err)
		return nil, err
	}

	return &exchangeRateResp, nil
}
