package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/NatBrian/Stockbit-Golang-Challenge/config"
	"github.com/NatBrian/Stockbit-Golang-Challenge/helper"
	"github.com/NatBrian/Stockbit-Golang-Challenge/kafka"
	"github.com/NatBrian/Stockbit-Golang-Challenge/model"
	__ "github.com/NatBrian/Stockbit-Golang-Challenge/pb"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type (
	StockService struct {
		KafkaProducer kafka.Producer
		Config        config.Config
		Redis         *redis.Client
		Context       context.Context
	}

	IStockService interface {
		CalculateOhlc(stockCodes []string, changeRecords []model.ChangeRecord) (map[string]model.Summary, error)
		ProduceRecords(changeRecords []model.ChangeRecordInput) error
		GetSummary() (map[string]model.Summary, error)
	}
)

func (ss *StockService) CalculateOhlc(stockCodes []string, changeRecords []model.ChangeRecord) (map[string]model.Summary, error) {
	stockCodesByteFromRedis, err := ss.Redis.Get(ss.Context, model.RedisKeyStockCode).Bytes()
	var stockCodesFromRedis []string
	if err == nil {
		err = json.Unmarshal(stockCodesByteFromRedis, &stockCodesFromRedis)
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: json.Unmarshal(stockCodesByte)")
			return nil, err
		}
	}

	stockCodes = append(stockCodes, stockCodesFromRedis...)
	stockCodes = helper.RemoveDuplicateString(stockCodes)

	stockCodesByte, err := json.Marshal(stockCodes)
	if err != nil {
		log.Error().Err(err).Msg("error CalculateOhlc: json.Marshal(stockCodes)")
	}
	err = ss.Redis.Set(ss.Context, model.RedisKeyStockCode, stockCodesByte, 0).Err()
	if err != nil {
		log.Error().Err(err).Msg("error CalculateOhlc: set StockCode")
		return nil, err
	}

	var (
		result = map[string]model.Summary{}
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
				openPrice, err := ss.Redis.Get(ss.Context, record.StockCode+model.RedisKeyOpenPrice).Int64()
				if err != nil {
					openPrice = 0
				}

				var (
					quantity int64
					price    int64
				)

				if record.Quantity != 0 {
					quantity = record.Quantity
				} else if record.ExecutedQuantity != 0 {
					quantity = record.ExecutedQuantity
				}

				if record.Price != 0 {
					price = record.Price
				} else if record.ExecutionPrice != 0 {
					price = record.ExecutionPrice
				}

				// assign prices
				// type A is only used by previous price
				if quantity == 0 {
					entry.PreviousPrice = price
				} else if quantity > 0 && openPrice == 0 && record.Type != "A" {
					entry.OpenPrice = price
					openPrice = price
					err = ss.Redis.Set(ss.Context, code+model.RedisKeyOpenPrice, openPrice, 0).Err()
					log.Info().Msg(fmt.Sprintf("======= set open price: %d", openPrice))
					if err != nil {
						log.Error().Err(err).Msg("error CalculateOhlc: set openPrice")
						return nil, err
					}
				} else {
					if record.Type != "A" {
						entry.ClosePrice = price
						if entry.HighestPrice < price {
							entry.HighestPrice = price
						}
						if entry.LowestPrice > price {
							entry.LowestPrice = price
						}
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

	for _, code := range stockCodes {
		indexCodeByte, err := ss.Redis.Get(ss.Context, code+model.RedisKeyIndexCode).Bytes()
		if err == nil {
			var indexCode []string
			err := json.Unmarshal(indexCodeByte, &indexCode)
			if err != nil {
				log.Error().Err(err).Msg("error marshal indexCode")
				return nil, err
			}

			indexCode = append(indexCode, result[code].IndexCode...)
			indexCode = helper.RemoveDuplicateString(indexCode)

			indexCodeByte, err = json.Marshal(indexCode)
			if err != nil {
				log.Error().Err(err).Msg("error CalculateOhlc: json.Marshal(indexCode)")
			}

			err = ss.Redis.Set(ss.Context, code+model.RedisKeyIndexCode, indexCodeByte, 0).Err()
			if err != nil {
				log.Error().Err(err).Msg("error CalculateOhlc: set IndexCode from indexCode")
				return nil, err
			}
		} else {
			indexCodeByte, err = json.Marshal(result[code].IndexCode)
			if err != nil {
				log.Error().Err(err).Msg("error CalculateOhlc: json.Marshal(result[code].IndexCode)")
			}

			err = ss.Redis.Set(ss.Context, code+model.RedisKeyIndexCode, indexCodeByte, 0).Err()
			if err != nil {
				log.Error().Err(err).Msg("error CalculateOhlc: set IndexCode from result[code].IndexCode")
				return nil, err
			}
		}

		// assuming that previous price only appear once
		previousPrice, err := ss.Redis.Get(ss.Context, code+model.RedisKeyPreviousPrice).Int64()
		if result[code].PreviousPrice != 0 || err != nil { // no previous price assigned to redis
			err = ss.Redis.Set(ss.Context, code+model.RedisKeyPreviousPrice, result[code].PreviousPrice, 0).Err()
			if err != nil {
				log.Error().Err(err).Msg("error CalculateOhlc: set PreviousPrice")
				return nil, err
			}
		} else { // set previous price = 0
			if previousPrice == 0 {
				previousPrice = result[code].PreviousPrice
			}
			err = ss.Redis.Set(ss.Context, code+model.RedisKeyPreviousPrice, previousPrice, 0).Err()
			if err != nil {
				log.Error().Err(err).Msg("error CalculateOhlc: set PreviousPrice")
				return nil, err
			}
		}

		//if result[code].OpenPrice != 0 {
		//	err = ss.Redis.Set(ss.Context, code+model.RedisKeyOpenPrice, result[code].OpenPrice, 0).Err()
		//	if err != nil {
		//		log.Error().Err(err).Msg("error CalculateOhlc: set OpenPrice")
		//		return nil, err
		//	}
		//}

		if result[code].HighestPrice != 0 {
			highestPrice, err := ss.Redis.Get(ss.Context, code+model.RedisKeyHighestPrice).Int64()
			if err != nil {
				highestPrice = 0
			}
			if result[code].HighestPrice > highestPrice {
				err = ss.Redis.Set(ss.Context, code+model.RedisKeyHighestPrice, result[code].HighestPrice, 0).Err()
				if err != nil {
					log.Error().Err(err).Msg("error CalculateOhlc: set HighestPrice")
					return nil, err
				}
			} else {
				err = ss.Redis.Set(ss.Context, code+model.RedisKeyHighestPrice, highestPrice, 0).Err()
				if err != nil {
					log.Error().Err(err).Msg("error CalculateOhlc: set HighestPrice")
					return nil, err
				}
			}
		}

		if result[code].LowestPrice != 0 {
			lowestPrice, err := ss.Redis.Get(ss.Context, code+model.RedisKeyLowestPrice).Int64()
			if err != nil {
				lowestPrice = result[code].LowestPrice
			}

			if result[code].LowestPrice < lowestPrice {
				err = ss.Redis.Set(ss.Context, code+model.RedisKeyLowestPrice, result[code].LowestPrice, 0).Err()
				if err != nil {
					log.Error().Err(err).Msg("error CalculateOhlc: set LowestPrice")
					return nil, err
				}
			} else {
				err = ss.Redis.Set(ss.Context, code+model.RedisKeyLowestPrice, lowestPrice, 0).Err()
				if err != nil {
					log.Error().Err(err).Msg("error CalculateOhlc: set LowestPrice")
					return nil, err
				}
			}
		}

		if result[code].ClosePrice != 0 {
			err = ss.Redis.Set(ss.Context, code+model.RedisKeyClosePrice, result[code].ClosePrice, 0).Err()
			if err != nil {
				log.Error().Err(err).Msg("error CalculateOhlc: set ClosePrice")
				return nil, err
			}
		}

		volume, err := ss.Redis.Get(ss.Context, code+model.RedisKeyVolume).Int64()
		if err != nil {
			volume = 0
		}
		volume += result[code].Volume
		err = ss.Redis.Set(ss.Context, code+model.RedisKeyVolume, volume, 0).Err()
		if err != nil {
			log.Error().Err(err).Msg("error CalculateOhlc: set Volume")
			return nil, err
		}

		value, err := ss.Redis.Get(ss.Context, code+model.RedisKeyValue).Int64()
		if err != nil {
			value = 0
		}
		value += result[code].Value
		err = ss.Redis.Set(ss.Context, code+model.RedisKeyValue, value, 0).Err()
		if err != nil {
			log.Error().Err(err).Msg("error CalculateOhlc: set Value")
			return nil, err
		}

		averagePrice := int64(math.Round(float64(value) / float64(volume)))
		err = ss.Redis.Set(ss.Context, code+model.RedisKeyAveragePrice, averagePrice, 0).Err()
		if err != nil {
			log.Error().Err(err).Msg("error CalculateOhlc: set AveragePrice")
			return nil, err
		}
	}

	return result, nil
}

func (ss *StockService) ProduceRecords(changeRecords []model.ChangeRecordInput) error {
	var (
		protoChangeRecords []*__.ChangeRecord
	)

	for _, record := range changeRecords {
		var (
			protoChangeRecord __.ChangeRecord
			quantity          int64
			executedQuantity  int64
			price             int64
			executionPrice    int64
			err               error
		)

		if record.Quantity != "" {
			quantity, err = strconv.ParseInt(record.Quantity, 10, 64)
			if err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("error strconv.ParseInt(record.Quantity, 10, 64) , record.Quantity: %s", record.Quantity))
				return err
			}
		} else if record.ExecutedQuantity != "" {
			executedQuantity, err = strconv.ParseInt(record.ExecutedQuantity, 10, 64)
			if err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("error strconv.ParseInt(record.ExecutedQuantity, 10, 64) , record.ExecutedQuantity: %s", record.ExecutedQuantity))
				return err
			}
		}

		if record.Price != "" {
			price, err = strconv.ParseInt(record.Price, 10, 64)
			if err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("error strconv.ParseInt(record.Price, 10, 64) , record.Price: %s", record.Price))
				return err
			}
		} else if record.ExecutionPrice != "" {
			executionPrice, err = strconv.ParseInt(record.ExecutionPrice, 10, 64)
			if err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("error strconv.ParseInt(record.ExecutionPrice, 10, 64) , record.ExecutionPrice: %s", record.ExecutionPrice))
				return err
			}
		}

		protoChangeRecord.Type = record.Type
		protoChangeRecord.OrderNumber = record.OrderNumber
		protoChangeRecord.OrderVerb = record.OrderVerb
		protoChangeRecord.Quantity = quantity
		protoChangeRecord.ExecutedQuantity = executedQuantity
		protoChangeRecord.OrderBook = record.OrderBook
		protoChangeRecord.Price = price
		protoChangeRecord.ExecutionPrice = executionPrice
		protoChangeRecord.StockCode = record.StockCode

		protoChangeRecords = append(protoChangeRecords, &protoChangeRecord)
	}

	message := __.ChangeRecords{
		ChangeRecords: protoChangeRecords,
	}

	payload, err := proto.Marshal(&message)
	if err != nil {
		log.Error().Err(err).Msg("error proto.Marshal message")
		return err
	}

	key := []byte(uuid.New().String())
	err = ss.KafkaProducer.Produce(context.Background(), ss.Config.KafkaConfig.Topic, kafka.Message{
		Key:   key,
		Value: payload,
	})
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error Produce changeRecords: %+v", changeRecords))
		return err
	}
	return nil
}

