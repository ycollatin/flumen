//       mot.go - Publicola

package main

import (
	"fmt"
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
	// ajout du genre pour les noms
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

func restostr(ans gocol.Res) string {
	var ll []string
	for _, an := range ans {
		lem := an.Lem.Gr
		var mm []string
		for _, m := range an.Morphos {
			mm = append(mm, m)
		}
		ll = append(ll, fmt.Sprintf("%s - %s", lem, strings.Join(mm, ", ")))
	}
	return strings.Join(ll, "\n")
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
	debog := m.gr=="iussu" && g.id=="n.gen"
	if debog {fmt.Println(" -estNoyau",m.gr,g.id,"ans2:",gocol.Restostring(m.ans2))}

	var ans3 gocol.Res
	// vérif du pos
	for _, a := range m.ans {
		for _, noy := range g.noyaux {
			//if debog {fmt.Println("  .estNoyau,noy",noy,"a.Lem.Pos",a.Lem.Pos)}
			if noy.vaPos(a.Lem.Pos) {
				ans3 = append(ans3, a)
				break
			}
		}
	}
	//if debog {fmt.Println("  .estNoyau, oka, len ans3", len(ans3))}
	// vérif lexicosyntaxique
	var ans4 gocol.Res
	for _, a := range ans3 {
		va := true
		for _, ls := range g.lexSynt {
			va = va && lexsynt(a.Lem.Gr[0], ls)
		}
		if va {
			ans4 = append(ans4, a)
		}
	}
	//if debog {fmt.Println("  .estNoyau, okb, len ans3",len(ans3))}
	if len(ans4) == 0 {
		return false
	}
	// vérif morpho. Si aucune n'est requise, renvoyer true
	if len(g.morph) == 0 {
		m.ans2 = ans4
		return true
	}

	var ans5 gocol.Res
	for _, sr := range ans4 {
		var morfos []string  // morphos de sr acceptées par g
		for _, morf := range sr.Morphos {
			//if debog {fmt.Println("  .estNoyau, morf",morf,"g.morph",g.morph)}
			if g.vaMorph(morf) {
				morfos = append(morfos, morf)
			}
		}
		//if debog {fmt.Println("  .estNoyau, morfos",len(morfos))}
		if len(morfos) > 0 {
			sr.Morphos = morfos
			ans5 = append(ans5, sr)
		}
	}
	if len(ans5) > 0 {
		m.ans2 = ans5
		return true
	}
	return false
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
func (m *Mot) estSub(sub *Sub, mn *Mot) gocol.Res {
	debog := sub.groupe.id=="n.prepAbl" && m.gr == "ex" && mn.gr=="luto"
	if debog {fmt.Println(" -estSub m",m.gr,"pos",m.pos,"sub",sub.groupe.id,"mn",mn.gr)}
	// signet motestSub
	var ans2 gocol.Res
	// vérification des pos
	if m.pos != "" {
		// 1. La pos définitif n'est pas encore fixé
		va := false
		for _, noy := range sub.noyaux {
			va = va || noy.vaPos(m.pos)
		}
		if !va {
			return ans2
		}
		ans2 = m.ans2
	} else {
		// 2. La pos du mot est définitive
		for _, an := range m.ans {
			for _, noy := range sub.noyaux {
				if noy.vaPos(an.Lem.Pos) {
					ans2 = append(ans2, an)
					break
				}
			}
		}
	}
	//if debog {fmt.Println("  .estSub, len(ans2)",len(ans2))}
	if len(ans2) == 0 {
		return ans2
	}
	if debog {fmt.Println("  .estSub1, oka, len mn.ans2",len(mn.ans2))}
	//morphologie
	var ans3 gocol.Res
	for _, an := range m.ans2 {
		if debog {fmt.Println("  .estSub2",an.Lem.Gr,"morphos",len(an.Morphos))}
		var lmorf []string
		for _, morfs := range an.Morphos {
			// pour toutes les morphos valides de m
			if sub.vaMorpho(morfs) {
				lmorf = append(lmorf, morfs)
			}
		}
		if len(lmorf) > 0 {
			an.Morphos = lmorf
			ans3 = append(ans3, an)
		}
	}

	// accord
	var ans4 gocol.Res
	// pour toutes les morphos valides de mn
	for _, ann := range mn.ans2 {
		for _, morfn := range ann.Morphos {
			// pour toutes les morphos valides de m
			for _, an := range ans3 {
				var lmorf []string
				for _, morfs := range an.Morphos {
					if accord(morfn, morfs, sub.accord) {
						lmorf = append(lmorf, morfs)
					}
				}
				if len(lmorf) > 0 {
					an.Morphos = lmorf
					ans4 = append(ans4, an)
				}
			}
		}
	}
	if debog {fmt.Println("  .estSub3, sortie ans3",len(ans3))}
	return ans4
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
	debog := g.id=="n.prepAbl" && m.gr == "luto"
	if debog {fmt.Println("noeud", m.gr, g.id,len(m.ans2),"ans2")}
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
	//if debog {fmt.Println("  .noeud okb",lante,"lante, r",r,nod.doc())}
	// reгcherche rétrograde des subs ante
	for ia := lante-1; ia > -1; ia-- {
		if r < 0 {
			// le rang du mot est < 0 : impossible
			return nil
		}
		//if debog {fmt.Println("  .noeud, oka, ia", ia,"r",r)}
		sub := g.ante[ia]
		ma := texte.phrase.mots[r]
		// passer les mots déjà subordonnés
		for ma.dejaSub() {
			r--
			if r < 0 {
				return nil
			}
			ma = texte.phrase.mots[r]
		}
		if debog {fmt.Println(" .noeud ma",ma.gr,"estSub",m.gr,"grup",sub.groupe.id,ma.estSub(sub, m))}
		// vérification de réciprocité, puis du lien lui-même
		if m.estSubDe(ma) || ma.estSub(sub, m) == nil {
			// réinitialiser lemme et morpho de ma
			return nil
		}
		ma.sub = sub
		nod.mma = append(nod.mma, ma)
		r--
		if debog {fmt.Println("    vu",ma.gr)}
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
		if m.estSubDe(mp) || mp.estSub(sub, m) == nil {
			// réinitialiser lemme et morpho de mp
			return nil
		}
		mp.sub = sub
		nod.mmp = append(nod.mmp, mp)
	}
	if len(nod.mma) + len(nod.mmp) > 0 {
		m.pos = g.id
		// TODO ? fixer lemme et morpho de tous les mots du nod
		if debog {fmt.Println("   .noeud", len(m.ans2),"ans2")}
		fmt.Println(g.id,"nod:",restostr(m.ans2))
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
