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
	idGr    string    // nom du groupe dont le mot est noyau
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
		br.photos[m.rang] = phm
		m.restmp = m.ans
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
	nb.photos = make(map[int]*PhotoMot)
	nb.photos = b.photos
	nb.veto = b.veto
	// créer une lemmatisation temporaire pour chaque mot
	for _, m := range mots {
		m.restmp = nb.photos[m.rang].res
	}
	return nb
}

func (b *Branche) dejasub(m *Mot) bool {
	for _, n := range b.nods {
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
	}
	return false
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
		n := bm.noeud(m, g)
		if n != nil {
			// Si le groupe a été exploré pour m dans une
			// autre branche, passer
			va := true
			for _, veto := range bm.veto[m.rang] {
				va = va && !n.egale(veto)
				if !va {
					break
				}
			}
			if !va {
				continue
			}
			// le noeud est accepté. créer une branche fille (bf)
			bf := bm.copie()
			bf.nods = append(bf.nods, n)
			for _, mph := range mots {
				if mph == m {
					ph := new(PhotoMot)
					ph.res = m.restmp
					ph.idGr = n.grp.id
					bf.photos[m.rang] = ph
					mph.restmp = ph.res
				}
				for _, ma := range n.mma {
					if mph == ma {
						ph := new(PhotoMot)
						ph.res = ma.restmp
						ph.idGr = bm.photos[ma.rang].idGr
						bf.photos[ma.rang] = ph
						mph.restmp = ph.res
					}
				}
				for _, mp := range n.mmp {
					if mph == mp {
						ph := new(PhotoMot)
						ph.res = mp.restmp
						ph.idGr = bm.photos[mp.rang].idGr
						bf.photos[mp.rang] = ph
						mph.restmp = ph.res
					}
				}
			}
			bm.filles = append(bm.filles, bf)
			bf.explore()
			bm.veto[m.rang] = append(bm.veto[m.rang], n)
			// rétablir toutes les lemmatisations temporaires
			for _, m := range mots {
				m.restmp = bm.photos[m.rang].res
			}
		}
	}
}

func (bm *Branche) explore() {
	// signet sexplore
	for _, m := range mots {
		// 1. groupes terminaux
		bm.exploreGroupes(m, grpTerm)
		// 2. groupes non terminaux
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

// id du Nod dont m est déjà le noyau
func (b *Branche) id(m *Mot) string {
	for _, nod := range b.nods {
		if nod.nucl.rang == m.rang {
			return nod.grp.id
		}
	}
	return ""
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
	//photo := b.photos[m.rang]
	//m.restmp = photo.res
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
		//resma := b.photos[ma.rang].res
		resma := ma.restmp
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
		resmp := mp.restmp
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
	m.restmp = b.photos[m.rang].res
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
		pos     string    // nom du groupe dont le mot est noyau
	*/
	photom := b.photos[m.rang]
	// vérif du pos
	if photom.idGr != "" {
		// 1. La pos définitif est fixée
		// noyaux exclus
		id := b.id(m)
		if g.estExclu(id) {
			return nil
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
						if noy.vaPos(photom.idGr) {
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
			if g.vaPos(a.Lem.Pos) {
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
func (b *Branche) resSub(m *Mot, sub *Sub, mn *Mot, res gocol.Res) (gocol.Res) {
	// signet sresub
	// si la fonction est déjà prise, renvoyer nil
	if !sub.multi && b.adeja(mn, sub.lien) {
		return nil
	}

	// photo m et mn pour la branche
	photom := b.photos[m.rang]
	// vérification du pos : id du noyau, ou pos du mot
	if photom.idGr != "" {
		// 1. La pos du mot est définitive
		id := b.id(m)
		if !sub.vaId(id) {
			return nil
		}
		va := false
		for _, noy := range sub.noyaux {
			va = va || noy.vaPos(photom.idGr)
		}
		if !va {
			return nil
		}
	} else {
		// 2. m n'est pas encore noyau : on vérifie lexicosyntaxe canon et pos 
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
