package geocoder

import "testing"

const (
	city       = "Seattle"
	state      = "WA"
	postalCode = "98164"
	seattleLat = 47.603832
	seattleLng = -122.330062
	antwerpLat = 51.221110
	antwerpLng = 4.399708
	beijingLat = 39.905963
	beijingLng = 116.391248
)

func TestGeocode(t *testing.T) {
	query := "Seattle WA"
	lat, lng, err := Geocode(query)
	if err != nil {
		t.Errorf("Seattle: Expected error to be nil ~ Received %v", err)
	}

	if lat != seattleLat || lng != seattleLng {
		t.Errorf("Seattle: Expected (%f, %f) ~ Received (%f, %f)", seattleLat, seattleLng, lat, lng)
	}
}

func TestReverseGeoCode(t *testing.T) {
	address, err := ReverseGeocode(seattleLat, seattleLng)
	if err != nil {
		t.Errorf("Seattle (reverse): Expected error to be nil ~ Received %v", err)
	}

	if address != nil && address.City != city || address.State != state || address.PostalCode != postalCode {
		t.Errorf("Seattle (reverse): Expected %s %s %s ~ Received %s %s %s",
			city, state, postalCode, address.City, address.State, address.PostalCode)
	}
}

func TestBatchGeocode(t *testing.T) {
	latLngs, err := BatchGeocode([]string{"Antwerp,Belgium", "Beijing,China"})
	if err != nil {
		t.Errorf("Seattle (reverse): Expected error to be nil ~ Received %v", err)
	}

	if len(latLngs) != 2 {
		t.Fatalf("Batch: Expected len(batch) to be 2, got %v", latLngs)
	}
	antwerp := latLngs[0]
	if antwerp.Lat != antwerpLat || antwerp.Lng != antwerpLng {
		t.Errorf("Antwerp: Expected %f, %f ~ Received %f, %f", antwerpLat, antwerpLng, antwerp.Lat, antwerp.Lng)
	}
	beijing := latLngs[1]
	if beijing.Lat != beijingLat || beijing.Lng != beijingLng {
		t.Errorf("Beijng: Expected %f, %f ~ Received %f, %f", beijingLat, beijingLng, beijing.Lat, beijing.Lng)
	}
}

func TestGeocodeShouldFail(t *testing.T) {
	query := "Seattle WA"
	// set a bad api key
	SetAPIKey("bad api key that doesn't exist, hopefully!")
	lat, lng, err := Geocode(query)
	SetAPIKey(apiKey)

	if err == nil {
		t.Errorf("Seattle: Expected error to not be nil ~ Received %v", err)
	}

	if lat != 0 || lng != 0 {
		t.Errorf("Seattle: Expected (0, 0) ~ Received (%f, %f)", lat, lng)
	}
}

func TestReverseGeoCodeShouldFail(t *testing.T) {
	// set a bad api key
	SetAPIKey("bad api key that doesn't exist, hopefully!")
	address, err := ReverseGeocode(seattleLat, seattleLng)
	SetAPIKey(apiKey)

	if err == nil {
		t.Errorf("Seattle (reverse): Expected error to not be nil ~ Received %v", err)
	}

	if address != nil {
		t.Errorf("Seattle (reverse): Expected address to be nil ~ Received %v",
			address)
	}
}

func TestBatchGeocodeShouldFail(t *testing.T) {
	// set a bad api key
	SetAPIKey("bad api key that doesn't exist, hopefully!")
	latLngs, err := BatchGeocode([]string{"Antwerp,Belgium", "Beijing,China"})
	SetAPIKey(apiKey)

	if err == nil {
		t.Errorf("Batch: Expected error to be nil ~ Received %v", err)
	}

	if len(latLngs) != 0 {
		t.Errorf("Batch: Expected len(batch) to be 0, got %v", latLngs)
	}
}