func (ss *StockService) GetSummary() (map[string]model.Summary, error) {
	stockCodesByte, err := ss.Redis.Get(ss.Context, model.RedisKeyStockCode).Bytes()
	if err != nil {
		log.Error().Err(err).Msg("error redis get " + model.RedisKeyStockCode)
		return nil, err
	}

	var stockCodes []string
	err = json.Unmarshal(stockCodesByte, &stockCodes)
	if err != nil {
		log.Error().Err(err).Msg("error GetSummary: json.Unmarshal(stockCodesByte)")
		return nil, err
	}

	summaries := make(map[string]model.Summary, len(stockCodes))

	for _, code := range stockCodes {
		var summary model.Summary

		summary.StockCode = code

		indexCodeBytes, err := ss.Redis.Get(ss.Context, code+model.RedisKeyIndexCode).Bytes()
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: IndexCode for stockCode: " + code)
			continue // skip for this stock code
		}
		err = json.Unmarshal(indexCodeBytes, &summary.IndexCode)
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: json.Unmarshal(indexCodeBytes)")
			continue // skip for this stock code
		}

		previousPrice, err := ss.Redis.Get(ss.Context, code+model.RedisKeyPreviousPrice).Int64()
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: PreviousPrice for stockCode: " + code)
			continue // skip for this stock code
		}
		summary.PreviousPrice = previousPrice

		openPrice, err := ss.Redis.Get(ss.Context, code+model.RedisKeyOpenPrice).Int64()
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: OpenPrice for stockCode: " + code)
			continue // skip for this stock code
		}
		summary.OpenPrice = openPrice

		highestPrice, err := ss.Redis.Get(ss.Context, code+model.RedisKeyHighestPrice).Int64()
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: HighestPrice for stockCode: " + code)
			continue // skip for this stock code
		}
		summary.HighestPrice = highestPrice

		lowestPrice, err := ss.Redis.Get(ss.Context, code+model.RedisKeyLowestPrice).Int64()
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: LowestPrice for stockCode: " + code)
			continue // skip for this stock code
		}
		summary.LowestPrice = lowestPrice

		closePrice, err := ss.Redis.Get(ss.Context, code+model.RedisKeyClosePrice).Int64()
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: ClosePrice for stockCode: " + code)
			continue // skip for this stock code
		}
		summary.ClosePrice = closePrice

		volume, err := ss.Redis.Get(ss.Context, code+model.RedisKeyVolume).Int64()
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: Volume for stockCode: " + code)
			continue // skip for this stock code
		}
		summary.Volume = volume

		value, err := ss.Redis.Get(ss.Context, code+model.RedisKeyValue).Int64()
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: Value for stockCode: " + code)
			continue // skip for this stock code
		}
		summary.Value = value

		avgPrice, err := ss.Redis.Get(ss.Context, code+model.RedisKeyAveragePrice).Int64()
		if err != nil {
			log.Error().Err(err).Msg("error GetSummary: AveragePrice for stockCode: " + code)
			continue // skip for this stock code
		}
		summary.AveragePrice = avgPrice

		summaries[code] = summary
	}

	return summaries, nil
}
