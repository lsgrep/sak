package numeric

import (
	"fmt"
	"math/big"
)

// `x` is string represented decimal real number having at most 8 digits after the decimal point.
// Besides, `x * 1e8` must fit in uint64.
func StringToUint64E8(x string) (uint64, error) {
	return stringToUint64WithPrecision(x, 1e8)
}

// `x` is string represented decimal real number having at most 10 digits after the decimal point.
// Besides, `x * 1e10` must fit in uint64.
func StringToUint64E10(x string) (uint64, error) {
	return stringToUint64WithPrecision(x, 1e10)
}

func stringToUint64WithPrecision(x string, precision int64) (uint64, error) {
	errConv := fmt.Errorf("Cannot convert %s to E:%v amount.", x, precision)

	r := new(big.Rat)
	if _, err := fmt.Sscan(x, r); err != nil {
		return 0, errConv
	}
	r.Mul(r, big.NewRat(precision, 1))
	if !r.IsInt() {
		return 0, errConv
	}

	val := r.Num()
	if !val.IsUint64() {
		return 0, errConv
	}
	return val.Uint64(), nil
}

// Converts a uint64-e8 number to string real number with 8 digits after the decimal point.
//
// For example, 154.321e8 -> "154.32100000"
func Uint64E8ToString(x uint64) string {
	integral := x / 1e8
	fraction := x - integral*1e8
	return fmt.Sprintf("%d.%08d", integral, fraction)
}

// Converts a uint64-e10 number to string real number with 10 digits after the decimal point.
//
// For example, 154.321e10 -> "154.3210000000"
func Uint64E10ToString(x uint64) string {
	integral := x / 1e10
	fraction := x - integral*1e10
	return fmt.Sprintf("%d.%010d", integral, fraction)
}
