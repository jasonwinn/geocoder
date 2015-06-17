/* Exposes (partially) the mapquest directions api.

Reference: http://open.mapquestapi.com/directions/

Example:

directions := NewDirections("Amsterdam,Netherlands", []string{"Antwerp,Belgium"})
directions.Unit = "k" // switch to km

url := directions.URL("json")
json := directions.Dump("json")
xml := directions.Dump("xml")
distance := directions.Distance("k")
route := directions.Get().Route

*/

package geocoder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	directionURL = "https://open.mapquestapi.com/directions/v2/route?inFormat=kvp&key="
)

// Directions provide information on how to get from one location
// to one (or more) other locations together with copyright and statuscode info.
// (style, shape and mapstate options are not implemented)
type Directions struct {
	// Starting location
	From string
	// Ending location
	To []string
	// type of units for calculating distance: m (Miles, default) or k (Km)
	Unit string
	// fastest(default), shortest, pedestrian, multimodal, bicycle
	RouteType string
	// If true(default), reverse geocode call will be performed even on lat/long.
	DoReverseGeocode bool
	// none, text(default), html, microformat
	NarrativeType string
	// Encompass extra advice such as intersections (default false)
	EnhancedNarrative bool
	// The maximum number of Link IDs to return for each maneuver (default is 0)
	MaxLinkID int
	// en_US(default), en_GB, fr_CA, fr_FR, de_DE, es_ES, es_MX, ru_RU
	Locale string
	// Limited Access, Toll Road, Ferry, Unpaved, Seasonal Closure, Country Crossing
	Avoids []string
	// Link IDs of roads to absolutely avoid. May cause some routes to fail.
	MustAvoidLinkIDs []int
	// Link IDs of roads to try to avoid during route calculation without guarantee.
	TryAvoidLinkIDs []int
	// Whether state boundary crossings will be displayed in narrative. (default true)
	StateBoundaryDisplay bool
	// Whether country boundary crossings will be displayed in narrative. (default true)
	CountryBoundaryDisplay bool
	// Whether "End at" destination maneuver will be displayed in narrative. (default true)
	DestinationManeuverDisplay bool
	// To return a route shape without a mapState. (default false)
	FullShape bool
	// A value of < 1 favors cycling on non-bike lane roads. [0.1..1(default)..100]
	CyclingRoadFactor float64
	// DEFAULT_STRATEGY (default), AVOID_UP_HILL, AVOID_DOWN_HILL,AVOID_ALL_HILLS,FAVOR_UP_HILL,FAVOR_DOWN_HILL,FAVOR_ALL_HILLS
	RoadGradeStrategy string
	// cautious, normal, aggressive
	DrivingStyle string
	// Fuel efficiency, given as miles per gallon. (0..235 mpg)
	HighwayEfficiency float64
	// If true (default), a small staticmap is displayed per maneuver
	ManMaps bool
	// Walking speed, always in miles per hour independent from unit (default 2.5)
	WalkingSpeed float64
	// Session id
	SessionID string
}

// NewDirections is a constructor to initialize a Directions struct
// with mapquest defaults.
func NewDirections(from string, to []string) *Directions {
	return &Directions{
		From:                       from,
		To:                         to,
		Unit:                       "m",
		RouteType:                  "fastest",
		DoReverseGeocode:           true,
		NarrativeType:              "text",
		EnhancedNarrative:          false,
		MaxLinkID:                  0,
		Locale:                     "en_US",
		StateBoundaryDisplay:       true,
		CountryBoundaryDisplay:     true,
		DestinationManeuverDisplay: true,
		FullShape:                  false,
		CyclingRoadFactor:          1,
		RoadGradeStrategy:          "DEFAULT_STRATEGY",
		DrivingStyle:               "normal",
		HighwayEfficiency:          22,
		ManMaps:                    true,
		WalkingSpeed:               -1, // 2.5
		SessionID:                  "",
	}
}

