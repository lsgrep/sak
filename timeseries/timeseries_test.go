package timeseries_test

import (
	"testing"
	"time"

	. "github.com/lsgrep/sak/timeseries"

	"github.com/stretchr/testify/suite"
)

var (
	sec = time.Duration(time.Second)
)

type TimeSeriesTestSuite struct {
	suite.Suite
}

func (ts *TimeSeriesTestSuite) SetupTest() {
}

// The hook of `go test`
func TestTimeSeriesTestSuite(t *testing.T) {
	suite.Run(t, new(TimeSeriesTestSuite))
}

func (ts *TimeSeriesTestSuite) TestLen() {
	now := time.Now()

	s := NewTimeSeries("foo", 10*sec, 3)
	ts.Equal(0, s.Len())

	s.Add(now.Add(1*sec), 1001)
	ts.Equal(1, s.Len())

	s.Add(now.Add(1*sec), 1002) // non-strictly later dot will be ignored
	ts.Equal(1, s.Len())

	s.Add(now.Add(2*sec), 1002)
	s.Add(now.Add(3*sec), 1003)
	ts.Equal(3, s.Len())

	s.Add(now.Add(4*sec), 1004) // drop oldest dos to fit MaxLen
	ts.Equal(3, s.Len())

	s.Add(now.Add(13*sec), 1013) // drop oldest dos to fit Window
	ts.Equal(2, s.Len())
}

func (ts *TimeSeriesTestSuite) TestGetLastChangeTime() {
	now := time.Now()

	s := NewTimeSeries("foo", 100*sec, 10)
	s.Add(now.Add(1*sec), 99)
	s.Add(now.Add(2*sec), 99)
	s.Add(now.Add(2*sec), 100) // ignored
	s.Add(now.Add(3*sec), 99)
	s.Add(now.Add(4*sec), 100)
	s.Add(now.Add(5*sec), 100)
	ts.Equal(now.Add(4*sec), s.GetLastChangeTime())

	s = NewTimeSeries("foo", 100*sec, 1)
	s.Add(now.Add(1*sec), 99)
	s.Add(now.Add(2*sec), 99)
	s.Add(now.Add(3*sec), 100)
	s.Add(now.Add(4*sec), 100)
	ts.Equal(1, s.Len())
	ts.Equal(now.Add(3*sec), s.GetLastChangeTime())
}

func (ts *TimeSeriesTestSuite) TestGetLastDotByTime() {
	s := NewTimeSeries("foo", 100*sec, 10)

	now := time.Now()
	dot, ok := s.GetLastDotByTime(now) // querying an empty series
	ts.False(ok)

	dot1 := Dot{now.Add(10 * sec), 101}
	dot2 := Dot{now.Add(20 * sec), 102}
	dot3 := Dot{now.Add(30 * sec), 103}
	dot4 := Dot{now.Add(40 * sec), 104}
	s.AddDot(dot1)
	s.AddDot(dot2)
	s.AddDot(dot3)
	s.AddDot(dot4)

	dot, ok = s.GetLastDotByTime(now.Add(9 * sec))
	ts.False(ok)
	dot, ok = s.GetLastDotByTime(now.Add(10 * sec))
	ts.True(ok)
	ts.Equal(dot1, dot)
	dot, ok = s.GetLastDotByTime(now.Add(11 * sec))
	ts.True(ok)
	ts.Equal(dot1, dot)
	dot, ok = s.GetLastDotByTime(now.Add(20 * sec))
	ts.True(ok)
	ts.Equal(dot2, dot)
	dot, ok = s.GetLastDotByTime(now.Add(30 * sec))
	ts.True(ok)
	ts.Equal(dot3, dot)
	dot, ok = s.GetLastDotByTime(now.Add(35 * sec))
	ts.True(ok)
	ts.Equal(dot3, dot)
	dot, ok = s.GetLastDotByTime(now.Add(40 * sec))
	ts.True(ok)
	ts.Equal(dot4, dot)
	dot, ok = s.GetLastDotByTime(now.Add(99 * sec))
	ts.True(ok)
	ts.Equal(dot4, dot)
	ts.Equal(s.Last(), dot)
}
