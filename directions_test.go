package geocoder

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

const (
	templDump      = "{\"route\":{\"hasTollRoad\":false,\"hasBridge\":false,\"distance\":0.27,\"shape\":{\"legIndexes\":[],\"maneuverIndexes\":[],\"shapePoints\":[51.529315,-0.269962,51.529087,-0.271016,51.52901,-0.271373,51.52885,-0.272208,51.528568,-0.274354,51.528472,-0.275138,51.528324,-0.27602]},\"hasTunnel\":false,\"hasHighway\":false,\"computedWaypoints\":[],\"routeError\":{\"errorCode\":-400,\"message\":\"\"},\"formattedTime\":\"00:00:41\",\"sessionId\":\"%s\",\"realTime\":-1,\"hasSeasonalClosure\":false,\"hasCountryCross\":false,\"fuelUsed\":0.02,\"legs\":[{\"hasTollRoad\":false,\"hasBridge\":false,\"destNarrative\":\"\",\"distance\":0.27,\"hasTunnel\":false,\"hasHighway\":false,\"index\":0,\"formattedTime\":\"00:00:41\",\"origIndex\":-1,\"hasSeasonalClosure\":false,\"hasCountryCross\":false,\"roadGradeStrategy\":[],\"destIndex\":-1,\"time\":41,\"hasUnpaved\":false,\"origNarrative\":\"\",\"maneuvers\":[{\"distance\":0.27,\"streets\":[\"Coronation Road\"],\"narrative\":\"\",\"turnType\":0,\"startPoint\":{\"lng\":-0.269962,\"lat\":51.529315},\"index\":0,\"formattedTime\":\"00:00:41\",\"directionName\":\"West\",\"maneuverNotes\":[],\"linkIds\":[91225298,91225281,91222079,91222078,91222062],\"signs\":[],\""
	testDistance   = 157
	testLegs       = 1
	testManeuvers  = 52
	testStatuscode = 0
	testTime       = 6085
	testUnit       = "m"
	testURL        = "https://open.mapquestapi.com/directions/v2/route?inFormat=kvp&key=Fmjtd%7Cluub256alu%2C7s%3Do5-9u82ur&outFormat=json&from=Amsterdam%2CNetherlands&to=Antwerp%2CBelgium&unit=m&routeTypefastest&narrativeType=text&enhancedNarrative=false&maxLinkId=0&locale=en_US&avoids=Ferry&mustAvoidLinkIds=5,7&stateBoundaryDisplay=true&countryBoundaryDisplay=true&destinationManeuverDisplay=true&fullShape=false&cyclingRoadFactor=1&roadGradeStrategy=DEFAULT_STRATEGY&drivingStyle=normal&highwayEfficiency=22&manMaps=true&walkingSpeed=-1"
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
	directions.Avoids = []string{"Ferry"}
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
	if err.Error() != "Error 402: We are unable to route with the given locations." {
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
	directions := NewDirections("1 Coronation Road, London, United Kingdom", []string{"3 Coronation Road, London, United Kingdom"})
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
