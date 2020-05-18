//       mot.go - Gentes

// signets :
// motadeja
// funcLemmatisation

// rappel de la lemmatisation dans gocol :
// type Sr struct {
//	Lem     *Lemme
//	Nmorph  []int
//	Morphos []string
// }
//
// type Res []Sr

package main

import (
	"github.com/ycollatin/gocol"
	"strings"
)

type Mot struct {
	gr     string    // graphie du mot
	rang   int       // rang du mot dans la phrase à partir de 0
	ans    gocol.Res // lemmatisations et id du groupe si le mot devient noyau
	restmp gocol.Res // lemmatisation de test d'un noeud
	interj bool		 // tous les lemmes du mot sont des interjections
}

func creeMot(m string) *Mot {
	mot := new(Mot)
	mot.gr = m
	var echec bool
	mot.ans, echec = gocol.Lemmatise(m)
	if echec {
		mot.ans, echec = gocol.Lemmatise(gocol.Majminmaj(m))
	}
	// ajout du genre pour les noms
	if !echec {
		for i, a := range mot.ans {
			mot.ans[i] = genus(a)
		}
	}

	// provisoire XXX
	// exclusions de mots rares faisant obstacle à des analyses importantes
	var nres gocol.Res
	mot.interj = true
	for _, an := range mot.ans {
		if !lexsynt(an.Lem, "excl") {
			nres = append(nres, an)
		}
		mot.interj = mot.interj && an.Lem.Pos == "intj"
	}
	mot.ans = nres
	return mot
}

// vrai si les lemmatisation lma et lmb sont accordées
// en cas (c) genre (g) et nombre (n)
func accord(lma, lmb, cgn string) bool {
	if strings.Contains(lmb, "inv.") {
		return false
	}
	va := true
	for i := 0; i < len(cgn); i++ {
		switch cgn[i] {
		case 'c':
			k := cas(lma)
			va = va && strings.Contains(lmb, k)
		case 'g':
			g := genre(lma)
			va = va && strings.Contains(lmb, g)
		case 'n':
			n := nombre(lma)
			va = va && strings.Contains(lmb, n)
		}
	}
	return va
}

// ajoute le genre à la morpho d'un nom
func genus(sr gocol.Sr) gocol.Sr {
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
