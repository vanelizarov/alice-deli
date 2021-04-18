package api

import (
	"encoding/json"
	"fmt"
	"sort"

	"deli/geo"
	"github.com/paulmach/orb/geojson"
)

const (
	DefaultNearbyKM = 0.3
)

var carModelMap = map[string]string{
	"hyundai solaris": "Hyundai Solaris",
	"vw polo": "Volkswagen Polo",
	"smart fortwo coupe": "Smart Fortwo Coupe",
	"bmw 320i": "BMW 320i",
	"smart forfour": "Smart Forfour",
	"mini cooper 5d": "MINI Cooper 5D",
	"fiat 500": "Fiat 500",
	"bmw 320i premium": "BMW 320i Premium",
	"mercedes-benz glc 250": "Mercedes-Benz GLC 250",
	"mercedes-benz e 200": "Mercedes-Benz E 200",
	"mini cooper 3d": "MINI Cooper 3D",
	"nissan qashqai": "Nissan Qashqai",
	"renault sandero": "Renault Sandero",
	"kia rio": "KIA Rio",
	"kia sportage": "KIA Sportage",
	"kia rio x-line": "KIA Rio X-Line",
	"vw polo vi": "Volkswagen Polo VI",
	"skoda rapid": "ŠKODA RAPID",
	"kia rio x": "KIA Rio X",
}

type Car struct {
	ID       int64
	Model    string
	Distance float64
}

func (c Car) String() string {
	var distanceStr string
	if c.Distance < 1.0 {
		distanceStr = fmt.Sprintf("%dм", int(c.Distance * 1000))
	} else {
		distanceStr = fmt.Sprintf("%.2fкм", c.Distance)
	}
	return fmt.Sprintf("%s — %s", carModelMap[c.Model], distanceStr)
}

type Availability struct {
	GeoJSON *geojson.FeatureCollection
}

func NewAvailability(regionID int64) (*Availability, error) {
	body, err := fetch(fmt.Sprintf("%s%s?regionId=%d", apiUrl, carsPath, regionID))
	if err != nil {
		return nil, err
	}

	carsRes := new(Availability)
	if err = json.Unmarshal(body, &carsRes); err != nil {
		return nil, err
	}

	return carsRes, nil
}

func (c *Availability) GetCarsNearby(userCoordinate *geo.Coordinate, radiusKM float64) []*Car {
	result := make([]*Car, 0, 10)
	for _, carFeature := range c.GeoJSON.Features {
		distance := userCoordinate.DistanceToPoint(carFeature.Point())
		if distance <= radiusKM {

			result = append(result, &Car{
				ID:       int64(carFeature.ID.(float64)),
				Model:    fmt.Sprintf("%v", carFeature.Properties["model"]),
				Distance: distance,
			})
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Distance < result[j].Distance
	})
	return result
}
