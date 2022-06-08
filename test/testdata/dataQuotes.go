package testdata

import "github.com/edx04/2022Q2GO-Bootcamp/internal/entity"

var Quotes = []entity.Quote{
	{Id: int64(0), Author: "Niklaus Wirth", Text: "Software gets slower faster than hardware gets faster."},
	{Id: int64(1), Author: "Anonymous", Text: "A few months writing code can save you a few hours in design."},
	{Id: int64(2), Author: "Alan Perlis", Text: "A year spent in artificial intelligence is enough to make one believe in God."},
}

var ApiQuotes = []entity.ApiQuote{
	{Author: "Niklaus Wirth", Text: "Software gets slower faster than hardware gets faster."},
	{Author: "Anonymous", Text: "A few months writing code can save you a few hours in design."},
	{Author: "Alan Perlis", Text: "A year spent in artificial intelligence is enough to make one believe in God."},
}
