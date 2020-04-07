package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/duyhtq/incognito-data-sync/config"
)

type HttpClient struct {
	*http.Client
	conf *config.Config
}

// NewHttpClient to get http client instance
func NewHttpClient(conf *config.Config) *HttpClient {
	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}
	return &HttpClient{
		httpClient,
		conf,
	}
}

func (client *HttpClient) RPCCall(
	method string,
	params interface{},
	rpcResponse interface{},
) (err error) {

	payload := map[string]interface{}{
		"method": method,
		"params": params,
		"id":     0,
	}
	payloadInBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Post(client.conf.Incognito.ChainEndpoint, "application/json", bytes.NewBuffer(payloadInBytes))

	if err != nil {
		return err
	}
	respBody := resp.Body
	defer respBody.Close()

	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		return err
	}

	fmt.Println("________________________________________________")
	fmt.Println("Response from BlockChain:", string(body))
	fmt.Println("________________________________________________")

	err = json.Unmarshal(body, rpcResponse)
	if err != nil {
		return err
	}
	return nil
}
