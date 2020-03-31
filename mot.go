//       mot.go - Publicola

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
	ans			gocol.Res	// ensemble des lemmatisations
	nl, nm		int			// n°s de Sr et de morpho, quand ils sont fixés
	tmpl, tmpm	int			// n°s provisoires de Sr et morpho
	pos			string		// id du groupe dont le mot est noyau
							// ou à défaut pos du mot, si elle est décidée
	lexsynt		[]string	// propriétés lexicosyntaxiques
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
		for i, a := range mot.ans {
			mot.ans[i] = genus(a)
		}
	}
	mot.nl = -1
	mot.nm = -1
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

func (m *Mot) elucide() bool {
	return m.nl > -1 && m.nm > -1
}

// teste si m peut être le noyau du groupe g
func (m *Mot) estNoyau(g *Groupe) bool {
	//debog := m.gr=="luto" && g.id=="n.prepAbl"
	//if debog {fmt.Println("estNoyau",m.gr,g.id,"nl/nm",m.nl,m.nm)}
	va := false

	//var nl, nm int

	// vérif du pos
	mnuclde := m.estNuclDe()
	//if debog {fmt.Println("  .estNoyau, mnuclde",mnuclde)}
	if len(mnuclde) == 0 {
		if m.elucide() {
			va = contient(g.pos, m.ans[m.nl].Lem.Pos)
		} else {
			for _, a := range m.ans {
				va = va || contient(g.pos, a.Lem.Pos)
			}
		}
	} else {
		for _, mnd := range mnuclde {
			va = va || contient(g.pos, mnd)
		}
	}
	if !va {
		return false
	}
	//if debog {fmt.Println(" .estNoyau, pos, nl/nm",m.nl,m.nm)}
	// vérif de la morpho
	if !m.elucide() {
		for i, an:= range m.ans {
			//if debog {fmt.Println("   .estNoyau >, morf",morf)}
			for j, gm := range an.Morphos {
				va = va || g.vaMorph(gm)
				if va {
					m.tmpl = i
					m.tmpm = j
					//if debog {fmt.Println("   .estNoyau, true, ml mn",m.nl,m.nm)}
					return true
				}
			}
		}
	} else {
		return g.vaMorph(m.ans[m.nl].Morphos[m.nm])
	}
	return false
}

// id des Nod dont m est déjà le noyau
func (m *Mot) estNuclDe() []string {
	var ret []string
	for _, nod := range phrase.nods {
		if nod.nucl == m {
			ret = append(ret, nod.grp.id)
		}
	}
	return ret
}

// vrai si m est compatible avec Sub et le noyau mn
// Sub : pos string, morpho []string, accord string
// gocol.Sr : Lem, Morphos []string
func (m *Mot) estSub(sub *Sub, mn *Mot) bool {
	//debog := m.gr=="ex" && mn.gr == "luto" && sub.groupe.id=="n.prepAbl"
	//if debog {fmt.Println("estSub m",m.gr,"mn",mn.gr,"grup",sub.groupe.id,"m.nl",m.nl)}
	//si le mot a déjà une lemmatisation fixée
	if m.elucide() {
		a := m.ans[m.nl]
		//if debog {fmt.Println(" .estSub alempos",a.Lem.Pos,"morfo",a.Morphos[m.nm])}
		if sub.vaPos(m.pos) && sub.vaMorpho(a.Morphos[m.nm]) {
			//if debog {fmt.Println(" .estsub, elucide", m.morphodef())}
			return true
		}
	} else {
	    // vérification de toutes les morphos	
		va := false
		var a gocol.Sr
		for i, an := range m.ans {
			if sub.vaPos(an.Lem.Pos) {
				va = true
				m.tmpl = i
				a = an
			}
		}
		va = false
		for i, morf := range a.Morphos {
			if sub.vaMorpho(morf) {
				va = true
				m.tmpm = i
			//if debog {fmt.Println("  .estsub, elucide2", m.morphodef())}
			}
		}
		if va {
			return true
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

func (ma *Mot) estSubDe(mb *Mot) bool {
	for _, n := range phrase.nods {
		if mb == n.nucl {
			for _, sub := range n.mma {
				if sub == ma {
					return true
				}
			}
			for _, sub := range n.mmp {
				if sub == ma {
					return true
				}
			}
		}
	}
	return false
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

// morpho définitive
func (m *Mot) morphodef() string {
	if !m.elucide() {
		return ""
	}
	return m.ans[m.nl].Morphos[m.nm]
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

// si m peut être noyau d'un gourpe g, un Nod est renvoyé, sinon nil.
func (m *Mot) noeud(g *Groupe) *Nod {
	//debog := g.id=="n.fam" && m.gr == "filius"
	//if debog {fmt.Println("noeud", m.gr, g.id)}
	rang := m.rang
	lante := len(g.ante)
	// mot de rang trop faible
	if rang < lante {
		return nil
	}
	// ou trop élevé
	if phrase.nbm() - rang < len(g.post) {
		return nil
	}
	//if debog {fmt.Println("   noeud oka, estNoyau",m.gr,g.id,m.estNoyau(g))}
	// m peut-il être noyau du groupe g ?
	if !m.estNoyau(g) {
		return nil
	}

	// création du noeud de retour
	nod := new(Nod)
	nod.grp = g
	nod.nucl = m
	nod.rang = rang
	// vérif des subs
	// ante
	r := rang - 1
	//if debog {fmt.Println("   noeud okb")}
	// reгcherche rétrograde des subs ante
	for ia := lante-1; ia > -1; ia-- {
		sub := g.ante[ia]
		ma := phrase.mots[r]
		// passer les mots
		for ma.dejaSub() && r > 0 {
			r--
			ma = phrase.mots[r]
		}
		//if debog {fmt.Println("  ma",ma.gr,"nl/nm",ma.nl,ma.nm,"estSub",m.gr,"id grup",sub.groupe.id,ma.estSub(sub, m))}
		if m.estSubDe(ma) || !ma.estSub(sub, m) {
			// réinitialiser lemme et morpho de ma
			return nil
		}
		ma.sub = sub
		nod.mma = append(nod.mma, ma)
		r--
		//if debog {fmt.Println("    vu",ma.gr)}
	}
	//if debog {fmt.Println("   okd",len(g.post),"g.post")}
	// post
	for ip, sub := range g.post {
		r := rang + ip + 1
		mp := phrase.mots[r]
		//if debog {fmt.Println("post, mp",mp.gr)}
		for mp.dejaSub() && r < len(phrase.mots) - 1 {
			r++
			mp = phrase.mots[r]
		}
		//if debog {fmt.Println("     mp", mp.gr,"estSub",m.gr,sub.groupe.id,mp.estSub(sub, m))}
		if m.estSubDe(mp) || !mp.estSub(sub, m) {
			// réinitialiser lemme et morpho de mp
			return nil
		}
		mp.sub = sub
		nod.mmp = append(nod.mmp, mp)
	}
	if len(nod.mma) + len(nod.mmp) > 0 {
		m.pos = g.id
		// fixer lemme et morpho de tous les mots du nod
		nod.valide()
		return nod
	}
	return nil
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

func (m *Mot) valideTmp() {
	m.nl = m.tmpl
	m.nm = m.tmpm
}
