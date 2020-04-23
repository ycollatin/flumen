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
	return p
}

func (b *Branche) copie() *Branche {
	nb := new(Branche)
	nb.gr = b.gr
	nb.nbmots = b.nbmots
	//copy(b.mots, nb.mots)
	for _, am := range b.mots {
		nm := am.copie()
		nb.mots = append(nb.mots, nm)
	}
	copy(nb.nods, b.nods)
	nb.mere = b
	nb.niveau = b.niveau + 1
	copy(b.filles, nb.filles)
	return nb
}

//func (p *Branche) arbre() ([]string, []string) {
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
			n := m.noeud(g)
			if n != nil {
				bf.nods = append(bf.nods, n)
				nbf := bf.copie()
				bm.filles = append(bm.filles, nbf)
			}
		}
	}
	// groupes non terminaux
	for _, g := range grp {
		for _, m := range bf.mots {
			n := m.noeud(g)
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
func (p *Branche) enClair() string {
	var lm []string
	for i := 0; i < len(p.mots); i++ {
		m := p.mots[i].gr
		if i == p.imot {
			m = rouge(m)
		}
		lm = append(lm, m)
	}
	return strings.Join(lm, " ") + "."
}

// affiche la Branche en colorant n mots à partir
// du mot n° d
func (p *Branche) exr(d, n int) (e string) {
	var gab string = "%s"
	for i := 0; i < len(p.mots); i++ {
		if e != "" {
			gab = " %s"
		}
		if i >= d && i < d+n {
			e += fmt.Sprintf(gab, rouge(p.mots[i].gr))
		} else {
			e += fmt.Sprintf(gab, p.mots[i].gr)
		}
	}
	return
}

func (p *Branche) motCourant() *Mot {
	return p.mots[p.imot]
}

func (p *Branche) teste() {
	if len(p.motCourant().ans) == 0 {
		texte.majPhrase()
	}
}
