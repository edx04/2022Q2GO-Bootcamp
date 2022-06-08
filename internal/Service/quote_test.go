package service

import (
	"log"
	"testing"

	"github.com/edx04/2022Q2GO-Bootcamp/internal/entity"
	"github.com/edx04/2022Q2GO-Bootcamp/test/testdata"
	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_ "github.com/stretchr/testify/mock"
)

type mockQuoteRepository struct {
	mock.Mock
}

func (mr *mockQuoteRepository) GetQuoteById(id int64) (*entity.Quote, error) {
	log.Printf("REPO MOCK: Get Quote with id %d", id)
	arg := mr.Called(id)
	return arg.Get(0).(*entity.Quote), arg.Error(1)
}

func (mr *mockQuoteRepository) InsertQuote(quote *entity.Quote) error {
	log.Printf("REPO MOCK: InsertQuote ")
	arg := mr.Called(quote)
	return arg.Error(0)
}

func (mr *mockQuoteRepository) QuoteWorkerPool(type_ int64, items int, itemsPerWorker int) (result []*entity.Quote, errors error) {
	log.Printf("REPO MOCK: InsertQuote ")
	arg := mr.Called(type_, items, itemsPerWorker)
	return arg.Get(0).([]*entity.Quote), arg.Error(1)
}

type mockClientRepository struct {
	mock.Mock
}

func (c *mockClientRepository) GetRandomQuote() (*entity.ApiQuote, error) {
	log.Printf("REPO MOCK: InsertQuote ")
	arg := c.Called()
	return arg.Get(0).(*entity.ApiQuote), arg.Error(1)
}
func Test_GetQuoteById(t *testing.T) {
	var testCases = []struct {
		name     string
		id       int64
		response *entity.Quote
		err      error
		// Repository
		repoRes *entity.Quote
		repoErr error
		//Client
		clientRes *entity.ApiQuote
		err2      error
	}{
		{
			"Should return 1 quote by id",
			0,
			&entity.Quote{Id: 0, Author: "Niklaus Wirth", Text: "Software gets slower faster than hardware gets faster."},
			nil,
			&testdata.Quotes[0],
			nil,
			&testdata.ApiQuotes[0],
			nil,
		},
		{
			"Should return 0 quotes by id that doesn´t exist",
			4,
			nil,
			ErrIdNotExist,
			nil,
			ErrIdNotExist,
			&testdata.ApiQuotes[0],
			nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// service initialize
			var svc QuoteService
			repo := &mockQuoteRepository{}
			client := &mockClientRepository{}
			repo.On("GetQuoteById", testCase.id).Return(testCase.repoRes, testCase.repoErr)
			svc = NewQuoteService(repo, client)

			// Run test
			quote, err := svc.FindQuoteById(testCase.id)
			t.Logf("Quote found: %v", quote)

			// Assert
			assert.Equal(t, testCase.response, quote)
			assert.Equal(t, testCase.err, err)
		})
	}
}
