/* Exposes (partially) the mapquest geocoding api.

Reference: http://open.mapquestapi.com/geocoding/

Example:

lat, lng := Geocode("Seattle WA")
address := ReverseGeocode(47.603561, -122.329437)

*/

package geocoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	geocodeURL        = "http://open.mapquestapi.com/geocoding/v1/address?inFormat=kvp&outFormat=json&location="
	reverseGeocodeURL = "http://open.mapquestapi.com/geocoding/v1/reverse?location="
	batchGeocodeURL   = "http://open.mapquestapi.com/geocoding/v1/batch?key="
)

// Geocode returns the latitude and longitude for a certain address
func Geocode(address string) (lat float64, lng float64) {
	// Query Provider
	resp, err := http.Get(geocodeURL + url.QueryEscape(address) + "&key=" + apiKey)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Decode our JSON results
	var result geocodingResults
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
		fmt.Sprintf("%f,%f&key=%s", lat, lng, apiKey))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Decode our JSON results
	var result geocodingResults
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

// BatchGeocode allows multiple locations to be geocoded at the same time.
// A limit of 100 locations exists for one call. Therefore the json is
// embedded as the body of an http post.
func BatchGeocode(addresses []string) (latLngs []LatLng) {
	var next, start, end int
	n := len(addresses)
	latLngs = make([]LatLng, n)
	batches := n/100 + 1
	next = 0
	for batch := 0; batch < batches; batch++ {
		start = next
		next = (batch + 1) * 100
		if n < next {
			end = n
		} else {
			end = next
		}
		bgb := batchGeocodeBody{
			Locations:  addresses[start:end],
			MaxResults: 1,
			ThumbMaps:  false,
		}
		b, err := json.Marshal(bgb)
		if err != nil {
			panic(err)
		}
		body := bytes.NewBuffer(b)
		resp, err := http.Post(batchGeocodeURL+apiKey, "application/json", body)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		var result geocodingResults
		err = decoder(resp).Decode(&result)
		if err != nil {
			panic(err)
		}
		for i, r := range result.Results {
			if len(r.Locations) == 0 {
				latLngs[start+i] = LatLng{Lat: 0, Lng: 0}
			} else {
				latLngs[start+i] = r.Locations[0].LatLng
			}
		}
	}
	return
}

// geocodingResults contains the locations of a geocoding request
// MapQuest providers more JSON fields than this but this is all we are interested in.
type geocodingResults struct {
	Results []struct {
		Locations []Location `json:"locations"`
	} `json:"results"`
}

// batchGeocodeBody will be marshalled as json to send in body with http post
type batchGeocodeBody struct {
	Locations  []string `json:"locations"`
	MaxResults int
	ThumbMaps  bool `json:"thumbMaps"`
}
