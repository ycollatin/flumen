//        mot.go - Publicola

package main

import (
	//"fmt"
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
	rang		int
	an			int			// n° de lemmatisation choisie
	ans			gocol.Res	// ensemble des lemmatisations
	lexsynt		[]string
	sub			*Sub		// sub qui lie le Mot à son noyau
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
	if len(mot.ans) > 1 {
		mot.an = -1
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
	//debog := m.gr=="filius" && g.id=="n.fam"
	//if debog {fmt.Println("   ",m.gr,"estNoyau",g.id)}
	for _, an := range m.ans {
		//if debog {fmt.Println("    oka",an.Lem.Gr, an.Morphos,"dejaNoy",m.dejaNoy())}
		//if debog {fmt.Println("    contient",g.pos, an.Lem.Pos,contient(g.pos,an.Lem.Pos))}
		// pos - m est-il déjà noyau de g ?
		if m.dejaNoy() {
			idgrp := phrase.estNuclDe(m)
			if !contient(g.pos, idgrp) {
				continue
			}
		} else if !contient(g.pos, an.Lem.Pos) {
				continue
		}
		//if debog {fmt.Println("   okb")}
		// lexSynt
		// morpho
		var va bool
		for im, morf := range an.Morphos {
			//if debog {fmt.Println("   morf",morf)}
			va = true
			//if debog {fmt.Println("    gmorf",g.morph)}
			for _, gmorf := range g.morph {
				ecl := strings.Split(gmorf, " ")
				for _, e := range ecl {
					va = va && strings.Contains(morf, e)
				}
				for _, lexs := range g.lexSynt {
					va = va && lexsynt(an.Lem.Gr[0], lexs)
				}
				if va {
					m.an = im
					break
				}
			}
			//if debog {fmt.Println("   va", va)}
		}
		//if debog {fmt.Println("   okc va")}
		if !va {
			m.an = -1
			continue
		}
		//if debog {fmt.Println("   okd")}
		if !va {
			continue
		}
		//if debog {fmt.Println("   oke")}
		return true
	}
	return false
}

// vrai si m est compatible avec Sub et le noyau mn
// Sub : pos string, morpho []string, accord string
// gocol.Sr : Lem, Morphos []string
func (m *Mot) estSub(sub *Sub, mn *Mot) bool {
	//debog := m.gr=="Iapeti" && mn.gr=="filius" && contient(sub.morpho, "gén")
	//if debog {fmt.Println("    estSub",m.gr, mn.gr, "sub.pos",sub.pos,"morpho",sub.morpho)}
	var respos, resmorf gocol.Res
	// pos
	if sub.terminal {
		//if debog {fmt.Println("    oka, ans",len(m.ans))}
		for _, an := range m.ans {
			if an.Lem.Pos == sub.pos {
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
	//if debog {fmt.Println("    estSub oka", len(respos),"respos")}
	// morpho
	for _, an := range respos {
		for _, anmorf := range an.Morphos {
			va := true
			// si le mot est déjà noyau, contrôler 
			// le nom de son groupe, (seul l'accord compte ?)
			idgrp := phrase.estNuclDe(m)
			//if debog {fmt.Println("    estSub, idgrp", idgrp)}
			//(idgrp = id du Nod dont m est le noyau)
			if idgrp > "" {
				if sub.generique {
					// il ne faut comparer que l'id générique du groupe
					ee := strings.Split(idgrp, ".")
					idgrp = ee[0]
					va = va && sub.posg == idgrp
				} else {
					va = va && sub.pos == idgrp
				}
			} else {
				//if debog {fmt.Println("   estSub okb")}
				// sinon, on vérifie la morpho du mot
				for _, trait := range sub.morpho {
					va = va && strings.Contains(anmorf, trait)
				}
			}
			//if debog {fmt.Println("   estSub okc")}
			if va {
				resmorf = append(resmorf, an)
			} else {
				continue
			}
			//if debog {fmt.Println("   estSub oke")}
			// accord
			if sub.accord != "" && !mn.accord(m, sub.accord) {
				continue
			}
			// lexsynt
			//if debog {fmt.Println("   estSub oke, sub.lexsynt", sub.lexsynt, an.Lem.Gr[0])}
			vals := true
			for _, lxs := range sub.lexsynt {
				vals = vals && lexsynt(an.Lem.Gr[0], lxs)
			}
			//if debog {fmt.Println("   estSub okf vals", vals)}
			if !vals {
				continue
			}
			if len(resmorf) > 0 {
				//if debog {fmt.Println("   okc, len(resmorf)",len(resmorf))}
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

// nombre de mots subs de m
func (m *Mot) nbSubs() int {
	if !m.dejaNoy() {
		return 0
	}
	var nbm int
	for _, mb := range phrase.mots {
		if mb == m {
			continue
		}
		if mb.subDe(m) {
			nbm++
		}
	}
	return nbm
}

// le mot m est il noyau d'un groupe d'id id ?
func (m *Mot) noyId(id string) bool {
	//debog := m.gr=="iussu" && id=="n.genabl"
	//if debog {fmt.Println("   noyId,ok")}
	for _, n := range phrase.nods {
		if n.nucl == m {
			ee := strings.Split(id, ".")
			//if debog {fmt.Println("   len ee",len(ee))}
			if len(ee) > 1 {
				//if debog {fmt.Println("   noyId, n.grp.id",n.grp.id,"id",id)}
				// si id contient '.', le nod doit avoir un id complet
				if n.grp.id == id {
					return true
				}
			} else {
				// sinon, l'id générique suffit
				eem := strings.Split(n.grp.id, "." )
				if eem[0] == id {
					return true
				}
			}
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
