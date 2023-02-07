package model

// ChangeRecordInput from ndjson
type ChangeRecordInput struct {
	Type             string `json:"type"`
	OrderNumber      string `json:"order_number"`
	OrderVerb        string `json:"order_verb"`
	Quantity         string `json:"quantity"`
	ExecutedQuantity string `json:"executed_quantity"`
	OrderBook        string `json:"order_book"`
	Price            string `json:"price"`
	ExecutionPrice   string `json:"execution_price"`
	StockCode        string `json:"stock_code"`
}

// ChangeRecord with correct types
type ChangeRecord struct {
	Type             string `json:"type"`
	OrderNumber      string `json:"order_number"`
	OrderVerb        string `json:"order_verb"`
	Quantity         int64  `json:"quantity"`
	ExecutedQuantity int64  `json:"executed_quantity"`
	OrderBook        string `json:"order_book"`
	Price            int64  `json:"price"`
	ExecutionPrice   int64  `json:"execution_price"`
	StockCode        string `json:"stock_code"`
}

type IndexMember struct {
	StockCode string
	IndexCode string
}

type Summary struct {
	StockCode     string   `json:"stock_code"`
	IndexCode     []string `json:"index_code"`
	PreviousPrice int64    `json:"previous_price"`
	OpenPrice     int64    `json:"open_price"`
	HighestPrice  int64    `json:"highest_price"`
	LowestPrice   int64    `json:"lowest_price"`
	ClosePrice    int64    `json:"close_price"`
	Volume        int64    `json:"volume"`
	Value         int64    `json:"value"`
	AveragePrice  int64    `json:"average_price"`
}

type SummaryResponse struct {
	Summaries map[string]Summary `json:"summaries"`
}

var (
	IndexMembers = []IndexMember{
		{
			StockCode: "BBCA",
			IndexCode: "IHSG",
		},
		{
			StockCode: "BBRI",
			IndexCode: "IHSG",
		},
		{
			StockCode: "ASII",
			IndexCode: "IHSG",
		},
		{
			StockCode: "GOTO",
			IndexCode: "IHSG",
		},
		{
			StockCode: "BBCA",
			IndexCode: "LQ45",
		},
		{
			StockCode: "BBRI",
			IndexCode: "LQ45",
		},
		{
			StockCode: "ASII",
			IndexCode: "LQ45",
		},
		{
			StockCode: "GOTO",
			IndexCode: "LQ45",
		},
		{
			StockCode: "BBCA",
			IndexCode: "KOMPAS100",
		},
		{
			StockCode: "BBRI",
			IndexCode: "KOMPAS100",
		},
	}
)

const (
	RedisKeyStockCode     = "stock-code"
	RedisKeyIndexCode     = "-index-code"
	RedisKeyPreviousPrice = "-previous-price"
	RedisKeyOpenPrice     = "-open-price"
	RedisKeyHighestPrice  = "-highest-price"
	RedisKeyLowestPrice   = "-lowest-price"
	RedisKeyClosePrice    = "-close-price"
	RedisKeyVolume        = "-volume"
	RedisKeyValue         = "-value"
	RedisKeyAveragePrice  = "-average-price"
)
