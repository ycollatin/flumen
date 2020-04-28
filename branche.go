//     Branche.go - Gentes

/*
Signets
exploregrou
scopie
sexplore
snoeud
snoyau
srecolte
sresub

*/

package main

import (
	"fmt"
	"github.com/ycollatin/gocol"
	"sort"
	"strings"
)

// Chaque branche modifie trois propriétés attachées
// au mot. Leur liste est donc enregistrée dans
// Branche.photos.
type PhotoMot struct {
	res     gocol.Res // lemmatisations réduites du mot
	dejasub bool      // appartenance du mot à un groupe
	pos     string    // nom du groupe dont le mot est noyau
}

var (
	mots   []*Mot
	nbmots int
)

type Branche struct {
	gr     string            // texte de la phrase
	imot   int               // rang du mot courant
	nods   []*Nod            // noeuds validés
	niveau int               // n° de la branche par rapport au tronc
	veto   map[int][]*Nod    // index : rang du mot; valeur : liste des liens interdits
	photos map[int]*PhotoMot // lemmatisations et appartenance de groupe propres à la branche
	filles []*Branche        // liste des branches filles
}

// le tronc est la branche de départ. Il
// est initialisé avec les mots lemmatisés
// par Collatinus
func creeTronc(t string) *Branche {
	mots = nil
	br := new(Branche)
	br.gr = t
	mm := gocol.Mots(t)
	for i, m := range mm {
		nm := creeMot(m)
		nm.rang = i
		mots = append(mots, nm)
	}
	nbmots = len(mots)
	br.photos = make(map[int]*PhotoMot) // l'index de la map est le numéro des mots
	br.veto = make(map[int][]*Nod)
	// peuplement des photos
	for _, m := range mots {
		phm := new(PhotoMot)
		phm.res = m.ans
		phm.dejasub = false
		br.photos[m.rang] = phm
	}
	return br
}

func (b *Branche) adeja(noyau *Mot, lien string) bool {
	for _, nod := range b.nods {
		if nod.nucl == noyau {
			for i, _ := range nod.mma {
				if nod.grp.ante[i].lien == lien {
					return true
				}
			}
			for i, _ := range nod.mmp {
				if nod.grp.post[i].lien == lien {
					return true
				}
			}
		}
	}
	return false
}

func (b *Branche) copie() *Branche {
	// signet scopie
	nb := new(Branche)
	nb.gr = b.gr
	nb.niveau = b.niveau + 1
	nb.nods = b.nods
	// les photos seront copiées après création
	// du noeud à l'origine de la copie
	nb.photos = make(map[int]*PhotoMot)
	nb.veto = b.veto
	return nb
}

func (b *Branche) dejasub(m *Mot) bool {
	return b.photos[m.rang].dejasub
}

func (b *Branche) domine(ma, mb *Mot) bool {
	mnoy := b.noyau(mb)
	for mnoy != nil {
		if mnoy == ma {
			return true
		}
		mnoy = b.noyau(mnoy)
	}
	return false
}

func (bm *Branche) exploreGroupes(m *Mot, grps []*Groupe) {
	// signet exploregrou
	for _, g := range grps {
		// Si le groupe a été exploré pour m dans une
		// autre branche, passer
		cont := false
		for _, lv := range bm.veto[m.rang] {
			// FIXME : cas des subs multi
			if m == lv.nucl && lv.grp.id == g.id && !lv.grp.multi() {
				cont = true
				break
			}
		}
		if cont {
			continue
		}
		n := bm.noeud(m, g)
		if n != nil {
			// le noeud est accepté. créer une branche fille (bf)
			bf := bm.copie()
			for _, mph := range mots {
				vu := false
				if mph == n.nucl {
					ph := new(PhotoMot)
					ph.res = mph.restmp
					ph.pos = n.grp.id
					bf.photos[mph.rang] = ph
					// interdire le groupe au noyau
					bm.veto[mph.rang] = append(bm.veto[mph.rang], n)
					vu = true
				}
				for _, ma := range n.mma {
					if mph == ma {
						ph := new(PhotoMot)
						ph.res = ma.restmp
						ph.dejasub = true
						ph.pos = bm.photos[ma.rang].pos
						bf.photos[mph.rang] = ph
						vu = true
					}
				}
				for _, mp := range n.mmp {
					if mph == mp {
						ph := new(PhotoMot)
						ph.res = mp.restmp
						ph.dejasub = true
						ph.pos = bm.photos[mp.rang].pos
						bf.photos[mph.rang] = ph
						vu = true
					}
				}
				if !vu {
					bf.photos[mph.rang] = bm.photos[mph.rang]
				}
			}
			bf.nods = append(bf.nods, n)
			bf.explore()
			bm.filles = append(bm.filles, bf)
		}
	}
}

