//     phrase.go - Publicola

package main

import (
	"fmt"
	"strings"
)

type Nod struct {
	mm		[]*Mot
	nucl	*Mot
	rangPr	int		// rang du 1er mot
	grp		*Groupe
}

// lignes graphviz du nœud
func (n *Nod) graf() (string) {
	var ll []string
	inod := phrase.rang(n.nucl)
	for _, m := range n.mm {
		ll = append(ll, fmt.Sprintf("%d -> %d", inod, phrase.rang(m)))
	}
	return strings.Join(ll, "\n")
}

type Phrase struct {
	gr		string
	mots	[]*Mot
	nods	[]*Nod
}

var phrase *Phrase

func (p *Phrase) append(m *Mot) {
	p.mots = append(p.mots, m)
}

func (p *Phrase) arbre() string {
	// groupes terminaux, recherche
	for _, m := range p.mots {
		for _, g := range grpTerm {
			n := p.noeud(m, g)
			if n != nil {
				p.nods = append(p.nods, n)
			}
		}
		// résolution des conflits
	}

	// groupes non terminaux
	// pour chaque mot
	for _, m := range p.mots {
		// si le mot est sub, passer
		if p.estSub(m) || p.estNucl(m) {
			continue
		}
		for _, g := range grp {
			// m noyau ?
			n := p.noeud(m, g)
			if n != nil {
				p.nods = append(p.nods, n)
			}
		}
	}

	// graphe
	var ll []string
	ll = append(ll, p.gr)
	for _, n := range p.nods {
		ll = append(ll, n.graf())
	}
	return strings.Join(ll, "\n")
}

func (p *Phrase) estNucl(m *Mot) bool {
	for _, nod := range p.nods {
		if nod.nucl == m {
			return true
		}
	}
	return false
}

func (p *Phrase) estNuclDe(m *Mot) string {
	for _, nod := range p.nods {
		if nod.nucl == m {
			return nod.grp.id
		}
	}
	return ""
}

func (p *Phrase) estSub(m *Mot) bool {
	for _, nod := range p.nods {
		for _, el := range nod.mm {
			if el == m {
				return true
			}
		}
	}
	return false
}

// extrait de la phrase p n mots à partir du mot
// n° d
func (p *Phrase) ex(d, n int) (e string) {
	var gab string = "%s"
	for i := 0; i<n; i++ {
		if e != "" {
			gab = " %s"
		}
		e += fmt.Sprintf(gab, p.mots[d+i].gr)
	}
	return
}

// affiche la phrase en colorant n mots à partir
// du mot n° d
func (p *Phrase) exr(d, n int) (e string) {
	var gab string = "%s"
	for i := 0; i<len(p.mots); i++ {
		if e != "" {
			gab = " %s"
		}
		if i >= d && i<d+n {
			e += fmt.Sprintf(gab, rouge(p.mots[i].gr))
		} else {
			e += fmt.Sprintf(gab, p.mots[i].gr)
		}
	}
	return
}

func (p *Phrase) enClair() (ec string) {
	for i:=0; i<len(p.mots); i++ {
		m := p.mots[i].gr
		if i == imot {
			m = rouge(m)
		}
		ec = fmt.Sprintf("%s %s", ec, m)
	}
	ec += "."
	return
}

func majPhrase() {
	phrase = texte.phrases[ip]
}

// nombre de mots
func (p *Phrase) nbm() int {
	return len(p.mots)
}

// renvoie le noeud dont m *est* le noyau
func (p *Phrase) nod(m *Mot) *Nod {
	for _, n := range p.nods {
		if n.nucl == m {
			return n
		}
	}
	return nil
}

// renvoie le noeud dont m peut être le noyau
func (p *Phrase) noeud(m *Mot, g *Groupe) *Nod {
	rang := p.rang(m)
	// mot de rang trop faible
	if rang < len(g.ante) {
		return nil
	}
	// ou trop élevé
	if p.nbm() - rang < len(g.post) {
		return nil
	}
	// vérif noyau
	if !m.estNoyau(g) {
		return nil
	}

	// création du noeud de retour
	nod := new(Nod)
	nod.rangPr = rang - len(g.ante)
	nod.grp = g
	nod.nucl = m
	// vérif des subs
	// ante
	for ia, sub := range g.ante {
		r := nod.rangPr + ia
		ma := p.mots[r]
		if !ma.estSub(sub) {
			return nil
		}
		nod.mm = append(nod.mm, ma)
	}
	// post
	for ip, sub := range g.post {
		r := rang + ip
		mp := p.mots[r]
		if !mp.estSub(sub) {
			return nil
		}
		nod.mm = append(nod.mm, mp)
	}
	if len(nod.mm) > 0 {
		return nod
	}
	return nil
}

func (p *Phrase) rang(m *Mot) int {
	for i, mot := range p.mots {
		if mot == m {
			return i
		}
	}
	return -1
}
