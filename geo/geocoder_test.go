package geo

import (
	"os"
	"testing"
)

func TestGetCoordinatesForAddress(t *testing.T) {
	_ = os.Setenv("GEOCODER_API_KEY", "278f85ec-4fdc-4107-aa20-9832af3d91e2")
	obj, err := GetGeocoderObjectForAddress("проспект мира 89")
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("Address not found")
	}

	t.Logf("%s; %s", obj.Name, obj.Point.Pos)
}

func BenchmarkGetGeocoderObjectForAddress(b *testing.B) {
	_ = os.Setenv("GEOCODER_API_KEY", "278f85ec-4fdc-4107-aa20-9832af3d91e2")
	var obj *Object
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		obj, _ = GetGeocoderObjectForAddress("проспект мира 89")
	}
	b.ReportAllocs()
	b.StopTimer()
	b.Logf("%s; %s", obj.Name, obj.Point.Pos)
}