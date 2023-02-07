package instructions

import (
	"encoding/json"
	"fmt"
	"math"
)

type ChangeRecord struct {
	StockCode string
	Price     int64
	Quantity  int64
}

type IndexMember struct {
	StockCode string
	IndexCode string
}

type Summary struct {
	StockCode string   `json:"stock_code"`
	IndexCode []string `json:"index_code"`
	Open      int64    `json:"open"`
	High      int64    `json:"high"`
	Low       int64    `json:"low"`
	Close     int64    `json:"close"`
	Prev      int64    `json:"prev"`
}

var result = map[string]Summary{}

func removeDuplicateString(oldList []string) []string {
	keys := make(map[string]bool)
	var newList []string

	for _, entry := range oldList {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			newList = append(newList, entry)
		}
	}
	return newList
}

func ohlc(changeRecords []ChangeRecord, indexMembers []IndexMember) map[string]Summary {
	var stockCodes []string
	for _, record := range changeRecords {
		stockCodes = append(stockCodes, record.StockCode)
	}

	stockCodes = removeDuplicateString(stockCodes)

	for _, code := range stockCodes {
		// not found then init code Summary
		entry, found := result[code]
		if !found {
			entry = Summary{
				StockCode: code,
				Low:       math.MaxInt64,
			}
		}

		for _, record := range changeRecords {
			if record.StockCode == code {
				var (
					quantity int64
					price    int64
				)

				quantity = record.Quantity
				price = record.Price

				// assign prices
				if quantity == 0 {
					entry.Prev = price
				} else if quantity > 0 && entry.Open == 0 {
					entry.Open = price
				} else {
					entry.Close = price
					if entry.High < price {
						entry.High = price
					}
					if entry.Low > price {
						entry.Low = price
					}
				}
			}
		}

		for _, indexMember := range indexMembers {
			if indexMember.StockCode == code {
				entry.IndexCode = append(entry.IndexCode, indexMember.IndexCode)
			}
		}

		result[code] = entry
	}

	return result
}

func main() {
	// x := []string{
	//	"BBCA", "BBRI", "ASII", "GOTO",
	// }
	w := []ChangeRecord{
		{
			StockCode: "BBCA",
			Price:     8783,
			Quantity:  0,
		},
		{
			StockCode: "BBRI",
			Price:     3233,
			Quantity:  0,
		},
		{
			StockCode: "ASII",
			Price:     1223,
			Quantity:  0,
		},
		{
			StockCode: "GOTO",
			Price:     321,
			Quantity:  0,
		},

		{
			StockCode: "BBCA",
			Price:     8780,
			Quantity:  1,
		},
		{
			StockCode: "BBRI",
			Price:     3230,
			Quantity:  1,
		},
		{
			StockCode: "ASII",
			Price:     1220,
			Quantity:  1,
		},
		{
			StockCode: "GOTO",
			Price:     320,
			Quantity:  1,
		},

		{
			StockCode: "BBCA",
			Price:     8800,
			Quantity:  1,
		},
		{
			StockCode: "BBRI",
			Price:     3300,
			Quantity:  1,
		},
		{
			StockCode: "ASII",
			Price:     1300,
			Quantity:  1,
		},
		{
			StockCode: "GOTO",
			Price:     330,
			Quantity:  1,
		},

		{
			StockCode: "BBCA",
			Price:     8600,
			Quantity:  1,
		},
		{
			StockCode: "BBRI",
			Price:     3100,
			Quantity:  1,
		},
		{
			StockCode: "ASII",
			Price:     1100,
			Quantity:  1,
		},
		{
			StockCode: "GOTO",
			Price:     310,
			Quantity:  1,
		},

		{
			StockCode: "BBCA",
			Price:     8785,
			Quantity:  1,
		},
		{
			StockCode: "BBRI",
			Price:     3235,
			Quantity:  1,
		},
		{
			StockCode: "ASII",
			Price:     1225,
			Quantity:  1,
		},
		{
			StockCode: "GOTO",
			Price:     325,
			Quantity:  1,
		},
	}
	p := []IndexMember{
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
	r := ohlc(w, p)
	for _, v := range r {
		jss, _ := json.Marshal(v)
		fmt.Println("summary: ", string(jss))
	}
}
