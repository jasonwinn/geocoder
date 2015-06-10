package geocoder

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

const (
	templDump      = "{\"route\":{\"hasTollRoad\":false,\"computedWaypoints\":[],\"fuelUsed\":0,\"shape\":{\"maneuverIndexes\":[],\"shapePoints\":[51.532039,-0.177344],\"legIndexes\":[]},\"hasUnpaved\":false,\"hasHighway\":false,\"realTime\":-1,\"distance\":0,\"time\":0,\"locationSequence\":[0,1],\"hasSeasonalClosure\":false,\"sessionId\":\"%s\",\"locations\":[{\"latLng\":{\"lng\":-0.177413,\"lat\":51.531965},\"adminArea4\":\"\",\"adminArea5Type\":\"City\",\"adminArea4Type\":\"County\",\"adminArea5\":\"London\",\"street\":\"3 Abbey Road\",\"adminArea1\":\"GB\",\"adminArea3\":\"England\",\"type\":\"s\",\"displayLatLng\":{\"lng\":-0.177413,\"lat\":51.531963},\"linkId\":82135400,\"postalCode\":\"NW8 9AY\",\"sideOfStreet\":\"N\",\"dragPoint\":false,\"adminArea1Type\":\"Country\",\"geocodeQuality\":\"POINT\",\"geocodeQualityCode\":\"P1AAX\",\"adminArea3Type\":\"State\"},{\"latLng\":{\"lng\":-0.177413,\"lat\":51.531965},\"adminArea4\":\"\",\"adminArea5Type\":\"City\",\"adminArea4Type\":\"County\",\"adminArea5\":\"London\",\"street\":\"3 Abbey Road\",\"adminArea1\":\"GB\",\"adminArea3\":\"England\",\"type\":\"s\",\"displayLatLng\":{\"lng\":-0.177413,\"lat\":51.531963},\"linkId\":82135400,\"postalCode\":\"NW8 9AY\",\"sideOfStreet\":\"N\",\"dragPoint\":false,\"adminArea1Type\":\"Country\",\"geocodeQuality\":\"POINT\",\"geocodeQualityCode\":\"P1AAX\",\"adminArea3Type\":\"State\"}],\"hasCountryCross\":false,\"legs\":[{\"hasTollRoad\":false,\"index\":0,\"roadGradeStrategy\":[],\"hasHighway\":false,\"hasUnpaved\":false,\"distance\":0,\"time\":0,\"origIndex\":-1,\"hasSeasonalClosure\":false,\"origNarrative\":\"\",\"hasCountryCross\":false,\"formattedTime\":\"00:00:00\",\"destNarrative\":\"\",\"destIndex\":-1,\"maneuvers\":[{\"signs\":[],\"index\":0,\"maneuverNotes\":[],\"direction\":5,\"narrative\":\"\",\"iconUrl\":\"http://content.mapquest.com/mqsite/turnsigns/icon-dirs-start_sm.gif\",\"distance\":0,\"time\":0,\"linkIds\":[82135400],\"streets\":[\"B507\",\"Abbey Road\"],\"attributes\":0,\"transportMode\":\"AUTO\",\"formattedTime\":\"00:00:00\",\"directionName\":\"Southeast\",\""
	testDistance   = 157
	testLegs       = 1
	testManeuvers  = 51
	testStatuscode = 0
	testTime       = 6085
	testUnit       = "m"
	testURL        = "http://open.mapquestapi.com/directions/v2/route?inFormat=kvp&key=Fmjtd%7Cluub256alu%2C7s%3Do5-9u82ur&outFormat=json&from=Amsterdam%2CNetherlands&to=Antwerp%2CBelgium&unit=m&routeTypefastest&narrativeType=text&enhancedNarrative=false&maxLinkId=0&locale=en_US&mustAvoidLinkIds=5,7&stateBoundaryDisplay=true&countryBoundaryDisplay=true&destinationManeuverDisplay=true&fullShape=false&cyclingRoadFactor=1&roadGradeStrategy=DEFAULT_STRATEGY&drivingStyle=normal&highwayEfficiency=22&manMaps=true&walkingSpeed=-1"
)

