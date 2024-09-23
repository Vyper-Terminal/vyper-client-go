package vyperclientgo

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type WalletAggregatedPnL struct {
	InvestedAmount       float64 `json:"investedAmount"`
	PnlPercent           float64 `json:"pnlPercent"`
	PnlUsd               float64 `json:"pnlUsd"`
	SoldAmount           float64 `json:"soldAmount"`
	TokensTraded         int     `json:"tokensTraded"`
	TotalPnlPercent      float64 `json:"totalPnlPercent"`
	TotalPnlUsd          float64 `json:"totalPnlUsd"`
	UnrealizedPnlPercent float64 `json:"unrealizedPnlPercent"`
	UnrealizedPnlUsd     float64 `json:"unrealizedPnlUsd"`
}

type WalletHolding struct {
	MarketId      string  `json:"marketId"`
	TokenHoldings float64 `json:"tokenHoldings"`
	TokenSymbol   string  `json:"tokenSymbol"`
	UsdValue      float64 `json:"usdValue"`
}

type WalletPnL struct {
	HolderSince     int64   `json:"holderSince"`
	InvestedAmount  float64 `json:"investedAmount"`
	InvestedTxns    int     `json:"investedTxns"`
	PnlPercent      float64 `json:"pnlPercent"`
	PnlUsd          float64 `json:"pnlUsd"`
	RemainingTokens float64 `json:"remainingTokens"`
	RemainingUsd    float64 `json:"remainingUsd"`
	SoldAmount      float64 `json:"soldAmount"`
	SoldTxns        int     `json:"soldTxns"`
}

type TopTrader struct {
	InvestedAmountTokens float64 `json:"investedAmount_tokens"`
	InvestedAmountUsd    float64 `json:"investedAmount_usd"`
	InvestedTxns         int     `json:"investedTxns"`
	PnlUsd               float64 `json:"pnlUsd"`
	RemainingTokens      float64 `json:"remainingTokens"`
	RemainingUsd         float64 `json:"remainingUsd"`
	SoldAmountTokens     float64 `json:"soldAmountTokens"`
	SoldAmountUsd        float64 `json:"soldAmountUsd"`
	SoldTxns             int     `json:"soldTxns"`
	WalletAddress        string  `json:"walletAddress"`
	WalletTag            string  `json:"walletTag,omitempty"`
}

type TokenSearchResult struct {
	ChainId           int     `json:"chainId"`
	MarketId          string  `json:"marketId"`
	CreatedTimestamp  int64   `json:"createdTimestamp"`
	Name              string  `json:"name"`
	Symbol            string  `json:"symbol"`
	TokenMint         string  `json:"tokenMint"`
	TokenType         string  `json:"tokenType"`
	PercentChange24h  float64 `json:"percentChange24h"`
	PooledAsset       float64 `json:"pooledAsset"`
	TokenLiquidityUsd float64 `json:"tokenLiquidityUsd"`
	TokenMarketCapUsd float64 `json:"tokenMarketCapUsd"`
	TokenPriceUsd     float64 `json:"tokenPriceUsd"`
	VolumeUsd         float64 `json:"volumeUsd"`
	Image             string  `json:"image,omitempty"`
	Telegram          string  `json:"telegram,omitempty"`
	Twitter           string  `json:"twitter,omitempty"`
	Website           string  `json:"website,omitempty"`
}

type TokenMarket struct {
	MarketCapUsd      float64 `json:"marketCapUsd"`
	MarketID          string  `json:"marketID"`
	TokenLiquidityUsd float64 `json:"tokenLiquidityUsd"`
	TokenType         string  `json:"tokenType"`
}

