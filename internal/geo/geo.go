package geo

import "math"

const (
	earthRadiusKm = 6371.0
	KmPerDegree   = 111.0
)

type Hub struct {
	Name       string
	Lat        float64
	Lon        float64
	DistanceKm float64
}

type BoundingBox struct {
	MinLat, MaxLat float64
	MinLon, MaxLon float64
}

func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// Returned BB fully encloses a circle of the given radius (in km) centred at (lat, lon)
// Can include points that are further than requested (in the corners),
// so additional filtering is still required.
func GetBoundingBox(lat, lon, distanceKm float64) BoundingBox {
	latOffset := distanceKm / KmPerDegree
	lonOffset := distanceKm / (KmPerDegree * math.Cos(degreesToRadians(lat)))

	return BoundingBox{
		MinLat: lat - latOffset,
		MaxLat: lat + latOffset,
		MinLon: lon - lonOffset,
		MaxLon: lon + lonOffset,
	}
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}
