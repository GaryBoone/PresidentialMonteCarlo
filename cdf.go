//
// cdf.go
//

package main

import (
	"math"
)

func prOverX(x, μ, σ float64) float64 {
	return 1.0 - cdf(x, μ, σ)
}

// cdf returns the probability of a random variable (μ, σ) being less than x.
func cdf(x, μ, σ float64) float64 {
	return 0.5 * (1.0 + erf((x-μ)/(σ*math.Sqrt(2.0))))
}

func erf(x float64) float64 {
	a := 0.14
	val := math.Sqrt(1.0 - math.Exp(-x*x*(4.0/math.Pi+a*x*x)/(1.0+a*x*x)))
	if x >= 0.0 {
		return val
	}
	return -val
}
