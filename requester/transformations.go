package requester

import (
	"arbitrage-bot-v2-sandbox/structures"
)

func siteTokenArray(source *[]structures.SiteToken) *structures.DexResponse {
	if source == nil {
		return nil
	}

	// Assuming you want to perform some transformation on DexResponse
	// For example, doubling the value of each token
	var transformedTokens []structures.SiteToken
	for _, siteToken := range *source {
		transformedToken := structures.SiteToken{
			Symbol: siteToken.Symbol,
			Value:  2 * siteToken.Value, // Example transformation
		}
		transformedTokens = append(transformedTokens, transformedToken)
	}

	return &structures.DexResponse{
		Data: transformedTokens,
	}
}
