package geocoder

import (
	"testing"
)

const (
	city       = "Seattle"
	state      = "WA"
	postalCode = "98104"
	seattleLat = 47.603561
	seattleLng = -122.329437
	antwerpLat = 51.221202
	antwerpLng = 4.399632
	beijingLat = 40.250538
	beijingLng = 116.625089
)

func TestGeocode(t *testing.T) {
	query := "Seattle WA"
	lat, lng := Geocode(query)

	if lat != seattleLat || lng != seattleLng {
		t.Errorf("Expected %f, %f ~ Received %f, %f", seattleLat, seattleLng, lat, lng)
	}
}

func TestReverseGeoCode(t *testing.T) {
	address := ReverseGeocode(seattleLat, seattleLng)

	if address.City != city || address.State != state || address.PostalCode != postalCode {
		t.Errorf("Expected %s %s %s ~ Received %s %s %s",
			city, state, postalCode, address.City, address.State, address.PostalCode)
	}
}

func TestBatchGeocode(t *testing.T) {
	latLngs := BatchGeocode([]string{"Antwerp,Belgium", "Beijing,China"})
	antwerp := latLngs[0]
	if antwerp.Lat != antwerpLat || antwerp.Lng != antwerpLng {
		t.Errorf("Antwerp: Expected %f, %f ~ Received %f, %f", antwerpLat, antwerpLng, antwerp.Lat, antwerp.Lng)
	}
	beijing := latLngs[1]
	if beijing.Lat != beijingLat || beijing.Lng != beijingLng {
		t.Errorf("Beijng: Expected %f, %f ~ Received %f, %f", beijingLat, beijingLng, beijing.Lat, beijing.Lng)
	}
}
