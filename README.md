PresidentialMonteCarlo is a Monte Carlo simulator of the 2016 presidential election written in Go. 

* It reads polling data from the Huffington Post API.
* It transforms the polling data into probabilities.
* It assumes that where there's no polling data, the state will vote as it did in 2012.
* It runs 25,000 simulations to determine the most likely electoral college total for Clinton.

## Update ##

This code was originally written for the 2012 Obama/Romney election. It has now been updated for the 2016 presidential election. It assumes Clinton and Trump are the Democratic and Republican candidates, although that isn't official as of this writing in early June. No matter. Clinton will win this election as surely as the sun will rise on November 9th, 2016. This simulation shows why.

## Usage ##

To build:
	
	$ go install github.com/GaryBoone/PresidentialMonteCarlo

To run:

	$ bin/PresidentialMonteCarlo

Try:
	$ bin/PresidentialMonteCarlo -minStdDev 0.05
  
    $ bin/PresidentialMonteCarlo -minStdDev 0.07 -acceptableSize 4000 -sims 50000
	Election 2016 Monte Carlo Simulation
	Run date: Monday Jun 13 2016 12:42:42
	There are 148 days until the election.

	Collecting survey data for the great states of NV, AL, LA, NE, MT, WY, CO, NY, ME, MA, IL, VT, SC, AR, AZ, KY, NJ, SD, OR, UT, AK, MN, NM, MI, MS, KS, ID, MD, FL, IN, WA, HI, GA, CT, DE, TX, OK, WV, DC, ND, TN, RI, MO, CA, VA, PA, NC, NH, OH, IA, WI.

	Swing States:
	NV has no polls yet, so it is assigned to Clinton based on 2012 outcome.
	CO has no polls yet, so it is assigned to Clinton based on 2012 outcome.
	Probability of Clinton winning FL: 58.24%
	Probability of Clinton winning VA: 73.41%
	Probability of Clinton winning PA: 75.68%
	Probability of Clinton winning NC: 52.53%
	Probability of Clinton winning NH: 65.12%
	Probability of Clinton winning OH: 61.48%
	Probability of Clinton winning IA: 59.49%
	Probability of Clinton winning WI: 84.31%
	35 states have no polls, so were assigned 2012 outcomes

	Clinton election probability: 92.12%
	Trump election probability: 7.88%
	Average electoral votes for Clinton: 311
	Average electoral votes for Trump: 227

For additional details about the data that was gathered and the simulation, see the the logfile.

	$ less logfile


## What This Simulation Shows ##

Why does the Monte Carlo simulation make the election appear more certain for Clinton than the national polls that show a close race?

First, you need to understand that the national polling reported in the media is irrelevant. 47% Clinton, 46% Trump? Irrelevant. Why? Because popular polls don't elect the president; the electoral college does. To predict the winner, you need to use the polls to simulate the various combinations of electoral college wins/losses for each candidate. This simulation uses 25,000 simulated elections. Most of the time, Clinton wins.

Second, you need to understand that if the polls are 52% Clinton in a state with a margin of error of 4%, then Clinton is 69% likely to win that state. You might think that 52% Clinton in the polling means that he's 52% percent likely to win the state. But that's not how the statistics work. If you have a coin that comes out heads 52% of the time with a 4% margin of error, then that means that the probability of exceeding 50% is about 70%. To see this, go to a [cumulative distribution function calculator](http://www.danielsoper.com/statcalc3/calc.aspx?id=53) and plug in 52, 5, 50 and then calculate 1 - the answer given. In Clinton's case, exceeding 50% is a win for that state. It turns out to have a probability of ~70%. 

Finally, realize that there is more certainty available in the data than is described in the polls. That's because 1) the electoral college elects presidents and 2) there are only so many ways that the electoral college can add up to a win for either candidate. For example, CA, HI, and NY will almost certainly vote for Clinton; no need to poll in those states. That's why we talk about swing states; they can change the outcome. But there are only a limited number of ways the swing states can combine to a victory for one candidate or the other. As it turns out, the current polling makes these combinations favor Clinton. It's nearly impossible for swing states to combine in ways that add up to a Trump win.

## Notes ##

* PresidentialMonteCarlo ignores Rasmussen polls (http://fivethirtyeight.blogs.nytimes.com/2010/11/04/rasmussen-polls-were-biased-and-inaccurate-quinnipiac-surveyusa-performed-strongly/)


## Poll Data Source ##

The state-by-state presidential polling data is provided by the [Pollster API](http://elections.huffingtonpost.com/pollster/api).



## Sources of Error ##

* Polling data is combined by pooling multiple samples into a single large sample.
* Undecides and Others are ignored. For Others, that's not a bad assumption, as the contest remains a two-person contest among the two leaders. For Undecideds, it implies that they'll allocate themselves according to the current proportions for each candidate. But historically, late deciders favor challengers.
* Only state-by-state polling data is considered. Professional modelers, such as the Nate Silver at the excellent [FiveThirtyEight blog](http://fivethirtyeight.blogs.nytimes.com/) consider many other factors, adjustments, and trendlines as well. See [FiveThirtyEight's methodology](http://fivethirtyeight.blogs.nytimes.com/methodology/). 
* On FiveThirtyEight, click on the tab titled "Presidential Now-cast" to see Nate's results if the elction were held right now. Those results correspond to the results generated by this simulation. 


## License ##
The code is available at github [GaryBoone/PresidentialMonteCarlo](https://github.com/GaryBoone/PresidentialMonteCarlo) under [MIT license](http://opensource.org/licenses/mit-license.php).
