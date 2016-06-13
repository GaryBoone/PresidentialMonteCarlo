//
// election.go
//
// An electoral college Monte Carlo simulation based on 2016 presidential polling.
//
// To build and run:
//     $ cd $GOPATH
//     $ mkdir -p src/github.com/GaryBoone/
//     $ cd src/github.com/GaryBoone/
//     $ git clone https://github.com/GaryBoone/PresidentialMonteCarlo
//     $ cd ../../..
//     $ go install github.com/GaryBoone/PresidentialMonteCarlo
//     $ bin/PresidentialMonteCarlo//
// Author:     Gary Boone     gary.boone@gmail.com
// History:
//             2012-09-25     • simulations in parallel
//             2012-09-24     • minimum σ
//                            • command line parameters
//                            • days until election countdown
//             2012-09-21     • cleanup, upload to github
//             2012-09-17     • initial version
// Notes:
//
// The state-by-state presidential polling data is provided by the Pollster API:
//   http://elections.huffingtonpost.com/pollster/api
// 
//   Example API call:
//   wget -O - 'http://elections.huffingtonpost.com/pollster/api/polls.json?topic=2016-president&state=OH'
//
// Read the logfile for details.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	democraticCandidate = "Clinton"// "Obama"
	republicanCandidate = "Trump"  // "Romney"
	electionYear = 2016 // 2012
	electionDay = 8 // 6
	swingStates = "CO,FL,IA,NC,NH,NV,OH,PA,VA,WI"
)

var (
	acceptableSize int
	numSimulations int
	min_σ          float64
	pollTopic = fmt.Sprintf("%d-president", electionYear)
	loc, _ = time.LoadLocation("America/New_York")
	electionDate = time.Date(electionYear, time.November, electionDay, 0, 0, 0, 0, loc)
)

func init() {
	const (
		acceptableSizeDefault = 2000
		numSimulationsDefault = 25000
		min_σDefault          = 0.0 // 0.0 => no minimum
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

		var democrat, republican, size int
		democrat, republican, size = parsePoll(state, poll, pollTopic)
		if democrat == 0 || republican == 0 {
			log.Printf("  Missing value (Democrat=%v, Republican=%v) for %v state poll by '%v'. Skipping.\n",
				democrat, republican, state, *poll.Pollster)
			continue
		}

		log.Printf("  adding %-30s %10s : Democrat(%v), Republican(%v), poll size(%v)\n",
			truncateString(pollster, 30), date[:10], democrat, republican, size)
		prob.update(democrat, republican, size)
		if prob.N > float64(acceptableSize) {
			return
		}
	}
	return
}

// for each state, flip a coin
func simulateObamaVotes(states []StateProbability, r *rand.Rand) int {
	votes := 0
	for _, state := range states {
		votes += state.simulateElection(r)
	}
	return votes
}

func loadProbability(state string) StateProbability {
	body := readPollingApi(pollTopic, state)
	polls := parseJson(body)

	msg := ""
	if strings.Contains(swingStates, state) {
		msg = ", a swing state"
	}	
	log.Printf("Found %v polls in %v%s.\n", len(polls), state, msg)

	prob := loadStateData(state, polls)
	prob.logStateProbability()
	return prob
}

func initializeSimulations() []StateProbability {
	results := make(chan StateProbability)
	// kick off all the polls
	for state, _ := range college {
		go func(state string) {
			results <- loadProbability(state)
		}(state)
	}

	stateProbabilities := make([]StateProbability, len(college))
	for i := range stateProbabilities {
		prob := <-results
		stateProbabilities[i] = prob
		if i == 0 {
			fmt.Printf("Collecting survey data for the great states of %v", prob.state)
		} else {
			fmt.Printf(", %v", prob.state)
		}
	}
	fmt.Printf(".\n")
	return stateProbabilities
}

type Result struct {
	votes, wins int
}

func doSome(n int, probs []StateProbability, c chan Result) {
	var voteSum, winSum int
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		votes := simulateObamaVotes(probs, r)
		if votes >= 270 {
			winSum++
		}
		voteSum += votes
	}
	c <- Result{voteSum, winSum}
}

func runSimulations(probs []StateProbability) (int, int) {
	nCPU := runtime.NumCPU()
	log.Printf("Using %v CPUs.\n", nCPU)
	runtime.GOMAXPROCS(nCPU)
	c := make(chan Result, nCPU)
	for i := 0; i < nCPU; i++ {
		go doSome(numSimulations/nCPU, probs, c)
	}

	var wins, votes int
	for i := 0; i < nCPU; i++ {
		res := <-c
		votes += res.votes
		wins += res.wins
	}
	return wins, votes
}

func daysUntilElection() int {
	now := time.Now()
	return int(math.Ceil(float64(electionDate.Sub(now)) / (24 * 60 * 60 * 1000000000.0)))
}

func reportProbalities(probs []StateProbability) {
	numStatesWithoutPolls := 0
	fmt.Println("\nSwing States:")
	for _, st := range probs {
		if st.N==0 {
			numStatesWithoutPolls++
			if strings.Contains(swingStates, st.state) {
				fmt.Printf("%s has no polls yet, so it is assigned to %s based on %d outcome.\n", st.state, democraticCandidate, electionYear-4)
			}
		} else {
			if strings.Contains(swingStates, st.state) {
				fmt.Printf("Probability of %s winning %v: %4.2f%%\n", democraticCandidate, st.state, 100.0*st.DemocratProbability)
			}
		}
	}
	fmt.Printf("%d states have no polls, so were assigned %d outcomes\n", numStatesWithoutPolls, electionYear-4)
}

func main() {
	flag.Parse()
	initializeLog()

	fmt.Println("Election %s Monte Carlo Simulation", electionYear)
	fmt.Printf("There are %v days until the election.\n\n", daysUntilElection())

	stateProbalities := initializeSimulations()
	reportProbalities(stateProbalities)
	
	wins, totalVotes := runSimulations(stateProbalities)

	demWinProb := 100.0*float64(wins)/float64(numSimulations)
	fmt.Printf("\n%s election probability: %.2f%%\n", democraticCandidate, demWinProb)
	fmt.Printf("%s election probability: %.2f%%\n", republicanCandidate, 100.0 - demWinProb)
	avgVotes := float64(totalVotes) / float64(numSimulations)
	roundedVotes := int(math.Floor(avgVotes + 0.5))
	fmt.Printf("Average electoral votes for %s: %v\n", democraticCandidate, roundedVotes)
	fmt.Printf("Average electoral votes for %s: %v\n", republicanCandidate, 538 - roundedVotes)
}