// URL constructs the mapquest directions url (format is json or xml)
func (directions Directions) URL(format string) string {
	// http://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go
	var (
		n        int
		routeURL bytes.Buffer
	)
	writeStringInts := func(label string, arrayInts []int) {
		n = len(arrayInts)
		if n > 0 {
			routeURL.WriteString(label)
			for i, number := range arrayInts {
				routeURL.WriteString(strconv.Itoa(number))
				if i < n-1 {
					routeURL.WriteString(",")
				}
			}
		}
	}
	routeURL.WriteString(directionURL)
	routeURL.WriteString(apiKey)
	routeURL.WriteString("&outFormat=" + format)
	routeURL.WriteString("&from=" + url.QueryEscape(directions.From))
	for _, to := range directions.To {
		routeURL.WriteString("&to=" + url.QueryEscape(to))
	}
	routeURL.WriteString("&unit=" + directions.Unit)
	routeURL.WriteString("&routeType" + directions.RouteType)
	routeURL.WriteString("&narrativeType=" + directions.NarrativeType)
	routeURL.WriteString("&enhancedNarrative=" + strconv.FormatBool(directions.EnhancedNarrative))
	routeURL.WriteString("&maxLinkId=" + strconv.Itoa(directions.MaxLinkID))
	routeURL.WriteString("&locale=" + directions.Locale)
	for _, avoids := range directions.Avoids {
		routeURL.WriteString("&avoids=" + url.QueryEscape(avoids))
	}
	writeStringInts("&mustAvoidLinkIds=", directions.MustAvoidLinkIDs)
	writeStringInts("&tryAvoidLinkIds=", directions.TryAvoidLinkIDs)
	routeURL.WriteString("&stateBoundaryDisplay=" + strconv.FormatBool(directions.StateBoundaryDisplay))
	routeURL.WriteString("&countryBoundaryDisplay=" + strconv.FormatBool(directions.CountryBoundaryDisplay))
	routeURL.WriteString("&destinationManeuverDisplay=" + strconv.FormatBool(directions.DestinationManeuverDisplay))
	routeURL.WriteString("&fullShape=" + strconv.FormatBool(directions.FullShape))
	routeURL.WriteString("&cyclingRoadFactor=" + strconv.FormatFloat(directions.CyclingRoadFactor, 'f', -1, 64))
	routeURL.WriteString("&roadGradeStrategy=" + directions.RoadGradeStrategy)
	routeURL.WriteString("&drivingStyle=" + directions.DrivingStyle)
	routeURL.WriteString("&highwayEfficiency=" + strconv.FormatFloat(directions.HighwayEfficiency, 'f', -1, 64))
	routeURL.WriteString("&manMaps=" + strconv.FormatBool(directions.ManMaps))
	routeURL.WriteString("&walkingSpeed=" + strconv.FormatFloat(directions.WalkingSpeed, 'f', -1, 64))
	if directions.SessionID != "" {
		routeURL.WriteString("&sessionId=" + directions.SessionID)
	}
	return routeURL.String()
}

// Dump directions as undecoded json or xml bytes
func (directions Directions) Dump(format string) (data []byte, err error) {
	resp, err := http.Get(directions.URL(format))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return
}

// Distance calculated in km or miles (unit is k [km] or m [miles])
func (directions Directions) Distance(unit string) (distance float64, err error) {
	// these changes are made on a copy and as such invisible to caller
	directions.ManMaps = false
	directions.NarrativeType = "none"
	directions.Unit = unit
	resp, err := http.Get(directions.URL("json"))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Decode our JSON results
	results := DistanceResults{}
	err = decoder(resp).Decode(&results)
	if err != nil {
		return
	}
	if results.Info.Statuscode != 0 {
		err = results.Info
		return
	}
	distance = results.Route.Distance
	return
}

// Private json struct to retrieve the distance
type DistanceResults struct {
	Route struct {
		Distance float64 `json:"distance"`
	} `json:"route"`
	Info Info `json:"info"`
}

