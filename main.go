package main

import (
	"encoding/json"
	"fmt"
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

func ohlc(x []string, w []ChangeRecord, p []IndexMember) map[string]Summary {
	for _, y := range x {
		found, not := result[y]
		if not {
			found = Summary{}
		}
		found.StockCode = y
		for _, u := range w {
			if u.StockCode == y {
				if u.Quantity == 0 {
					found.Prev = u.Price
					fmt.Println("done")
					fmt.Println("price updated")
					result[y] = found
				} else if u.Quantity > 0 && result[y].Open == 0 {
					found.Open = u.Price
					fmt.Println("done")
					fmt.Println("price updated")
					result[y] = found
				} else {
					found.Close = u.Price
					if found.High < u.Price {
						found.High = u.Price
					}
					if found.Low > u.Price {
						found.Low = u.Price
					}
					fmt.Println("done")
					fmt.Println("price updated")
					result[y] = found
				}
			} else {
				fmt.Println("done")
				fmt.Println("price updated")
				result[y] = found
			}
		}
		for _, i := range p {
			if i.StockCode == y {
				found.IndexCode = append(found.IndexCode, i.IndexCode)
				fmt.Println("index updated")
				result[y] = found
			} else {
				fmt.Println("index updated")
				result[y] = found
			}
		}
	}

	return result
}

func main() {
	x := []string{
		"BBCA", "BBRI", "ASII", "GOTO",
	}
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
	r := ohlc(x, w, p)
	for _, v := range r {
		jss, _ := json.Marshal(v)
		fmt.Println("summary: ", string(jss))
	}
}
