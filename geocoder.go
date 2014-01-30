package geocoder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// The ApiKey is probably okay to use for testing. But for production,
	// you should create your own key at http://mapquestapi.com

	apiKey     = "Fmjtd%7Cluub256alu%2C7s%3Do5-9u82ur"
	format     = "json"
	geocodeUrl = "http://www.mapquestapi.com/geocoding/v1/address" +
		"?key=" + apiKey +
		"&inFormat=kvp" +
		"&outFormat=" + format +
		"&location="
	reverseGeocodeUrl = "http://www.mapquestapi.com/geocoding/v1/reverse" +
		"?key=" + apiKey +
		"&location="
)

func Geocode(query string) (lat float64, lng float64) {
	// Query Provider
	resp, err := http.Get(geocodeUrl + url.QueryEscape(query))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Decode our JSON results
	result := new(ProviderResult)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)

	if err != nil {
		panic(err)
	}

	if len(result.Results[0].Locations) > 0 {
		lat = result.Results[0].Locations[0].LatLng.Lat
		lng = result.Results[0].Locations[0].LatLng.Lng
	}

	return
}

func ReverseGeocode(lat float64, lng float64) *GeoAddress {
	// Query Provider
	resp, err := http.Get(reverseGeocodeUrl + fmt.Sprintf("%f", lat) + "," + fmt.Sprintf("%f", lng))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Decode our JSON results
	result := new(ReverseProviderResult)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)

	if err != nil {
		panic(err)
	}

	var address GeoAddress

	// Assign the results to the GeoAddress struct
	if len(result.Results[0].Locations) > 0 {
		address = result.Results[0].Locations[0]
	}

	return &address
}
