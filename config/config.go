package config

import (
	"goinvest/models"
	"sort"

	"fmt"

	viper "github.com/spf13/viper"
)

const (
	ConfigFile = "config.yml"
)

type Config struct {
	OutputFilename string             `yaml:"outputFilename"`
	Portfolios     []models.Portfolio `yaml:"portfolios"`
	Balances       []models.Balance   `yaml:"balances"`
}

// LoadConfig reads from a provided yaml-formatted configuration filename
func LoadConfig() (conf Config, err error) {
	err = viper.BindEnv("config")
	if err != nil {
		return conf, fmt.Errorf("failed to bind config env: %v", err.Error())
	}
	configName := viper.GetString("config")
	viper.SetConfigName(configName) // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./res")    // optionally look for config in the working directory
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	err = viper.ReadInConfig()      // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		return conf, fmt.Errorf("error config file: %s", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		return conf, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return conf, nil
}

// GetAllSymbols uses a map to retrieve all unique ticker symbols
// across all portfolios from the config, and then returns them
func (conf *Config) GetAllSymbols() (symbols []string) {
	uniqueSymbols := make(map[string]string)

	for _, portfolio := range conf.Portfolios {
		for _, symbol := range portfolio.Symbols {
			uniqueSymbols[symbol.Symbol] = symbol.Symbol
		}
	}

	for symbol := range uniqueSymbols {
		symbols = append(symbols, symbol)
	}

	sort.Strings(symbols)

	return
}

// GetPortfolio attempts to retrieve a portfolio by name. If it cannot find one by the provided name, it will return an error.
func (conf *Config) GetPortfolio(name string) (result models.Portfolio, err error) {
	for _, portfolio := range conf.Portfolios {
		if portfolio.Name == name {
			return portfolio, nil
		}
	}

	return result, fmt.Errorf(
		"failed to find a portfolio by name %v",
		name,
	)
}
