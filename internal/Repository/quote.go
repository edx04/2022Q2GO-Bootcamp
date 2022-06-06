package repository

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/edx04/2022Q2GO-Bootcamp/internal/entity"
)

type QuoteRepository interface {
	InsertQuote(quote *entity.Quote) error
	GetQuoteById(id int64) (*entity.Quote, error)
}

type quoteRepository struct {
	filepath string
}

func NewQuoteRepo(file string) QuoteRepository {
	return &quoteRepository{file}
}

func (qr *quoteRepository) InsertQuote(quote *entity.Quote) error {
	//"data/data.csv"
	csvFile := qr.filepath
	file, err := os.OpenFile(csvFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		return ErrOpenCsv
	}

	defer file.Close()

	csvWriter := csv.NewWriter(file)

	err = csvWriter.Write([]string{strconv.FormatInt(quote.Id, 10), quote.Author, quote.Text})

	if err != nil {
		fmt.Println(err)
		return ErrWriteCsv
	}

	csvWriter.Flush()

	return nil
}

func (qr *quoteRepository) GetQuoteById(_id int64) (*entity.Quote, error) {
	//Open CSV
	file, err := os.Open(qr.filepath)
	if err != nil {
		log.Println(err)
		return nil, ErrOpenCsv
	}

	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()

	if err != nil {
		log.Println(err)
		return nil, ErrReadCsv
	}

	for _, row := range rows {
		id, _ := strconv.ParseInt(row[0], 0, 0)
		if id == _id {
			return &entity.Quote{
				Id:     id,
				Author: row[1],
				Text:   row[2],
			}, nil
		}
	}

	return nil, ErrIdNotExist
}
