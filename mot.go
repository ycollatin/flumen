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
			mot.ans[i] = genus(an)
		}
	}
	return mot
}

func (ma *Mot) accord(mb *Mot, cgn string) bool {
	for _, sra := range ma.ans {
		for _, srb := range mb.ans {
			for _, morfa := range sra.Morphos {
				for _, morfb := range srb.Morphos {
					va := true
					for i:=0; i<len(cgn); i++ {
						switch cgn[i] {
						case 'c':
							k := cas(morfa)
							va = va && strings.Contains(morfb, k)
						case 'g':
							g := genre(morfa)
							va = va && strings.Contains(morfb, g)
						case 'n':
							n := nombre(morfa)
							va = va && strings.Contains(morfb, n)
						}
					}
					if va {
						return true
					}
				}
			}
		}
	}
	return false
}

func (m *Mot) dejaNoy() bool {
	for _, n := range phrase.nods {
		if n.nucl == m {
			return true
		}
	}
	return false
}

func (m *Mot) dejaSub() bool {
	for _, n := range phrase.nods {
		if m.elDe(n) {
			return true
		}
	}
	return false
}

func (m *Mot) elDe(n *Nod) bool {
	for _, ma := range n.mma {
		if ma == m {
			return true
		}
	}
	for _, mp := range n.mmp {
		if mp == m {
			return true
		}
	}
	return false
}

// teste si m peut être le noyau du groupe g
func (m *Mot) estNoyau(g *Groupe) bool {
	/*
	grp:P.sv
	pos:GV.objprep
	morph:indic 3;subj 3
	ag:GN;sujet;nomin;n
	*/
	for _, an := range m.ans {
		// pos
		if m.dejaNoy() {
			idgrp := phrase.estNuclDe(m)
			if !contient(g.pos, idgrp) {
				continue
			}
		} else if !contient(g.pos, an.Lem.Pos) {
			continue
		}
		// morpho
		var va bool
		for _, morf := range an.Morphos {
			va = true
			for _, gmorf := range g.morph {
				ecl := strings.Split(gmorf, " ")
				for _, e := range ecl {
					va = va && strings.Contains(morf, e)
				}
				if va {
					break
				}
			}
		}
		if !va {
			continue
		}
		for _, ls := range(g.lexSynt) {
			va = va && contient(m.lexsynt, ls)
		}
		if !va {
			continue
		}
		return true
	}
	return false
}

// vrai si m est compatible avec Sub et le noyau mn
// Sub : pos string, morpho []string, accord string
// gocol.Sr : Lem, Morphos []string
func (m *Mot) estSub(sub *Sub, mn *Mot) bool {
	var respos, resmorf gocol.Res
	// pos
	if sub.terminal {
		for _, an := range m.ans {
			if an.Lem.Pos == an.Lem.Pos {
				respos = append(respos, an)
			}
		}
	} else {
		if m.noyId(sub.pos) {
			respos = m.ans
		} else {
			return false
		}
	}
	// morpho
	for _, an := range respos {
		for _, anmorf := range an.Morphos {
			va := true
			// si le mot est déjà noyau, contrôler 
			// le nom de son groupe, (seul l'accord compte ?)
			idgrp := phrase.estNuclDe(m)
			if idgrp > "" {
				va = va && sub.pos == idgrp
			} else {
				// sinon, on vérifie la morpho du mot
				for _, trait := range sub.morpho {
					va = va && strings.Contains(anmorf, trait)
				}
			}
			if va {
				resmorf = append(resmorf, an)
			} else {
				continue
			}
			// accord
			if sub.accord != "" && !mn.accord(m, sub.accord) {
				continue
			}
			if len(resmorf) > 0 {
				return true
			}
		}
	}
	return false
}

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

// le mot m est il noyau d'un groupe d'id id ?
func (m *Mot) noyId(id string) bool {
	for _, n := range phrase.nods {
		if n.grp.id == id && n.nucl == m {
			return true
		}
	}
	return false
}

func (ma *Mot) subDe(mb *Mot) bool {
	// chercher le groupe dont mb est noyau
	for _, n := range phrase.nods {
		if mb == n.nucl {
			return ma.elDe(n)
		}
	}
	return false
}
