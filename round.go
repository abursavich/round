// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package round provides a few rounding utility functions.
package round

import "time"

const (
	hour   = int64(time.Hour)
	minute = int64(time.Minute)
	second = int64(time.Second)
)

var pow10tab [20]uint64

func init() {
	pow10tab[0] = 1
	for e := 1; e < len(pow10tab); e++ {
		pow10tab[e] = 10 * pow10tab[e-1]
	}
}

// Duration returns the result of rounding d to the nearest multiple of n.
// The rounding behavior for halfway values is to round up. If n <= 0,
// Round returns d unchanged.
func Duration(d, n time.Duration) time.Duration {
	return time.Duration(Int64(int64(d), int64(n)))
}

// DurationN returns the result of rounding d to n significant decimal figures
// for standard string formatting. The rounding behavior for halfway values
// is to round up. If n <= 0, DurationN returns d unchanged.
//
// Examples (representing time.Duration values as strings):
//	DurationN("1h35m42.567s", 1) == "2h0m0s"
//	DurationN("1h35m42.567s", 2) == "1h40m0s"
//	DurationN("1h35m42.567s", 3) == "1h36m0s"
//	DurationN("1h35m42.567s", 4) == "1h35m40s"
//	DurationN("3.789s", 1) == "4s"
//	DurationN("1.567ms", 3) == "1.57ms"
//	DurationN("-41.5ms", 2) == "-42ms"
func DurationN(d time.Duration, n int) time.Duration {
	if n <= 0 {
		return d
	}
	if d < 0 {
		return -DurationN(-d, n)
	}
	v := int64(d)
	if v >= hour {
		k := i64digits(v / hour)
		if k >= n {
			return time.Duration(Int64(v, i64pow10(hour, k-n)))
		}
		n -= k
		k = i64digits(v % hour / minute)
		if k >= n {
			return time.Duration(Int64(v, i64pow10(minute, k-n)))
		}
		return time.Duration(Int64(v, i64pow10(100*second, k-n)))
	}
	if v >= minute {
		k := i64digits(v / minute)
		if k >= n {
			return time.Duration(Int64(v, i64pow10(minute, k-n)))
		}
		return time.Duration(Int64(v, i64pow10(100*second, k-n)))
	}
	return time.Duration(Int64N(v, n))
}

// Int64 returns the result of rounding v to the nearest multiple of n.
// The rounding behavior for halfway values is to round up. If n <= 0,
// Int64 returns v unchanged.
func Int64(v, n int64) int64 {
	if n <= 0 {
		return v
	}
	neg := v < 0
	if neg {
		v = -v
	}
	if r := v % n; r+r < n {
		v = v - r
	} else {
		v = v + n - r
	}
	if neg {
		return -v
	}
	return v
}

// Int64N returns the result of rounding v to n significant decimal figures.
// The rounding behavior for halfway values is to round up. If n <= 0,
// Int64N returns v unchanged.
func Int64N(v int64, n int) int64 {
	if n <= 0 {
		return v
	}
	if e := i64digits(v) - n; e > 0 {
		return Int64(v, i64pow10(1, e))
	}
	return v
}

// Uint64 returns the result of rounding v to the nearest multiple of n.
// The rounding behavior for halfway values is to round up. If n == 0,
// Uint64 returns v unchanged.
func Uint64(v, n uint64) uint64 {
	if n == 0 {
		return v
	}
	r := v % n
	if r+r < n {
		return v - r
	}
	return v + n - r
}

// Uint64N returns the result of rounding v to n significant decimal figures.
// The rounding behavior for halfway values is to round up. If n == 0,
// Uint64N returns v unchanged.
func Uint64N(v uint64, n int) uint64 {
	if n == 0 {
		return v
	}
	if e := u64digits(v) - n; e > 0 {
		return Uint64(v, u64pow10(1, e))
	}
	return v
}

// return range: [1,19]
func i64digits(v int64) int {
	if v < 0 {
		v = -v
	}
	// return u64digits(uint64(v))
	d := 1
	for v > 9 {
		v /= 10
		d++
	}
	return d
}

// return range: [1,20]
func u64digits(v uint64) int {
	d := 1
	for v > 9 {
		v /= 10
		d++
	}
	return d
}

// e range: [-19,19]
func i64pow10(b int64, e int) int64 {
	if e < 0 {
		return b / int64(pow10tab[-e])
	}
	return b * int64(pow10tab[e])
}

// e range: [-20,20]
func u64pow10(b uint64, e int) uint64 {
	if e < 0 {
		return b / pow10tab[-e]
	}
	return b * pow10tab[e]
}
