//
// cumDist_test.go
//
// go test election2012.go states.go api.go cumDist.go cumDist_test.go
//
//
package main

import (
	"math"
	"testing"
)

const TOL = 1e-14

func checkFloat64(x, y, tol float64, test string, t *testing.T) {
	if math.Abs(x-y) > math.Abs(x*tol) {
		t.Errorf("Found %v, but expected %v for test %v", x, y, test)
	}
}

func TestErf(t *testing.T) {
	checkFloat64(erf(3), 0.999979214581871, TOL, "cumDist", t)
	checkFloat64(erf(-3), -0.999979214581871, TOL, "cumDist", t)
}

func TestCumDist(t *testing.T) {
	checkFloat64(cumDist(1.0, 0.0, 1.0), 0.841384034263321, TOL, "cumDist", t)
	checkFloat64(cumDist(40.0, 47.0, 10.0), 0.24195429670945612, TOL, "cumDist", t)
	checkFloat64(cumDist(12.0, 10.0, 2.5), 0.7881610565888237, TOL, "cumDist", t)
}

func TestPrOverX(t *testing.T) {
	checkFloat64(prOverX(1.0, 0.0, 1.0), 1.0-0.841384034263321, TOL, "cumDist", t)
	checkFloat64(prOverX(40.0, 47.0, 10.0), 1.0-0.24195429670945612, TOL, "cumDist", t)
	checkFloat64(prOverX(12.0, 10.0, 2.5), 1.0-0.7881610565888237, TOL, "cumDist", t)
}
