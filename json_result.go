package geocoder

// Structs for JSON results.
// MapQuest providers more JSON fields than this but this is all we are interested in.

type ProviderResult struct {
	Results []struct {
		Locations []struct {
			LatLng struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"latLng"`
		} `json:"locations"`
	} `json:"results"`
}
type ReverseProviderResult struct {
	Results []struct {
		Locations []struct {
			Street      string `json:"street"`
			City        string `json:"adminArea5"`
			State       string `json:"adminArea3"`
			PostalCode  string `json:"postalCode"`
			County      string `json:"adminArea4"`
			CountryCode string `json:"adminArea1"`
		} `json:"locations"`
	} `json:"results"`
}
