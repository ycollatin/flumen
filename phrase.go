//     phrase.go - Publicola

package main

import (
	"fmt"
)

type Phrase struct {
	mots	[]*Mot
}

type Nod struct {
	mm		[]*Mot
	nod		int		// rang du 1er mot
	grp		*Groupe
	graf	[]string
}

var (
	phrase *Phrase
)

func (p *Phrase) append(m *Mot) {
	p.mots = append(p.mots, m)
}

func (p *Phrase) arbre() string {
	var noeuds []*Nod
	// groupes terminaux, recherche
	for _, m := range p.mots {
		for _, g := range grpTerm {
			n := p.noeud(m, g)
			if n != nil {
				noeuds = append(noeuds, n)
			}
		}
		// résolution des conflits
	}
	// groupes non terminaux
	return "incomplet"
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
	nod.mm = append(nod.mm, m)
	nod.nod = rang - len(g.ante)
	nod.grp = g
	// vérif des subs
	// ante
	for ia, sub := range g.ante {
		r := rang - len(g.ante) + ia
		ma := p.mots[r]
		if !ma.estSub(sub) {
			return nil
		}
		nod.graf = append(nod.graf, fmt.Sprintf("%s -> %s", m.gr, ma.gr))
	}
	// post
	for ip, sub := range g.post {
		r := rang + ip
		mp := p.mots[r]
		if !mp.estSub(sub) {
			return nil
		}
		nod.graf = append(nod.graf, fmt.Sprintf("%s -> %s", m.gr, mp.gr))
	}

	return nod
}

func (p *Phrase) rang(m *Mot) int {
	for i, mot := range p.mots {
		if mot == m {
			return i
		}
	}
	return -1
}