func unexpected(err error, t *testing.T) bool {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return true
	}
	return false
}

func TestUrl(t *testing.T) {
	directions := NewDirections("Amsterdam,Netherlands", []string{"Antwerp,Belgium"})
	directions.MustAvoidLinkIDs = []int{5, 7}
	routeURL := directions.URL("json")
	if testURL != routeURL {
		t.Errorf("Expected %s ~ Received %s", testURL, routeURL)
	}
}

func TestDistance(t *testing.T) {
	directions := NewDirections("Amsterdam,Netherlands", []string{"Antwerp,Belgium"})
	distance, err := directions.Distance("k")
	if unexpected(err, t) {
		return
	}
	if testDistance != int(distance) {
		t.Errorf("Expected %d ~ Received %d", testDistance, int(distance))
	}
	// distance function may not alter the original unit
	if directions.Unit != "m" {
		t.Errorf("Unit: Expected %s ~ Received %s", testUnit, directions.Unit)
	}
}

func TestDistanceError(t *testing.T) {
	directions := NewDirections("Amsterdam,Netherlands", []string{"sdfa"})
	_, err := directions.Distance("k")
	if err.Error() != "Error 400: We are unable to route with the given locations." {
		t.Errorf("Expected %s ~ Received %s", "", err.Error())
	}
}

func TestDirections(t *testing.T) {
	directions := NewDirections("Amsterdam,Netherlands", []string{"Antwerp,Belgium"})
	directions.Unit = "k" // switch to km
	directions.ManMaps = false
	//fmt.Println(string(directions.Dump("json")))
	results, err := directions.Get()
	if unexpected(err, t) {
		return
	}
	route := results.Route
	if route.HasTollRoad {
		t.Errorf("Route.HasTollRoad: Expected false ~ Received %v", route.HasUnpaved)
	}
	if route.HasUnpaved {
		t.Errorf("Route.HasUnpaved: Expected false ~ Received %v", route.HasUnpaved)
	}
	if !route.HasHighway {
		t.Errorf("Route.HasHighway: Expected true ~ Received %v", route.HasHighway)
	}
	distance := int(route.Distance)
	if testDistance != distance {
		t.Errorf("Route.Distance: Expected %d ~ Received %d", testDistance, distance)
	}
	if math.Abs(float64(testTime-route.Time)) > 60.0 { // tolerate minute difference
		t.Errorf("Route.Time: Expected %d ~ Received %d", testTime, route.Time)
	}
	legs := route.Legs
	if testLegs != len(legs) {
		t.Errorf("len(Route.Legs): Expected %d ~ Received %d", testLegs, len(legs))
	} else if testManeuvers != len(legs[0].Maneuvers) {
		t.Errorf("Route.Legs[0].Maneuvers: Expected %d ~ Received %d", testManeuvers, len(legs[0].Maneuvers))
	}
	statuscode := results.Info.Statuscode
	if testStatuscode != statuscode {
		t.Errorf("Info.Statuscode: Expected %d ~ Received %d", testStatuscode, statuscode)
	}
}

func TestDump(t *testing.T) {
	directions := NewDirections("3 Abbey Road, London, United Kingdom", []string{"3 Abbey Road, London, United Kingdom"})
	directions.NarrativeType = "none"
	directions.ManMaps = false
	results, err := directions.Get()
	if unexpected(err, t) {
		return
	}
	sessionID := results.Route.SessionID
	directions.SessionID = sessionID
	testDump := fmt.Sprintf(templDump, sessionID)
	dumpBytes, err := directions.Dump("json")
	if unexpected(err, t) {
		return
	}
	// avoid mapUrl as it contains a random sequence which is always different
	dump := strings.SplitN(string(dumpBytes), "mapUrl", 2)[0]
	if testDump != dump {
		t.Errorf("Expected\n%s\nReceived\n%s", testDump, dump)
	}
}
