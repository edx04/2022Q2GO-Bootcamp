package repository

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/edx04/2022Q2GO-Bootcamp/internal/entity"
)

type ClientRepository interface {
	GetRandomQuote() (*entity.ApiQuote, error)
}

type clientRepository struct {
	endpoint string
}

func NewClientRepo(endpoint string) ClientRepository {
	return &clientRepository{endpoint}
}

func (c *clientRepository) GetRandomQuote() (*entity.ApiQuote, error) {
	res, err := http.DefaultClient.Get(c.endpoint)

	if err != nil {
		return nil, err
	}

	quote := entity.ApiQuote{}
	res2, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(res2, &quote)
	return &quote, nil

}