func (bm *Branche) explore() {
	// signet sexplore
	// 1. groupes terminaux
	for _, m := range mots {
		bm.exploreGroupes(m, grpTerm)
	}
	// 2. groupes non terminaux
	for _, m := range mots {
		bm.exploreGroupes(m, grp)
	}
}

// texte de la Branche, le mot courant surligné en rouge
func (b *Branche) enClair() string {
	var lm []string
	for i := 0; i < len(mots); i++ {
		m := mots[i].gr
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
	for i := 0; i < len(mots); i++ {
		if e != "" {
			gab = " %s"
		}
		if i >= d && i < d+n {
			e += fmt.Sprintf(gab, rouge(mots[i].gr))
		} else {
			e += fmt.Sprintf(gab, mots[i].gr)
		}
	}
	return
}

// id des Nod dont m est déjà le noyau
func (b *Branche) ids(m *Mot) []string {
	var ret []string
	for _, nod := range b.nods {
		if nod.nucl.rang == m.rang {
			ret = append(ret, nod.grp.id)
		}
	}
	return ret
}

func (b *Branche) motCourant() *Mot {
	return mots[b.imot]
}

// si m peut être noyau d'un gourpe g, un Nod est renvoyé, sinon nil.
func (b *Branche) noeud(m *Mot, g *Groupe) *Nod {
	// signet snoeud

	// utilistation des photos
	//mot     *Mot      // liaison avec le mot
	//res     gocol.Res // lemmatisations réduites du mot
	//dejasub bool      // appartenance du mot à un groupe
	//pos     string    // nom du groupe dont le mot est noyau

	// vérification de rang
	rang := m.rang
	lante := len(g.ante)
	// mot de rang trop faible
	if rang-lante < 0 {
		return nil
	}
	// ou trop élevé
	if rang+len(g.post)-1 >= nbmots {
		return nil
	}

	// m peut-il être noyau du groupe g ?
	photo := b.photos[m.rang]
	m.restmp = photo.res
	res := b.resNoyau(m, g, m.restmp)
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
		ma := mots[r]
		// passer les mots déjà subordonnés
		for b.dejasub(ma) {
			r--
			if r < 0 {
				return nil
			}
			ma = mots[r]
		}
		// vérification de réciprocité, puis du lien lui-même
		if b.domine(ma, m) {
			return nil
		}
		sub := g.ante[ia]
		resma := b.photos[ma.rang].res
		resma = b.resSub(ma, sub, m, resma)
		if resma == nil {
			return nil
		}
		ma.restmp = resma
		nod.mma = append(nod.mma, ma)
		r--
	}

	// reгcherche des subs post
	r = rang + 1
	for _, sub := range g.post {
		if r >= nbmots {
			break
		}
		if sub.lien == "" {
			continue
		}
		mp := mots[r]
		for b.dejasub(mp) {
			r++
			if r >= nbmots {
				return nil
			}
			mpn := b.noyau(mp)
			if mpn != nil && mpn.rang < m.rang {
				return nil
			}
			mp = mots[r]
		}
		// réciprocité
		if b.domine(mp, m) {
			return nil
		}
		resmp := b.photos[mp.rang].res
		resmp = b.resSub(mp, sub, m, resmp)
		if resmp == nil {
			return nil
		}
		mp.restmp = resmp
		nod.mmp = append(nod.mmp, mp)
		r++
	}

	if len(nod.mma)+len(nod.mmp) > 0 {
		return nod
	}
	return nil
}