// Get the Direction Results (Route & Info)
func (directions Directions) Get() (results *DirectionsResults, err error) {
	resp, err := http.Get(directions.URL("json"))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Decode our JSON results
	results = &DirectionsResults{}
	err = decoder(resp).Decode(results)
	if err != nil {
		return
	}
	if results.Info.Statuscode != 0 {
		err = results.Info
	}
	return
}

// DirectionsResults can be decoded from the Directions route response
type DirectionsResults struct {
	Route Route `json:"route"`
	Info  Info  `json:"info"`
}

// Info about copyright and statuscode
type Info struct {
	Copyright struct {
		Text         string `json:"text"`         // "© 2014 MapQuest, Inc."
		ImageURL     string `json:"imageUrl"`     // "http://api.mqcdn.com/res/mqlogo.gif"
		ImageAltText string `json:"imageAltText"` // "© 2014 MapQuest, Inc."
	} `json:"copyright"`
	Statuscode int      `json:"statuscode"`
	Messages   []string `json:"messages"`
}

func (err Info) Error() string {
	if err.Statuscode != 0 {
		return fmt.Sprintf(
			"Error %d: %s",
			err.Statuscode, strings.Join(err.Messages, "; "))
	}
	return "No error"
}

// Route provides information on how to get from one location
// to one (or more) other locations.
type Route struct {
	// Returns true if at least one leg contains a Toll Road.
	HasTollRoad bool `json:"hasTollRoad"`
	// Returns true if at least one leg contains a Limited Access/Highway.
	HasHighway bool `json:"hasHighway"`
	// Returns true if at least one leg contains a Seasonal Closure.
	HasSeasonalClosure bool `json:"hasSeasonalClosure"`
	// Returns true if at least one leg contains an Unpaved.
	HasUnpaved bool `json:"hasUnpaved"`
	// Returns true if at least one leg contains a Country Crossing.
	HasCountryCross bool `json:"hasCountryCross"`
	// Returns lat/lng bounding rectangle of all points
	BoundingBox struct {
		UpperLeft  LatLng `json:"ul"`
		LowerRight LatLng `json:"lr"`
	} `json:"BoundingBox"`
	// Returns the calculated elapsed time in seconds for the route.
	Time int `json:"time"`
	// Returns the calculated elapsed time as formatted text in HH:MM:SS format.
	FormattedTime string `json:"formattedTime"`
	// Returns the calculated distance of the route.
	Distance float64 `json:"distance"`
	// The estimated amount of fuel used during the route
	FuelUsed float64 `json:"fuelUsed"`
	// A collection of leg objects, one for each "leg" of the route.
	Legs []Leg `json:"legs"`
	// Error report
	RouteError struct {
		Message   string `json:"message"`
		ErrorCode int    `json:"errorCode"`
	} `json:"routeError"`
	// A collection of locations in the form of an address.
	Locations []Location `json:"locations"`
	// Location sequence
	LocationSequence []int `json:"locationSequence"`
	// A unique identifier used to refer to a session
	SessionID string `json:"sessionId"`
	/* // Routing Options (not necessary as same as request)
	Options Options `json:"options"` */
}

