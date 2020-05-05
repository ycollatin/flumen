//     Branche.go - Gentes

/*
Signets
exploregrou
scopie
sexplore
snoeud
srecolte
sresEl
*/

package main

import (
	"fmt"
	"github.com/ycollatin/gocol"
	"sort"
	"strings"
)

var (
	mots    []*Mot
	nbmots  int
	journal []string
)

type Branche struct {
	gr     string            // texte de la phrase
	imot   int               // rang du mot courant
	nods   []*Nod            // noeuds validés
	niveau int               // n° de la branche par rapport au tronc
	veto   map[int][]*Nod    // index : rang du mot; valeur : liste des liens interdits
	photos map[int]gocol.Res // lemmatisations et appartenance de groupe propres à la branche
	filles []*Branche        // liste des branches filles
}

// le tronc est la branche de départ. Il
// est initialisé avec les mots lemmatisés
// par Collatinus. var tronc est déclarée dans texte.go
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
	br.photos = make(map[int]gocol.Res) // l'index de la map est le numéro des mots
	br.veto = make(map[int][]*Nod)
	// peuplement des photos
	for _, m := range mots {
		phm := m.ans
		br.photos[m.rang] = phm
	}
	journal = nil
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
	nb.photos = make(map[int]gocol.Res)
	nb.photos = b.photos
	nb.veto = make(map[int][]*Nod)
	nb.veto = b.veto
	return nb
}

