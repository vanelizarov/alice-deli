package geo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type GeocoderResponse struct {
	Response Response
}

type Response struct {
	ObjectCollection ObjectCollection `json:"GeoObjectCollection"`
}

type ObjectCollection struct {
	FeatureMembers []*FeatureMember `json:"featureMember"`
}

type FeatureMember struct {
	Object Object `json:"GeoObject"`
}

type Object struct {
	Name string
	Description string
	Point Point
}

type Point struct {
	Pos string
}

func (p Point) ToCoordinate() (*Coordinate, error) {
	coords := strings.Split(p.Pos, " ")
	var lon, lat float64
	lon, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, err
	}
	lat, err = strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, err
	}
	return NewCoordinate(lon, lat), nil
}

func GetGeocoderObjectForAddress(addr string) (*Object, error) {
	apiKey, ok := os.LookupEnv("GEOCODER_API_KEY")
	if !ok {
		return nil, errors.New("GEOCODER_API_KEY not found")
	}

	req, err := http.NewRequest("GET", "https://geocode-maps.yandex.ru/1.x", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("geocode", addr)
	q.Add("apikey", apiKey)
	q.Add("format", "json")
	q.Add("results", "1")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(){
		_ = res.Body.Close()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	r := new(GeocoderResponse)
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if len(r.Response.ObjectCollection.FeatureMembers) == 0 {
		return nil, nil
	}

	return &r.Response.ObjectCollection.FeatureMembers[0].Object, nil
}