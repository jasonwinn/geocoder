// Golang mapquest api

package geocoder

import (
	"encoding/json"
	"net/http"
)

var apiKey = "Fmjtd%7Cluub256alu%2C7s%3Do5-9u82ur"

// SetAPIKey lets you set your own api key.
// The default api key is probably okay to use for testing.
// But for production, you should create your own key at http://mapquestapi.com
func SetAPIKey(key string) {
	apiKey = key
}

// Shortcut for creating a json decoder out of a response
func decoder(resp *http.Response) *json.Decoder {
	return json.NewDecoder(resp.Body)
}

// LatLng specifies a point with latitude and longitude
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Location is specified by its address and coordinates
type Location struct {
	Street      string `json:"street"`
	City        string `json:"adminArea5"`
	State       string `json:"adminArea3"`
	PostalCode  string `json:"postalCode"`
	County      string `json:"adminArea4"`
	CountryCode string `json:"adminArea1"`
	LatLng      LatLng `json:"latLng"`
	Type        string `json:"type"`
	DragPoint   bool   `json:"dragPoint"`
}
