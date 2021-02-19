package helpers

import (
	"fmt"
	"goinvest/config"
	"goinvest/models"
	"goinvest/quote"
	"math"
)

func BalanceAccount(conf config.Config, bal models.Balance, quotes []quote.Quote) (result map[string]map[string]models.CalculatedAllocation, err error) {
	// divide a balance according to the portfolio
	portfolio, err := conf.GetPortfolio(bal.Portfolio)
	if err != nil {
		return result, fmt.Errorf(
			"failed to balance: %v",
			err.Error(),
		)
	}

	// first, group symbols according to their classification
	groups := make(map[string]map[string]models.CalculatedAllocation)
	for _, symbol := range portfolio.Symbols {
		if groups[symbol.Type] == nil {
			groups[symbol.Type] = make(map[string]models.CalculatedAllocation)
		}
		groups[symbol.Type][symbol.Symbol] = models.CalculatedAllocation{}
	}

	// now that we have all symbols grouped, proceed
	// to apply the allocations
	for group, symbols := range groups {
		groupAllocation := (portfolio.Allocations[group] / float64(100.0)) * bal.Balance
		allocPercentageFromTotal := portfolio.Allocations[group] / float64(len(symbols))
		allocPerSymbol := groupAllocation / float64(len(symbols))

		// the balance for this type of investment has been established,
		// so proceed to skim over each symbol associated with this type
		// of investment and find out how many shares to buy
		for symbol := range symbols {
			for _, quote := range quotes {
				if quote.Symbol == symbol {
					shares := int64(math.Floor(allocPerSymbol / quote.Price))
					totalAllocated := float64(shares) * quote.Price
					groups[group][symbol] = models.CalculatedAllocation{
						Shares:                          shares,
						SharePrice:                      quote.Price,
						Remainder:                       allocPerSymbol - totalAllocated,
						TotalAllocated:                  totalAllocated,
						IdealAllocation:                 allocPerSymbol,
						IdealGroupAllocationPercentage:  portfolio.Allocations[group],
						IdealSymbolAllocationPercentage: allocPercentageFromTotal,
					}
				}
			}
		}
	}

	return groups, nil
}
