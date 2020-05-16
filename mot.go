//       mot.go - Gentes

// signets :
// motadeja

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
	"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

type Mot struct {
	gr     string    // graphie du mot
	rang   int       // rang du mot dans la phrase à partir de 0
	ans    gocol.Res // lemmatisations et id du groupe si le mot devient noyau
	restmp gocol.Res // lemmatisation de test d'un noeud
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
	for _, an := range mot.ans {
		if !lexsynt(an.Lem, "excl") {
			nres = append(nres, an)
		}
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

func (m *Mot) lemmatisation(sol Sol) string {
	var ll []string
	for _, nod := range sol.nods {
		if nod.nucl == m {
			for _, l := range restoll(nod.rnucl) {
				ll = append(ll, rouge(l))
			}
		}
		for i, ma := range nod.mma {
			if ma == m {
				for _, l := range restoll(nod.rra[i]) {
					ll = append(ll, rouge(l))
				}
			}
		}
		for i, mp := range nod.mmp {
			if mp == m {
				for _, l := range restoll(nod.rrp[i]) {
					ll = append(ll, rouge(l))
				}
			}
		}
	}
	for _, l := range restoll(m.ans) {
		if !contient(ll, l) {
			ll = append(ll, l)
		}
	}
	return strings.Join(ll, "\n")
}

func restoll(an gocol.Res) []string {
	var lr []string
	for _, rl := range an {
		if rl.Lem == nil {
			continue
		}
		l := fmt.Sprintf("   %s, %s [%s]: %s",
			strings.Join(rl.Lem.Grq, " "),
			rl.Lem.Indmorph,
			rl.Lem.Pos,
			rl.Lem.Traduction)
		lr = append(lr, l)
		for _, m := range rl.Morphos {
			lr = append(lr, "      "+m)
		}
	}
	return lr
}

// retourne une chaîne humainement lisible des
// lemmatisations de an.
func restostring(an gocol.Res) string {
	return strings.Join(restoll(an), "\n")
}
