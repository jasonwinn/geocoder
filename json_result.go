// Structs for JSON results.
// MapQuest providers more JSON fields than this but this is all we are interested in.

package geocoder

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

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

type GeocodingResults struct {
	Results []struct {
		Locations []Location `json:"locations"`
	} `json:"results"`
}
