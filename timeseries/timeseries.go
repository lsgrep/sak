/*
This package implements a struct TimeSeries for managing a time series. It provides APIs for
querying various information of a time series.
*/
package timeseries

import "time"

type TimeSeries struct {
	//--------------- Constants --------------------
	Name   string
	Window time.Duration
	MaxLen int
	//----------------------------------------------

	lastValChangeTime time.Time

	dots []Dot
}

type Dot struct {
	Time  time.Time
	Value float64
}

// Oldest dots will be dropped (when adding new dot or when calling `DropToFit`) until the remaining
// dots are strictly within `window` and no more than `maxLen`.
//
// Precondition: `window` > 0 and `maxLen` > 0
func NewTimeSeries(name string, window time.Duration, maxLen int) *TimeSeries {
	return &TimeSeries{
		Name:   name,
		Window: window,
		MaxLen: maxLen,

		dots: nil,
	}
}

// Return the length of the series, i.e. the number of dots in the series (NOT including the dropped
// ones).
func (ts *TimeSeries) Len() int {
	return len(ts.dots)
}

// Get the last dot in the time series.
//
// Precondition: `ts` has at least one dot.
func (ts *TimeSeries) Last() Dot {
	return ts.dots[len(ts.dots)-1]
}

// Add a new dot to the time series. Remove old dots when necessary.
//
// Return false and ignore the new dot if `time` is not strictly later than the last dot in the
// time series.
func (ts *TimeSeries) Add(time time.Time, value float64) bool {
	if ts.Len() > 0 && !time.After(ts.Last().Time) {
		return false
	}

	if ts.Len() > 0 && ts.Last().Value != value {
		ts.lastValChangeTime = time
	}

	newDot := Dot{time, value}
	ts.add(newDot)
	return true
}

// A shortcut of `ts.Add(dot.Time, dot.Value)`
func (ts *TimeSeries) AddDot(dot Dot) bool {
	return ts.Add(dot.Time, dot.Value)
}

// Add a new dot without checking its time and then pop out oldest dots when necessary.
func (ts *TimeSeries) add(newDot Dot) {
	ts.dots = append(ts.dots, newDot)
	ts.DropToFit(newDot.Time)
}

// Drop oldest dots immediately until the remaining dots are strictly within `window` and no more
// than `maxLen`.
func (ts *TimeSeries) DropToFit(timeline time.Time) {
	for ts.Len() > ts.MaxLen {
		ts.dots = ts.dots[1:]
	}

	cutLine := timeline.Add(-ts.Window)
	for ts.Len() > 0 && !ts.dots[0].Time.After(cutLine) {
		ts.dots = ts.dots[1:]
	}
}

// Get the time when the last time that a new dot has a different value than its predecessor.
func (ts *TimeSeries) GetLastChangeTime() time.Time {
	return ts.lastValChangeTime
}

// Get the last dot on or before `time`. Return `ok` as false if no such a dot.
// Runtime O(log(Len))
func (ts *TimeSeries) GetLastDotByTime(time time.Time) (result Dot, ok bool) {
	head, tail := 0, ts.Len()-1
	for head <= tail {
		mid := (head + tail) >> 1
		if ts.dots[mid].Time.After(time) {
			tail = mid - 1
		} else {
			head = mid + 1
		}
	}
	if tail < 0 {
		return Dot{}, false
	}
	return ts.dots[tail], true
}
