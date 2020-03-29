//       mot.go - Publicola

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
	rang		int
	nl, nm		int			// n°s de Sr et de morpho
	ans			gocol.Res	// ensemble des lemmatisations
	lexsynt		[]string	// propriétés lexicosyntaxiques
	idgr		string		// id du groupe dont le mot est le noyau
	sub			*Sub		// sub qui lie le Mot à son noyau
}

func creeMot(m string) *Mot {
	//debog := m == "luto"
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
	mot.nl = -1
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
	//debog := m.gr=="Prometheus" && g.id=="n.appFam"
	//if debog {fmt.Println("   estNoyau",m.gr,g.id)}
	for nam, an := range m.ans {
		//if debog {fmt.Println("   .estNoyau",an.Lem.Gr, an.Morphos,"dejaNoy",m.dejaNoy())}
		//if debog{fmt.Println("   .estNoyau g.pos",g.pos,"an.Lem.Pos",an.Lem.Pos)}
		// pos - m est-il déjà noyau de g ?
		if m.dejaNoy() {
			//if debog {fmt.Println("    .estNoyau idgrp",phrase.estNuclDe(m))}
			idgrp := phrase.estNuclDe(m)
			if !contient(g.pos, idgrp) {
				continue
			}
		} else if !contient(g.pos, an.Lem.Pos) {
				continue
		}
		//if debog {fmt.Println("    .estNoyau okb")}
		// morpho
		var va bool
		for im, morf := range an.Morphos {
			if m.nl == nam && im != m.nm {
				// la lemmatisation mot est déjà fixée à m.an
				continue
			}
			va = true
			for _, gmorf := range g.morph {
				//if debog {fmt.Println("   gmorf",gmorf)}
				ecl := strings.Split(gmorf, " ")
				for _, e := range ecl {
					//if debog {fmt.Println("     ",morf,"contains",e,strings.Contains(morf,e))}
					va = va && strings.Contains(morf, e)
				}
				if va {
					// XXX Mot.ans désigne le nombre de lemmes.
					//     chaque lemme a plusieurs morphos !
					m.nl = nam
					m.nm = im
					break
				}
			}
			//if debog {fmt.Println("   va", va)}
			for _, lexs := range g.lexSynt {
				va = va && lexsynt(an.Lem.Gr[0], lexs)
			}
		}
		//if debog {fmt.Println("    .estNoyau okc va", va)}
		if !va {
			m.nl = -1
			m.nm = -1
			continue
		}
		//if debog {fmt.Println("   okd, return true")}
		return true
	}
	//if debog {fmt.Println("   oke false")}
	return false
}

