// Golang mapquest api

package geocoder

import (
	"encoding/json"
	"net/http"
)

var apiKey = "Fmjtd%7Cluub256alu%2C7s%3Do5-9u82ur"

// SetApiKey lets you set your own api key.
// The default api key is probably okay to use for testing.
// But for production, you should create your own key at http://mapquestapi.com
func SetApiKey(key string) {
	apiKey = key
}

// Shortcut for creating a json decoder out of a response
func decoder(resp *http.Response) *json.Decoder {
	return json.NewDecoder(resp.Body)
}
