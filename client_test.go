package vyperclientgo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewVyperClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewVyperClient(apiKey)

	if client.ApiKey != apiKey {
		t.Errorf("Expected ApiKey to be %s, but got %s", apiKey, client.ApiKey)
	}

	if client.BaseURL != "https://api.vyper.trade" {
		t.Errorf("Expected BaseURL to be https://api.vyper.trade, but got %s", client.BaseURL)
	}

	if client.HttpClient == nil {
		t.Error("Expected HttpClient to be initialized, but it's nil")
	}
}

func TestGetChainIds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/chain/ids" {
			t.Errorf("Expected to request '/api/v1/chain/ids', got: %s", r.URL.Path)
		}
		response := struct {
			Status  string         `json:"status"`
			Message string         `json:"message"`
			Data    map[string]int `json:"data"`
		}{
			Status:  "success",
			Message: "Chain IDs retrieved successfully",
			Data: map[string]int{
				"solana":   900,
				"tron":     1000,
				"ethereum": 1,
				"base":     8453,
				"arbitrum": 42161,
				"bsc":      56,
				"blast":    81457,
			},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	chainIds, err := client.GetChainIds()
	if err != nil {
		t.Fatalf("GetChainIds returned an error: %v", err)
	}

	expected := map[string]int{
		"solana":   900,
		"tron":     1000,
		"ethereum": 1,
		"base":     8453,
		"arbitrum": 42161,
		"bsc":      56,
		"blast":    81457,
	}

	if !reflect.DeepEqual(chainIds, expected) {
		t.Errorf("Expected chain IDs to be %+v, but got %+v", expected, chainIds)
	}
}

func TestGetTokenAth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/token/ath" {
			t.Errorf("Expected to request '/api/v1/token/ath', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("chainID") != "1" || r.URL.Query().Get("marketID") != "test-market" {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		response := APIResponse{
			Status:  "success",
			Message: "Token ATH retrieved successfully",
			Data: TokenATH{
				MarketCapUsd:      1000000,
				Timestamp:         1625097600,
				TokenLiquidityUsd: 500000,
			},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Error encoding response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	ath, err := client.GetTokenAth(1, "test-market")
	if err != nil {
		t.Fatalf("GetTokenAth returned an error: %v", err)
	}

	expected := &TokenATH{
		MarketCapUsd:      1000000,
		Timestamp:         1625097600,
		TokenLiquidityUsd: 500000,
	}

	if !reflect.DeepEqual(ath, expected) {
		t.Errorf("Expected TokenATH to be %v, but got %v", expected, ath)
	}
}

func TestGetTokenMarket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/token/market/test-market" {
			t.Errorf("Expected to request '/api/v1/token/market/test-market', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("chainID") != "1" || r.URL.Query().Get("interval") != "1d" {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		response := APIResponse{
			Status:  "success",
			Message: "Token market data retrieved successfully",
			Data: map[string]interface{}{
				"marketId":           "test-market",
				"tokenPriceUsd":      1.5,
				"tokenLiquidityUsd":  1000000.0,
				"tokenMarketCapUsd":  5000000.0,
				"volumeUsd":          500000.0,
				"priceChangePercent": -2.5,
			},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Error encoding response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	market, err := client.GetTokenMarket("test-market", 1, "1d")
	if err != nil {
		t.Fatalf("GetTokenMarket returned an error: %v", err)
	}

	expected := &TokenPair{
		MarketId:           "test-market",
		TokenPriceUsd:      1.5,
		TokenLiquidityUsd:  1000000.0,
		TokenMarketCapUsd:  5000000.0,
		VolumeUsd:          500000.0,
		PriceChangePercent: -2.5,
	}

	if !reflect.DeepEqual(market, expected) {
		t.Errorf("Expected TokenPair to be %+v, but got %+v", expected, market)
	}
}

