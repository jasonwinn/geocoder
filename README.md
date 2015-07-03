Mapquest geocoder and directions for Go (golang)
================================================

## What does it do

* Returns a Longitude and Latitude for a given string query
* Returns an address for a Longitude and Longitude
* Returns directions between two or more points. (JSON or XML)

## API Key
Get a free API Key at [http://mapquestapi.com](http://mapquestapi.com)

## Why MapQuest API?
Google Maps Geocoding API has a limitation that prohibits querying their geocoding API unless you will be displaying the results on a Google Map. Google directions is limited to 2 requests per second.

## Install

* go get "github.com/jasonwinn/geocoder"
* import "github.com/jasonwinn/geocoder"
 

## Examples

### Set API Key

You'll want to set an api key for the Mapquest API to go into production.
```go
// this is the testing key used in `go test`
SetAPIKey("Fmjtd%7Cluub256alu%2C7s%3Do5-9u82ur")
```

### Geocode
```go
  query := "Seattle WA"
  lat, lng, err := geocoder.Geocode(query)
  if err != nil {
    panic("THERE WAS SOME ERROR!!!!!")
  }
  
  // 47.6064, -122.330803
 
```

### Reverse Geocode
```go
  address, err := geocoder.ReverseGeocode(47.6064, -122.330803)
  if err != nil {
    panic("THERE WAS SOME ERROR!!!!!")
  }

  address.Street 	        // 542 Marion St   
  address.City 		        // Seattle
  address.State 	        // WA
  address.PostalCode 	    // 98104 
  address.County 	        // King
  address.CountryCode       // US 

```

### Directions
```go
  directions := NewDirections("Amsterdam,Netherlands", []string{"Antwerp,Belgium"})
  results, err := directions.Get()
  if err != nil {
    panic("THERE WAS SOME ERROR!!!!!")
  }

  route := results.Route
  time:= route.Time
  legs:= route.Legs
  distance:= route.Distance

  // or get distance with this shortcut
  //
  // use "k" to return result in km
  // use "m" to return result in miles
  distance, err := directions.Distance("k")
  if err != nil {
    panic("THERE WAS SOME ERROR!!!!!")
  }
```

## Documentation

[https://godoc.org/github.com/jasonwinn/geocoder](https://godoc.org/github.com/jasonwinn/geocoder)

