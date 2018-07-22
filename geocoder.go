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

// Complete geocoding result
type GeocodingResult struct {
	Info    Info `json:"info"`
	Options struct {
		MaxResults        int  `json:"maxResults"`
		ThumbMaps         bool `json:"thumbMaps"`
		IgnoreLatLngInput bool `json:"ignoreLatLngInput"`
	} `json:"options"`
	Results []struct {
		ProvidedLocation struct {
			Location string `json:"location"`
		} `json:"providedLocation"`
		Locations []struct {
			Street string `json:"street"`
			// Neighborhood
			AdminArea6     string `json:"adminArea6"`
			AdminArea6Type string `json:"adminArea6Type"`
			// City
			AdminArea5     string `json:"adminArea5"`
			AdminArea5Type string `json:"adminArea5Type"`
			// County
			AdminArea4     string `json:"adminArea4"`
			AdminArea4Type string `json:"adminArea4Type"`
			// State
			AdminArea3     string `json:"adminArea3"`
			AdminArea3Type string `json:"adminArea3Type"`
			// Country
			AdminArea1         string `json:"adminArea1"`
			AdminArea1Type     string `json:"adminArea1Type"`
			PostalCode         string `json:"postalCode"`
			GeocodeQualityCode string `json:"geocodeQualityCode"`
			// ex: "NEIGHBORHOOD", "CITY", "COUNTY"
			GeocodeQuality string `json:"geocodeQuality"`
			DragPoint      bool   `json:"dragPoint"`
			SideOfStreet   string `json:"sideOfStreet"`
			LinkId         string `json:"linkId"`
			UnknownInput   string `json:"unknownInput"`
			Type           string `json:"type"`
			LatLng         LatLng `json:"latLng"`
			DisplayLatLng  LatLng `json:"displayLatLng"`
			MapUrl         string `json:"mapUrl"`
		} `json:"locations"`
	} `json:"results"`
}
