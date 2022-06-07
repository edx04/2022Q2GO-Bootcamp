package repository

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/edx04/2022Q2GO-Bootcamp/internal/entity"
)

type QuoteRepository interface {
	InsertQuote(quote *entity.Quote) error
	GetQuoteById(id int64) (*entity.Quote, error)
	QuoteWorkerPool(type_ int64, items int, itemsPerWorker int) (result []*entity.Quote, errors error)
}

type quoteRepository struct {
	filepath string
}

type GoroutinePool struct {
	queue  chan []string
	result chan *entity.Quote
	wg     sync.WaitGroup
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

func (qr *quoteRepository) QuoteWorkerPool(type_ int64, items int, itemsPerWorker int) (result []*entity.Quote, errors error) {

	numWorkers := items / itemsPerWorker

	pool := &GoroutinePool{
		queue:  make(chan []string),
		result: make(chan *entity.Quote),
		wg:     sync.WaitGroup{},
	}
	pool.wg.Add(numWorkers)

	file, err := os.Open(qr.filepath)
	if err != nil {
		log.Println(err)
		return nil, ErrOpenCsv
	}

	for i := 0; i < numWorkers; i++ {
		go func(workerId int) {

			func(queue <-chan []string, results chan<- *entity.Quote, maxItemsPerWorker int) {
				count := 0
				for queue := range queue {
					id, _ := strconv.ParseInt(queue[0], 0, 0)
					results <- &entity.Quote{Id: id,
						Author: queue[1],
						Text:   queue[2],
					}

					count++
					if count >= maxItemsPerWorker {
						break
					}

				}
			}(pool.queue, pool.result, itemsPerWorker)
			defer pool.wg.Done()
		}(i)
	}

	go func(file *os.File, items int, type_ int64, pool *GoroutinePool) {
		reader := csv.NewReader(file)
		count := 0
		for {
			row, err := reader.Read()

			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
				break
			}

			id, _ := strconv.ParseInt(row[0], 0, 0)

			if type_ == 0 {
				if id%2 == 0 {
					pool.queue <- row
					count++
				}

			} else if type_ == 1 {
				if id%2 != 0 {
					pool.queue <- row
					count++
				}

			} else {
				break
			}

			if count >= items {
				break
			}

		}
		close(pool.queue)
	}(file, items, type_, pool)

	go func() {
		pool.wg.Wait()
		close(pool.result)
	}()

	for quote := range pool.result {
		result = append(result, quote)
	}

	return
}
