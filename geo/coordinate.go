package geo

import (
	"fmt"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

const (
	earthRadiusKM = 6371.009
)

type Coordinate struct {
	Point orb.Point
}

func NewCoordinate(lon, lat float64) *Coordinate {
	return &Coordinate{
		Point: orb.Point{lon, lat},
	}
}

func (c Coordinate) String() string {
	return fmt.Sprintf("Coordinate{Lon: %f, Lat: %f}", c.Point.Lon(), c.Point.Lat())
}

func (c Coordinate) DistanceToPoint(c2 orb.Point) float64 {
	lat1 := degToRad(c.Point.Lat())
	lat2 := degToRad(c2.Lat())
	dlat := lat2 - lat1
	dlon := degToRad(c2.Lon() - c.Point.Lon())

	a := sinSqr(dlat*0.5) + math.Cos(lat1)*math.Cos(lat2)*sinSqr(dlon*0.5)

	return 2 * earthRadiusKM * math.Asin(math.Sqrt(a))
}

func (c Coordinate) DistanceTo(c2 *Coordinate) float64 {
	return c.DistanceToPoint(c2.Point)
}

func (c Coordinate) IsInsidePolygon(fc *geojson.FeatureCollection) bool {
	for _, feature := range fc.Features {
		mp, isMulti := feature.Geometry.(orb.MultiPolygon)
		if isMulti {
			if planar.MultiPolygonContains(mp, c.Point) {
				return true
			}
		} else {
			poly, isPoly := feature.Geometry.(orb.Polygon)
			if isPoly {
				if planar.PolygonContains(poly, c.Point) {
					return true
				}
			}
		}
	}
	return false
}

func degToRad(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func sinSqr(radians float64) float64 {
	sin := math.Sin(radians)
	return sin * sin
}