package numeric

import (
	"github.com/lsgrep/sak/utest"
	"testing"
)

type TestNumericSuite struct {
	utest.RequireSuite
}

// The hook of `go test`
func TestRun_TestNumericSuite(t *testing.T) {
	utest.Run(t, new(TestNumericSuite))
}

func (t *TestNumericSuite) TestStringToUint64E8() {
	t.StringToUint64E8_Equal("0.00000", 0)
	t.StringToUint64E8_Equal("0.00000001", 1)
	t.StringToUint64E8_Equal("0.0001", 0.0001e8)
	t.StringToUint64E8_Equal("1234.56780", 123456780000)
	t.StringToUint64E8_Equal("184467440737.09551615", 0xFFFFFFFFFFFFFFFF)

	t.StringToUint64E8_Error("-1.2")
	t.StringToUint64E8_Error("0.000000005")
	t.StringToUint64E8_Error("184467440737.09551616")
}

func (t *TestNumericSuite) StringToUint64E8_Equal(x string, expectedE8 uint64) {
	actualE8, err := StringToUint64E8(x)
	t.NoError(err)
	t.Equal(expectedE8, actualE8)
}

func (t *TestNumericSuite) StringToUint64E8_Error(x string) {
	valE8, err := StringToUint64E8(x)
	t.Error(err)
	t.Equal(uint64(0), valE8)
}

func (t *TestNumericSuite) TestStringToUint64E10() {
	t.StringToUint64E10_Equal("0.00000", 0)
	t.StringToUint64E10_Equal("0.00000001", 100)
	t.StringToUint64E10_Equal("0.0001", 0.01e8)
	t.StringToUint64E10_Equal("1234.56780", 12345678000000)
	t.StringToUint64E10_Equal("1844674407.3709551615", 0xFFFFFFFFFFFFFFFF)

	t.StringToUint64E10_Error("-1.2")
	t.StringToUint64E10_Error("0.00000000005")
	t.StringToUint64E10_Error("1844674407.3709551616")
}

func (t *TestNumericSuite) StringToUint64E10_Equal(x string, expectedE10 uint64) {
	actualE10, err := StringToUint64E10(x)
	t.NoError(err)
	t.Equal(expectedE10, actualE10)
}

func (t *TestNumericSuite) StringToUint64E10_Error(x string) {
	valE10, err := StringToUint64E10(x)
	t.Error(err)
	t.Equal(uint64(0), valE10)
}
