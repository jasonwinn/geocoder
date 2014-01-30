package geocoder

import (
	"fmt"
	"testing"
)

const (
	city       = "Seattle"
	state      = "WA"
	postalCode = "98104"
	seattleLat = 47.603561
	seattleLng = -122.329437
)

func TestGeocode(t *testing.T) {
	query := "Seattle WA"
	lat, lng := Geocode(query)

	if lat != seattleLat || lng != seattleLng {
		t.Error(fmt.Sprintf("Expected %f, %f ~ Received %f, %f", seattleLat, seattleLng, lat, lng))
	}
}

func TestReverseGeoCode(t *testing.T) {
	address := ReverseGeocode(seattleLat, seattleLng)

	if address.City != city || address.State != state || address.PostalCode != postalCode {
		t.Error(fmt.Sprintf("Expected %s %s %s ~ Received %s %s %s",
			city, state, postalCode, address.City, address.State, address.PostalCode))
	}
}
