package main

import (
	"fmt"

	"deli/alice"
	"deli/api"
	aerrors "deli/errors"
	"deli/geo"
)

const (
	MeasureUnitMeter     = "meter"
	MeasureUnitKiloMeter = "kilometer"
)

const (
	IntentNearby = "cars.nearby"
)

const (
	SlotAddress = "address"
	SlotRadius  = "radius"
	SlotMeasure = "measure"
)

const HelloText = "Привет! Я помогу вам найти каршеринг поблизости. Просто скажите \"найди машины рядом\" и адрес. Еще вы можете задать радиус поиска, для этого после адреса назовите расстояние, например 300м или 5км"

var ErrorNoAddress = aerrors.NewError(aerrors.ErrorNoAddress)

func getAddressFromIntentSlot(slot *alice.IntentSlot) (string, error) {
	switch slot.Type {
	case alice.EntityYandexString:
		return fmt.Sprintf("%v", slot.Value), nil
	case alice.EntityYandexGeo:
		val := slot.Value.(map[string]interface{})
		var addr string
		if street, exists := val["street"]; exists {
			addr += fmt.Sprintf("%v", street)
		}
		if houseNumber, exists := val["house_number"]; exists {
			addr += fmt.Sprintf(" %v", houseNumber)
		}
		return addr, nil
	}
	return "", aerrors.NewError(aerrors.ErrorNoAddress)
}

func getRadiusKMFromIntent(intent *alice.Intent) float64 {
	radius := api.DefaultNearbyKM
	radiusSlot, radiusExists := intent.Slots[SlotRadius]
	if radiusExists {
		radius, _ = radiusSlot.Value.(float64)
	}

	measure := 1.0
	measureSlot, measureExists := intent.Slots[SlotMeasure]
	if measureExists {
		if measureSlot.Value.(string) == MeasureUnitMeter {
			measure = 0.001
		}
	}

	return radius * measure
}

func getAvailableCars(req *alice.Request) ([]*api.Car, int, *aerrors.AliceError) {
	for name, intent := range req.Request.NLU.Intents {
		switch name {
		case IntentNearby:
			addrSlot, addrExists := intent.Slots[SlotAddress]
			if !addrExists {
				return nil, 0, ErrorNoAddress
			}
			addr, err := getAddressFromIntentSlot(addrSlot)
			if err != nil {
				return nil, 0, ErrorNoAddress
			}

			geocoderObj, err := geo.GetGeocoderObjectForAddress(addr)
			if err != nil {
				return nil, 0, ErrorNoAddress
			}

			geoZone, err := api.NewGeoZone()
			if err != nil {
				return nil, 0, ErrorNoAddress
			}

			userPoint, err := geocoderObj.Point.ToCoordinate()
			if err != nil {
				return nil, 0, ErrorNoAddress
			}

			userRegion, err := geoZone.GetUserRegion(userPoint)
			if err != nil {
				return nil, 0, aerrors.NewError(aerrors.ErrorNoCarsheringInUserRegion)
			}

			regionAvailability, err := api.NewAvailability(userRegion.ID)
			if err != nil {
				return nil, 0, aerrors.NewError(aerrors.ErrorNoCarsInUserRegion)
			}

			var radius = getRadiusKMFromIntent(intent)
			carsNearby := regionAvailability.GetCarsNearby(userPoint, radius)
			carsNearbyLen := len(carsNearby)
			if carsNearbyLen == 0 {
				return nil, 0, aerrors.NewError(aerrors.ErrorNoCarsNearby)
			}
			return carsNearby, carsNearbyLen, nil
		}
	}
	return nil, 0, aerrors.NewError(aerrors.ErrorNoIntent)
}

func declOfNum(number int, titles []string) string {
	cases := []int{2, 0, 1, 1, 1, 2}
	var currentCase int
	if number%100 > 4 && number%100 < 20 {
		currentCase = 2
	} else if number%10 < 5 {
		currentCase = cases[number%10]
	} else {
		currentCase = cases[5]
	}
	return titles[currentCase]
}

func Handler(req *alice.Request) (*alice.Response, error) {
	var resText string
	if len(req.Request.NLU.Intents) > 0 {
		cars, carsLen, err := getAvailableCars(req)
		if err == nil {
			carsWord := declOfNum(carsLen, []string{"машина", "машины", "машин"})
			resText = fmt.Sprintf("Рядом с вами нашлась %d %s", carsLen, carsWord)
			for _, car := range cars {
				resText += fmt.Sprintf("\n%s", car)
			}
		} else {
			resText = err.Error()
		}
	} else {
		resText = HelloText
	}

	return &alice.Response{
		Response: alice.ResponsePayload{
			Text:       resText,
			EndSession: false,
		},
		Session: alice.ResponseSession{
			SessionID: req.Session.SessionID,
			MessageID: req.Session.MessageID,
			UserID:    req.Session.UserID,
		},
		Version: req.Version,
	}, nil
}
