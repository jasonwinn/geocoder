/* Exposes (partially) the mapquest geocoding api.

Reference: http://open.mapquestapi.com/geocoding/

Example:

lat, lng := Geocode("Seattle WA")
location, err := ReverseGeocode(47.603561, -122.329437)

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
	geocodeURL        = "https://open.mapquestapi.com/geocoding/v1/address?inFormat=kvp&outFormat=json&location="
	reverseGeocodeURL = "https://open.mapquestapi.com/geocoding/v1/reverse?location="
	batchGeocodeURL   = "https://open.mapquestapi.com/geocoding/v1/batch?key="
)

// Returns the latitude and longitude of the best location match
// for the specified query.
func Geocode(address string) (float64, float64, error) {
	// Query Provider
	resp, err := http.Get(geocodeURL + url.QueryEscape(address) + "&key=" + apiKey)

	if err != nil {
		return 0, 0, fmt.Errorf("Error geocoding address: <%v>", err)
	}

	defer resp.Body.Close()

	// Decode our JSON results
	var result geocodingResults
	err = decoder(resp).Decode(&result)

	if err != nil {
		return 0, 0, fmt.Errorf("Error decoding geocoding result: <%v>", err)
	}

	var lat float64
	var lng float64
	if len(result.Results[0].Locations) > 0 {
		lat = result.Results[0].Locations[0].LatLng.Lat
		lng = result.Results[0].Locations[0].LatLng.Lng
	}

	return lat, lng, nil
}

// Returns the full geocoding response including all of the matches
// as well as reverse-geocoded for each match location.
func FullGeocode(address string) (*GeocodingResult, error) {
	// Query Provider
	resp, err := http.Get(geocodeURL + url.QueryEscape(address) + "&key=" + apiKey)

	if err != nil {
		return nil, fmt.Errorf("Error geocoding address: <%v>", err)
	}

	defer resp.Body.Close()

	// Decode our JSON results
	var result GeocodingResult
	err = decoder(resp).Decode(&result)

	if err != nil {
		return nil, fmt.Errorf("Error decoding geocoding result: <%v>", err)
	}

	return &result, nil
}

// Returns the address for a latitude and longitude.
func ReverseGeocode(lat float64, lng float64) (*Location, error) {
	// Query Provider
	resp, err := http.Get(reverseGeocodeURL +
		fmt.Sprintf("%f,%f&key=%s", lat, lng, apiKey))

	if err != nil {
		return nil, fmt.Errorf("Error reverse geocoding lat, long pair: <%v>", err)
	}

	defer resp.Body.Close()

	// Decode our JSON results
	var result geocodingResults
	err = decoder(resp).Decode(&result)

	if err != nil {
		return nil, fmt.Errorf("Error decoding reverse geocoding result: <%v>", err)
	}

	var location Location

	// Assign the results to the Location struct
	if len(result.Results[0].Locations) > 0 {
		location = result.Results[0].Locations[0]
	}

	return &location, nil
}

// Geocodes multiple locations with a single API request.
// Up to 100 locations per call may be provided.
func BatchGeocode(addresses []string) ([]LatLng, error) {
	var next, start, end int
	n := len(addresses)
	latLngs := make([]LatLng, n)
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
			return nil, err
		}
		body := bytes.NewBuffer(b)
		resp, err := http.Post(batchGeocodeURL+apiKey, "application/json", body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var result geocodingResults
		err = decoder(resp).Decode(&result)
		if err != nil {
			return nil, err
		}
		for i, r := range result.Results {
			if len(r.Locations) == 0 {
				latLngs[start+i] = LatLng{Lat: 0, Lng: 0}
			} else {
				latLngs[start+i] = r.Locations[0].LatLng
			}
		}
	}
	return latLngs, nil
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
