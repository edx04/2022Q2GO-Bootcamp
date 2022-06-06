package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/edx04/2022Q2GO-Bootcamp/internal/entity"
	"github.com/gorilla/mux"
)

type QuoteHandler interface {
	GetQuoteByIdHandlers() http.HandlerFunc
	GenerateQuoteHandlers() http.HandlerFunc
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
		fmt.Println("id : ", params["id"])

		id, _ := strconv.ParseInt(params["id"], 0, 0)
		quote, err := q.service.FindQuoteById(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetQuoteByIdResponse{
			Author: quote.Author,
			Text:   quote.Text,
			Status: true,
		})

	}
}

func (q *quoteHandler) GenerateQuoteHandlers() http.HandlerFunc {
	fmt.Println("heree")
	return func(w http.ResponseWriter, r *http.Request) {
		quote, err := q.service.GenerateQuote()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(GenerateQuoteResponse{
			Id:     quote.Id,
			Author: quote.Author,
			Text:   quote.Text,
			Status: true,
		})

	}
}
