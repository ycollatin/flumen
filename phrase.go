//     phrase.go - Publicola

package main

import (
	"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

type Phrase struct {
	gr     string
	imot   int
	nbmots int
	mots   []*Mot
	nods   []*Nod
	ar     []string // arbre de la phrase
	src    []string // source de l'arbre
}

func creePhrase(t string) *Phrase {
	p := new(Phrase)
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

func (p *Phrase) arbre() ([]string, []string) {
	if len(p.ar) > 0 {
		return p.ar, p.src
	}
	var lexpl []string
	var ll []string
	p.nods = nil
	nba := 0 //
	nbn := 0 // pour permettre le départ de boucle
	debut := true
	for nbn > nba || debut {
		// recherche des noyaux
		// groupes terminaux, à n'essayer qu'à la première passe
		nba = nbn
		if nba == 0 {
			for _, g := range grpTerm {
				for _, m := range p.mots {
					if m.dejaNoy() {
						continue
					}
					n := m.noeud(g)
					if n != nil {
						p.nods = append(p.nods, n)
						lexpl = append(lexpl, n.grp.id)
					}
				}
			}
		}
		// groupes non terminaux
		for _, g := range grp {
			for _, m := range p.mots {
				n := m.noeud(g)
				if n != nil {
					p.nods = append(p.nods, n)
					lexpl = append(lexpl, n.grp.id)
				}
			}
		}
		nbn = len(p.nods)
		debut = false
	}

	// graphe
	ll = append(ll, p.gr)
	for _, n := range p.nods {
		ll = append(ll, n.graf()...)
	}
	p.ar = ll
	p.src = graphe(ll)
	return ll, lexpl
}

// texte de la phrase, le mot courant surligné en rouge
func (p *Phrase) enClair() string {
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

// affiche la phrase en colorant n mots à partir
// du mot n° d
func (p *Phrase) exr(d, n int) (e string) {
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

func (p *Phrase) motCourant() *Mot {
	return p.mots[p.imot]
}

func (p *Phrase) teste() {
	if len(p.motCourant().ans) == 0 {
		texte.majPhrase()
	}
}
