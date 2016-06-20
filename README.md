PresidentialMonteCarlo is a Monte Carlo simulator of the 2016 presidential election written in Go. 

This code uses a technique called Monte Carlo analysis which is a method for calculating probabilities when no simpler, closed-form, formula is known. It allows us to calculate probabilities for complex systems such as presidential elections. 

Specifically, the process is:

* It reads polling data from the Huffington Post API for each state.
* It transforms the polling data into probabilities.
* It assumes that where there's no polling data, the state will vote as it did in 2012.
* It uses the probabilities to simulate an election for a state, doing so for all 50 states and DC to create a trial election. The electoral college votes are counted to determine the winning candidate.
* It runs 25,000 simulations or more, counting how many times each candidate wins. 

The percentage of these wins is then the final probability of the candidate winning. So if Clinton wins 92% of 50,000 trial elections, we conclude she has a 92% chance of winning in November.

What Monte Carlo simulation does is allow you to find certainty (eg, 92% certainty) from many uncertain components, such as close state polls. 

It doesn't have to be this way. It could be that the overall simulation of combining uncertain components also leads to uncertain outcomes. But what's so interesting in the case of the presidential election is that the simulation shows that the outcome of the 2016 Presidential is certain. Hillary Clinton will win.
 
## Update ##

This code was originally written for the 2012 Obama/Romney election. It predicted a 99.86% probability of Obama winning. It predicted 314 Electoral College votes for Obama. He actually received 332. The simulation was wrong about Indiana (actually 54.3% Romney / 43.8% Obama, 11 votes) and Florida (actually 50.0% Obama / 49.1% Romney, 29 votes).

It has now been updated for the 2016 presidential election. It assumes Clinton and Trump are the Democratic and Republican candidates, although that isn't official as of this writing in mid June.

## Usage ##

Assuming a standard Go development directory layout, to build:
	
	$ go install github.com/GaryBoone/PresidentialMonteCarlo

To run:

	$ bin/PresidentialMonteCarlo

Try:

	$ bin/PresidentialMonteCarlo -minStdDev 0.07

or:
  
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

First, you need to understand that the _national_ polling reported in the media is irrelevant. 47% Clinton, 46% Trump? Irrelevant. Why? Because popular polls don't elect the president; the electoral college does based on _state_ elections. To predict the winner, you need to use the polls to simulate the various combinations of electoral college wins/losses for each candidate. This simulation uses 25,000 simulated elections by default. Most of the time, Clinton wins.

Second, if the polls are 52% Clinton in a state with a margin of error of 4%, then Clinton is 69% likely to win that state. You might think that 52% Clinton in the polling means that she's 52% percent likely to win the state. But that's not how the statistics work. If you estimate an outcome to occur 52% of the time with a 4% margin of error, then that means that the probability of that outcome of exceeding 50% is about 70%. To see this, go to a [cumulative distribution function calculator](http://www.danielsoper.com/statcalc3/calc.aspx?id=53) and plug in 52, 5, 50 and then calculate 1 - the answer given. In Clinton's case, exceeding 50% is a win for that state. It turns out to have a probability of ~70%. 

Third, while any given poll may have a margin of error of say 4%, we can combine polls and reduce the error. Polling is expensive, so polling companies only call people until they have a representative sample with the desired margin of error. But we can increase certainty for our simulations by aggregating polls. 

Finally, realize that there is more certainty available in the data than is described in the polls. That's because there are only so many ways that the electoral college can add up to a win for either candidate. For example, CA, HI, and NY will almost certainly vote for Clinton; no need to poll in those states. That's why we talk about swing states; they can change the outcome. But there are only a limited number of ways the swing states can combine to a victory for one candidate or the other. As it turns out, the current polling makes these combinations favor Clinton. It's nearly impossible for swing states to combine in ways that add up to a Trump win.


## Poll Data Source ##

The state-by-state presidential polling data is provided by the [Pollster API](http://elections.huffingtonpost.com/pollster/api).


## Notes and Sources of Error ##

* If a state has no recent polls, the simulation assumes it will vote the  way it did in 2012. That's the best choice because the reason these states aren't being polled is that they're so likely to vote as they did in 2012.
* Polls are combined most-recent-first to meet the requested aggregation size.
* The _-minStdDev_ parameter allows you to force uncertainty into the simulation. By default, there's no minimum standard deviation. The aggregation of polls reduces their uncertainty as the statistics dictate. That, combined with the other factors above, causes the simulation to report Clinton winning 100% of the time. To see more of how the states can vary add _-minStdDev_ to the run.
* PresidentialMonteCarlo ignores pollsters that are marked banned in [FiveThirtyEight's Pollster Ratings](http://projects.fivethirtyeight.com/pollster-ratings/)
* The simulation uses only polls of 'likely' voters.
* _Undecideds_ and _Others_ are ignored. For _Others_, that's not a bad assumption, as the contest remains a two-person contest among the two leaders. For _Undecideds_, it implies that they'll allocate themselves according to the current proportions for each candidate. But historically, late deciders favor challengers.
* Only state-by-state polling data is considered. Professional modelers, such as the Nate Silver at the excellent [FiveThirtyEight](http://fivethirtyeight.com/politics/) consider many other factors, adjustments, and trendlines as well. See [FiveThirtyEight's methodology](http://fivethirtyeight.blogs.nytimes.com/methodology/). 


## License ##
The code is available at github [GaryBoone/PresidentialMonteCarlo](https://github.com/GaryBoone/PresidentialMonteCarlo) under [MIT license](http://opensource.org/licenses/mit-license.php).
