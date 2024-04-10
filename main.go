package main

import (
	"arbitrage-bot-v2-sandbox/requester"
	"arbitrage-bot-v2-sandbox/structures"
	"fmt"
	"reflect"
	"sync"
)

// main is the entry point of the program.
func main() { //nolint:typecheck
	structures.Setup()
	fmt.Println(structures.Exchanges)

	var wg sync.WaitGroup
	results := make(chan structures.DexResponse, len(structures.Exchanges))

	for _, dex := range structures.Exchanges {
		wg.Add(1)
		params := &requester.Params[[]structures.SiteToken]{ // Specify the type here
			Dex:                     dex,
			Method:                  "GET",
			ExpectedAPIResponseType: reflect.TypeOf([]structures.SiteToken{}), // Using reflect correctly
		}
		go requester.Make(params, &wg, results) //nolint:errcheck
	}

	wg.Wait()
	close(results)

	allResults := make([]structures.DexResponse, 0, len(structures.Exchanges))
	for result := range results {
		allResults = append(allResults, result)
	}

	compareDexes(allResults)
}

func compareDexes(allResults []structures.DexResponse) {
	fmt.Println("allResults:", allResults, "\n")
	tokenMap := make(map[string][]structures.SiteTokenData)

	// Aggregate all tokens by their symbols
	for _, dex := range allResults {
		for _, token := range dex.Data {
			tokenMap[token.Symbol] = append(tokenMap[token.Symbol], structures.SiteTokenData{
				DexName: dex.Name,
				Value:   token.Value,
			})
		}
	}

	fmt.Println(tokenMap)
	// Now compare all tokens within each symbol group
	for symbol, tokens := range tokenMap {
		fmt.Printf("Comparisons for symbol: %s\n", symbol)
		for i := 0; i < len(tokens); i++ {
			for j := 0; j < len(tokens); j++ {
				if tokens[i].DexName != tokens[j].DexName {
					result := compare(tokens[i].Value, tokens[j].Value)
					fmt.Printf("%s: %s (%d) vs %s (%d) = %d\n",
						symbol, tokens[i].DexName, tokens[i].Value, tokens[j].DexName, tokens[j].Value, result)

					if result == -1 {
						trade(tokens[i], tokens[j])
					} else if result == 1 {
						trade(tokens[j], tokens[i])
					}
				}
			}
		}
	}
}

// compare returns -1 if a > b, 0 if a == b, 1 if a < b.
func compare(a, b int) int {
	if a > b {
		return -1
	} else if a < b {
		return 1
	}
	return 0
}

func trade(to structures.SiteTokenData, from structures.SiteTokenData) {
	fmt.Println("To", to, "From", from)
}
