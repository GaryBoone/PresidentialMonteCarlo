//
// cdf_test.go
//
// go test election2012.go state.go api.go cdf.go parse.go college.go cdf_test.go
//

package main

import (
	"math"
	"testing"
)

const TOL = 1e-14

func checkFloat64(found, expected, tol float64, test string, t *testing.T) {
	if math.Abs(found-expected) > math.Abs(found*tol) {
		t.Errorf("Found %v, but expected %v for test %v", found, expected, test)
	}
}

func TestErf(t *testing.T) {
	checkFloat64(erf(3), 0.999979214581871, TOL, "erf", t)
	checkFloat64(erf(-3), -0.999979214581871, TOL, "erf", t)
}

func TestCdf(t *testing.T) {
	checkFloat64(cdf(1.0, 0.0, 1.0), 0.841384034263321, TOL, "cdf", t)
	checkFloat64(cdf(40.0, 47.0, 10.0), 0.24195429670945612, TOL, "cdf", t)
	checkFloat64(cdf(12.0, 10.0, 2.5), 0.7881610565888237, TOL, "cdf", t)
}

func TestPrOverX(t *testing.T) {
	checkFloat64(prOverX(1.0, 0.0, 1.0), 1.0-0.841384034263321, TOL, "prOverX", t)
	checkFloat64(prOverX(40.0, 47.0, 10.0), 1.0-0.24195429670945612, TOL, "prOverX", t)
	checkFloat64(prOverX(12.0, 10.0, 2.5), 1.0-0.7881610565888237, TOL, "prOverX", t)
}
