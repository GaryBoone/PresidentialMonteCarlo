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
	"strings"
)

func main() {
	fmt.Println("Election 2012\n")
	for state, _ := range states {
		fmt.Printf("Collecting survey data for the great state of %v\n", state)
		body := readPollingApi(state)
		polls := parseJson(body)
		fmt.Printf("  Found %v polls.\n", len(polls))

		for _, poll := range polls {
			fmt.Printf(" id: %v\n", poll.Id)
			fmt.Printf(" end date: %v\n", poll.End_date)
			for _, question := range poll.Questions {
				if question.Topic != nil && strings.EqualFold(*question.Topic, "2012-president") {
					fmt.Printf(" topic: %v\n", *question.Topic)
				}
			}
		}
	}
}
