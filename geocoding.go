/* Exposes (partially) the mapquest geocoding api.

Reference: http://open.mapquestapi.com/geocoding/

Example:

lat, lng := Geocode("Seattle WA")
address := ReverseGeocode(47.603561, -122.329437)

*/

package geocoder

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	geocodeURL        = "http://www.mapquestapi.com/geocoding/v1/address?inFormat=kvp&outFormat=json&location="
	reverseGeocodeURL = "http://www.mapquestapi.com/geocoding/v1/reverse?location="
)

// Geocode returns the latitude and longitude for a certain address
func Geocode(query string) (lat float64, lng float64) {
	// Query Provider
	resp, err := http.Get(geocodeURL + url.QueryEscape(query) + "&key=" + apiKey)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Decode our JSON results
	result := new(GeocodingResults)
	err = decoder(resp).Decode(&result)

	if err != nil {
		panic(err)
	}

	if len(result.Results[0].Locations) > 0 {
		lat = result.Results[0].Locations[0].LatLng.Lat
		lng = result.Results[0].Locations[0].LatLng.Lng
	}

	return
}

// ReverseGeocode returns the address for a certain latitude and longitude
func ReverseGeocode(lat float64, lng float64) *Location {
	// Query Provider
	resp, err := http.Get(reverseGeocodeURL +
		fmt.Sprintf("%f", lat) + "," + fmt.Sprintf("%f", lng) +
		"&key=" + apiKey)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Decode our JSON results
	result := new(GeocodingResults)
	err = decoder(resp).Decode(&result)

	if err != nil {
		panic(err)
	}

	var location Location

	// Assign the results to the Location struct
	if len(result.Results[0].Locations) > 0 {
		location = result.Results[0].Locations[0]
	}

	return &location
}
