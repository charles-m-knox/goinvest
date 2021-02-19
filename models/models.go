package models

type Symbol struct {
	Symbol string `yaml:"symbol"`
	Type   string `yaml:"type"`
}

type Portfolio struct {
	Name        string             `yaml:"name"`
	Symbols     []Symbol           `yaml:"symbols"`
	Allocations map[string]float64 `yaml:"allocations"`
}

type Balance struct {
	Name      string  `yaml:"name"`
	Balance   float64 `yaml:"balance"`
	Portfolio string  `yaml:"portfolio"`
}

type CalculatedAllocation struct {
	Shares                          int64
	SharePrice                      float64
	Remainder                       float64
	TotalAllocated                  float64
	IdealAllocation                 float64
	IdealGroupAllocationPercentage  float64
	IdealSymbolAllocationPercentage float64
}
