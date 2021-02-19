package main

import (
	"encoding/csv"
	"fmt"
	"goinvest/config"
	"goinvest/helpers"
	"goinvest/quote"
	"os"
	"strings"

	"log"

	"github.com/go-resty/resty/v2"
)

func main() {
	// load config
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err.Error())
	}

	f, err := os.Create(conf.OutputFilename)
	if err != nil {
		log.Fatalf(
			"failed to open %v for writing: %v",
			conf.OutputFilename,
			err.Error(),
		)
	}

	symbols := conf.GetAllSymbols()
	symbolsCSV := strings.Join(symbols, ",")
	log.Printf("symbols: %v", symbolsCSV)

	// url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com&symbols=%s", symbolsString)
	client := resty.New()
	quotes := quote.GetQuotes(*client, symbols)
	realQuotes := quotes()
	for _, quote := range realQuotes {
		log.Printf("quote: %v %v", quote.Symbol, quote.Price)
	}

	// write CSV headers
	headers := []string{
		"Name",
		"Symbol",
		"Type",
		"Shares",
		"Share Price",
		"Purchase Price",
		"Allocated",
		"Remainder",
		"Symbol Allocation %",
		"Group Allocation %",
		"From Balance",
	}

	w := csv.NewWriter(f)

	err = w.Write(headers)
	if err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, balance := range conf.Balances {
		// for each balance in conf.Balances, proceed to apply the portfolio
		groups, err := helpers.BalanceAccount(conf, balance, realQuotes)
		if err != nil {
			log.Fatalf("failed to balance: %v", err.Error())
		}

		for group, symbols := range groups {
			for symbol, s := range symbols {
				err := w.Write([]string{
					balance.Name,                            // "Name"
					symbol,                                  // "Symbol"
					group,                                   // "Type"
					fmt.Sprintf("%v", s.Shares),             // "Shares"
					fmt.Sprintf("$%v", s.SharePrice),        // "Shares"
					fmt.Sprintf("$%.2f", s.TotalAllocated),  // "Purchase Price"
					fmt.Sprintf("$%.2f", s.IdealAllocation), // "Allocated"
					fmt.Sprintf("$%.2f", s.Remainder),       // "Remainder"
					fmt.Sprintf("%.2f%%", s.IdealSymbolAllocationPercentage), // "Allocation %"
					fmt.Sprintf("%.2f%%", s.IdealGroupAllocationPercentage),  // "Allocation %"
					fmt.Sprintf("$%.2f", balance.Balance),                    // "From Balance"
				})
				if err != nil {
					log.Fatalf("failed to write newRecord: %v", err.Error())
				}
			}
		}
	}

	w.Flush()

	err = w.Error()
	if err != nil {
		log.Fatalf("failed to write csv: %v", err.Error())
	}

	log.Printf("done: finished writing  to %v", conf.OutputFilename)
}
