package main

import (
	"encoding/json"
	"fmt"

	config "github.com/duyhtq/incognito-data-sync/config"
	pg "github.com/duyhtq/incognito-data-sync/databases/postgresql"
	"github.com/duyhtq/incognito-data-sync/entities"
	_ "github.com/lib/pq"
)

func main() {

	if true {
		return
	}

	conf := config.GetConfig()

	db, err := pg.Init(conf)
	if err != nil {
		return
	}

	txStore, err := pg.NewTransactionsStore(db)
	if err != nil {
		return
	}

	txs, err := txStore.GetTransactions()

	fmt.Println("txs count: ", len(txs))

	// if true {
	// 	return
	// }

	for index, tx := range txs {

		fmt.Println("index: ", index)

		type MetaData struct {
			Type int `json:"Type"`
		}
		metaData := &MetaData{Type: 0}
		if tx.Metadata != nil {
			json.Unmarshal([]byte(*tx.Metadata), metaData)
		}

		// Begin custom data =========================================
		var listSerialNumber, listCoinCommitment, listPublicKey []string

		proofDetailString := *tx.ProofDetail

		proofDetailObject := &entities.ProofDetail{}
		err = json.Unmarshal([]byte(proofDetailString), proofDetailObject)

		if err == nil {
			// 1. for ProofDetail:
			// 1.a. Get InputCoins in SerialNumber:
			if proofDetailObject.InputCoins != nil {
				for _, coinDetail := range proofDetailObject.InputCoins {
					if coinDetail.CoinDetails.SerialNumber != "" {
						listSerialNumber = append(listSerialNumber, coinDetail.CoinDetails.SerialNumber)
					}
				}
			}
			// 1.b. Get CoinCommitment, PublicKey in OutputCoins:
			if proofDetailObject.OutputCoins != nil {
				for _, coinDetail := range proofDetailObject.OutputCoins {
					if coinDetail.CoinDetails.PublicKey != "" {
						listPublicKey = append(listPublicKey, coinDetail.CoinDetails.PublicKey)
					}
					if coinDetail.CoinDetails.CoinCommitment != "" {
						listCoinCommitment = append(listCoinCommitment, coinDetail.CoinDetails.CoinCommitment)
					}
				}
			}
		}

		privacyCustomTokenProofDetailString := *tx.TransactedPrivacyCoinProofDetail

		transactedPrivacyCoinProofDetailObject := &entities.ProofDetail{}
		err = json.Unmarshal([]byte(privacyCustomTokenProofDetailString), transactedPrivacyCoinProofDetailObject)

		if err == nil {
			// 2. for PrivacyCustomTokenProofDetail:
			// 2.a. Get SerialNumber in InputCoins:
			if transactedPrivacyCoinProofDetailObject.InputCoins != nil {
				for _, coinDetail := range transactedPrivacyCoinProofDetailObject.InputCoins {
					if coinDetail.CoinDetails.SerialNumber != "" {
						listSerialNumber = append(listSerialNumber, coinDetail.CoinDetails.SerialNumber)
					}
				}
			}
			// 2.b. Get CoinCommitment, PublicKey in OutputCoins:
			if transactedPrivacyCoinProofDetailObject.OutputCoins != nil {
				for _, coinDetail := range transactedPrivacyCoinProofDetailObject.OutputCoins {
					if coinDetail.CoinDetails.PublicKey != "" {
						listPublicKey = append(listPublicKey, coinDetail.CoinDetails.PublicKey)
					}
					if coinDetail.CoinDetails.CoinCommitment != "" {
						listCoinCommitment = append(listCoinCommitment, coinDetail.CoinDetails.CoinCommitment)
					}
				}
			}
		}

		// fmt.Println("SerialNumberList, CoinCommitmentList, PublickeyList", listSerialNumber, listCoinCommitment, listPublicKey)

		// update
		err = txStore.UpdateCustomFiledTransaction(tx.ID, listSerialNumber, listPublicKey, listCoinCommitment, metaData.Type)
		if err == nil {
			fmt.Println("update success!", tx.ID)
		}

		fmt.Println("-----------------------")
	}

}

// 1: 153418 -> 225222
//CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
