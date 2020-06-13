//     Branche.go - Gentes
// Une Branche essaie d'intégrer dans des groupes les mots de la phrase qui
// sont encore isolés. Dès qu'elle a trouvé une solution, elle passe la
// main à une branche fille, qui fait la même chose, puis lui rend la main
// quand elle est parvenue à la fin de son exploration.
// Si elle ne trouve aucune solution, elle rend la main à sa branche mère, si
// elle existe.

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
	"strings"

	"github.com/ycollatin/gocol"
)

type Sol struct {
	nods   []*Nod
	nbarcs int
}

type Branche struct {
	filles   []*Branche        // liste des branches filles
	gr       string            // texte de la phrase
	nods     []*Nod            // noeuds validés
	niveau   int               // n° de la branche par rapport au tronc
	photos   map[int]gocol.Res // lemmatisations et appartenance de groupe propres à la branche
	vendange []Sol             // résultat de la récolte
	veto     map[int][]*Nod    // index : rang du mot; valeur : liste des liens interdits
}

var (
	mots    []*Mot
	nbmots  int
	journal []string
	tronc   *Branche
)

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

// vrai si le *Mot noyau a déjà un sub de lien "lien"
func (b *Branche) adeja(noyau *Mot, lien string) bool {
	for _, nod := range b.nods {
		if nod.nucl == noyau {
			for i, _ := range nod.mma {
				if nod.groupe.ante[i].lien == lien {
					return true
				}
			}
			for i, _ := range nod.mmp {
				if nod.groupe.post[i].lien == lien {
					return true
				}
			}
		}
	}
	return false
}

func (b *Branche) aSubLien(lien string, mot *Mot) bool {
	for _, nod := range b.nods {
		for k, v := range nod.lla {
			if k == mot.rang && v == lien {
				return true
			}
		}
		for k, v := range nod.llp {
			if k == mot.rang && v == lien {
				return true
			}
		}
	}
	return false
}

// copie branche mère  - branche fille
func (b *Branche) copie() *Branche {
	// signet scopie
	nb := new(Branche)
	nb.gr = b.gr
	nb.niveau = b.niveau + 1
	for _, nbm := range b.nods {
		nb.nods = append(nb.nods, nbm.copie())
	}
	nb.photos = make(map[int]gocol.Res)
	for i, r := range b.photos {
		for _, sr := range r {
			var nr gocol.Sr
			nr.Lem = sr.Lem
			for _, morf := range sr.Morphos {
				nr.Morphos = append(nr.Morphos, morf)
			}
			nb.photos[i] = append(nb.photos[i], nr)
		}
	}
	nb.veto = make(map[int][]*Nod)
	for i, v := range b.veto {
		for _, nod := range v {
			nb.veto[i] = append(nb.veto[i], nod.copie())
		}
	}
	return nb
}

// Faux si m є à un groupe
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

// mb є au groupe ma, directement ou indirectement
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

// supprime les branches qui ne sont pas assez productive
func (b *Branche) elague() {
	// le nombre maximal d'arcs est immédiatement inférieur
	// au nombre de mots.
	// recherche du nombre maximal de branches dans les
	// solutions enregistrées dans b.vendange
	var max int
	for _, sol := range b.vendange {
		if sol.nbarcs > max {
			max = sol.nbarcs
		}
	}
	var nv []Sol
	for _, sol := range b.vendange {
		if sol.nbarcs == max {
			nv = append(nv, sol)
		}
	}
	b.vendange = nv
}

// explore toutes les possibilités d'une branche
func (b *Branche) explore() {
	// signet sexplore
	for _, m := range mots {
		// 1. groupes terminaux
		b.exploreGroupes(m, grpTerm)
		// 2. groupes non terminaux
		b.exploreGroupes(m, grp)
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
		if n != nil && n.valide {
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
				}
				for _, ma := range n.mma {
					// photos des éléments antéposés
					if mph == ma {
						bf.photos[ma.rang] = ma.restmp
					}
				}
				for _, mp := range n.mmp {
					// photos des éléments postposés
					if mph == mp {
						bf.photos[mp.rang] = mp.restmp
					}
				}
			}
			// màj journal
			indent := strings.Repeat("  ", bm.niveau)
			journal = append(journal, fmt.Sprintf("%s %d. %s", indent, bm.niveau, n.doc(false)))
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
			// réinitialiser les photos
			for _, mot := range mots {
				mot.restmp = bm.photos[mot.rang]
			}
		}
	}
}

// id des groupes dont m est noyau
func (b *Branche) ids(m *Mot) (lids []string) {
	for _, nod := range b.nods {
		if nod.nucl == m {
			lids = append(lids, nod.groupe.id)
		}
	}
	return
}

