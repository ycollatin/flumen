//     phrase.go - Publicola

package main

import (
	"fmt"
	"strings"
)

type Nod struct {
	mma,mmp	[]*Mot
	nucl	*Mot
	grp		*Groupe
}

// lignes graphviz du nœud
func (n *Nod) graf() (string) {
	var ll []string
	inod := phrase.rang(n.nucl)
	for i, m := range n.mma {
		ll = append(ll, fmt.Sprintf("%d -> %d [%s]", inod, phrase.rang(m), n.grp.ante[i].lien))
	}
	for i, m := range n.mmp {
		ll = append(ll, fmt.Sprintf("%d -> %d [%s]", inod, phrase.rang(m), n.grp.post[i].lien))
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

func (p *Phrase) arbre() []string {
	// groupes terminaux, recherche
	for _, m := range p.mots {
		if m.dejaNoy() {
			continue
		}
		for _, g := range grpTerm {
			n := p.noeud(m, g)
			if n != nil {
				p.nods = append(p.nods, n)
			}
		}
		// résolution des conflits (à écrire)
	}

	// groupes non terminaux
	// recherche des noyaux
	// pour chaque mot
	for _, m := range p.mots {
		// pour chaque déf. de groupe non terminal
		for _, g := range grp {
			// m noyau ?
			n := p.noeud(m, g)
			if n != nil {
				p.nods = append(p.nods, n)
			}
		}
		// résolution des conflits (à écrire)
	}

	// graphe
	var ll []string
	ll = append(ll, p.gr)
	for _, n := range p.nods {
		ll = append(ll, n.graf())
	}
	return ll
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
		for _, el := range nod.mma {
			if el == m {
				return true
			}
		}
		for _, el := range nod.mmp {
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
	lante := len(g.ante)
	// mot de rang trop faible
	if rang < lante {
		return nil
	}
	// ou trop élevé
	if p.nbm() - rang < len(g.post) {
		return nil
	}
	// m peut-il être noyau ?
	if !m.estNoyau(g) {
		return nil
		}

	// création du noeud de retour
	nod := new(Nod)
	nod.grp = g
	nod.nucl = m
	// vérif des subs
	// ante
	r := rang - 1
	for ia := lante-1; ia > -1; ia-- {
		sub := g.ante[ia]
		ma := p.mots[r]
		for ma.dejaSub() && r > 0 {
			r--
			ma = p.mots[r]
		}
		if !ma.estSub(sub, m) {
			return nil
		}
		nod.mma = append(nod.mma, ma)
		if r <= 0 {
			break
		}
		r--
	}
	// post
	for ip, sub := range g.post {
		r := rang + ip + 1
		mp := p.mots[r]
		for mp.dejaSub() && r < len(p.mots) - 1 {
			r++
			mp = p.mots[r]
		}
		if !mp.estSub(sub, m) {
			return nil
		}
		nod.mmp = append(nod.mmp, mp)
	}
	if len(nod.mma) + len(nod.mmp) > 0 {
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
