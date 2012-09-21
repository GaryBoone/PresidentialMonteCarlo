// 
// state.go
// 
// 
package main

import (
	"math"
	"math/rand"
)

type StateProbability struct {
	state            string
	Obama            float64
	Romney           float64
	N                float64
	ObamaProbability float64
}

// Update state data with a new poll. The new N is calculated with the actual 
// number of votes for Obama and Romney, not the N of the poll. The effect 
// is to not count undecideds and Others. Essentially, the poll is reduced 
// to a new poll between the two potential winners. In both cases, that's
// what actually happens. Because the N is reduced, the uncertainty is 
// increased as it should be.
func (s *StateProbability) update(oPerc, rPerc, pollSize int) {
	obamaVotes := float64(oPerc) * float64(pollSize) / 100.0
	romneyVotes := float64(rPerc) * float64(pollSize) / 100.0
	s.Obama += obamaVotes
	s.Romney += romneyVotes
	s.N += obamaVotes + romneyVotes

	obamaPerc := s.Obama / s.N
	σ := math.Sqrt((obamaPerc - obamaPerc*obamaPerc) / s.N)
	s.ObamaProbability = prOverX(0.50, obamaPerc, σ)

	// fmt.Printf("%v, O%%=%6.4f, σ=%6.4f, Pr(Obama)=%6.4f, votes=%d\n",
	// 	s.state, obamaPerc, σ, s.ObamaProbability, s.N)
}

func (s *StateProbability) simulateElection() int {
	if s.N != 0 {
		if rand.Float64() < s.ObamaProbability {
			return college[s.state].votes
		}
	} else {
		// give state to 2008 winner
		if college[s.state].dem2008 {
			return college[s.state].votes
		}
	}
	return 0
}