// si m peut être noyau d'un Groupe g, un Nod est renvoyé, sinon nil.
func (b *Branche) noeud(m *Mot, g *Groupe) *Nod {
	// snoeud, signet

	// noeud nnul pour le retour d'échec
	var nnul *Nod

	// vérification de rang
	rang := m.rang
	lante := len(g.ante)
	// mot de rang trop faible
	if rang-lante < 0 {
		return nnul
	}
	// ou trop élevé
	if rang+len(g.post)-1 >= nbmots {
		return nnul
	}

	// validation du noyau
	m.restmp = b.photos[m.rang]
	m.restmp = b.resEl(m, g.nucl, m, m.restmp)
	if m.restmp == nil {
		return nnul
	}

	// création du noeud de retour
	nod := new(Nod)
	nod.rra = make(map[int]gocol.Res)
	nod.rrp = make(map[int]gocol.Res)
	nod.lla = make(map[int]string)
	nod.llp = make(map[int]string)
	nod.groupe = g
	nod.nucl = m
	nod.rnucl = m.restmp
	nod.rang = rang
	nod.nbsubs = g.nbsubs

	// reгcherche rétrograde des subs ante
	r := rang - 1
	var nonLies int
	for ia := lante - 1; ia > -1; ia-- {
		if r < 0 {
			// le rang du mot est < 0 : impossible
			return nnul
		}
		ma := mots[r]
		// passer les mots déjà subordonnés
		for b.dejasub(ma) || ma.interj {
			r--
			if r < 0 {
				return nnul
			}
			ma = mots[r]
		}
		// vérification de réciprocité, puis du lien lui-même
		if b.domine(ma, m) && g.ante[ia].lien > "" {
			return nnul
		}
		sub := g.ante[ia]
		// réinitialisation des lemmatisations de test
		ma.restmp = b.photos[ma.rang]
		// validation du noyau
		ma.restmp = b.resEl(ma, sub, m, ma.restmp)
		if ma.restmp == nil {
			return nnul
		}
		if sub.lien > "" {
			nod.mma = append(nod.mma, ma)
			nod.rra[ma.rang] = ma.restmp
			nod.lla[ma.rang] = sub.lien
		} else {
			// si le lien est muet, c'est qu'il est étranger au
			// groupe. Il y a hyperbate.
			nonLies++
		}
		r--
	}
	// si les ante ne sont pas au complet, renvoyer nnul
	// Tenir compte des éléments non liés :
	// hyperbates, vocatifs et interjections
	if len(nod.mma)+nonLies < lante {
		return nnul
	}

	// reгcherche des subs post
	r = rang + 1
	for ip, sub := range g.post {
		if r >= nbmots {
			break
		}
		mp := mots[r]
		for b.dejasub(mp) || b.domine(mp, m) || mp.interj {
			r++
			if r >= nbmots {
				return nnul
			}
			mpn := b.noyau(mp)
			if mpn != nil && mpn.rang < m.rang {
				return nnul
			}
			mp = mots[r]
		}
		mp.restmp = b.photos[mp.rang]
		mp.restmp = b.resEl(mp, sub, m, mp.restmp)
		if mp.restmp == nil {
			return nnul
		}
		// cf. supra nod.mma
		if g.post[ip].lien > "" {
			nod.mmp = append(nod.mmp, mp)
			nod.rrp[mp.rang] = mp.restmp
			nod.llp[mp.rang] = sub.lien
		} else {
			nonLies++
		}
		r++
	}
	// le noeud est valide si tous les post ont été trouvés
	if len(nod.mmp)+nonLies < len(g.post) {
		return nnul
	}
	nod.valide = true
	return nod
}

func (b *Branche) noyau(m *Mot) *Mot {
	for _, n := range b.nods {
		for i, msub := range n.mma {
			if n.groupe.ante[i].lien > "" && msub == m {
				return n.nucl
			}
		}
		for i, msub := range n.mmp {
			if n.groupe.post[i].lien > "" && msub == m {
				return n.nucl
			}
		}
	}
	return nil
}

// récolte tous les noeuds terminaux d'un arbre
func (b *Branche) recolte() {
	b.vendange = nil
	for _, f := range b.filles {
		if f.terminale() {
			var (
				nods []*Nod
				nbn  int
			)
			for _, n := range f.nods {
				nods = append(nods, n)
				nbn += len(n.mma) + len(n.mmp)
			}
			if nbn > 0 {
				b.vendange = append(b.vendange, Sol{nods, nbn})
			}
		} else {
			f.recolte()
			b.vendange = append(b.vendange, f.vendange...)
		}
	}
}

