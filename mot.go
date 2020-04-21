//       mot.go - Gentes

// signets :
//
// motnoeud
// motresnoyau
// motestNoyauDeGroupe
// motresSub
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
	ans, ans2  gocol.Res // ensemble des lemmatisations, ans2 réduit par chaque noeud créé.
	restmp     gocol.Res // analyses temporaires du mot pendand le calcul d'un noeud
	dejasub    bool      // le mot est déjà l'élément d'n nœud
	llm        []Lm      // liste des lemmes ٍ+ morpho possibles
	tmpl, tmpm int       // n°s provisoires de Sr et morpho
	pos        string    // id du groupe dont le mot est noyau
	// ou à défaut pos du mot, si elle est décidée
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

	// provisoire XXX
	// exclusions de mots rares faisant obstacle à des analyses importantes
	var nres gocol.Res
	for _, an := range mot.ans {
		if !lexsynt(an.Lem, "excl") {
			nres = append(nres, an)
		}
	}
	mot.ans2 = nres
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

func (m *Mot) adeja(sub *Sub) bool {
	// signet motadeja
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

	// vérification de rang
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
	m.restmp = res

	// création du noeud de retour
	nod := new(Nod)
	nod.grp = g
	nod.nucl = m
	nod.rang = rang

	// reгcherche rétrograde des subs ante
	r := rang - 1
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
		res := ma.resSub(sub, m, ma.restmp)
		if res == nil {
			return nil
		}
		ma.restmp = res
		nod.mma = append(nod.mma, ma)
		r--
	}

	// reгcherche des subs post
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
		// la pos du noyau devient celle du groupe
		m.pos = g.id
		// restriction des lemmatisations des antéposés
		for _, ms := range nod.mma {
			ms.dejasub = true
			ms.ans2 = ms.restmp
			ms.restmp = nil
		}
		//restriction des lemmatisations du noyau
		m.ans2 = m.restmp
		m.restmp = nil
		// restriction des lemmatisations des postposés
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

// renvoie quelles lemmatisations de m lui permettent d'être le noyau du groupe g
func (m *Mot) resNoyau(g *Groupe, res gocol.Res) gocol.Res {
	//signet motresnoyau
	// vérif du pos
	if m.pos != "" {
		// 1. La pos définitif est fixée
		va := false
		for _, noy := range g.noyaux {
			if noy.canon > "" {
				for _, a := range res {
					va = va || noy.vaSr(a)
				}
			} else {
				va = va || noy.vaPos(m.pos)
			}
		}
		if !va {
			return nil
		}
	} else {
		// Le mot est encore isolé
		var nres gocol.Res
		for _, a := range res {
			va := false
			for _, noy := range g.noyaux {
				if noy.canon > "" {
					va = va || noy.vaSr(a)
				} else {
					va = va || noy.vaPos(a.Lem.Pos)
				}
			}
			if va {
				nres = append(nres, a)
			}
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres
	}

	// vérif lexicosyntaxique
	var nres gocol.Res
	for _, a := range res {
		va := true
		for _, ls := range g.lexsynt {
			va = va && lexsynt(a.Lem, ls)
		}
		if va {
			nres = append(nres, a)
		}
	}
	if len(nres) == 0 {
		return nil
	}
	res = nres

	// verif des exclusions lexicosyntaxiques
	nres = nil
	for _, a := range res {
		va := true
		for _, ls := range g.exclls {
			va = va && !lexsynt(a.Lem, ls)
		}
		if va {
			nres = append(nres, a)
		}
	}
	if len(nres) == 0 {
		return nil
	}
	res = nres

	// vérif morpho.
	// Si aucune n'est requise, renvoyer true
	if len(g.morph) == 0 {
		return res
	}

	nres = nil
	for _, sr := range res {
		var morfos []string // morphos de sr acceptées par g
		for _, morf := range sr.Morphos {
			if g.vaMorph(morf) {
				morfos = append(morfos, morf)
			}
		}
		if len(morfos) > 0 {
			sr.Morphos = morfos
			nres = append(nres, sr)
		}
	}
	// pour faire comme pour les autres vérifs :
	res = nres
	return res
}

// vrai si m est compatible avec Sub et le noyau mn
func (m *Mot) resSub(sub *Sub, mn *Mot, res gocol.Res) (vares gocol.Res) {
	// signet motresSub
	// si la fonction est déjà prise, renvoyer nil
	if !sub.multi && mn.adeja(sub) {
		return nil
	}
	// vérification des pos
	// FIXME legatos decernis : avec v.obj, seul legagos pp est sélectionné par vaPos
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
		var nres gocol.Res
		// lexicosyntaxe
		for _, an := range res {
			va := true
			for _, ls := range sub.lexsynt {
				va = va && lexsynt(an.Lem, ls)
			}
			if va {
				nres = append(nres, an)
			}
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres

		// canon et POS
		nres = nil
		for _, an := range res {
			va := false
			for _, noy := range sub.noyaux {
				if noy.canon > "" {
					va = va || noy.vaSr(an)
				} else {
					va = va || noy.vaPos(an.Lem.Pos)
				}
			}
			if va {
				nres = append(nres, an)
			}
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres
	}

	//morphologie
	// si aucune morpho n'est requise, passer
	if len(sub.morpho) > 0 {
		var nres gocol.Res
		for _, an := range res {
			var lmorf []string
			for _, morfs := range an.Morphos {
				// pour toutes les morphos valides de m
				if strings.Contains(morfs, "inv.") || sub.vaMorpho(morfs) {
					lmorf = append(lmorf, morfs)
				}
			}
			if len(lmorf) > 0 {
				an.Morphos = lmorf
				nres = append(nres, an)
			}
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres
	}

	// accord
	// pour chaque an.
	if sub.accord > "" {
		var nres gocol.Res
		for _, an := range res {
			va := false
			for _, anoy := range mn.restmp {
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
				}
			}
			if va {
				nres = append(nres, an)
			}
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres
	}
	return res
}
