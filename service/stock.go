package service

import (
	"log"
	"strconv"

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
	var (
		result = map[string]model.Summary{}
		err    error
	)

	for _, y := range stockCodes {
		found, not := result[y]
		if not {
			found = model.Summary{}
		}
		found.StockCode = y
		for _, u := range changeRecords {
			if u.StockCode == y {
				quantity := 0
				if u.Quantity != "" {
					quantity, err = strconv.Atoi(u.Quantity)
					if err != nil {
						log.Println("error strconv.Atoi(u.Quantity) ", u.Quantity, err)
						return nil, err
					}
				}

				priceInt, err := strconv.Atoi(u.Price)
				if err != nil {
					log.Println("error strconv.Atoi(u.Price) ", u.Price, err)
					return nil, err
				}
				price := int64(priceInt)

				if quantity == 0 {
					found.Prev = price
					//fmt.Println("done")
					//fmt.Println("price updated")
					result[y] = found
				} else if quantity > 0 && result[y].Open == 0 {
					found.Open = price
					//fmt.Println("done")
					//fmt.Println("price updated")
					result[y] = found
				} else {
					found.Close = price
					if found.High < price {
						found.High = price
					}
					if found.Low > price {
						found.Low = price
					}
					//fmt.Println("done")
					//fmt.Println("price updated")
					result[y] = found
				}
			} else {
				//fmt.Println("done")
				//fmt.Println("price updated")
				result[y] = found
			}
		}
		for _, i := range model.IndexMembers {
			if i.StockCode == y {
				found.IndexCode = append(found.IndexCode, i.IndexCode)
				//fmt.Println("index updated")
				result[y] = found
			} else {
				//fmt.Println("index updated")
				result[y] = found
			}
		}
	}

	return result, nil
}
