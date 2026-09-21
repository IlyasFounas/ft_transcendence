package haversine

import (
	"math"
)

type Coordinates struct {
	Lat float64
	Lon float64
}

// Constants for the radius of the Earth in meters and feet
const (
	EarthDiameterMeters = 12742176 // median diameter of the Earth in meters
	EarthRadiusMeters = 6371088 // median radius of the Earth in meters
	EarthRadiusFeet = 20902520 // 6371088 * 3.28084 > meters to feet conversion
)

// Converts from degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Distance calculates the distance in meters and feet between two coordinates using the Haversine formula
func Distance(p1, p2 Coordinates) (m, ft float64) {
	lat1 := degreesToRadians(p1.Lat)
	lon1 := degreesToRadians(p1.Lon)
	lat2 := degreesToRadians(p2.Lat)
	lon2 := degreesToRadians(p2.Lon)
	latDiff := lat2 - lat1
	lonDiff := lon2 - lon1

	// Haversine formula
	// hav0 = sin²(Δφ/2) + cos φ1 * cos φ2 * sin²(Δλ/2)
	sinHalfLat := math.Sin(latDiff / 2)
	sinHalfLon := math.Sin(lonDiff/2)
	hav0 := (sinHalfLat * sinHalfLat) + math.Cos(lat1)*math.Cos(lat2)*(sinHalfLon * sinHalfLon)

	// radians to central angle: 
	centralAngle := 2 * math.Atan2(math.Sqrt(hav0), math.Sqrt(1-hav0))
 
	m = centralAngle * EarthRadiusMeters
	ft = centralAngle * EarthRadiusFeet

	return m, ft
}