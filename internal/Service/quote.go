package service

import (
	"github.com/edx04/2022Q2GO-Bootcamp/internal/entity"
)

type QuoteService interface {
	GenerateQuote() (*entity.Quote, error)
	FindQuoteById(id int64) (*entity.Quote, error)
}

type QuoteRepository interface {
	InsertQuote(quote *entity.Quote) error
	GetQuoteById(id int64) (*entity.Quote, error)
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