type TokenMetadata struct {
	Image    string `json:"image,omitempty"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Telegram string `json:"telegram,omitempty"`
	Twitter  string `json:"twitter,omitempty"`
	Website  string `json:"website,omitempty"`
}

type TokenSymbol struct {
	Symbol string `json:"symbol"`
}

type TokenHolder struct {
	PercentOwned  float64 `json:"percentOwned"`
	TokenHoldings float64 `json:"tokenHoldings"`
	UsdHoldings   float64 `json:"usdHoldings"`
	WalletAddress string  `json:"walletAddress"`
	WalletTag     string  `json:"walletTag,omitempty"`
}

type TokenATH struct {
	MarketCapUsd      float64 `json:"marketCapUsd"`
	Timestamp         int64   `json:"timestamp"`
	TokenLiquidityUsd float64 `json:"tokenLiquidityUsd"`
}

type MigrationState struct {
	DurationMinutes    int     `json:"durationMinutes"`
	Makers             int     `json:"makers"`
	MigrationTimestamp int64   `json:"migrationTimestamp"`
	Volume             float64 `json:"volume"`
}

type TokenPair struct {
	Abused                 *bool           `json:"abused,omitempty"`
	BondingCurvePercentage *float64        `json:"bondingCurvePercentage,omitempty"`
	BuyTxnCount            int             `json:"buyTxnCount"`
	ChainId                int             `json:"chainId"`
	ContractCreator        string          `json:"contractCreator"`
	CreatedTimestamp       int64           `json:"createdTimestamp"`
	Description            string          `json:"description,omitempty"`
	FreezeAuthority        *bool           `json:"freezeAuthority,omitempty"`
	Image                  string          `json:"image,omitempty"`
	InitialAssetLiquidity  float64         `json:"initialAssetLiquidity"`
	InitialUsdLiquidity    float64         `json:"initialUsdLiquidity"`
	IsMigrated             *bool           `json:"isMigrated,omitempty"`
	LpBurned               bool            `json:"lpBurned"`
	LpCreator              string          `json:"lpCreator"`
	MarketId               string          `json:"marketId"`
	MetadataUri            string          `json:"metadataUri,omitempty"`
	MigratedMarketId       string          `json:"migratedMarketId,omitempty"`
	MigrationState         *MigrationState `json:"migrationState,omitempty"`
	MintAuthority          *bool           `json:"mintAuthority,omitempty"`
	Name                   string          `json:"name"`
	PooledAsset            float64         `json:"pooledAsset"`
	PooledToken            float64         `json:"pooledToken"`
	PriceChangePercent     float64         `json:"priceChangePercent"`
	SellTxnCount           int             `json:"sellTxnCount"`
	Symbol                 string          `json:"symbol"`
	Telegram               string          `json:"telegram,omitempty"`
	TokenLiquidityAsset    float64         `json:"tokenLiquidityAsset"`
	TokenLiquidityUsd      float64         `json:"tokenLiquidityUsd"`
	TokenMarketCapAsset    float64         `json:"tokenMarketCapAsset"`
	TokenMarketCapUsd      float64         `json:"tokenMarketCapUsd"`
	TokenMint              string          `json:"tokenMint"`
	TokenPriceAsset        float64         `json:"tokenPriceAsset"`
	TokenPriceUsd          float64         `json:"tokenPriceUsd"`
	TokenType              string          `json:"tokenType"`
	Top10HoldingPercent    float64         `json:"top10HoldingPercent"`
	TotalSupply            float64         `json:"totalSupply"`
	TransactionCount       int             `json:"transactionCount"`
	Twitter                string          `json:"twitter,omitempty"`
	VolumeAsset            float64         `json:"volumeAsset"`
	VolumeUsd              float64         `json:"volumeUsd"`
	Website                string          `json:"website,omitempty"`
}

type TokenPairs struct {
	HasNext bool        `json:"hasNext"`
	Pairs   []TokenPair `json:"pairs"`
}

type ChainAction struct {
	Signer              string  `json:"signer"`
	TokenAccount        string  `json:"tokenAccount,omitempty"`
	TransactionId       string  `json:"transactionId"`
	TokenMint           string  `json:"tokenMint,omitempty"`
	MarketId            string  `json:"marketId"`
	ActionType          string  `json:"actionType"`
	TokenAmount         float64 `json:"tokenAmount"`
	AssetAmount         float64 `json:"assetAmount"`
	TokenPriceUsd       float64 `json:"tokenPriceUsd"`
	TokenPriceAsset     float64 `json:"tokenPriceAsset"`
	SwapTotalUsd        float64 `json:"swapTotalUsd,omitempty"`
	SwapTotalAsset      float64 `json:"swapTotalAsset,omitempty"`
	TokenMarketCapAsset float64 `json:"tokenMarketCapAsset"`
	TokenMarketCapUsd   float64 `json:"tokenMarketCapUsd"`
	TokenLiquidityAsset float64 `json:"tokenLiquidityAsset"`
	TokenLiquidityUsd   float64 `json:"tokenLiquidityUsd"`
	PooledToken         float64 `json:"pooledToken"`
	PooledAsset         float64 `json:"pooledAsset"`
	ActionTimestamp     int64   `json:"actionTimestamp"`
	BondingCurvePercent float64 `json:"bondingCurvePercentage,omitempty"`
	BotUsed             string  `json:"botUsed,omitempty"`
}

type TokenPairsParams struct {
	AtLeastOneSocial    *bool    `json:"atLeastOneSocial,omitempty"`
	BuysMax             *int     `json:"buysMax,omitempty"`
	BuysMin             *int     `json:"buysMin,omitempty"`
	ChainIds            []int    `json:"chainIds,omitempty"`
	FreezeAuthDisabled  *bool    `json:"freezeAuthDisabled,omitempty"`
	InitialLiquidityMax *float64 `json:"initialLiquidityMax,omitempty"`
	InitialLiquidityMin *float64 `json:"initialLiquidityMin,omitempty"`
	Interval            string   `json:"interval,omitempty"`
	LiquidityMax        *float64 `json:"liquidityMax,omitempty"`
	LiquidityMin        *float64 `json:"liquidityMin,omitempty"`
	LpBurned            *bool    `json:"lpBurned,omitempty"`
	MarketCapMax        *float64 `json:"marketCapMax,omitempty"`
	MarketCapMin        *float64 `json:"marketCapMin,omitempty"`
	MintAuthDisabled    *bool    `json:"mintAuthDisabled,omitempty"`
	Page                *int     `json:"page,omitempty"`
	SellsMax            *int     `json:"sellsMax,omitempty"`
	SellsMin            *int     `json:"sellsMin,omitempty"`
	Sorting             string   `json:"sorting,omitempty"`
	SwapsMax            *int     `json:"swapsMax,omitempty"`
	SwapsMin            *int     `json:"swapsMin,omitempty"`
	TokenTypes          []string `json:"tokenTypes,omitempty"`
	Top10Holders        *bool    `json:"top10Holders,omitempty"`
	VolumeMax           *float64 `json:"volumeMax,omitempty"`
	VolumeMin           *float64 `json:"volumeMin,omitempty"`
}
