// 
// state.go
// 

package main

import (
	"log"
	"math"
	"math/rand"
)

type StateProbability struct {
	state            string
	Obama            float64
	Romney           float64
	N                float64
	obamaPerc        float64
	σ                float64
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

	s.obamaPerc = s.Obama / s.N
	s.σ = math.Sqrt((s.obamaPerc - s.obamaPerc*s.obamaPerc) / s.N)
	if min_σ != 0.0 && s.σ < min_σ {
		s.σ = min_σ
	}
	s.ObamaProbability = prOverX(0.50, s.obamaPerc, s.σ)
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

func (s *StateProbability) logStateProbability() {
	if s.N != 0 {
		log.Printf("  %v: Obama polling=%6.4f, N=%d, σ=%6.4f --> Pr(Obama)=%6.4f\n",
			s.state, s.obamaPerc, int(s.N), s.σ, s.ObamaProbability)
	} else {
		if college[s.state].dem2008 {
			log.Printf("  %s voted Democratic in 2008.\n", s.state)
		} else {
			log.Printf("  %s voted Republican in 2008.\n", s.state)
		}
	}
}
