package main

import (
	"arbitrage-bot-v2/constants"
	"arbitrage-bot-v2/structures"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp" //nolint:goimports
	"log"
	"strconv"
)

func main() {
	startServers(constants.NumExchanges, constants.PartialPort)
	select {}
}

func startServers(number int, partialPort string) {
	for i := 0; i < number; i++ {
		go func(i int) { // Launch each server in its own goroutine
			strNum := strconv.Itoa(i)
			handler := func(ctx *fasthttp.RequestCtx) {
				switch string(ctx.Path()) {
				case "/":
					tokens := structures.DevGenerateToken(constants.MAXRAND)
					jsonData, err := json.Marshal(tokens)
					if err != nil {
						ctx.Error("Unsupported path", fasthttp.StatusInternalServerError)
						return
					}
					ctx.SetContentType("Content-Type; application/json")
					// Respond with a simple message
					fmt.Fprintf(ctx, string(jsonData))
				default:
					// Handle 404: Not Found
					ctx.Error("Unsupported path", fasthttp.StatusNotFound)
				}
			}

			// Start the server
			port := ":" + partialPort + strNum
			fmt.Println("Server listening on", port)
			if err := fasthttp.ListenAndServe(port, handler); err != nil {
				log.Fatalf("Error in ListenAndServe on port %s: %s", port, err)
			}
		}(i) // Pass the current value of i into the goroutine
	}
}
