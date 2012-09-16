//
// election2012.go
//
// $ go run election2012.go
//
// wget -O - 'http://elections.huffingtonpost.com/pollster/api/polls.json?topic=2012-president&state=OH' | underscore print --color
//

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const stateStr = "http://elections.huffingtonpost.com/pollster/api/polls.json?topic=2012-president&state="

// states and their electoral college votes
var states = map[string]int{"AL": 9, "AK": 3, "AZ": 11, "AR": 6, "CA": 55, "CO": 9, "CT": 7,
	"DE": 3, "DC": 3, "FL": 29, "GA": 16, "HI": 4, "ID": 4, "IL": 20, "IN": 11, "IA": 6, "KS": 6,
	"KY": 8, "LA": 8, "ME": 4, "MD": 10, "MA": 11, "MI": 16, "MN": 10, "MS": 6, "MO": 10, "MT": 3,
	"NE": 5, "NV": 6, "NH": 4, "NJ": 14, "NM": 5, "NY": 29, "NC": 15, "ND": 3, "OH": 18, "OK": 7,
	"OR": 7, "PA": 20, "RI": 4, "SC": 9, "SD": 3, "TN": 11, "TX": 38, "UT": 6, "VT": 3, "VA": 13,
	"WA": 12, "WV": 5, "WI": 10, "WY": 3}

type Responses struct {
	Choice     *string
	Value      *int
	First_name *string
	Last_name  *string
	Party      *string
	Incumbent  *bool
}

type Subpopulations struct {
	Name            *string
	Observations    *int
	Margin_of_error *float64
	Responses       []Responses
}

type Question struct {
	Name           *string
	Chart          *string
	Topic          *string
	State          *string
	Subpopulations []Subpopulations
}

type SurveyHouse struct {
	Name *string
}

type Sponsor struct {
	Name *string
}

type Poll struct {
	Id            int
	Pollster      *string
	Start_date    *string
	End_date      *string
	Method        *string
	Source        *string
	Last_updated  *string
	Survey_houses []SurveyHouse
	Sponsors      []Sponsor
	Questions     []Question
}

func readFromWeb(state string) {
	resp, err := http.Get(stateStr + state)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(body))

	// var response interface{}
	var response []Poll
	err = json.Unmarshal(body, &response)
	if err != nil {
		// log.Fatal(err)
		fmt.Printf(" we had an error: %v\n", err)
		// fmt.Println(string(body))
		var error map[string][]string
		err = json.Unmarshal(body, &error)
		if err == nil {
			for _, e := range error["errors"] {
				fmt.Printf(" error: %v\n", e)
			}
		}
	}
	// //fmt.Printf("** resp = %v\n", response)
	// switch t := response.(type) {
	// default:
	// 	fmt.Printf(" ** foo type %T", t)
	// }
	//switch on type?
}

func main() {
	fmt.Println("Election 2012\n")
	for state, _ := range states {
		fmt.Println("Collecting survey data for the great state of " + state)
		readFromWeb(state)
	}
}
