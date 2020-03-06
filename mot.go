//        mot.go - Publicola

package main

import (
	"fmt"
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

func (m *Mot) estNoyau(g *Groupe) bool {
	debog := m.gr == "Prometheus" && g.id == "GN.0"
	if debog {
		fmt.Println(m.gr, "estNoyau",g.id)
	}
	for _, an := range m.ans {
		if debog {
			fmt.Println(" an.Lem",an.Lem.Gr, an.Lem.Pos,"g.pos",g.pos)
		}
		// pos
		if !contient(g.pos, an.Lem.Pos) {
			return false
		}
		if debog {
			fmt.Println("   OKa")
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
		if debog {
			fmt.Println("   OKb")
		}
		for _, ls := range(g.lexSynt) {
			va = va && contient(m.lexsynt, ls)
		}
		if !va {
			return false
		}
		if debog {
			fmt.Println("   OKc")
		}
	}
	return true
}

func (m *Mot) estSub(sub *Sub) bool {
	for _, an := range m.ans {
		// pos
		if sub.pos != an.Lem.Pos {
			return false
		}
		// morpho
		var va bool
		for _, morf := range an.Morphos {
			va = true
			for _, gmorf := range sub.morpho {
				va = va && strings.Contains(morf, gmorf)
			}
		}
		if !va {
			return false
		}
		/*
		for _, ls := range(g.lexSynt) {
			va = va && contient(m.lexsynt, ls)
		}
		if !va {
			return false
		}
		*/
	}
	return true
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
