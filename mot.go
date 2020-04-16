//       mot.go - Gentes

// signets :
//
// motnoeud
// motresnoyau
// motestNoyauDeGroupe
// motresSub

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

type Lm struct {
	l *gocol.Lemme
	m string
}

type Mot struct {
	gr         string    // graphie du mot
	rang       int       // rang du mot dans la phrase à partir de 0
	ans, ans2  gocol.Res // ensemble des lemmatisations, ans provisoire
	dejasub    bool      // le mot est déjà l'élément d'n nœud
	llm        []Lm      // liste des lemmes ٍ+ morpho possibles
	tmpl, tmpm int       // n°s provisoires de Sr et morpho
	pos        string    // id du groupe dont le mot est noyau
	// ou à défaut pos du mot, si elle est décidée
	lexsynt []string // propriétés lexicosyntaxiques
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

func (m *Mot) dejaNoy() bool {
	for _, n := range texte.phrase.nods {
		if n.nucl == m {
			return true
		}
	}
	return false
}

func (ma *Mot) domine(mb *Mot) bool {
	//mnoy := noyDe(mb)
	mnoy := mb.noyau()
	for mnoy != nil {
		if mnoy == ma {
			return true
		}
		mnoy = mnoy.noyau()
		//noyDe(mnoy)
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
func (m *Mot) resSub(sub *Sub, mn *Mot) gocol.Res {
	// signet motresSub
	var ans2 gocol.Res
	// vérification des pos
	if m.pos != "" {
		// 1. La pos du mot est définitive
		// noyaux exclus
		veto := false
		lgr := m.estNuclDe()
		for _, noy := range sub.noyexcl {
			veto = veto || contient(lgr, noy.id)
		}
		if veto {
			return ans2
		}
		// noyaux possibles
		va := false
		for _, noy := range sub.noyaux {
			va = va || noy.vaPos(m.pos)
		}
		if !va {
			return ans2
		}
		ans2 = m.ans
	} else {
		// 2. La pos définitif n'est pas encore fixée
		for _, an := range m.ans {
			// lexicosyntaxe
			va := true
			for _, ls := range sub.lexsynt {
				va = va && lexsynt(an.Lem.Gr[0], ls)
			}
			if !va {
				continue
			}
			// canon et POS
			for _, noy := range sub.noyaux {
				if noy.canon > "" {
					if noy.vaSr(an) {
						ans2 = append(ans2, an)
						break
					}
				} else {
					if noy.vaPos(an.Lem.Pos) {
						ans2 = append(ans2, an)
						break
					}
				}
			}
		}
	}
	if len(ans2) == 0 {
		return ans2
	}

	//morphologie
	var ans3 gocol.Res
	for _, an := range ans2 {
		var lmorf []string
		for _, morfs := range an.Morphos {
			// pour toutes les morphos valides de m
			if strings.Contains(morfs, "inv.") || sub.vaMorpho(morfs) {
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
	return ans4
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

// si m peut être noyau d'un gourpe g, un Nod est renvoyé, sinon nil.
func (m *Mot) noeud(g *Groupe) *Nod {
	// signet motnoeud
	rang := m.rang
	lante := len(g.ante)
	// mot de rang trop faible
	if rang-lante < 0 {
		return nil
	}
	// ou trop élevé
	if rang+len(g.post)-1 >= texte.phrase.nbmots {
		return nil
	}
	// m peut-il être noyau du groupe g ?
	res2 := m.resNoyau(g)
	if len(res2) == 0 {
		return nil
	}
	m.ans2 = res2
	// création du noeud de retour
	nod := new(Nod)
	nod.grp = g
	nod.nucl = m
	nod.rang = rang
	// vérif des subs
	// ante
	r := rang - 1
	// reгcherche rétrograde des subs ante
	for ia := lante - 1; ia > -1; ia-- {
		if r < 0 {
			// le rang du mot est < 0 : impossible
			return nil
		}
		ma := texte.phrase.mots[r]
		// passer les mots déjà subordonnés
		for ma.dejasub {
			r--
			if r < 0 {
				return nil
			}
			ma = texte.phrase.mots[r]
		}
		// vérification de réciprocité, puis du lien lui-même
		if ma.domine(m) {
			return nil
		}
		sub := g.ante[ia]
		res3 := ma.resSub(sub, m)
		if len(res3) == 0 {
			return nil
		}
		ma.ans2 = res3
		nod.mma = append(nod.mma, ma)
		r--
	}
	// post
	for ip, sub := range g.post {
		r := rang + ip + 1
		if r >= texte.phrase.nbmots {
			break
		}
		if sub.lien == "" {
			continue
		}
		mp := texte.phrase.mots[r]
		for mp.dejasub {
			r++
			if r >= texte.phrase.nbmots {
				return nil
			}
			mpn := mp.noyau()
			if mpn != nil && mpn.rang < m.rang {
				return nil
			}
			mp = texte.phrase.mots[r]
		}
		// réciprocité
		if mp.domine(m) {
			return nil
		}
		res4 := mp.resSub(sub, m)
		if len(res4) == 0 {
			return nil
		}
		mp.ans2 = res4
		nod.mmp = append(nod.mmp, mp)
		r++
	}
	// fixer les pos et sub des mots du noeud
	if len(nod.mma)+len(nod.mmp) > 0 {
		m.pos = g.id
		for _, m := range nod.mma {
			m.dejasub = true
		}
		for _, m := range nod.mmp {
			m.dejasub = true
		}
		return nod
	}
	return nil
}

func (m *Mot) noyau() *Mot {
	for _, n := range texte.phrase.nods {
		for _, msub := range n.mma {
			if msub == m {
				return n.nucl
			}
		}
		for _, msub := range n.mmp {
			if msub == m {
				return n.nucl
			}
		}
	}
	return nil
}

// renvoie quelles lemmatisations de m lui permettent d'être le noyau du groupe g
func (m *Mot) resNoyau(g *Groupe) gocol.Res {
	//signet motresnoyau

	var ans3 gocol.Res
	// vérif du pos
	if m.pos != "" {
		// 1. La pos définitif est fixée
		va := false
		for _, noy := range g.noyaux {
			va = va || noy.vaPos(m.pos)
		}
		if !va {
			return ans3
		}
		// vérification du Pos des lemmatisations sélectionnées
		va = false
		for _, noy := range g.noyaux {
			for _, an := range m.ans2 {
				va = va || noy.vaPos(an.Lem.Pos)
			}
		}
		if !va {
			return ans3
		}
		ans3 = m.ans2
	} else {
		// Le mot est encore isolé
		for _, a := range m.ans {
			for _, noy := range g.noyaux {
				if noy.canon > "" {
					if noy.vaSr(a) {
						ans3 = append(ans3, a)
						break
					}
				} else {
					if noy.vaPos(a.Lem.Pos) {
						ans3 = append(ans3, a)
						break
					}
				}
			}
		}
	}

	// vérif lexicosyntaxique
	var ans4 gocol.Res
	for _, a := range ans3 {
		va := true
		for _, ls := range g.lexsynt {
			va = va && lexsynt(a.Lem.Gr[0], ls)
		}
		for _, ls := range g.exclls {
			va = va && !lexsynt(a.Lem.Gr[0], ls)
		}
		if va {
			ans4 = append(ans4, a)
		}
	}

	// vérif morpho.
	// Si aucune n'est requise, renvoyer true
	if len(g.morph) == 0 {
		return ans4
	}

	var ans5 gocol.Res
	for _, sr := range ans4 {
		var morfos []string // morphos de sr acceptées par g
		for _, morf := range sr.Morphos {
			if g.vaMorph(morf) {
				morfos = append(morfos, morf)
			}
		}
		if len(morfos) > 0 {
			sr.Morphos = morfos
			ans5 = append(ans5, sr)
		}
	}
	return ans5
}
