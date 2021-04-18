package geo

import (
	"testing"
)

func BenchmarkDistanceTo(b *testing.B) {
	me := NewCoordinate(
		37.63632017059521,
		55.80415798653096,
	)
	car := NewCoordinate(
		37.6305,
		55.8046,
	)
	var d float64
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d = me.DistanceTo(car)
	}
	b.ReportAllocs()
	b.StopTimer()
	b.Log(d)
}
