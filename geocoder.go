package geocoder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type GeoAddress struct {
	Street      string
	City        string
	State       string
	PostalCode  string
	County      string
	CountryCode string
}

func Geocode(query string) (lat float64, lng float64) {
	// Query Provider
	resp, err := http.Get(geocodeUrl + url.QueryEscape(query))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Decode our JSON results
	result := new(ProviderResult)
	json.Unmarshal(body, &result)

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
	body, _ := ioutil.ReadAll(resp.Body)

    // Decode our JSON results
	result := new(ReverseProviderResult)
	json.Unmarshal(body, &result)

	address := &GeoAddress{}

    // Assign the results to the GeoAddress struct
	if len(result.Results[0].Locations) > 0 {
		l := result.Results[0].Locations[0]
		address.Street = l.Street
		address.City = l.City
		address.State = l.State
		address.PostalCode = l.PostalCode
		address.County = l.County
		address.CountryCode = l.CountryCode
	}


	return address
}
