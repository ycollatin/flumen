//       mot.go - Gentes

// signets :
//
// motnoeud
// motresnoyau
// motestNoyauDeGroupe
// motresSub
// fotesr

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
	ans, ans2  gocol.Res // ensemble des lemmatisations, ans2 réduite par la syntaxe
	restmp     gocol.Res // analyses temporaires du mot pendand le calcul d'un noeud
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
	mot.ans2 = mot.ans
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

func (m *Mot) adeja(sub *Sub) bool {
	sublien := sub.lien
	for _, nod := range texte.phrase.nods {
		if nod.nucl == m {
			for i, _ := range nod.mma {
				if nod.grp.ante[i].lien == sublien {
					return true
				}
			}
			for i, _ := range nod.mmp {
				if nod.grp.post[i].lien == sublien {
					return true
				}
			}
		}
	}
	return false
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
	m.restmp = m.ans2
	res := m.resNoyau(g, m.restmp)
	if res == nil {
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
		//ma.restmp = ma.ans2
		res := ma.resSub(sub, m, ma.restmp)
		if res == nil {
			return nil
		}
		ma.restmp = res
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
		//mp.restmp = mp.ans2
		res := mp.resSub(sub, m, mp.restmp)
		if res == nil {
			return nil
		}
		mp.restmp = res
		nod.mmp = append(nod.mmp, mp)
		r++
	}
	// fixer les pos et sub des mots du noeud
	if len(nod.mma)+len(nod.mmp) > 0 {
		m.pos = g.id
		m.ans2 = m.restmp
		for _, ms := range nod.mma {
			ms.dejasub = true
			ms.ans2 = ms.restmp
			ms.restmp = nil
		}
		for _, ms := range nod.mmp {
			ms.dejasub = true
			ms.ans2 = ms.restmp
			ms.restmp = nil
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

func oteSr(res gocol.Res, n int) gocol.Res {
	// signet foteSr
	var restmp gocol.Res
	for i, sr := range res {
		if i != n {
			restmp = append(restmp, sr)
		}
	}
	return restmp
}

// renvoie quelles lemmatisations de m lui permettent d'être le noyau du groupe g
func (m *Mot) resNoyau(g *Groupe, res gocol.Res) gocol.Res {
	//signet motresnoyau
	// vérif du pos
	if m.pos != "" {
		// 1. La pos définitif est fixée
		va := false
		for _, noy := range g.noyaux {
			va = va || noy.vaPos(m.pos)
		}
		if !va {
			return nil
		}
		/*
		// vérification du Pos des lemmatisations sélectionnées
		var aoter []int
		for _, noy := range g.noyaux {
			for i, an := range res {
				if !noy.vaPos(an.Lem.Pos) {
					aoter = append(aoter, i)
				}
			}
		}
		for ao := len(aoter) -1; ao > -1; ao-- {
			res = oteSr(res, aoter[ao])
		}
		if len(res) == 0 {
			return nil
		}
		*/
	} else {
		// Le mot est encore isolé
		var aoter []int
		for i, a := range res {
			va := false
			for _, noy := range g.noyaux {
				if noy.canon > "" {
					va = va || noy.vaSr(a)
				} else {
					va = va || noy.vaPos(a.Lem.Pos)
				}
			}
			if !va {
				//res = oteSr(res, i)
				aoter = append(aoter, i)
			}
		}
		for ao := len(aoter) -1; ao > -1; ao-- {
			res = oteSr(res, aoter[ao])
		}
	}
	if len(res) == 0 {
		return nil
	}

	// vérif lexicosyntaxique
	var aoter []int
	for _, ls := range g.lexsynt {
		for i, a := range res {
			if !lexsynt(a.Lem.Gr[0], ls) {
				aoter = append(aoter, i)
			}
		}
	}
	for ao := len(aoter) -1; ao > -1; ao-- {
		res = oteSr(res, aoter[ao])
	}
	if len(res) == 0 {
		return nil
	}
	// verif des exclusions lexicosyntaxiques
	aoter = nil
	for _, ls := range g.exclls {
		for i, a := range res {
			if lexsynt(a.Lem.Gr[0], ls) {
				aoter = append(aoter, i)
			}
		}
	}
	for ao := len(aoter) -1; ao > -1; ao-- {
		res = oteSr(res, aoter[ao])
	}
	if len(res) == 0 {
		return nil
	}

	// vérif morpho.
	// Si aucune n'est requise, renvoyer true
	if len(g.morph) == 0 {
		return res
	}

	aoter = nil
	for i, sr := range res {
		var morfos []string // morphos de sr acceptées par g
		for _, morf := range sr.Morphos {
			if g.vaMorph(morf) {
				morfos = append(morfos, morf)
			}
		}
		if len(morfos) == 0 {
			aoter = append(aoter, i)
			//res = oteSr(res, i)
		}
		sr.Morphos = morfos
	}
	for ao := len(aoter) -1; ao > -1; ao-- {
		res = oteSr(res, aoter[ao])
	}
	return res
}

// vrai si m est compatible avec Sub et le noyau mn
func (m *Mot) resSub(sub *Sub, mn *Mot, res gocol.Res) (vares gocol.Res) {
	// signet motresSub
	// si la fonction est déjà prise, renvoyer nil
	if mn.adeja(sub) {
		return nil
	}
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
			return nil
		}
		// noyaux possibles
		va := false
		for _, noy := range sub.noyaux {
			va = va || noy.vaPos(m.pos)
		}
		if !va {
			return nil
		}
	} else {
		// 2. La pos définitif n'est pas encore fixée
		var aoter []int
		for i, an := range res {
			// lexicosyntaxe
			va := true
			for _, ls := range sub.lexsynt {
				va = va && lexsynt(an.Lem.Gr[0], ls)
			}
			if !va {
				aoter = append(aoter, i)
			}
		}
		for i := len(aoter)-1; i > -1; i-- {
			res = oteSr(res, i)
		}
		if len(res) == 0 {
			return nil
		}

		// canon et POS
		aoter = nil
		va := false
		for i, an := range res {
			for _, noy := range sub.noyaux {
				if noy.canon > "" {
					va = va || noy.vaSr(an)
				} else {
					va = va || noy.vaPos(an.Lem.Pos)
				}
			}
			if !va {
				aoter = append(aoter, i)
				//res = oteSr(res, i)
			}
		}
		for i := len(aoter)-1; i > -1; i-- {
			res = oteSr(res, i)
		}
	}
	if len(res) == 0 {
		return nil
	}

	//morphologie
	// si aucune morpho n'est requise, passer
	if len(sub.morpho) > 0 {
		var aoter []int
		for i, an := range res {
			var lmorf []string
			for _, morfs := range an.Morphos {
				// pour toutes les morphos valides de m
				if strings.Contains(morfs, "inv.") || sub.vaMorpho(morfs) {
					lmorf = append(lmorf, morfs)
				}
			}
			if len(lmorf) == 0 {
				aoter = append(aoter, i)
			} else {
				res[i].Morphos = lmorf
			}
		}
		for i := len(aoter)-1; i > -1; i-- {
			res = oteSr(res, i)
		}
	}
	// accord
	// pour toutes les morphos valides de mn
	if sub.accord > "" {
		var aoter []int
		for i, an := range res {
			va := false
			for _, anoy := range mn.ans2 {
				// pour toutes les morphos valides de m
				var lmorf []string
				for _, morfn := range anoy.Morphos {
					for _, morfs := range an.Morphos {
						if accord(morfn, morfs, sub.accord) {
							lmorf = append(lmorf, morfs)
							va = true
						}
					}
				}
				if len(lmorf) > 0 {
					an.Morphos = lmorf
					// XXX à vérifier : à placer + bas ?
					res[i] = an
				}
			}
			if !va {
				aoter = append(aoter, i)
			}
		}
		for i := len(aoter)-1; i > -1; i-- {
			res = oteSr(res, i)
		}
	}
	return res
}