// teste m comme élément du groupe de mn comme défini par el, retourne et modifie en conséquence
// sa lemmatisation res. Renvoie nil si m ne peut pas être élément du groupe el.groupe
func (b *Branche) resEl(m *Mot, el *El, mn *Mot, res gocol.Res) gocol.Res {
	// signet sresEl

	// contraintes de groupe
	if !el.multi && b.adeja(mn, el.lien) {
		return nil
	}

	if el.lienNexcl > "" && b.aSubLien(el.lienNexcl, mn) {
		return nil
	}

	if el.lienN > "" && !b.aSubLien(el.lienN, mn) {
		return nil
	}

	// vérification du pos : id du noyau, ou pos du mot
	var va bool
	ids := b.ids(m)
	if len(ids) > 0 {
		for _, id := range ids {
			// familles
			pel := PrimEl(id, ".")
			if len(el.famexcl) > 0 && contient(el.famexcl, pel) {

				/*
					! FIXME : !si inopérant dans /Haec te scire uolui./
					grp:v.infobj
					n:@v;;act;;infobj
					ag:@v !si;objet;inf
				*/
				return nil
			}
			if len(el.idsexcl) > 0 && contient(el.idsexcl, id) {
				return nil
			}
		}
		if len(el.familles) > 0 {
			var vafam bool
			for _, id := range ids {
				for _, elf := range el.familles {
					if elf == PrimEl(id, ".") {
						vafam = true
					}
				}
			}
			if !vafam {
				return nil
			}
		}
		// id des groupes dont m est noyau
		if len(el.ids) > 0 {
			var vaids bool
			for _, id := range ids {
				for _, idel := range el.ids {
					if idel == id {
						vaids = true
					}
				}
			}
			if !vaids {
				return nil
			}
		}
		va = true
	}
	if el.generique() {
		// l'élément n'a aucune des propriétés requises pour un mot isolé,
		// la recherche est terminée.
		if va {
			return res
		} else {
			return nil
		}
	}

	// terminer si aucune des propriétés suivantes n'est requise
	if len(el.poss)+len(el.cles)+len(el.morpho)+len(el.morphexcl)+len(el.lexsynt) == 0 {
		return nil
	}

	var nres gocol.Res

	// contraintes de lemmatisation
	// Commencer par les exclusions
	// lexicosyntaxe, exclus
	if len(el.lsexcl) > 0 {
		nres = nil
		for _, an := range res {
			va := true
			for _, excl := range el.lsexcl {
				va = va && !lexsynt(an.Lem, excl)
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

	// canons exclus
	if len(el.clesexcl) > 0 {
		nres = nil
		for _, an := range res {
			if contient(el.clesexcl, an.Lem.Cle) {
				continue
			}
			nres = append(nres, an)
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres
	}
	// pos exclus
	if len(el.posexcl) > 0 {
		nres = nil
		for _, an := range res {
			if contient(el.posexcl, an.Lem.Pos) {
				continue
			}
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres
	}
	// morpho exclues
	if len(el.morphexcl) > 0 {
		for _, mexcl := range el.morphexcl {
			for _, an := range res {
				for _, morfs := range an.Morphos {
					// pour toutes les morphos valides de m
					if strings.Contains(morfs, mexcl) {
						return nil
					}
				}
			}
		}
	}

	// groupes exclus

	// ensuite les possibles
	// lexicosyntaxe
	if len(el.lexsynt) > 0 {
		nres = nil
		for _, an := range res {
			va := true
			for _, ls := range el.lexsynt {
				va = va && lexsynt(an.Lem, ls) || (ls == "que" && m.que)
			}
			// suffixe -que
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
			for _, cle := range el.cles {
				if an.Lem.Cle == cle {
					nres = append(nres, an)
				}
			}
		}
		if len(nres) == 0 {
			return nil
		}
		res = nres
	}

	// pos
	if len(el.poss) > 0 {
		nres = nil
		for _, an := range res {
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
	if len(el.morpho) > 0 {
		var nres gocol.Res
		for _, an := range res {
			for _, morfs := range an.Morphos {
				// pour toutes les morphos valides de m
				if el.vaMorpho(morfs) {
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
	acc := el.accord
	if acc > "" {
		var nres gocol.Res
		for _, an := range res {
			for _, anoy := range mn.restmp {
				if strings.Contains(acc, "g") && lexsynt(anoy.Lem, "mf") {
					acc = strings.Replace(acc, "g", "", 1)
				}
				for _, morfs := range an.Morphos {
					for _, morfn := range anoy.Morphos {
						// cas des noms f. et m.
						if accord(morfn, morfs, acc) {
							//nres = append(nres, an)
							nres = gocol.AddRes(nres, an.Lem, morfs, 0)
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

// Une branche est terminale si elle n'a pas de filles
// la récolte récupère alors tous les liens des noeuds
// dont elle a hérité.
func (b *Branche) terminale() bool {
	return len(b.filles) == 0
}
