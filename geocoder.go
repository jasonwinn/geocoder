// Golang mapquest api

package geocoder

import (
	"encoding/json"
	"net/http"
)

const (
	// The ApiKey is probably okay to use for testing. But for production,
	// you should create your own key at http://mapquestapi.com

	apiKey = "Fmjtd%7Cluub256alu%2C7s%3Do5-9u82ur"
)

// Shortcut for creating a json decoder out of a response
func decoder(resp *http.Response) *json.Decoder {
	return json.NewDecoder(resp.Body)
}
