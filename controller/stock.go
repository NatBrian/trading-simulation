package controller

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/NatBrian/Stockbit-Golang-Challenge/config"
	"github.com/NatBrian/Stockbit-Golang-Challenge/helper"
	"github.com/NatBrian/Stockbit-Golang-Challenge/kafka"
	"github.com/NatBrian/Stockbit-Golang-Challenge/model"
	__ "github.com/NatBrian/Stockbit-Golang-Challenge/pb"
	"github.com/NatBrian/Stockbit-Golang-Challenge/service"
	"github.com/gabriel-vasile/mimetype"
	"github.com/olivere/ndjson"
	"google.golang.org/protobuf/proto"
)

type (
	StockController struct {
		Config        config.Config
		StockService  service.StockService
		Context       context.Context
		KafkaConsumer kafka.Consumer
	}

	IStockController interface {
		UploadTransaction(w http.ResponseWriter, r *http.Request)
		GetSummary(w http.ResponseWriter, r *http.Request)

		ConsumeRecords(msg kafka.Message) error
	}
)

func (sc *StockController) UploadTransaction(w http.ResponseWriter, r *http.Request) {
	// limit Mb
	err := r.ParseMultipartForm(int64(sc.Config.Constants.MaxZipMb) << 20)
	if err != nil {
		errorMessage := fmt.Sprintf("fileInZip above max size limit: %d mb", sc.Config.Constants.MaxZipMb)
		log.Error().Err(err).Msg(errorMessage)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusBadRequest, resp)
		return
	}

	// read zip fileInZip
	formZipFile, formZipFileHeader, err := r.FormFile("file")
	if err != nil {
		errorMessage := "error FormFile"
		log.Error().Err(err).Msg(errorMessage)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	defer func(formZipFile multipart.File) {
		err = formZipFile.Close()
		if err != nil {
			errorMessage := "error formZipFile.ClosePrice()"
			log.Error().Err(err).Msg(errorMessage)
			resp := model.NewErrorResponse(errorMessage, err)
			helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
			return
		}
	}(formZipFile)

	// check mime type
	bReader := bufio.NewReader(formZipFile)
	mType, err := mimetype.DetectReader(bReader)
	if err != nil {
		errorMessage := "error mimetype.DetectReader"
		log.Error().Err(err).Msg(errorMessage)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	if mType.Extension() != ".zip" {
		errorMessage := "wrong fileInZip type, need .zip"
		log.Error().Err(err).Msg(errorMessage)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusBadRequest, resp)
		return
	}

	// open fileInZip
	zipFile, err := formZipFileHeader.Open()
	if err != nil {
		errorMessage := "error OpenPrice fileInZip"
		log.Error().Err(err).Msg(errorMessage)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	defer func(zipFile multipart.File) {
		err = zipFile.Close()
		if err != nil {
			errorMessage := "error fileInZip.ClosePrice()"
			log.Error().Err(err).Msg(errorMessage)
			resp := model.NewErrorResponse(errorMessage, err)
			helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
			return
		}
	}(zipFile)

	// get fileInZip size
	zipFileSize, err := zipFile.Seek(0, 2)
	if err != nil {
		errorMessage := "error fileInZip.Seek"
		log.Error().Err(err).Msg(errorMessage)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	zipReader, err := zip.NewReader(zipFile, zipFileSize)
	if err != nil {
		errorMessage := "error zip.NewReader"
		log.Error().Err(err).Msg(errorMessage)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	log.Info().Msg("Reading Zip")
	for _, fileInZip := range zipReader.File {
		log.Info().Msg(fmt.Sprintf("Reading fileInZip: %s", fileInZip.Name))

		// open ndjson File
		ndjsonFile, err := fileInZip.Open()
		if err != nil {
			errorMessage := "error OpenPrice ndjsonFile"
			log.Error().Err(err).Msg(errorMessage)
			resp := model.NewErrorResponse(errorMessage, err)
			helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
			return
		}

		defer func(ndjsonFile io.ReadCloser) {
			err = ndjsonFile.Close()
			if err != nil {
				errorMessage := "error ndjsonFile.ClosePrice()"
				log.Error().Err(err).Msg(errorMessage)
				resp := model.NewErrorResponse(errorMessage, err)
				helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
				return
			}
		}(ndjsonFile)

		ndjsonReader := ndjson.NewReader(ndjsonFile)

		var (
			changeRecords []model.ChangeRecordInput
		)

		// read each json line in ndjson
		for ndjsonReader.Next() {
			var changeRecord model.ChangeRecordInput
			if err = ndjsonReader.Decode(&changeRecord); err != nil {
				errorMessage := "error Decode changeRecord"
				log.Error().Err(err).Msg(errorMessage)
				// continue to next fileInZip which has correct struct
			}

			changeRecords = append(changeRecords, changeRecord)
		}
		if err = ndjsonReader.Err(); err != nil {
			errorMessage := "error ndjsonReader"
			log.Error().Err(err).Msg(errorMessage)
			resp := model.NewErrorResponse(errorMessage, err)
			helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
			return
		}

		err = sc.StockService.ProduceRecords(changeRecords)
		if err != nil {
			errorMessage := "error ProduceRecords"
			log.Error().Err(err).Msg(errorMessage)
			resp := model.NewErrorResponse(errorMessage, err)
			helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
			return
		}

	}

	helper.ResponseFormatter(w, http.StatusCreated, nil)
}

func (sc *StockController) GetSummary(w http.ResponseWriter, r *http.Request) {
	summaries, err := sc.StockService.GetSummary()
	if err != nil {
		errorMessage := "error GetSummary"
		log.Error().Err(err).Msg(errorMessage)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	helper.ResponseFormatter(w, http.StatusOK, model.SummaryResponse{
		Summaries: summaries,
	})
}

func (sc *StockController) ConsumeRecords(msg kafka.Message) error {
	var (
		kafkaPayload __.ChangeRecords
		data         []model.ChangeRecord
		stockCodes   []string
	)

	err := proto.Unmarshal(msg.Value, &kafkaPayload)
	if err != nil {
		return err
	}

	for _, payload := range kafkaPayload.ChangeRecords {
		data = append(data, model.ChangeRecord{
			Type:             payload.Type,
			OrderNumber:      payload.OrderNumber,
			OrderVerb:        payload.OrderVerb,
			Quantity:         payload.Quantity,
			ExecutedQuantity: payload.ExecutedQuantity,
			OrderBook:        payload.OrderBook,
			Price:            payload.Price,
			ExecutionPrice:   payload.ExecutionPrice,
			StockCode:        payload.StockCode,
		})
	}

	for _, d := range data {
		stockCodes = append(stockCodes, d.StockCode)
	}

	summary, err := sc.StockService.CalculateOhlc(stockCodes, data)
	if err != nil {
		errorMessage := "error CalculateOhlc"
		log.Error().Err(err).Msg(errorMessage)
		return err
	}

	log.Info().Msg(fmt.Sprintf("ConsumeRecords success. Actual correct summary can be retrieved via API GET summary. temporary summary: %+v", summary))

	return nil
}
