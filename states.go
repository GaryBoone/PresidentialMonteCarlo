//
// states.go
// 
// 
package main

type State struct {
	votes   int
	dem2008 bool
}

var states = map[string]State{
	"AL": State{9, false},
	"AK": State{3, false},
	"AZ": State{11, false},
	"AR": State{6, false},
	"CA": State{55, true},
	"CO": State{9, true},
	"CT": State{7, true},
	"DE": State{3, true},
	"DC": State{3, true},
	"FL": State{29, true},
	"GA": State{16, false},
	"HI": State{4, true},
	"ID": State{4, false},
	"IL": State{20, true},
	"IN": State{11, true},
	"IA": State{6, true},
	"KS": State{6, false},
	"KY": State{8, false},
	"LA": State{8, false},
	"ME": State{4, true},
	"MD": State{10, true},
	"MA": State{11, true},
	"MI": State{16, true},
	"MN": State{10, true},
	"MS": State{6, false},
	"MO": State{10, false},
	"MT": State{3, false},
	"NE": State{5, false},
	"NV": State{6, true},
	"NH": State{4, true},
	"NJ": State{14, true},
	"NM": State{5, true},
	"NY": State{29, true},
	"NC": State{15, true},
	"ND": State{3, false},
	"OH": State{18, true},
	"OK": State{7, false},
	"OR": State{7, true},
	"PA": State{20, true},
	"RI": State{4, true},
	"SC": State{9, false},
	"SD": State{3, false},
	"TN": State{11, false},
	"TX": State{38, false},
	"UT": State{6, false},
	"VT": State{3, true},
	"VA": State{13, true},
	"WA": State{12, true},
	"WV": State{5, false},
	"WI": State{10, true},
	"WY": State{3, false},
}
