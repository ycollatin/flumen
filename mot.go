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
	//debog := m.gr=="ignem" && g.id=="v.prepobj"
	//if debog {fmt.Println("estNoyau",m.gr,g.id)}
	var ans3 gocol.Res
	// vérif du pos
	for _, a := range m.ans {
		if contient(g.pos, a.Lem.Pos) {
			ans3 = append(ans3, a)
		}
	}
	//if debog {fmt.Println("  .estNoyau, oka")}
	// vérif lexicosyntaxique
	for i, a := range ans3 {
		va := true
		for _, ls := range g.lexSynt {
			va = va && lexsynt(a.Lem.Gr[0], ls)
		}
		if !va {
			ans3 = supprSr(ans3, i)
		}
	}
	//if debog {fmt.Println("  .estNoyau, ans3",len(ans3))}
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
		//if debog {fmt.Println("  .estNoyau, morfos",len(morfos))}
		if len(morfos) > 0 {
			sr.Morphos = morfos
			m.ans2 = append(m.ans2, sr)
		}
	}
	return len(m.ans2) > 0
}

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
	//debog := sub.groupe.id=="v.prepobj" && m.gr == "immortalibus" //&& mn.gr=="petebant"
	//if debog {fmt.Println("estSub m",m.gr,"pos",m.pos,"sub",sub.groupe.id,"mn",mn.gr)}
	// signet motestSub
	var ans2 gocol.Res
	// vérification des pos
	if m.pos != "" && !sub.vaPos(m.pos) {
		//if debog {fmt.Println("  .estSub,",sub.groupe.id, "(sub) vaPos("+m.pos+")",sub.vaPos(m.pos))}
		return false
	} else {
		for _, an := range m.ans {
			//if debog {fmt.Println("  .estSub,an.Lem.Pos",an.Lem.Pos,"vapos",sub.vaPos(an.Lem.Pos))}
			if sub.vaPos(an.Lem.Pos) {
				ans2 = append(ans2, an)
			}
		}
	}
	//if debog {fmt.Println("  .estSub, len(ans2)",len(ans2))}
	if len(ans2) == 0 {
		return false
	}
	//if debog {fmt.Println("  .estSub, oka")}
	// accord et morphologie
	// pour toutes les morphos valides de mn
	for _, ann := range mn.ans2 {
		for _, morfn := range ann.Morphos {
			// pour toutes les morphos valides de m
			for i, ansub := range ans2 {
				var lmorf []string
				for _, morfs := range ansub.Morphos {
					//if debog {fmt.Println("  .estSub,morpho morfs",morfs,"sub",sub.morpho,sub.vaMorpho(morfs))}
					if accord(morfn, morfs, sub.accord) && sub.vaMorpho(morfs) {
						lmorf = append(lmorf, morfs)
					}
				}
				if len(lmorf) == 0 {
					ans2 = supprSr(ans2, i)
				} else {
					ans2[i].Morphos = lmorf
				}
			}
		}
	}
	if len(ans2) > 0 {
		m.ans2 = ans2
		return true
	}
	return false
}

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
	//debog := g.id=="v.prepobj" && m.gr == "ignem"
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
	//if debog {fmt.Println("  .noeud oka, estNoyau",m.gr,g.id,m.estNoyau(g),"lante",lante)}
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
	//if debog {fmt.Println("  .noeud okb",lante,"lante, r",r)}
	// reгcherche rétrograde des subs ante
	for ia := lante-1; ia > -1; ia-- {
		if r < 0 {
			break
		}
		//if debog {fmt.Println("  .noeud, oka, ia", ia,"r",r)}
		sub := g.ante[ia]
		ma := texte.phrase.mots[r]
		//if debog {fmt.Println(" .noeud ma0",ma.gr,"dejasub",ma.dejaSub(),"r",r)}
		// passer les mots
		for ma.dejaSub() {
			r--
			if r < 0 {
				return nil
			}
			ma = texte.phrase.mots[r]
		}
		//if debog {fmt.Println(" .noeud ma",ma.gr,"estSub",m.gr,"grup",sub.groupe.id,ma.estSub(sub, m))}
		if m.estSubDe(ma) || !ma.estSub(sub, m) {
			// réinitialiser lemme et morpho de ma
			return nil
		}
		ma.sub = sub
		nod.mma = append(nod.mma, ma)
		r--
		//if debog {fmt.Println("    vu",ma.gr)}
	}
	//if debog {fmt.Println("  .noeud okd",len(g.post),"g.post, rang",rang,"nbmots",texte.phrase.nbmots)}
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
		// TODO ? fixer lemme et morpho de tous les mots du nod
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

func supprSr(res gocol.Res, p int) gocol.Res {
	var ret gocol.Res
	for i, an := range res {
		if i != p {
			ret = append(ret, an)
		}
	}
	return ret
}