// vrai si m est compatible avec Sub et le noyau mn
// Sub : pos string, morpho []string, accord string
// gocol.Sr : Lem, Morphos []string
func (m *Mot) estSub(sub *Sub, mn *Mot) bool {
	debog := m.gr=="filius" && mn.gr=="Prometheus" && sub.groupe.id=="n.appFam"
	if debog {fmt.Println("   .estSub",m.gr, mn.gr, "sub.groupe:"+sub.groupe.id+"-morpho",sub.morpho)}
	var respos, resmorf gocol.Res
	// le sub a-t-il le bon pos ?
	//if debog {fmt.Println("    . estSub, sub.terminal",sub.terminal)}
	if sub.terminal {
		//if debog {fmt.Println("    . oka, ans",len(m.ans))}
		// le sub a plusieurs pos, ex. "NP n"
		for _, an := range m.ans {
			if sub.vaPos(an) {
				respos = append(respos, an)
			}
		}
		//if debog {fmt.Println("    . estSub okb", len(respos),"respos")}
		if len(respos) == 0 {
			return false
		}
	} else {
		if debog {fmt.Println("   . estSub m.idgr",m.idgr,"==",sub.lien)}
		// le mot m est il noyau d'un groupe sub
		//if contient(sub.idgr(), m.idgr) {
		if sub.vaId(m.idgr) {
			respos = m.ans
		} else {
			if debog {fmt.Println("    . estSub false")}
			return false
		}
	}
	if debog {fmt.Println("    . estSub, respos",len(respos))}
	// le sub a-t-il la bonne morpho ?
	for _, an := range respos {
		for _, anmorf := range an.Morphos {
			//if debog {fmt.Println("    . anmorf",anmorf)}
			va := true
			// si le mot est déjà noyau, contrôler 
			// le nom de son groupe, (seul l'accord compte ?)
			//idgrp := phrase.estNuclDe(m)
			idgrp := m.estNuclDuGroupe()
			if debog {fmt.Println("    . estSub, idgrp", idgrp,"generique",sub.generique)}
			//(idgrp = id du Nod dont m est le noyau)
			if idgrp > "" {
				if sub.generique {
					// il ne faut comparer que l'id générique du groupe
					ee := strings.Split(idgrp, ".")
					idgrp = ee[0]
					va = va && sub.vaId(idgrp)
					if debog {fmt.Println("   estSub, va", va)}
				} else {
					va = va && sub.vaId(idgrp)
					if debog {fmt.Println("   . estSub, va",va)}
				}
			}
			//if debog {fmt.Println("    . estSub okc, va",va)}
			// vérification de la morpho du mot
			for _, trait := range sub.morpho {
				//if debog {fmt.Println("   estSub anmorf",anmorf,"trait",trait)}
				va = va && strings.Contains(anmorf, trait)
			}
			//if debog {fmt.Println("   estSub okc")}
			if va {
				resmorf = append(resmorf, an)
			} else {
				continue
			}
			//if debog {fmt.Println("   . estSub okd")}
			// accord
			if sub.accord != "" && !mn.accord(m, sub.accord) {
				continue
			}
			// lexsynt
			//if debog {fmt.Println("   . estSub oke, sub.lexsynt",sub.lexsynt,an.Lem.Gr[0])}
			vals := true
			for _, lxs := range sub.lexsynt {
				vals = vals && lexsynt(an.Lem.Gr[0], lxs)
			}
			//if debog {fmt.Println("   . estSub okf vals", vals)}
			if !vals {
				continue
			}
			if len(resmorf) > 0 {
				//if debog {fmt.Println("   return true, len(resmorf)",len(resmorf))}
				return true
			}
		}
	}
	return false
}

// id du groupe dont m est le noyau
// XXX nod.grp est-il le groupe dont nod est le noyau
func (m *Mot) estNuclDuGroupe() string {
	for _, nod := range phrase.nods {
		if nod.nucl == m {
			return nod.grp.id
		}
	}
	return ""
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

// le mot m est il noyau d'un groupe g
func (m *Mot) noyId(sub *Sub) bool {
	if m.nl < 0 {
		return false
	}
	for _, n := range phrase.nods {
		//if (n.nucl == m) {fmt.Println("debog noyId",m.gr,"len ans",len(m.ans),"nl",m.nl)}
		if n.nucl == m {
			return true
		}
	}
	return false
}

/*
// le mot m est il noyau d'un groupe sub
func (m *Mot) noyId(id string) bool {
	//debog := m.gr=="rem" && id=="n.prepDetAcc"
	//if debog {fmt.Println("   ..noyId",m.gr,id,"nb noeuds",len(phrase.nods))}
	for _, n := range phrase.nods {
		//if debog {fmt.Println("   n.nucl",n.nucl.gr,"m",m.gr)}
		if n.nucl == m {
			ee := strings.Split(id, ".")
			//if debog {fmt.Println("    ..noyId len ee",len(ee),"n.grp.id",n.grp.id)}
			if len(ee) > 1 {
				//if debog {fmt.Println("    ..noyId>1 n.grp.id",n.grp.id,"id",id)}
				// si id contient '.', le nod doit avoir un id complet
				if n.grp.id == id {
					//if debog {fmt.Println("   ..noyId true")}
					return true
				}
			} else {
				// sinon, l'id générique suffit
				eem := strings.Split(n.grp.id, "." )
				//if debog {fmt.Println("   noyId else, n.grp.id",n.grp.id,"eem0",eem[0])}
				if eem[0] == id {
					return true
				}
			}
		}
	}
	return false
}
*/

func (ma *Mot) subDe(mb *Mot) bool {
	// chercher le groupe dont mb est noyau
	for _, n := range phrase.nods {
		if mb == n.nucl {
			return ma.elDe(n)
		}
	}
	return false
}
