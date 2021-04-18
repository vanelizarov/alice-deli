package api

import (
	"encoding/json"

	"deli/geo"
	"github.com/paulmach/orb/geojson"
)

type GeoJSON struct {
	Driving *geojson.FeatureCollection
	Parking *geojson.FeatureCollection
}

type Region struct {
	ID int64
	Name string
	Title string
	TitleRu string
	GeoJSON GeoJSON
}

type GeoZone struct {
	Data []*Region
}

func NewGeoZone() (*GeoZone, error) {
	body, err := fetch(apiUrl + geoZonesPath)
	if err != nil {
		return nil, err
	}

	gz := new(GeoZone)
	if err = json.Unmarshal(body, &gz); err != nil {
		return nil, err
	}

	return gz, nil
}

func (g *GeoZone) GetUserRegion(userCoordinate *geo.Coordinate) (*Region, error) {
	for _, r := range g.Data {
		if isInside := userCoordinate.IsInsidePolygon(r.GeoJSON.Driving); isInside {
			return r, nil
		}
	}

	return nil, nil
}