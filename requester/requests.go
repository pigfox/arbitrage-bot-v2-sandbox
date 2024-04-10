package requester

import (
	"arbitrage-bot-v2-sandbox/constants"
	"arbitrage-bot-v2-sandbox/structures"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"
	"time"
	//nolint:goimports
)

type Params[T structures.CustomData] struct {
	Dex                     structures.Dex
	Method                  string
	Headers                 map[string]string
	APIKey                  string
	Body                    string
	ExpectedAPIResponseType reflect.Type //nolint:gofmt    // Holds the type, not an instance
}

type Resp struct {
	Code int         `json:"code,omitempty"` //nolint:gofmt
	Body interface{} `json:"body,omitempty"`
}

func Make[T structures.CustomData](params *Params[T], wg *sync.WaitGroup, results chan<- structures.DexResponse) {
	defer wg.Done()
	var req *http.Request
	var err error
	if params.Body == "" {
		req, err = http.NewRequest(params.Method, params.Dex.URL, http.NoBody)
	} else {
		req, err = http.NewRequest(params.Method, params.Dex.URL, bytes.NewBufferString(params.Body))
	}
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}

	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.TIMEOUT*time.Millisecond)
	defer cancel()

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Println("request failed: %w", err)
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("failed to read response body: %w", err)
		return
	}

	// Handle non-OK status
	if res.StatusCode != http.StatusOK {
		fmt.Println("HTTP request failed with status code %d, body: %s", res.StatusCode, string(resBody)) //nolint:govet
		return
	}

	// Initialize the expected response type
	responseInstance := reflect.New(params.ExpectedAPIResponseType).Interface()
	if err = json.Unmarshal(resBody, &responseInstance); err != nil {
		fmt.Println("error unmarshalling response: %w", err)
		return
	}

	transformedBody, err := transformResponse(responseInstance)
	if err != nil {
		fmt.Println("error transforming response: %w", err)
		return
	}

	transformedBody.DexMetadata = params.Dex.DexMetadata
	results <- transformedBody
}

func transformResponse(resp interface{}) (structures.DexResponse, error) { //nolint:gofmt
	switch v := resp.(type) {
	case *[]structures.SiteToken:
		return *siteTokenArray(v), nil
	default:
		return structures.DexResponse{}, fmt.Errorf("unsupported response type: %T", resp)
	}
}
