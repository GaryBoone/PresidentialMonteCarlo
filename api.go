//
// api.go
//

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// sort polls by most recent first so we use the most recent data
const apiUrl = "http://elections.huffingtonpost.com/pollster/api/polls.json?sort=updated&topic=%s&state=%s"

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

// readPollingApi reads the data from the Pollster API.
func readPollingApi(topic, state string) []byte {
	url := fmt.Sprintf(apiUrl, topic, state)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func parseJson(body []byte) []Poll {
	var response []Poll
	err := json.Unmarshal(body, &response)
	if err != nil {
		// try unmarshalling into an error map
		var error map[string][]string
		err2 := json.Unmarshal(body, &error)
		if err2 == nil {
			for _, e := range error["errors"] {
				fmt.Printf(" API error: %v\n", e)
			}
		} else {
			fmt.Printf(" JSON parsing error: %v\n", err)
		}
	}
	return response
}
