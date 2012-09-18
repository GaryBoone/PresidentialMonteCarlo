//
// cumDist.go
//
// https://1e47a410-a-62cb3a1a-s-sites.googlegroups.com/site/winitzki/sergei-winitzkis-files/erf-approx.pdf
// Sergei Winitzki
//

package main

import (
	"math"
)

func prOverX(x, μ, σ float64) float64 {
	return 1.0 - cumDist(x, μ, σ)
}

// cumDist returns the probability of a random variable (μ, σ) being less than x.
func cumDist(x, μ, σ float64) float64 {
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