func (b *Branche) noyau(m *Mot) *Mot {
	for _, n := range b.nods {
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
	// signet snoyau
	// valeurs variables de m pour la branche
	/*
		// utilistation des photos
		mot     *Mot      // liaison avec le mot
		res     gocol.Res // lemmatisations réduites du mot
		dejasub bool      // appartenance du mot à un groupe
		pos     string    // nom du groupe dont le mot est noyau
	*/
	photom := b.photos[m.rang]
	// vérif du pos
	if photom.pos != "" {
		// 1. La pos définitif est fixée
		// noyaux exclus
		ids := b.ids(m)
		for _, id := range ids {
			if g.estExclu(id) {
				return nil
			}
		}
		var nres gocol.Res
		// noyaux admis
		for _, noy := range g.noyaux {
			if noy.canon > "" {
				for _, a := range res {
					for _, morf := range a.Morphos {
						if noy.vaSr(a) {
							nres = gocol.AddRes(nres, a.Lem, morf, 0)
						}
					}
				}
			} else {
				for _, a := range res {
					for _, morf := range a.Morphos {
						if noy.vaPos(photom.pos) {
							nres = gocol.AddRes(nres, a.Lem, morf, 0)
						}
					}
				}
			}
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres
	} else { // Le mot est encore isolé
		// vérif des pos
		var nres gocol.Res
		for _, a := range res {
			for _, noy := range g.noyaux {
				if noy.canon > "" && noy.vaSr(a) {
					nres = append(nres, a)
				} else {
					if noy.vaPos(a.Lem.Pos) {
						nres = append(nres, a)
					}
				}
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

	// vérif morpho.
	// Si aucune n'est requise, renvoyer true
	if len(g.morph) == 0 {
		return res
	}

	nres = nil
	for _, sr := range res {
		for _, morf := range sr.Morphos {
			if g.vaMorph(morf) {
				nres = gocol.AddRes(nres, sr.Lem, morf, 0)
			}
		}
	}
	// pour faire comme pour les autres vérifs :
	res = nres
	return res
}

// récolte tous les noeuds terminaux d'un arbre
func (b *Branche) recolte() (rec [][]*Nod) {
	// signet srecolte
	if len(b.filles) == 0 {
		rec = append(rec, b.nods)
		return rec
	}
	for _, f := range b.filles {
		nrec := f.recolte()
		rec = append(rec, nrec...)
	}
	sort.Slice(rec, func(i, j int) bool {
		return len(rec[i]) > len(rec[j])
	})
	return rec
}

// vrai si m est compatible avec Sub et le noyau mn
func (b *Branche) resSub(m *Mot, sub *Sub, mn *Mot, res gocol.Res) (vares gocol.Res) {
	// signet sresub
	// si la fonction est déjà prise, renvoyer nil
	if !sub.multi && b.adeja(mn, sub.lien) {
		return nil
	}

	// photo m et mn pour la branche
	photom := b.photos[m.rang]
	// vérification des pos
	if photom.pos != "" {
		// 1. La pos du mot est définitive
		// noyaux exclus
		excl := false
		lgr := b.ids(m)
		for _, noy := range sub.noyexcl {
			excl = excl || contient(lgr, noy.id)
		}
		if excl {
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
			for _, morfs := range an.Morphos {
				// pour toutes les morphos valides de m
				if strings.Contains(morfs, "inv.") || sub.vaMorpho(morfs) {
					nres = gocol.AddRes(nres, an.Lem, morfs, 0)
				}
			}
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres
	}

	// accord
	if sub.accord > "" {
		var nres gocol.Res
		for _, an := range res {
			for _, anoy := range mn.restmp {
				for _, morfs := range an.Morphos {
					for _, morfn := range anoy.Morphos {
						if accord(morfn, morfs, sub.accord) {
							nres = gocol.AddRes(nres, anoy.Lem, morfs, 0)
						}
					}
				}
			}
		}
		if len(nres) == 0 {
			return nil
		}
	}
	return res
}