func TestGetTokenHolders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/token/holders" {
			t.Errorf("Expected to request '/api/v1/token/holders', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("marketID") != "test-market" || r.URL.Query().Get("chainID") != "1" {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		holdersData := struct {
			Holders      []TokenHolder `json:"holders"`
			TotalHolders int           `json:"total_holders"`
		}{
			Holders: []TokenHolder{
				{
					WalletAddress: "0x123...",
					TokenHoldings: 1000.0,
					UsdHoldings:   1500.0,
					PercentOwned:  0.1,
				},
				{
					WalletAddress: "0x456...",
					TokenHoldings: 500.0,
					UsdHoldings:   750.0,
					PercentOwned:  0.05,
				},
			},
			TotalHolders: 2,
		}
		holdersBytes, err := json.Marshal(holdersData)
		if err != nil {
			t.Fatalf("Failed to marshal holders data: %v", err)
		}
		response := APIResponse{
			Status:  "success",
			Message: "Token holders retrieved successfully",
			Data:    string(holdersBytes),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	holders, totalHolders, err := client.GetTokenHolders("test-market", 1)
	if err != nil {
		t.Fatalf("GetTokenHolders returned an error: %v", err)
	}

	expectedHolders := []TokenHolder{
		{
			WalletAddress: "0x123...",
			TokenHoldings: 1000.0,
			UsdHoldings:   1500.0,
			PercentOwned:  0.1,
		},
		{
			WalletAddress: "0x456...",
			TokenHoldings: 500.0,
			UsdHoldings:   750.0,
			PercentOwned:  0.05,
		},
	}

	if !reflect.DeepEqual(holders, expectedHolders) {
		t.Errorf("Expected holders to be %+v, but got %+v", expectedHolders, holders)
	}

	if totalHolders != 2 {
		t.Errorf("Expected totalHolders to be 2, but got %d", totalHolders)
	}
}

func TestGetTokenMarkets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/token/markets" {
			t.Errorf("Expected to request '/api/v1/token/markets', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("tokenMint") != "test-token" || r.URL.Query().Get("chainID") != "1" {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		marketsData := []TokenMarket{
			{
				MarketID:          "market1",
				TokenLiquidityUsd: 1000000.0,
				MarketCapUsd:      5000000.0,
				TokenType:         "SPL",
			},
			{
				MarketID:          "market2",
				TokenLiquidityUsd: 2000000.0,
				MarketCapUsd:      10000000.0,
				TokenType:         "ERC20",
			},
		}
		marketsBytes, err := json.Marshal(marketsData)
		if err != nil {
			t.Fatalf("Failed to marshal markets data: %v", err)
		}
		response := APIResponse{
			Status:  "success",
			Message: "Token markets retrieved successfully",
			Data:    string(marketsBytes),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	markets, err := client.GetTokenMarkets("test-token", 1)
	if err != nil {
		t.Fatalf("GetTokenMarkets returned an error: %v", err)
	}

	expectedMarkets := []TokenMarket{
		{
			MarketID:          "market1",
			TokenLiquidityUsd: 1000000.0,
			MarketCapUsd:      5000000.0,
			TokenType:         "SPL",
		},
		{
			MarketID:          "market2",
			TokenLiquidityUsd: 2000000.0,
			MarketCapUsd:      10000000.0,
			TokenType:         "ERC20",
		},
	}

	if !reflect.DeepEqual(markets, expectedMarkets) {
		t.Errorf("Expected markets to be %+v, but got %+v", expectedMarkets, markets)
	}
}

func TestGetWalletHoldings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/wallet/holdings" {
			t.Errorf("Expected to request '/wallet/holdings', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("walletAddress") != "0x123..." || r.URL.Query().Get("chainID") != "1" {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		holdingsData := []WalletHolding{
			{
				MarketId:      "market1",
				TokenHoldings: 1000.0,
				UsdValue:      1500.0,
				TokenSymbol:   "TKN1",
			},
			{
				MarketId:      "market2",
				TokenHoldings: 500.0,
				UsdValue:      750.0,
				TokenSymbol:   "TKN2",
			},
		}
		holdingsBytes, err := json.Marshal(holdingsData)
		if err != nil {
			t.Fatalf("Failed to marshal holdings data: %v", err)
		}
		response := APIResponse{
			Status:  "success",
			Message: "Wallet holdings retrieved successfully",
			Data:    string(holdingsBytes),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	holdings, err := client.GetWalletHoldings("0x123...", 1)
	if err != nil {
		t.Fatalf("GetWalletHoldings returned an error: %v", err)
	}

	expectedHoldings := []WalletHolding{
		{
			MarketId:      "market1",
			TokenHoldings: 1000.0,
			UsdValue:      1500.0,
			TokenSymbol:   "TKN1",
		},
		{
			MarketId:      "market2",
			TokenHoldings: 500.0,
			UsdValue:      750.0,
			TokenSymbol:   "TKN2",
		},
	}

	if !reflect.DeepEqual(holdings, expectedHoldings) {
		t.Errorf("Expected holdings to be %+v, but got %+v", expectedHoldings, holdings)
	}
}

