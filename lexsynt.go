//      lexsynt.go -- gentes
// 		samedi 21 mars 2020 

package main

import (
	"github.com/ycollatin/gocol"
	"strings"
)

var llexs map[string][]string

func lisLexsynt() {
	llexs = make(map[string][]string)
	ll := gocol.Lignes(chData+"lexsynt.la")
	for _, l := range ll {
		ecl := strings.Split(l, ":")
		ecls := strings.Split(ecl[1], ",")
		llexs[ecl[0]] = ecls
	}
}

func lexsynt(lem, el string) bool {
	ls := llexs[lem]
	if ls == nil {
		return false
	}
	return contient(ls, el)
}
