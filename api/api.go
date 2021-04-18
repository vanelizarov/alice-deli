package api

import (
	"io/ioutil"
	"net/http"
)

/*
1. Получить список геозон - https://api.delimobil.ru/api/geozones
2. Получить список машин - https://api.delimobil.ru/api/cars?regionId=<id из п.1>
3. Получить информацию по авто - https://api.delimobil.ru/api/cars/<id>
 */

const (
	apiUrl = "https://api.delimobil.ru/api"
	geoZonesPath = "/geozones"
	carsPath = "/cars"
)

func fetch(url string) ([]byte, error) {
	apiRes, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(){
		_ = apiRes.Body.Close()
	}()

	body, err := ioutil.ReadAll(apiRes.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}