// Leg contains the maneuvers describing how to get from one location
// to the next location. Multiple legs belong to a route.
type Leg struct {
	// The shape point index which starts a specific route segment.
	Index int `json:"index"`
	// Returns true if at least one maneuver contains a Toll Road.
	HasTollRoad bool `json:"hasTollRoad"`
	// Returns true if at least one maneuver contains a Limited Access/Highway.
	HasHighway bool `json:"hasHighway"`
	// Returns true if at least one maneuver contains a Seasonal Closure.
	HasSeasonalClosure bool `json:"hasSeasonalClosure"`
	// Returns true if at least one maneuver contains an Unpaved.
	HasUnpaved bool `json:"hasUnpaved"`
	// Returns true if at least one maneuver contains a Country Crossing.
	HasCountryCross bool `json:"hasCountryCross"`
	// Returns the calculated elapsed time in seconds for the leg.
	Time int `json:"time"`
	// Returns the calculated elapsed time as formatted text in HH:MM:SS format.
	FormattedTime string `json:"formattedTime"`
	// Returns the calculated distance of the leg.
	Distance float64 `json:"distance"`
	// A collection of Maneuver objects
	Maneuvers []Maneuver `json:"maneuvers"`
	// Road grade avoidance strategies
	RoadGradeStrategy [][]int `json:"roadGradeStrategy"`

	// Collapsed Narrative Parameters if the user is familiar with the area

	// The origin index is the index of the first non-collapsed maneuver.
	OrigIndex int `json:"origIndex"`
	//  The rephrased origin narrative string for the first non-collapsed maneuver.
	OrigNarrative string `json:"origNarrative"`
	// The destination index is the index of the last non-collapsed maneuver.
	DestIndex int `json:"destIndex"`
	// The rephrased destination narrative string for the destination maneuver.
	DestNarrative string `json:"destNarrative"`
	/* RoadGradeStrategy ??? `json:"roadGradeStrategy"` */
}

// Maneuver describes each one step in a route narrative.
// Multiple maneuvers belong to a leg.
type Maneuver struct {
	Index int `json:"index"`
	// Returns the calculated elapsed time in seconds for the maneuver.
	Time int `json:"time"`
	// Returns the calculated elapsed time as formatted text in HH:MM:SS format.
	FormattedTime string `json:"formattedTime"`
	// Returns the calculated distance of the maneuver.
	Distance float64 `json:"distance"`
	// Text name, extra text, type (road shield), direction and image
	Signs []Sign `json:"signs"`
	// An URL to a static map of this maneuver.
	MapURL string `json:"mapUrl"`
	// Textual driving directions for a particular maneuver.
	Narrative string `json:"narrative"`
	/* // A collection of maneuverNote objects, one for each maneuver.
	   ManeuverNotes []ManeuverNote `json:"maneuverNotes"` */
	// none=0,north=1,northwest=2,northeast=3,south=4,southeast=5,southwest=6,west=7,east=8
	Direction int `json:"direction"`
	// Name of the direction
	DirectionName string `json:"directionName"`
	// Collection of street names this maneuver applies to
	Streets []string `json:"streets"`
	// none=0,portions toll=1,portions unpaved=2,possible seasonal road closure=4,gate=8,ferry=16,avoid id=32,country crossing=64,limited access (highways)=128
	Attributes int `json:"attributes"`
	// straight=0,slight right=1,right=2,sharp right=3,reverse=4,sharp left=5,left=6,slight left=7,right u-turn=8,
	// left u-turn=9,right merge=10,left merge=11,right on ramp=12,left on ramp=13,right off ramp=14,left off ramp=15,right fork=16,left fork=17,straight fork=18
	TurnType int `json:"direction"`
	// 1st shape point latLng for a particular maneuver (eg for zooming)
	StartPoint LatLng `json:"startPoint"`
	// Icon
	IconURL string `json:"iconUrl"`
	// "AUTO", "WALKING", "BICYCLE", "RAIL", "BUS" (future use)
	TransportMode string `json:"transportMode"`
	// Link Ids of roads
	LinkIDs []int `json:"linkIds"`
}

// Sign specifies information for a particular maneuver.
type Sign struct {
	Text      string `json:"text"`
	ExtraText string `json:"extraText"`
	Direction int    `json:"direction"`
	// road shield
	Type int `json:"type"`
	// Image
	URL string `json:"url"`
}

/* TODO
// Maneuver notes can exist for Timed Turn Restrictions, Timed Access Roads, HOV Roads, Seasonal Closures, and Timed Direction of Travel.
type ManeuverNote struct {
	RuleId int `json:"ruleId"`
	// ... incomplete
}
*/
