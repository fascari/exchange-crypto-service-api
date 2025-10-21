package testdata

const (
	SmallAmount  = 1.0
	MediumAmount = 10.0
	LargeAmount  = 100.0
)

func ExcessiveAmount(baseBalance float64) float64 {
	return baseBalance + LargeAmount
}
