//       mot.go - Publicola

package main

import (
	//"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

// signets :
//
// motnoeud
// motestnoyau
// motestSub

// rappel de la lemmatisation dans gocol :
// type Sr struct {
//	Lem     *Lemme
//	Nmorph  []int
//	Morphos []string
// }
//
// type Res []Sr

type Lm struct {
	l	*gocol.Lemme
	m	string
}

type Mot struct {
	gr			string
	rang		int
	ans, ans2	gocol.Res	// ensemble des lemmatisations
	llm			[]Lm		// liste des lemmes ٍ+ morpho possibles
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
	return mot
}

func accord(lma, lmb, cgn string) bool {
	va := true
	for i:=0; i<len(cgn); i++ {
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

func (m *Mot) dejaNoy() bool {
	for _, n := range texte.phrase.nods {
		if n.nucl == m {
			return true
		}
	}
	return false
}

func (m *Mot) dejaSub() bool {
	for _, n := range texte.phrase.nods {
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

// teste si m peut être le noyau du groupe groupe g
func (m *Mot) estNoyau(g *Groupe) bool {
//signet motestnoyau
	var ans3 gocol.Res
	//debog := m.gr=="currum" && g.id=="n.gen"
	//if debog {fmt.Println("estNoyau",m.gr,g.id,"nl/nm",m.nl,m.nm,"eluc.",m.elucide())}
	// vérif du pos
	for _, a := range m.ans {
		if contient(g.pos, a.Lem.Pos) {
			ans3 = append(ans3, a)
		}
		if len(ans3) == 0 {
			return false
		}
		// vérif lexicosyntaxique
		for _, an := range ans3 {
			for _, ls := range g.lexSynt {
				if !lexsynt(an.Lem.Gr[0], ls) {
					continue
				}
			}
		}
		if len(ans3) == 0 {
			return false
		}
		// vérif morpho
		for _, sr := range ans3 {
			var morfos []string  // morphos de sr acceptées par g
			for _, morf := range sr.Morphos {
				if g.vaMorph(morf) {
					morfos = append(morfos, morf)
				}
			}
			if len(morfos) > 0 {
				m.ans2 = append(m.ans2, sr)
			}
		}
	}
	return len(m.ans2) > 0
}

/*
// teste si m peut être le noyau du groupe groupe g
func (m *Mot) estNoyau(g *Groupe) bool {
	var lemmes []*gocol.Lemme
	//debog := m.gr=="currum" && g.id=="n.gen"
	//if debog {fmt.Println("estNoyau",m.gr,g.id,"nl/nm",m.nl,m.nm,"eluc.",m.elucide())}
	// vérif du pos
	va := false
	for _, a := range m.ans {
		if contient(g.pos, a.Lem.Pos) {
			lemmes = append(lemmes, a.Lem)
		}
		if len(lemmes) == 0 {
			return false
		}
		// vérif lexicosyntaxique
		for i, l := range lemmes {
			for _, ls := range g.lexSynt {
				if !lexsynt(l.Gr[0], ls) {
					if len(lemmes) == 1 {
						return false
					}
					lemmes[i] = lemmes[len(lemmes) - 1]
					lemmes[len(lemmes)-1]  = nil
					lemmes = lemmes[:len(lemmes)-1]
				}
			}
		}
		// vérif morpho
		for i, sr := range m.ans {
			// si le lemme de sr n'est pas dans lemmes, continuer
			for _, l := range lemmes {
				va = va || l == sr.Lem
			}
			if !va {
				continue
			}
			va := false
			for j, morf := range sr.Morphos {
				if g.vaMorph(morf) {
					var nlm Lm
					nlm.l = sr.Lem
					nlm.m = morf
					m.llm = append(m.llm, nlm)
					va := true
				}
			}
		}
	}
	return va
}
*/

// id des Nod dont m est déjà le noyau
func (m *Mot) estNuclDe() []string {
	var ret []string
	for _, nod := range texte.phrase.nods {
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
	// signet motestSub
	for lm := range m.llm {
		// vérification de toutes les morphos	
		var a gocol.Sr
		va := false
		for i, an := range m.ans {
			//if debog {fmt.Println("  .estSub, i",i,"an.lem.pos",an.Lem.Pos)}
			if sub.vaPos(an.Lem.Pos) {
				va = true
				m.tmpl = i
				a = an
				break
			}
		}
		if !va {
			return false
		}
		for i, morf := range a.Morphos {
			// accord
			// XXX
			//if debog {fmt.Println("  .estSub,i morf", i, morf)}
			if sub.vaMorpho(morf) {
				m.tmpm = i
				return true
			}
		}
	}
	return false
}

/*
func (m *Mot) estSub(sub *Sub, mn *Mot) bool {
	// signet motestSub
	for lm := range m.llm {
		// accord
		if sub.accord > "" {
			if !accord(m, sub.accord) {
				return false
			}
		}
		// vérification de toutes les morphos	
		var a gocol.Sr
		va := false
		for i, an := range m.ans {
			//if debog {fmt.Println("  .estSub, i",i,"an.lem.pos",an.Lem.Pos)}
			if sub.vaPos(an.Lem.Pos) {
				va = true
				m.tmpl = i
				a = an
				// FIXME : la première solution est prise. Une autre pourrait être la bonne !
				break
			}
		}
		if !va {
			return false
		}
		for i, morf := range a.Morphos {
			//if debog {fmt.Println("  .estSub,i morf", i, morf)}
			if sub.vaMorpho(morf) {
				m.tmpm = i
				return true
			}
		}
		return false
	}
}
*/

// id du groupe dont m est le noyau
func (m *Mot) estNuclDuGroupe() string {
	for _, nod := range texte.phrase.nods {
		if nod.nucl == m {
			return nod.grp.id
		}
	}
	return ""
}

func (ma *Mot) estSubDe(mb *Mot) bool {
	for _, n := range texte.phrase.nods {
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

/*
// morpho définitive
func (m *Mot) morphodef() string {
	if !m.elucide() {
		return ""
	}
	return m.ans[m.nl].Morphos[m.nm]
}
*/

// nombre de mots subs de m
func (m *Mot) nbSubs() int {
	if !m.dejaNoy() {
		return 0
	}
	var nbm int
	for _, mb := range texte.phrase.mots {
		if mb == m {
			continue
		}
		if mb.subDe(m) {
			nbm++
		}
	}
	return nbm
}

// signet motnoeud
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
	if texte.phrase.nbmots - rang < len(g.post) {
		return nil
	}
	//if debog {fmt.Println("  .noeud oka, estNoyau",m.gr,g.id,m.estNoyau(g))}
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
	//if debog {fmt.Println("  .noeud okb",lante,"lante")}
	// reгcherche rétrograde des subs ante
	for ia := lante-1; ia > -1; ia-- {
		if r < 0 {
			break
		}
		//if debog {fmt.Println("  .noeud, oka")}
		sub := g.ante[ia]
		ma := texte.phrase.mots[r]
		// passer les mots
		for ma.dejaSub() && r > 0 {
			r--
			ma = texte.phrase.mots[r]
		}
		//if debog {fmt.Println(" ma",ma.gr,"nl/nm",ma.nl,ma.nm,"estSub",m.gr,"grup",sub.groupe.id,ma.estSub(sub, m))}
		if m.estSubDe(ma) || !ma.estSub(sub, m) {
			// réinitialiser lemme et morpho de ma
			return nil
		}
		ma.sub = sub
		nod.mma = append(nod.mma, ma)
		r--
		//if debog {fmt.Println("    vu",ma.gr)}
	}
	//if debog {fmt.Println("  .noeud okd",len(g.post),"g.post, rang",rang,"nbmots",phrase.nbmots)}
	// post
	for ip, sub := range g.post {
		r := rang + ip + 1
		if r >= texte.phrase.nbmots {
			break
		}
		mp := texte.phrase.mots[r]
		//if debog {fmt.Println("post, mp",mp.gr)}
		for mp.dejaSub() && r < len(texte.phrase.mots) - 1 {
			r++
			mp = texte.phrase.mots[r]
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
		return nod
	}
	return nil
}

// vrai si ma est sub de mb
func (ma *Mot) subDe(mb *Mot) bool {
	// chercher le groupe dont mb est noyau
	for _, n := range texte.phrase.nods {
		if mb == n.nucl {
			return ma.elDe(n)
		}
	}
	return false
}
