//
// election2012.go
//
// An electoral college Monte Carlo simulation based on 2012 presidential polling.
//
// To run:
//   $ go run election2012.go state.go api.go cdf.go parse.go college.go
//
// Author:     Gary Boone     gary.boone@gmail.com
// History:    2012-09-17     • initial version
//             2012-09-21     • cleanup, upload to github
// Notes:
//
// The state-by-state presidential polling data is provided by the Pollster API:
//   http://elections.huffingtonpost.com/pollster/api
// 
//   Example API call:
//   wget -O - 'http://elections.huffingtonpost.com/pollster/api/polls.json?topic=2012-president&state=OH'
//
// Read the logfile for details.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	acceptableSize int
	numSimulations int
	min_σ          float64
)

func init() {
	const (
		acceptableSizeDefault = 2000
		numSimulationsDefault = 25000
		min_σDefault          = 0.0
	)
	flag.IntVar(&acceptableSize, "acceptableSize", acceptableSizeDefault,
		"Don't add more polls after this many samples are obtained")
	flag.IntVar(&numSimulations, "sims", numSimulationsDefault, "number of simulations to run")
	flag.Float64Var(&min_σ, "minStdDev", min_σDefault, "minimum standard deviation")
}

func initializeLog() {
	f, err := os.Create("logfile")
	if err != nil {
		fmt.Printf("can't open file: %v\n", err)
	}
	log.SetOutput(f)
	log.Println("Simulation Parameters:")
	log.Printf("  Stop adding new polls when we have more than %v samples. "+
		"(param: -acceptableSize x)\n", acceptableSize)
	log.Printf("  Run %v simulations. (param: -numSimulations x)\n", numSimulations)
	log.Printf("  Don't allow the standard deviation to shrink below %v. "+
		"(0=no limit, param: -minStdDev x)\n", min_σ)
}

func truncateString(inStr string, length int) string {
	if len(inStr) < length || length < 4 {
		return inStr
	}
	return inStr[:length-3] + "..."
}

func loadStateData(state string, polls []Poll) (prob StateProbability) {
	prob.state = state
	for _, poll := range polls {

		pollster := parsePollster(poll)
		date := parseDateAsString(poll)

		// skip systemically biased polls
		// http://fivethirtyeight.blogs.nytimes.com/2010/11/04/rasmussen-polls-were-biased-and-inaccurate-quinnipiac-surveyusa-performed-strongly/
		if strings.EqualFold(pollster, "Rasmussen") {
			continue
		}

		var obama, romney, size int
		obama, romney, size = parsePoll(state, poll)
		if obama == 0 || romney == 0 {
			log.Printf("  Missing value (Obama=%v, Romney=%v) for %v state poll by '%v'. Skipping.\n",
				obama, romney, state, *poll.Pollster)
			continue
		}

		log.Printf("  adding %-30s %10s : O(%v), R(%v), N(%v)\n",
			truncateString(pollster, 30), date[:10], obama, romney, size)
		prob.update(obama, romney, size)
		if prob.N > float64(acceptableSize) {
			return
		}
	}
	return
}

// for each state, flip a coin
func simulateObamaVotes(stateProbalities []StateProbability) int {
	votes := 0
	for _, prob := range stateProbalities {
		votes += prob.simulateElection()
	}
	return votes
}

func initializeSimulations() []StateProbability {
	var stateProbalities []StateProbability
	i := 0
	for state, _ := range college {
		if i == 0 {
			fmt.Printf("Collecting survey data for the great state of %v", state)
		} else {
			fmt.Printf(", %v", state)
		}
		i++
		body := readPollingApi(state)
		polls := parseJson(body)
		log.Printf("Found %v polls in %v.\n", len(polls), state)
		prob := loadStateData(state, polls)
		prob.logStateProbability()
		stateProbalities = append(stateProbalities, prob)
	}
	fmt.Printf(".\n")
	return stateProbalities
}

func runSimulations(stateProbalities []StateProbability) (int, int) {
	wins := 0
	totalVotes := 0
	for i := 0; i < numSimulations; i++ {
		votes := simulateObamaVotes(stateProbalities)
		if votes >= 270 {
			wins++
		}
		totalVotes += votes
	}
	return wins, totalVotes
}

// Let's say election day begins on midnight Eastern Time on Nov 6, 2012
func daysUntilElection() int {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println("*** problem loading timezone information. Using UTC. ***")
		location = time.UTC
	}
	now := time.Now().In(location)
	electionDay := time.Date(2012, time.November, 6, 0, 0, 0, 0, location).In(location)
	return int(electionDay.Sub(now) / (24 * 60 * 60 * 1000000000))
}

func main() {
	flag.Parse()
	initializeLog()

	fmt.Println("Election 2012 Monte Carlo Simulation")
	fmt.Printf("There are %v days until the election.\n\n", daysUntilElection())

	stateProbalities := initializeSimulations()
	wins, totalVotes := runSimulations(stateProbalities)

	fmt.Printf("\nObama win probability: %.2f%%. Average votes: %.2f\n", 100.0*float64(wins)/float64(numSimulations),
		float64(totalVotes)/float64(numSimulations))

}
