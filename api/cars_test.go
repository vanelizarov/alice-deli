package api

import (
	"testing"

	"deli/geo"
)

func BenchmarkGetCarsNearby(b *testing.B) {
	me := geo.NewCoordinate(
		37.63632017059521,
		55.80415798653096,
	)

	gz, err := NewGeoZone()
	if err != nil {
		b.Fatal(err)
	}

	region, err := gz.GetUserRegion(me)
	if err != nil {
		b.Fatal(err)
	}
	//b.Logf("Region for %s: %s (ID: %d)", me, region.TitleRu, region.ID)

	a, err := NewAvailability(region.ID)
	if err != nil {
		b.Fatal(err)
	}
	//b.Logf("Cars for %s (ID: %d): %v", region.TitleRu, region.ID, a.Cars)

	var carsNearby []*Car
	var carsNearbyLen int
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		carsNearby = a.GetCarsNearby(me, DefaultNearbyKM)
		carsNearbyLen = len(carsNearby)
		//if carsNearbyLen == 0 {
		//	b.Fatalf("No cars found in radius %f", nearbyKM)
		//}
	}
	b.ReportAllocs()
	b.StopTimer()
	b.Logf("Found %d cars in radius %fkm: %v", carsNearbyLen, DefaultNearbyKM, carsNearby)
}