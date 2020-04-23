//     Branche.go - Publicola

package main

import (
	"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

type Branche struct {
	gr     string
	imot   int
	nbmots int
	mots   []*Mot
	nods   []*Nod
	niveau int
	photos map[*Mot]*PhotoMot
	mere   *Branche
	filles []*Branche
}

func creeBranche(t string) *Branche {
	p := new(Branche)
	p.gr = t
	mm := gocol.Mots(t)
	for i, m := range mm {
		nm := creeMot(m)
		nm.rang = i
		p.mots = append(p.mots, nm)
	}
	p.nbmots = len(p.mots)
	p.photos = make(map[*Mot]*PhotoMot)
	return p
}

func (b *Branche) copie() *Branche {
	nb := new(Branche)
	nb.gr = b.gr
	nb.nbmots = b.nbmots
	for _, am := range b.mots {
		nm := am.copie()
		nb.mots = append(nb.mots, nm)
	}
	copy(nb.nods, b.nods)
	nb.mere = b
	nb.niveau = b.niveau + 1
	copy(b.filles, nb.filles)
	nb.photos = make(map[*Mot]*PhotoMot)
	return nb
}

func (b *Branche) dejasub(m *Mot) bool {
	return b.photos[m].dejasub
}

func (b *Branche) domine(ma, mb *Mot) bool {
//func (b *Branche) noyau(m *Mot) *Mot {
	mnoy := b.noyau(mb)
	for mnoy != nil {
		if mnoy == ma {
			return true
		}
		mnoy = b.noyau(mnoy)
	}
	return false
}

// FIXME : la photo de dea devrait avoir un pos
// et la photo de Isis devrait être dejasub
func (bm *Branche) explore() {
	// on copie la branche mère pour la rendre
	// indépendante et en faire une fille possible
	bf := bm.copie()
	// recherche des noyaux
	// groupes terminaux
	for _, g := range grpTerm {
		for _, m := range bf.mots {
			if m.dejaNoy() {
				continue
			}
			// XXX b.photos[m] n'existe pas
			n := bm.noeud(m, g)
			if n != nil {
				bf.nods = append(bf.nods, n)
				nbf := bf.copie()
				// calcul des photos pour nbf
				for _, mph := range bm.mots {
					photo := new(PhotoMot)
					photo.mot = mph
					if n.inclut(mph) {
						photo.res = mph.restmp
						if mph == n.nucl {
							photo.pos = n.grp.id
						}
					} else {
						photo.res = bm.photos[mph].res
					}
				}
				bm.filles = append(bm.filles, nbf)
			}
		}
	}
	// groupes non terminaux
	for _, g := range grp {
		for _, m := range bf.mots {
			n := bm.noeud(m, g)
			if n != nil {
				bf.nods = append(bf.nods, n)
				nbf := bf.copie()
				nbf.mere = bm
				bm.filles = append(bm.filles, nbf)
			}
		}
	}
	if len(bm.filles) > 0 {
		for _, f := range bm.filles {
			f.explore()
		}
	}
}

// texte de la Branche, le mot courant surligné en rouge
func (b *Branche) enClair() string {
	var lm []string
	for i := 0; i < len(b.mots); i++ {
		m := b.mots[i].gr
		if i == b.imot {
			m = rouge(m)
		}
		lm = append(lm, m)
	}
	return strings.Join(lm, " ") + "."
}

// affiche la Branche en colorant n mots à partir
// du mot n° d
func (b *Branche) exr(d, n int) (e string) {
	var gab string = "%s"
	for i := 0; i < len(b.mots); i++ {
		if e != "" {
			gab = " %s"
		}
		if i >= d && i < d+n {
			e += fmt.Sprintf(gab, rouge(b.mots[i].gr))
		} else {
			e += fmt.Sprintf(gab, b.mots[i].gr)
		}
	}
	return
}

func (b *Branche) initTronc() {
	for _, m := range b.mots {
		p := new(PhotoMot)
		p.mot = m
		p.res = m.ans
		p.dejasub = false
		b.photos[m] = p
	}
}

func (b *Branche) motCourant() *Mot {
	return b.mots[b.imot]
}

// si m peut être noyau d'un gourpe g, un Nod est renvoyé, sinon nil.
func (b *Branche) noeud(m *Mot, g *Groupe) *Nod {
	// signet motnoeud

	// vérification de rang
	rang := m.rang
	lante := len(g.ante)
	// mot de rang trop faible
	if rang-lante < 0 {
		return nil
	}
	// ou trop élevé
	if rang+len(g.post)-1 >= texte.tronc.nbmots {
		return nil
	}

	// m peut-il être noyau du groupe g ?
	m.restmp = cloneRes(m.ans)
	res := b.resNoyau(m, g, m.restmp)
	if res == nil {
		return nil
	}
	res = cloneRes(m.restmp)

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
		ma := texte.tronc.mots[r]
		// passer les mots déjà subordonnés
		for b.dejasub(ma) {
			r--
			if r < 0 {
				return nil
			}
			ma = texte.tronc.mots[r]
		}
		// vérification de réciprocité, puis du lien lui-même
		if b.domine(ma, m) {
			return nil
		}
		sub := g.ante[ia]
		res := b.resSub(ma, sub, m, ma.restmp)
		if res == nil {
			return nil
		}
		ma.restmp = cloneRes(res)
		nod.mma = append(nod.mma, ma)
		r--
	}

	// reгcherche des subs post
	for ip, sub := range g.post {
		r := rang + ip + 1
		if r >= texte.tronc.nbmots {
			break
		}
		if sub.lien == "" {
			continue
		}
		mp := texte.tronc.mots[r]
		for b.dejasub(mp) {
			r++
			if r >= texte.tronc.nbmots {
				return nil
			}
			mpn := b.noyau(mp)
			if mpn != nil && mpn.rang < m.rang {
				return nil
			}
			mp = texte.tronc.mots[r]
		}
		// réciprocité
		if b.domine(mp, m) {
			return nil
		}
		mp.restmp = cloneRes(mp.ans)
		res := b.resSub(mp, sub, m, mp.restmp)
		if res == nil {
			return nil
		}
		mp.restmp = cloneRes(res)
		nod.mmp = append(nod.mmp, mp)
		r++
	}

	// fixer les pos et sub des mots du noeud
	if len(nod.mma)+len(nod.mmp) > 0 {
		// la pos du noyau devient celle du groupe
		photo := b.photos[m]
		photo.pos = g.id
		b.photos[m] = photo
		// restriction des lemmatisations des antéposés
		for _, ms := range nod.mma {
			photo := b.photos[ms]
			photo.dejasub = true
			photo.res = cloneRes(ms.restmp)
			ms.restmp = nil
		}
		//restriction des lemmatisations du noyau
		m.ans = m.restmp
		m.restmp = nil
		// restriction des lemmatisations des postposés
		for _, ms := range nod.mmp {
			photo := b.photos[ms]
			photo.dejasub = true
			photo.res = cloneRes(ms.restmp)
			ms.restmp = nil
		}
		return nod
	}
	return nil
}

func (b *Branche) noyau(m *Mot) *Mot {
	for _, n := range texte.tronc.nods {
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
func (b *Branche) resNoyau(m *Mot, g *Groupe, res gocol.Res) gocol.Res {
	// valeurs variable de m pour la branche
	photom := b.photos[m]
	// vérif du pos
	if photom.pos != "" {
		// 1. La pos définitif est fixée
		va := false
		for _, noy := range g.noyaux {
			if noy.canon > "" {
				for _, a := range res {
					va = va || noy.vaSr(a)
				}
			} else {
				va = va || noy.vaPos(photom.pos)
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
func (b *Branche) resSub(m *Mot, sub *Sub, mn *Mot, res gocol.Res) (vares gocol.Res) {
	// signet motresSub
	// si la fonction est déjà prise, renvoyer nil
	if !sub.multi && mn.adeja(sub) {
		return nil
	}

	// photo m et mn pour la branche
	photom := b.photos[m]
	// vérification des pos
	// FIXME legatos decernis : avec v.obj, seul legagos pp est sélectionné par vaPos
	if photom.pos != "" {
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
			va = va || noy.vaPos(photom.pos)
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