func TestGetWalletAggregatedPnl(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/wallet/aggregated-pnl" {
			t.Errorf("Expected to request '/wallet/aggregated-pnl', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("walletAddress") != "0x123..." || r.URL.Query().Get("chainID") != "1" {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		pnlData := WalletAggregatedPnL{
			InvestedAmount:       1000.0,
			PnlPercent:           10.5,
			PnlUsd:               105.0,
			SoldAmount:           500.0,
			TokensTraded:         5,
			TotalPnlPercent:      15.0,
			TotalPnlUsd:          150.0,
			UnrealizedPnlPercent: 5.0,
			UnrealizedPnlUsd:     50.0,
		}
		pnlBytes, err := json.Marshal(pnlData)
		if err != nil {
			t.Fatalf("Failed to marshal PNL data: %v", err)
		}
		response := APIResponse{
			Status:  "success",
			Message: "Wallet aggregated PNL retrieved successfully",
			Data:    string(pnlBytes),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	pnl, err := client.GetWalletAggregatedPnl("0x123...", 1)
	if err != nil {
		t.Fatalf("GetWalletAggregatedPnl returned an error: %v", err)
	}

	expected := &WalletAggregatedPnL{
		InvestedAmount:       1000.0,
		PnlPercent:           10.5,
		PnlUsd:               105.0,
		SoldAmount:           500.0,
		TokensTraded:         5,
		TotalPnlPercent:      15.0,
		TotalPnlUsd:          150.0,
		UnrealizedPnlPercent: 5.0,
		UnrealizedPnlUsd:     50.0,
	}

	if !reflect.DeepEqual(pnl, expected) {
		t.Errorf("Expected PNL to be %+v, but got %+v", expected, pnl)
	}
}

func TestGetWalletPnl(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/wallet/pnl" {
			t.Errorf("Expected to request '/wallet/pnl', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("walletAddress") != "0x123..." || r.URL.Query().Get("marketID") != "market1" || r.URL.Query().Get("chainID") != "1" {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		pnlData := WalletPnL{
			HolderSince:     1625097600,
			InvestedAmount:  1000.0,
			InvestedTxns:    5,
			PnlPercent:      10.5,
			PnlUsd:          105.0,
			RemainingTokens: 100.0,
			RemainingUsd:    150.0,
			SoldAmount:      500.0,
			SoldTxns:        3,
		}
		pnlBytes, err := json.Marshal(pnlData)
		if err != nil {
			t.Fatalf("Failed to marshal PNL data: %v", err)
		}
		response := APIResponse{
			Status:  "success",
			Message: "Wallet PNL retrieved successfully",
			Data:    string(pnlBytes),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	pnl, err := client.GetWalletPnl("0x123...", "market1", 1)
	if err != nil {
		t.Fatalf("GetWalletPnl returned an error: %v", err)
	}

	expected := &WalletPnL{
		HolderSince:     1625097600,
		InvestedAmount:  1000.0,
		InvestedTxns:    5,
		PnlPercent:      10.5,
		PnlUsd:          105.0,
		RemainingTokens: 100.0,
		RemainingUsd:    150.0,
		SoldAmount:      500.0,
		SoldTxns:        3,
	}

	if !reflect.DeepEqual(pnl, expected) {
		t.Errorf("Expected PNL to be %+v, but got %+v", expected, pnl)
	}
}

func TestGetTokenMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/token/metadata" {
			t.Errorf("Expected to request '/api/v1/token/metadata', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("chainID") != "1" || r.URL.Query().Get("tokenMint") != "0xabc..." {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		metadataData := TokenMetadata{
			Image:    "https://example.com/token.png",
			Name:     "Example Token",
			Symbol:   "EXT",
			Telegram: "https://t.me/exampletoken",
			Twitter:  "https://twitter.com/exampletoken",
			Website:  "https://exampletoken.com",
		}
		metadataBytes, err := json.Marshal(metadataData)
		if err != nil {
			t.Fatalf("Failed to marshal metadata: %v", err)
		}
		response := APIResponse{
			Status:  "success",
			Message: "Token metadata retrieved successfully",
			Data:    string(metadataBytes),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	metadata, err := client.GetTokenMetadata(1, "0xabc...")
	if err != nil {
		t.Fatalf("GetTokenMetadata returned an error: %v", err)
	}

	expected := &TokenMetadata{
		Image:    "https://example.com/token.png",
		Name:     "Example Token",
		Symbol:   "EXT",
		Telegram: "https://t.me/exampletoken",
		Twitter:  "https://twitter.com/exampletoken",
		Website:  "https://exampletoken.com",
	}

	if !reflect.DeepEqual(metadata, expected) {
		t.Errorf("Expected metadata to be %+v, but got %+v", expected, metadata)
	}
}

func TestGetTokenSymbol(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/token/symbol" {
			t.Errorf("Expected to request '/api/v1/token/symbol', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("chainID") != "1" || r.URL.Query().Get("tokenMint") != "0xabc..." {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		symbolData := TokenSymbol{
			Symbol: "EXT",
		}
		symbolBytes, err := json.Marshal(symbolData)
		if err != nil {
			t.Fatalf("Failed to marshal symbol data: %v", err)
		}
		response := APIResponse{
			Status:  "success",
			Message: "Token symbol retrieved successfully",
			Data:    string(symbolBytes),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	symbol, err := client.GetTokenSymbol(1, "0xabc...")
	if err != nil {
		t.Fatalf("GetTokenSymbol returned an error: %v", err)
	}

	expected := &TokenSymbol{
		Symbol: "EXT",
	}

	if !reflect.DeepEqual(symbol, expected) {
		t.Errorf("Expected symbol to be %+v, but got %+v", expected, symbol)
	}
}

func TestGetTopTraders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/token/top-traders" {
			t.Errorf("Expected to request '/api/v1/token/top-traders', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("marketID") != "market1" || r.URL.Query().Get("chainID") != "1" {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		tradersData := []TopTrader{
			{
				InvestedAmountTokens: 1000.0,
				InvestedAmountUsd:    1500.0,
				InvestedTxns:         5,
				PnlUsd:               200.0,
				RemainingTokens:      800.0,
				RemainingUsd:         1200.0,
				SoldAmountTokens:     200.0,
				SoldAmountUsd:        300.0,
				SoldTxns:             2,
				WalletAddress:        "0x123...",
				WalletTag:            "Trader1",
			},
		}
		tradersBytes, err := json.Marshal(tradersData)
		if err != nil {
			t.Fatalf("Failed to marshal traders data: %v", err)
		}
		response := APIResponse{
			Status:  "success",
			Message: "Top traders retrieved successfully",
			Data:    string(tradersBytes),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	traders, err := client.GetTopTraders("market1", 1)
	if err != nil {
		t.Fatalf("GetTopTraders returned an error: %v", err)
	}

	expected := []TopTrader{
		{
			InvestedAmountTokens: 1000.0,
			InvestedAmountUsd:    1500.0,
			InvestedTxns:         5,
			PnlUsd:               200.0,
			RemainingTokens:      800.0,
			RemainingUsd:         1200.0,
			SoldAmountTokens:     200.0,
			SoldAmountUsd:        300.0,
			SoldTxns:             2,
			WalletAddress:        "0x123...",
			WalletTag:            "Trader1",
		},
	}

	if !reflect.DeepEqual(traders, expected) {
		t.Errorf("Expected traders to be %+v, but got %+v", expected, traders)
	}
}

func TestSearchTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/token/search" {
			t.Errorf("Expected to request '/api/v1/token/search', got: %s", r.URL.Path)
		}
		if r.URL.Query().Get("criteria") != "test" || r.URL.Query().Get("chainID") != "1" {
			t.Errorf("Unexpected query parameters: %v", r.URL.Query())
		}
		searchData := []TokenSearchResult{
			{
				ChainId:           1,
				MarketId:          "market1",
				CreatedTimestamp:  1625097600,
				Name:              "Test Token",
				Symbol:            "TEST",
				TokenMint:         "0xabc...",
				TokenType:         "ERC20",
				PercentChange24h:  5.5,
				PooledAsset:       1000.0,
				TokenLiquidityUsd: 10000.0,
				TokenMarketCapUsd: 100000.0,
				TokenPriceUsd:     1.5,
				VolumeUsd:         50000.0,
			},
		}
		searchBytes, err := json.Marshal(searchData)
		if err != nil {
			t.Fatalf("Failed to marshal search data: %v", err)
		}
		response := APIResponse{
			Status:  "success",
			Message: "Token search completed successfully",
			Data:    string(searchBytes),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	chainId := 1
	results, err := client.SearchTokens("test", &chainId)
	if err != nil {
		t.Fatalf("SearchTokens returned an error: %v", err)
	}

	expected := []TokenSearchResult{
		{
			ChainId:           1,
			MarketId:          "market1",
			CreatedTimestamp:  1625097600,
			Name:              "Test Token",
			Symbol:            "TEST",
			TokenMint:         "0xabc...",
			TokenType:         "ERC20",
			PercentChange24h:  5.5,
			PooledAsset:       1000.0,
			TokenLiquidityUsd: 10000.0,
			TokenMarketCapUsd: 100000.0,
			TokenPriceUsd:     1.5,
			VolumeUsd:         50000.0,
		},
	}

	if !reflect.DeepEqual(results, expected) {
		t.Errorf("Expected search results to be %+v, but got %+v", expected, results)
	}
}

func TestGetTokenPairs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/token/pairs" {
			t.Errorf("Expected to request '/token/pairs', got: %s", r.URL.Path)
		}
		pairsData := TokenPairs{
			HasNext: true,
			Pairs: []TokenPair{
				{
					MarketId:           "market1",
					ChainId:            1,
					Name:               "Test Token",
					Symbol:             "TEST",
					TokenMint:          "0xabc...",
					TokenType:          "ERC20",
					TokenPriceUsd:      1.5,
					TokenLiquidityUsd:  10000.0,
					TokenMarketCapUsd:  100000.0,
					VolumeUsd:          50000.0,
					PriceChangePercent: 5.5,
					// ... more fields ...
				},
			},
		}
		pairsBytes, err := json.Marshal(pairsData)
		if err != nil {
			t.Fatalf("Failed to marshal pairs data: %v", err)
		}
		response := APIResponse{
			Status:  "success",
			Message: "Token pairs retrieved successfully",
			Data:    string(pairsBytes),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &VyperClient{
		BaseURL:    server.URL,
		ApiKey:     "test-api-key",
		HttpClient: server.Client(),
	}

	params := TokenPairsParams{
		ChainIds: []int{1},
		Sorting:  "volume",
	}
	pairs, err := client.GetTokenPairs(params)
	if err != nil {
		t.Fatalf("GetTokenPairs returned an error: %v", err)
	}

	expected := &TokenPairs{
		HasNext: true,
		Pairs: []TokenPair{
			{
				MarketId:           "market1",
				ChainId:            1,
				Name:               "Test Token",
				Symbol:             "TEST",
				TokenMint:          "0xabc...",
				TokenType:          "ERC20",
				TokenPriceUsd:      1.5,
				TokenLiquidityUsd:  10000.0,
				TokenMarketCapUsd:  100000.0,
				VolumeUsd:          50000.0,
				PriceChangePercent: 5.5,
				// ... more fields ...
			},
		},
	}

	if !reflect.DeepEqual(pairs, expected) {
		t.Errorf("Expected token pairs to be %+v, but got %+v", expected, pairs)
	}
}
