package service

import (
	"log"
	"math"
	"strconv"

	"github.com/NatBrian/Stockbit-Golang-Challenge/helper"
	"github.com/NatBrian/Stockbit-Golang-Challenge/model"
)

type (
	StockService struct {
	}

	IStockService interface {
		CalculateOhlc(stockCodes []string, changeRecords []model.ChangeRecord) (map[string]model.Summary, error)
	}
)

func (ss *StockService) CalculateOhlc(stockCodes []string, changeRecords []model.ChangeRecord) (map[string]model.Summary, error) {
	stockCodes = helper.RemoveDuplicateString(stockCodes)

	var (
		result = map[string]model.Summary{}
		err    error
	)

	for _, code := range stockCodes {

		// not found then init code Summary
		entry, found := result[code]
		if !found {
			entry = model.Summary{
				StockCode:   code,
				LowestPrice: math.MaxInt64,
			}
		}

		for _, record := range changeRecords {
			if record.StockCode == code {
				var (
					quantity int64
					price    int64
				)

				if record.Quantity != "" {
					quantity, err = strconv.ParseInt(record.Quantity, 10, 64)
					if err != nil {
						log.Println("error strconv.ParseInt(record.Quantity, 10, 64) ", record.Quantity, err)
						return nil, err
					}
				} else if record.ExecutedQuantity != "" {
					quantity, err = strconv.ParseInt(record.ExecutedQuantity, 10, 64)
					if err != nil {
						log.Println("error strconv.ParseInt(record.ExecutedQuantity, 10, 64) ", record.ExecutedQuantity, err)
						return nil, err
					}
				}

				if record.Price != "" {
					price, err = strconv.ParseInt(record.Price, 10, 64)
					if err != nil {
						log.Println("error strconv.ParseInt(record.Price, 10, 64) ", record.Price, err)
						return nil, err
					}
				} else if record.ExecutionPrice != "" {
					price, err = strconv.ParseInt(record.ExecutionPrice, 10, 64)
					if err != nil {
						log.Println("error strconv.ParseInt(record.ExecutionPrice, 10, 64) ", record.ExecutionPrice, err)
						return nil, err
					}
				}

				// assign prices
				if quantity == 0 {
					entry.PreviousPrice = price
				} else if quantity > 0 && entry.OpenPrice == 0 {
					entry.OpenPrice = price
				} else {
					entry.ClosePrice = price
					if entry.HighestPrice < price {
						entry.HighestPrice = price
					}
					if entry.LowestPrice > price {
						entry.LowestPrice = price
					}
				}

				// calculate volume and value
				if record.Type == "E" || record.Type == "P" {
					entry.Volume += quantity
					entry.Value += quantity * price
				}

				result[code] = entry
			}
		}
	}

	for _, code := range stockCodes {
		entry, _ := result[code]

		// assign index
		for _, indexMember := range model.IndexMembers {
			if indexMember.StockCode == code {
				entry.IndexCode = append(entry.IndexCode, indexMember.IndexCode)
			}
		}

		// calculate average price
		entry.AveragePrice = int64(math.Round(float64(entry.Value) / float64(entry.Volume)))

		result[code] = entry
	}

	return result, nil
}
