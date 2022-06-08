package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/edx04/2022Q2GO-Bootcamp/internal/entity"
	"github.com/gorilla/mux"
)

type QuoteHandler interface {
	GetQuoteByIdHandlers() http.HandlerFunc
	GenerateQuoteHandlers() http.HandlerFunc
	ConcurrentlyQuotesHandlers() http.HandlerFunc
	responseWriter(w http.ResponseWriter, httpStatus int, response ...interface{})
}

type GetQuoteByIdRequest struct {
	Id string `json:"id"`
}

type GetQuoteByIdResponse struct {
	Author string `json:"message"`
	Text   string `json:"quote"`
	Status bool   `json:"status"`
}

type GenerateQuoteResponse struct {
	Id     int64  `json:"id"`
	Author string `json:"message"`
	Text   string `json:"quote"`
	Status bool   `json:"status"`
}

type QuoteService interface {
	GenerateQuote() (*entity.Quote, error)
	FindQuoteById(id int64) (*entity.Quote, error)
	GetQuoteWorkerPool(type_ string, items string, itemsPerWorker string) (result []*entity.Quote, errors error)
}

type quoteHandler struct {
	service QuoteService
}

func NewHandlerQuote(service QuoteService) QuoteHandler {
	return &quoteHandler{
		service: service,
	}
}

func (q *quoteHandler) GetQuoteByIdHandlers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Println("id : ", params["id"])

		id, err := strconv.ParseInt(params["id"], 0, 0)

		if err != nil {
			log.Println(err)
			http.Error(w, ErrParamId.Error(), http.StatusBadRequest)
			return
		}

		quote, err := q.service.FindQuoteById(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		q.responseWriter(w, http.StatusOK, GetQuoteByIdResponse{
			Author: quote.Author,
			Text:   quote.Text,
			Status: true,
		})

	}
}

func (q *quoteHandler) GenerateQuoteHandlers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quote, err := q.service.GenerateQuote()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		q.responseWriter(w, http.StatusCreated, GenerateQuoteResponse{
			Id:     quote.Id,
			Author: quote.Author,
			Text:   quote.Text,
			Status: true,
		})

	}
}

func (q *quoteHandler) ConcurrentlyQuotesHandlers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		params := []string{"type", "items", "items_per_workers"}
		value := []string{}

		//params validation
		for _, param := range params {
			val, ok := query[param]
			if !ok {
				http.Error(w, fmt.Sprintf("missing %s parameter", param), http.StatusBadRequest)
				return
			}
			value = append(value, val[0])
		}

		type_ := value[0]
		items := value[1]
		items_per_worker := value[2]

		type_ = strings.ToUpper(type_)

		if !(type_ == "ODD" || type_ == "EVEN") {
			http.Error(w, "The parameter type can only be odd or even", http.StatusBadRequest)
			return
		}

		quotes, err := q.service.GetQuoteWorkerPool(type_, items, items_per_worker)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		quotes_res := []GenerateQuoteResponse{}

		for _, quote := range quotes {

			quotes_res = append(quotes_res, GenerateQuoteResponse{
				Id:     quote.Id,
				Author: quote.Author,
				Text:   quote.Text,
				Status: true,
			})

		}

		q.responseWriter(w, http.StatusOK, quotes_res)

	}

}

func (q *quoteHandler) responseWriter(w http.ResponseWriter, httpStatus int, response ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(response)

}
