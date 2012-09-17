//
// election2012.go
//
// $ go run election2012.go api.go states.go
//
// wget -O - 'http://elections.huffingtonpost.com/pollster/api/polls.json?topic=2012-president&state=OH' | underscore print --color
//

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Sim struct {
	state  string
	Obama  int
	Romney int
	N      int
}

// Update sim data with a new poll. The new N is calculated with the actual 
// number of votes for Obama and Romney, not the N of the poll. The effect 
// is to not count undecideds and Others. Essentially, the poll is reduced 
// to a new poll between the two potential winners. Because the N is reduced, 
// the uncertainty is increased as it should be.
func (s *Sim) update(obamaPerc, romneyPerc, pollSize int) {
	obamaVotes := int(obamaPerc * pollSize / 100.0)
	romneyVotes := int(romneyPerc * pollSize / 100.0)
	s.Obama += obamaVotes
	s.Romney += romneyVotes
	s.N += obamaVotes + romneyVotes
}

func parseResponses(state string, poll Poll, responses []Responses) (obama, romney int) {
	for _, resp := range responses {
		if resp.Choice == nil {
			log.Printf("No Choice for %v state poll by '%v'. Skipping.\n",
				state, *poll.Pollster)
			continue
		}
		if resp.Value == nil {
			log.Printf("No Value for %v state poll by '%v'. Skipping.\n",
				state, *poll.Pollster)
			continue
		}
		if strings.EqualFold(*resp.Choice, "obama") {
			obama = *resp.Value
		}
		if strings.EqualFold(*resp.Choice, "romney") {
			romney = *resp.Value
		}
	}
	return
}

func parseSubpopulation(state string, poll Poll, sub Subpopulations) (obama, romney, size int) {
	if sub.Observations == nil {
		log.Printf("No N for %v state poll by '%v'. Skipping.\n",
			state, *poll.Pollster)
		return
	}

	size = *sub.Observations
	obama, romney = parseResponses(state, poll, sub.Responses)
	return
}

func parsePoll(state string, poll Poll) (obama, romney, size int) {
	for _, question := range poll.Questions {
		if question.Topic != nil && strings.EqualFold(*question.Topic, "2012-president") {
			// given multiple subpopulations, prefer likely voters
			switch len(question.Subpopulations) {
			case 1:
				obama, romney, size = parseSubpopulation(state, poll, question.Subpopulations[0])
			default:
				foundLikelyVoters := false
				for _, sub := range question.Subpopulations {
					if sub.Name != nil && strings.EqualFold(*sub.Name, "Likely Voters") {
						obama, romney, size = parseSubpopulation(state, poll, sub)
						foundLikelyVoters = true
					}
				}
				if !foundLikelyVoters {
					log.Printf("No Likely voters in multi-subpopulation poll for "+
						"%v state poll by '%v'. Skipping.\n", state, *poll.Pollster)
				}
			}
		}
	}
	return
}

func loadSimulations(state string, polls []Poll) Sim {
	sim := Sim{state, 0, 0, 0}
	for _, poll := range polls {

		// skip systemically biased polls
		if poll.Pollster != nil && strings.EqualFold(*poll.Pollster, "Rasmussen") {
			continue
		}

		// TODO: use date to weight down older polls
		// pollDate, err := time.Parse("2006-01-02", poll.End_date) // will assume UTC
		// if err != nil {
		// 	fmt.Printf(" parsing error for date '%v'. Error: %v\n", poll.End_date, err)
		// 	continue
		// }

		var obama, romney, size int
		obama, romney, size = parsePoll(state, poll)
		if obama == 0 || romney == 0 {
			log.Printf("Missing value (Obama=%v, Romney=%v) for %v state poll by '%v'. Skipping.\n",
				obama, romney, state, *poll.Pollster)
			continue
		}
		// fmt.Printf(" adding O(%v), R(%v), N(%v)\n", obama, romney, size)
		sim.update(obama, romney, size)
	}
	return sim
}

func doSimulation(simulations []Sim) {

}

func main() {
	f, err := os.Create("logfile")
	if err != nil {
		fmt.Printf("can't open file: %v\n", err)
	}
	log.SetOutput(f)

	fmt.Println("Election 2012\n")
	var simulations []Sim
	for state, _ := range states {
		fmt.Printf("Collecting survey data for the great state of %v\n", state)
		body := readPollingApi(state)
		polls := parseJson(body)
		log.Printf("  Found %v polls.\n", len(polls))
		sim := loadSimulations(state, polls)
		simulations = append(simulations, sim)
	}

	fmt.Printf("num simulations: %v\n", len(simulations))

	doSimulation(simulations)

}
