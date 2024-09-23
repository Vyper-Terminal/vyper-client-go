package vyperclientgo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type VyperClient struct {
	BaseURL    string
	ApiKey     string
	HttpClient *http.Client
}

func NewVyperClient(apiKey string) *VyperClient {
	return &VyperClient{
		BaseURL: "https://api.vyper.trade",
		ApiKey:  apiKey,
		HttpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *VyperClient) request(method, endpoint string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, c.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.ApiKey)
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, &VyperApiError{
				Message:    fmt.Sprintf("HTTP error: %s", resp.Status),
				StatusCode: resp.StatusCode,
			}
		}
		return nil, &VyperApiError{
			Message:    apiResp.Message,
			StatusCode: resp.StatusCode,
			Response:   apiResp,
		}
	}

	return body, nil
}

func (c *VyperClient) GetChainIds() (map[string]int, error) {
	body, err := c.request("GET", "/api/v1/chain/ids", nil)
	if err != nil {
		return nil, err
	}

	var result map[string]int
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *VyperClient) GetTokenAth(chainId int, marketId string) (*TokenATH, error) {
	params := map[string]string{
		"chainID":  fmt.Sprintf("%d", chainId),
		"marketID": marketId,
	}
	body, err := c.request("GET", "/api/v1/token/ath", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result TokenATH
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataBytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *VyperClient) GetTokenMarket(marketId string, chainId int, interval string) (*TokenPair, error) {
	params := map[string]string{
		"chainID":  fmt.Sprintf("%d", chainId),
		"interval": interval,
	}
	body, err := c.request("GET", fmt.Sprintf("/api/v1/token/market/%s", marketId), params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result TokenPair
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataBytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *VyperClient) GetTokenHolders(marketId string, chainId int) ([]TokenHolder, int, error) {
	params := map[string]string{
		"marketID": marketId,
		"chainID":  fmt.Sprintf("%d", chainId),
	}
	body, err := c.request("GET", "/api/v1/token/holders", params)
	if err != nil {
		return nil, 0, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, 0, err
	}

	var result struct {
		Holders      []TokenHolder `json:"holders"`
		TotalHolders int           `json:"total_holders"`
	}

	dataStr, ok := apiResp.Data.(string)
	if !ok {
		return nil, 0, fmt.Errorf("unexpected data type for API response")
	}

	err = json.Unmarshal([]byte(dataStr), &result)
	if err != nil {
		return nil, 0, err
	}

	return result.Holders, result.TotalHolders, nil
}

func (c *VyperClient) GetTokenMarkets(tokenMint string, chainId int) ([]TokenMarket, error) {
	params := map[string]string{
		"tokenMint": tokenMint,
		"chainID":   fmt.Sprintf("%d", chainId),
	}
	body, err := c.request("GET", "/api/v1/token/markets", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result []TokenMarket
	dataStr, ok := apiResp.Data.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for API response")
	}

	err = json.Unmarshal([]byte(dataStr), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *VyperClient) GetWalletHoldings(walletAddress string, chainId int) ([]WalletHolding, error) {
	params := map[string]string{
		"walletAddress": walletAddress,
		"chainID":       fmt.Sprintf("%d", chainId),
	}
	body, err := c.request("GET", "/wallet/holdings", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result []WalletHolding
	dataStr, ok := apiResp.Data.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for API response")
	}

	err = json.Unmarshal([]byte(dataStr), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *VyperClient) GetWalletAggregatedPnl(walletAddress string, chainId int) (*WalletAggregatedPnL, error) {
	params := map[string]string{
		"walletAddress": walletAddress,
		"chainID":       fmt.Sprintf("%d", chainId),
	}
	body, err := c.request("GET", "/wallet/aggregated-pnl", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result WalletAggregatedPnL
	dataStr, ok := apiResp.Data.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for API response")
	}
	err = json.Unmarshal([]byte(dataStr), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *VyperClient) GetWalletPnl(walletAddress string, marketId string, chainId int) (*WalletPnL, error) {
	params := map[string]string{
		"walletAddress": walletAddress,
		"marketID":      marketId,
		"chainID":       fmt.Sprintf("%d", chainId),
	}
	body, err := c.request("GET", "/wallet/pnl", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result WalletPnL
	dataStr, ok := apiResp.Data.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for API response")
	}
	err = json.Unmarshal([]byte(dataStr), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *VyperClient) GetTokenMetadata(chainId int, tokenMint string) (*TokenMetadata, error) {
	params := map[string]string{
		"chainID":   fmt.Sprintf("%d", chainId),
		"tokenMint": tokenMint,
	}
	body, err := c.request("GET", "/api/v1/token/metadata", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result TokenMetadata
	dataStr, ok := apiResp.Data.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for API response")
	}
	err = json.Unmarshal([]byte(dataStr), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *VyperClient) GetTokenSymbol(chainId int, tokenMint string) (*TokenSymbol, error) {
	params := map[string]string{
		"chainID":   fmt.Sprintf("%d", chainId),
		"tokenMint": tokenMint,
	}
	body, err := c.request("GET", "/api/v1/token/symbol", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result TokenSymbol
	dataStr, ok := apiResp.Data.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for API response")
	}
	err = json.Unmarshal([]byte(dataStr), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *VyperClient) GetTopTraders(marketId string, chainId int) ([]TopTrader, error) {
	params := map[string]string{
		"marketID": marketId,
		"chainID":  fmt.Sprintf("%d", chainId),
	}
	body, err := c.request("GET", "/api/v1/token/top-traders", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result []TopTrader
	dataStr, ok := apiResp.Data.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for API response")
	}
	err = json.Unmarshal([]byte(dataStr), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *VyperClient) SearchTokens(criteria string, chainId *int) ([]TokenSearchResult, error) {
	params := map[string]string{
		"criteria": criteria,
	}
	if chainId != nil {
		params["chainID"] = fmt.Sprintf("%d", *chainId)
	}
	body, err := c.request("GET", "/api/v1/token/search", params)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result []TokenSearchResult
	dataStr, ok := apiResp.Data.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for API response")
	}
	err = json.Unmarshal([]byte(dataStr), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *VyperClient) GetTokenPairs(params TokenPairsParams) (*TokenPairs, error) {
	queryParams := make(map[string]string)

	v := reflect.ValueOf(params)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.IsZero() {
			tag := t.Field(i).Tag.Get("json")
			if tag != "" {
				tag = strings.Split(tag, ",")[0]
				var value string
				switch field.Kind() {
				case reflect.Ptr:
					value = fmt.Sprintf("%v", field.Elem().Interface())
				case reflect.String:
					value = field.String()
				case reflect.Slice:
					switch field.Type().Elem().Kind() {
					case reflect.Int:
						intSlice := field.Interface().([]int)
						strSlice := make([]string, len(intSlice))
						for i, v := range intSlice {
							strSlice[i] = strconv.Itoa(v)
						}
						value = strings.Join(strSlice, ",")
					case reflect.String:
						value = strings.Join(field.Interface().([]string), ",")
					}
				}
				queryParams[tag] = value
			}
		}
	}

	body, err := c.request("GET", "/token/pairs", queryParams)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var result TokenPairs
	dataStr, ok := apiResp.Data.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for API response")
	}
	err = json.Unmarshal([]byte(dataStr), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
