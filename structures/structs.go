package structures

var Exchanges []Dex

type SiteTokenData struct {
	DexName string // Name of the Dex from which this token data was obtained
	Value   int    // The numeric value of the token
}

type Dex struct {
	DexMetadata
	URL    string
	Active bool
}

type DexMetadata struct {
	Name string
}

type DexResponse struct {
	DexMetadata
	Code int         `json:"code"`
	Data []SiteToken `json:"data"`
}

type SiteToken struct {
	Symbol string
	Value  int
}

type CustomData interface {
	DexResponse | []SiteToken
}
