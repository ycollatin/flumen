//        mot.go - Publicola

package main

import (
	"github.com/ycollatin/gocol"
	"strings"
)

// rappel de la lemmatisation dans gocol :
// type Sr struct {
//	Lem     *Lemme
//	Nmorph  []int
//	Morphos []string
// }
//
// type Res []Sr

type Mot struct {
	gr			string
	an			gocol.Sr	// lemmatisation choisie
	ans			gocol.Res	// ensemble des lemmatisations
	lexsynt		[]string
}

func creeMot(m string) *Mot {
	mot := new(Mot)
	mot.gr = m
	var echec bool
	mot.ans, echec = gocol.Lemmatise(m)
	if echec {
		mot.ans, echec = gocol.Lemmatise(gocol.Majminmaj(m))
	}
	if !echec {
		for i, an := range mot.ans {
			mot.ans[i] = genre(an)
		}
	}
	return mot
}

// teste si m peut être le noyau du groupe g
func (m *Mot) estNoyau(g *Groupe) bool {
	for _, an := range m.ans {
		// pos
		if !contient(g.pos, an.Lem.Pos) {
			return false
		}
		// morpho
		var va bool
		for _, morf := range an.Morphos {
			va = true
			for _, gmorf := range g.morph {
				va = va && strings.Contains(morf, gmorf)
			}
		}
		if !va {
			return false
		}
		for _, ls := range(g.lexSynt) {
			va = va && contient(m.lexsynt, ls)
		}
		if !va {
			return false
		}
	}
	return true
}

// vrai si m est compatible avec Sub
// Sub : pos string, morpho []string, accord string
// gocol.Sr : Lem, Morphos []string
func (m *Mot) estSub(sub *Sub) bool {
	var respos, resmorf gocol.Res
	// pos
	for _, an := range m.ans {
		if an.Lem.Pos == an.Lem.Pos {
			respos = append(respos, an)
		}
	}
	// morpho
	for _, an := range respos {
		for _, anmorf := range an.Morphos {
			va := true
			// si le mot est déjà noyau, contrôler 
			// le nom de son groupe
			idgrp := phrase.estNuclDe(m)
			if idgrp > "" {
				va = va && sub.pos == idgrp
			} else {
				for _, trait := range sub.morpho {
					va = va && strings.Contains(anmorf, trait)
				}
			}
			if va {
				resmorf = append(resmorf, an)
			}
		}
	}
	// accord
	// . . .
	return len(resmorf) > 0
}

func genre(sr gocol.Sr) gocol.Sr {
	if sr.Lem.Pos != "n" && sr.Lem.Pos != "NP" {
		return sr
	}
	inc := 12
	switch sr.Lem.Genre {
	case "féminin":
		inc += 12
	case "neutre":
		inc += 24
	}
	for i, _ := range sr.Nmorph {
		sr.Nmorph[i] += inc
	}
	return sr
}
