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
	Democrat         float64
	Republican       float64
	N                float64
	democratPerc     float64
	σ                float64
	DemocratProbability float64
}

// Update state data with a new poll. The new N is calculated with the actual 
// number of votes for the Democrat and Republican, not the N of the poll. The effect 
// is to not count undecideds and Others. Essentially, the poll is reduced 
// to a new poll between the two potential winners. In both cases, that's
// what actually happens. Because the N is reduced, the uncertainty is 
// increased as it should be.
func (s *StateProbability) update(oPerc, rPerc, pollSize int) {
	democratVotes := float64(oPerc) * float64(pollSize) / 100.0
	republicanVotes := float64(rPerc) * float64(pollSize) / 100.0
	s.Democrat += democratVotes
	s.Republican += republicanVotes
	s.N += democratVotes + republicanVotes

	s.democratPerc = s.Democrat / s.N
	s.σ = math.Sqrt((s.democratPerc - s.democratPerc*s.democratPerc) / s.N)
	if min_σ != 0.0 && s.σ < min_σ {
		s.σ = min_σ
	}
	s.DemocratProbability = prOverX(0.50, s.democratPerc, s.σ)
}

func (s *StateProbability) simulateElection(r *rand.Rand) int {
	if s.N != 0 {
		if r.Float64() < s.DemocratProbability {
			return college[s.state].votes
		}
	} else {
		// give state to winner of last election
		if college[s.state].lastElection {
			return college[s.state].votes
		}
	}
	return 0
}

func (s *StateProbability) logStateProbability() {
	if s.N != 0 {
		log.Printf("  %v: Democrat polling=%6.4f, N=%d, σ=%6.4f --> Pr(Democrat)=%6.4f\n",
			s.state, s.democratPerc, int(s.N), s.σ, s.DemocratProbability)
	} else {
		if college[s.state].lastElection {
			log.Printf("  %s voted Democratic in the last election. Assuming %d votes for the Democrat.\n", s.state, college[s.state].votes)
		} else {
			log.Printf("  %s voted Republican in the last election. Assuming %d votes for the Republican\n", s.state, college[s.state].votes)
		}
	}
}