func (b *Branche) copieRestmp() {
	for _, m := range mots {
	    m.restmp = b.photos[m.rang]
	}
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

// explore toutes les possibilités d'une branche
func (bm *Branche) explore() {
	// signet sexplore
	for _, m := range mots {
		photo := bm.photos[m.rang]
		// 1. groupes terminaux
		bm.exploreGroupes(m, grpTerm)
		// 2. groupes non terminaux
		bm.exploreGroupes(m, grp)
		bm.photos[m.rang] = photo
	}
}

// essaye toutes les règles de groupes de grps où m pourrait
// être noyau
func (bm *Branche) exploreGroupes(m *Mot, grps []*Groupe) {
	// signet exploregrou
	for _, g := range grps {
		// tester la possibilité de création noeud de type g
		// dont le noyau est m
		n := bm.noeud(m, g)
		if n != nil {
			// Si le groupe a été exploré pour m dans une
			// autre branche, passer
			// XXX Pourrait être fait dans noeud(m, g) avant calcul.
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
			//créer une branche fille (bf)
			// copiée de la mère (bm)
			// où le noeud sera obligatoire. Dans la branche mère,
			// il sera interdit.
			bf := bm.copie()
			bf.nods = append(bf.nods, n)
			for _, mph := range mots {
				if mph == m {
					// photo du noyau
					bf.photos[mph.rang] = m.restmp
					n.rnucl = mph.restmp
				}
				for _, ma := range n.mma {
					// photos des éléments antéposés
					if mph == ma {
						bf.photos[ma.rang] = ma.restmp
						n.rra[ma.rang] = mph.restmp
					}
				}
				for _, mp := range n.mmp {
					// photos des éléments postposés
					if mph == mp {
						bf.photos[mp.rang] = mp.restmp
						n.rrp[mp.rang] = mph.restmp
					}
				}
			}
			// màj journal
			indent := strings.Repeat("  ", bm.niveau)
			journal = append(journal, fmt.Sprintf("%s %d. %s", indent, bm.niveau, n.doc()))
			// la fille est explorée récursivement
			journal = append(journal, fmt.Sprintf("%s %d montée au niveau %d", indent, bm.niveau, bf.niveau))
			bf.explore()
			if len(bf.filles) == 0 {
				journal = append(journal, fmt.Sprintf("%s   %d branche terminale, %d arc(s)",
					indent, bf.niveau, len(bf.nods)))
			}
			journal = append(journal, fmt.Sprintf("%s retour niveau %d. niveau %d, %d fille(s)",
				indent, bm.niveau, bf.niveau, len(bf.filles)))
			// ajout à la liste des filles
			bm.filles = append(bm.filles, bf)
			// le noeud est désormais interdit chez
			// la mère et ses futures filles
			bm.veto[m.rang] = append(bm.veto[m.rang], n)
		}
	}
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

func (b *Branche) idgr(m *Mot) (id string) {
	var max int
	for _, nod := range b.nods {
		nbe := nod.nbEl()
		if nod.nucl == m && nbe > max {
			id = nod.grp.id
		}
	}
	return
}

/*
// id du Nod dont m est déjà le noyau
func (b *Branche) ids(m *Mot) (ii []string) {
	for _, nod := range b.nods {
		if nod.nucl.rang == m.rang {
			ii = append(ii, nod.grp.id)
		}
	}
	return
}
*/

func (b *Branche) motCourant() *Mot {
	return mots[b.imot]
}

// si m peut être noyau d'un gourpe g, un Nod est renvoyé, sinon nil.
func (b *Branche) noeud(m *Mot, g *Groupe) *Nod {
	// signet snoeud

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
	m.restmp = b.photos[m.rang]

	m.restmp = b.resEl(m, g.nucl, m, m.restmp)
	if m.restmp == nil {
		return nil
	}

	// création du noeud de retour
	nod := new(Nod)
	nod.rra = make(map[int]gocol.Res)
	nod.rrp = make(map[int]gocol.Res)
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
		// passer les mots déjà subordonnés | b branche.go:336 cond 2 m.gr=="rapuit" && g.id=="v.conjet"
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
		ma.restmp = b.photos[ma.rang]
		ma.restmp = b.resEl(ma, sub, m, ma.restmp)
		if ma.restmp == nil {
			return nil
		}
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
		mp.restmp = b.photos[mp.rang]
		mp.restmp = b.resEl(mp, sub, m, mp.restmp)
		if mp.restmp == nil {
			return nil
		}
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

// récolte tous les noeuds terminaux d'un arbre
func (b *Branche) recolte() (rec [][]*Nod) {
	// signet srecolte
	if b.terminale() {
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
func (b *Branche) resEl(m *Mot, el *El, mn *Mot, res gocol.Res) gocol.Res {
	// signet sresEl
	// si la fonction est déjà prise, renvoyer nil
	if !el.multi && b.adeja(mn, el.lien) {
		return nil
	}

	// vérification du pos : id du noyau, ou pos du mot
	id := b.idgr(m)
	if id > "" {
		// familles
		pel := PrimEl(id, ".")
		if contient(el.famexcl, pel) {
			return nil
		}
		if contient(el.idsexcl, id) {
			return nil
		}
	}
	var nres gocol.Res
	// 2. m n'est pas encore noyau : on vérifie lexicosyntaxe canon et pos
	// lexicosyntaxe, exclus
	if len(el.lsexcl) > 0 {
		nres = nil
		for _, excl := range el.lsexcl {
			for _, an := range res {
				if !lexsynt(an.Lem, excl) {
					nres = append(nres, an)
				}
			}
		}
		res = nres
	}

	// possibles
	if len(el.lexsynt) > 0 {
		nres = nil
		for _, an := range res {
			va := true
			for _, ls := range el.lexsynt {
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
	}


	// canons
	if len(el.cles) > 0 {
		nres = nil
		for _, an := range res {
			if contient(el.clesexcl, an.Lem.Cle) {
				return nil
			}
			for _, cle := range el.cles {
				if an.Lem.Cle == cle {
					nres = append(nres, an)
				}
			}
			if len(nres) == 0 {
				return nil
			}
			res = nres
		}
	}

	if len(el.poss) > 0 {
		// pos
		nres = nil
		for  _, an := range res {
			if contient(el.posexcl, an.Lem.Pos) {
				continue
			}
			if contient(el.poss, an.Lem.Pos) {
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
	if len(el.morpho) > 0 {
		var nres gocol.Res
		for _, an := range res {
			for _, morfs := range an.Morphos {
				// pour toutes les morphos valides de m
				if strings.Contains(morfs, "inv.") || el.vaMorpho(morfs) {
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
	if el.accord > "" {
		var nres gocol.Res
		for _, an := range res {
			for _, anoy := range mn.restmp {
				for _, morfs := range an.Morphos {
					for _, morfn := range anoy.Morphos {
						if accord(morfn, morfs, el.accord) {
							nres = gocol.AddRes(nres, anoy.Lem, morfs, 0)
						}
					}
				}
			}
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres
	}
	return res
}

func (b *Branche) terminale() bool {
	return len(b.filles) == 0
}
