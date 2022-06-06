package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	handler "github.com/edx04/2022Q2GO-Bootcamp/internal/Handler"
	repository "github.com/edx04/2022Q2GO-Bootcamp/internal/Repository"
	service "github.com/edx04/2022Q2GO-Bootcamp/internal/Service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config/.env")

	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	PORT := os.Getenv("PORT")
	ENDPOINT_URL := os.Getenv("Endpoint")

	// s, err := server.NewServer(context.Background(), &server.Config{
	// 	Port:      PORT,
	// 	JWTSecret: JWT_SECRET,
	// })

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("data/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File created successfully")
	defer file.Close()

	repo := repository.NewQuoteRepo("data/data.csv")
	client := repository.NewClientRepo(ENDPOINT_URL)

	service := service.NewQuoteService(repo, client)

	handler := handler.NewHandlerQuote(service)

	handler.GetQuoteByIdHandlers()

	r := mux.NewRouter()

	Routes(r, handler)

	srv := &http.Server{
		Handler: r,
		Addr:    PORT,
		// Good practice: enforce timeouts for servers you create!
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting server Port", PORT)
	log.Fatal(srv.ListenAndServe())
}

func Routes(r *mux.Router, handler handler.QuoteHandler) {
	r.HandleFunc("/Quote/{id}", handler.GetQuoteByIdHandlers()).Methods("GET")
	r.HandleFunc("/RandomQuote", handler.GenerateQuoteHandlers()).Methods("POST")

}
