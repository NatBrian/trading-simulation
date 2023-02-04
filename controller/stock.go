package controller

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/NatBrian/Stockbit-Golang-Challenge/config"
	"github.com/NatBrian/Stockbit-Golang-Challenge/helper"
	"github.com/NatBrian/Stockbit-Golang-Challenge/model"
	"github.com/NatBrian/Stockbit-Golang-Challenge/service"
	"github.com/gabriel-vasile/mimetype"
	"github.com/olivere/ndjson"
)

type (
	StockController struct {
		Config       config.Config
		StockService service.StockService
	}

	IStockController interface {
		UploadTransaction(w http.ResponseWriter, r *http.Request)
	}
)

func (sc *StockController) UploadTransaction(w http.ResponseWriter, r *http.Request) {
	// limit Mb
	err := r.ParseMultipartForm(int64(sc.Config.MaxZipMb) << 20)
	if err != nil {
		errorMessage := fmt.Sprintf("fileInZip above max size limit: %d mb", sc.Config.MaxZipMb)
		log.Println(errorMessage, err)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusBadRequest, resp)
		return
	}

	// read zip fileInZip
	formZipFile, formZipFileHeader, err := r.FormFile("file")
	if err != nil {
		errorMessage := "error FormFile"
		log.Println(errorMessage, err)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	defer func(formZipFile multipart.File) {
		err := formZipFile.Close()
		if err != nil {
			errorMessage := "error formZipFile.ClosePrice()"
			log.Println(errorMessage, err)
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
		log.Println(errorMessage, err)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	if mType.Extension() != ".zip" {
		errorMessage := "wrong fileInZip type, need .zip"
		log.Println(errorMessage, err)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusBadRequest, resp)
		return
	}

	// open fileInZip
	zipFile, err := formZipFileHeader.Open()
	if err != nil {
		errorMessage := "error OpenPrice fileInZip"
		log.Println(errorMessage, err)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	defer func(zipFile multipart.File) {
		err := zipFile.Close()
		if err != nil {
			errorMessage := "error fileInZip.ClosePrice()"
			log.Println(errorMessage, err)
			resp := model.NewErrorResponse(errorMessage, err)
			helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
			return
		}
	}(zipFile)

	// get fileInZip size
	zipFileSize, err := zipFile.Seek(0, 2)
	if err != nil {
		errorMessage := "error fileInZip.Seek"
		log.Println(errorMessage, err)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	zipReader, err := zip.NewReader(zipFile, zipFileSize)
	if err != nil {
		errorMessage := "error zip.NewReader"
		log.Println(errorMessage, err)
		resp := model.NewErrorResponse(errorMessage, err)
		helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
		return
	}

	// Read all the files in zip archive
	var (
		summaries []map[string]model.Summary
	)

	log.Println("Reading Zip: ", zipReader.File)
	for _, fileInZip := range zipReader.File {
		log.Println("Reading fileInZip:", fileInZip.Name)

		// open ndjson File
		ndjsonFile, err := fileInZip.Open()
		if err != nil {
			errorMessage := "error OpenPrice ndjsonFile"
			log.Println(errorMessage, err)
			resp := model.NewErrorResponse(errorMessage, err)
			helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
			return
		}

		defer func(ndjsonFile io.ReadCloser) {
			err := ndjsonFile.Close()
			if err != nil {
				errorMessage := "error ndjsonFile.ClosePrice()"
				log.Println(errorMessage, err)
				resp := model.NewErrorResponse(errorMessage, err)
				helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
				return
			}
		}(ndjsonFile)

		ndjsonReader := ndjson.NewReader(ndjsonFile)

		var (
			changeRecords []model.ChangeRecord
			stockCodes    []string
		)

		// read each json line in ndjson
		for ndjsonReader.Next() {
			var changeRecord model.ChangeRecord
			if err := ndjsonReader.Decode(&changeRecord); err != nil {
				errorMessage := "error Decode changeRecord"
				log.Println(errorMessage, err)
				// continue to next fileInZip which has correct struct
			}

			changeRecords = append(changeRecords, changeRecord)
			stockCodes = append(stockCodes, changeRecord.StockCode)
		}
		if err := ndjsonReader.Err(); err != nil {
			errorMessage := "error ndjsonReader"
			log.Println(errorMessage, err)
			resp := model.NewErrorResponse(errorMessage, err)
			helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
			return
		}

		// calculate records
		summary, err := sc.StockService.CalculateOhlc(stockCodes, changeRecords)
		if err != nil {
			errorMessage := "error CalculateOhlc"
			log.Println(errorMessage, err)
			resp := model.NewErrorResponse(errorMessage, err)
			helper.ResponseFormatter(w, http.StatusInternalServerError, resp)
			return
		}

		summaries = append(summaries, summary)
	}

	helper.ResponseFormatter(w, http.StatusOK, summaries)
}
