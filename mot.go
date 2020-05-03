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

type An struct {
	res  gocol.Res // lemmatisations réduites du mot
	idGr string    // nom du groupe dont le mot est noyau
}

type Mot struct {
	gr     string    // graphie du mot
	rang   int       // rang du mot dans la phrase à partir de 0
	ans		An		 // lemmatisations et id du groupe si le mot devient noyau
	restmp gocol.Res // lemmatisation de test d'un noeud
}

func creeMot(m string) *Mot {
	mot := new(Mot)
	mot.gr = m
	var echec bool
	mot.ans.res, echec = gocol.Lemmatise(m)
	if echec {
		mot.ans.res, echec = gocol.Lemmatise(gocol.Majminmaj(m))
	}
	// ajout du genre pour les noms
	if !echec {
		for i, a := range mot.ans.res {
			mot.ans.res[i] = genus(a)
		}
	}

	// provisoire XXX
	// exclusions de mots rares faisant obstacle à des analyses importantes
	var nres gocol.Res
	for _, an := range mot.ans.res {
		if !lexsynt(an.Lem, "excl") {
			nres = append(nres, an)
		}
	}
	mot.ans.res = nres
	return mot
}

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

func restostring(an An) string {
	var lr []string
	for _, rl := range an.res {
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
	return strings.Join(lr, "\n")
}
