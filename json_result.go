// Structs for JSON results.
// MapQuest providers more JSON fields than this but this is all we are interested in.

package geocoder

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

// GeocodingResults contains the locations of a geocoding request
type GeocodingResults struct {
	Results []struct {
		Locations []Location `json:"locations"`
	} `json:"results"`
}
