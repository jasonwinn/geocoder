package geocoder

import "testing"

func TestSetAPIKey(t *testing.T) {
	key := apiKey
	SetAPIKey("foo")
	if apiKey != "foo" {
		t.Errorf("Expected foo ~ Received %s", apiKey)
	}
	SetAPIKey(key)
}
