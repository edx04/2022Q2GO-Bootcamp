package service

import (
	"strconv"

	"github.com/edx04/2022Q2GO-Bootcamp/internal/entity"
)

type QuoteService interface {
	GenerateQuote() (*entity.Quote, error)
	FindQuoteById(id int64) (*entity.Quote, error)
	GetQuoteWorkerPool(type_ string, items string, itemsPerWorker string) (result []*entity.Quote, errors error)
}

type QuoteRepository interface {
	InsertQuote(quote *entity.Quote) error
	GetQuoteById(id int64) (*entity.Quote, error)
	QuoteWorkerPool(type_ int64, items int, itemsPerWorker int) (result []*entity.Quote, errors error)
}

type ClientRepository interface {
	GetRandomQuote() (*entity.ApiQuote, error)
}

type quoteService struct {
	id     int
	repo   QuoteRepository
	client ClientRepository
}

func NewQuoteService(repo QuoteRepository, client ClientRepository) QuoteService {
	return &quoteService{0, repo, client}
}

func (qs *quoteService) GenerateQuote() (*entity.Quote, error) {
	quote, err := qs.client.GetRandomQuote()

	if err != nil {
		return nil, err
	}

	quoteWithId := &entity.Quote{Id: int64(qs.id),
		Author: quote.Author,
		Text:   quote.Text}

	err = qs.repo.InsertQuote(quoteWithId)

	if err != nil {
		return nil, err
	}

	qs.id++

	return quoteWithId, nil
}

func (qs *quoteService) FindQuoteById(id int64) (*entity.Quote, error) {
	quote, err := qs.repo.GetQuoteById(id)

	if err != nil {
		return nil, err
	}

	return quote, nil
}

func (qs *quoteService) GetQuoteWorkerPool(type_ string, items string, itemsPerWorker string) (result []*entity.Quote, errors error) {
	var type_int int64
	if type_ == "odd" || type_ == "ODD" {
		type_int = 0
	} else {
		type_int = 1
	}

	items_int, err := strconv.ParseInt(items, 0, 0)

	if err != nil {
		return nil, ErrParamItems
	}

	itemsPerWorker_int, err := strconv.ParseInt(items, 0, 0)

	if err != nil {
		return nil, ErrItemsPerWorker
	}

	quotes, err := qs.repo.QuoteWorkerPool(type_int, int(items_int), int(itemsPerWorker_int))
	if err != nil {
		return nil, err
	}

	return quotes, err
}
