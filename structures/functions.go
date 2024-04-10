package structures

import (
	"arbitrage-bot-v2-sandbox/constants"
	"strconv"
)

func Setup() {
	devGenerateDexes(constants.NumExchanges, constants.PartialPort)
}

func devGenerateDexes(number int, partialPort string) {
	for i := 0; i < number; i++ {
		strNum := strconv.Itoa(i)
		name := "Dex" + strNum
		url := "http://localhost:" + partialPort + strNum
		dex := Dex{
			DexMetadata: DexMetadata{
				Name: name,
			},
			URL:    url,
			Active: true,
		}
		Exchanges = append(Exchanges, dex)
	}
}
