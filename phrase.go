//     phrase.go - Publicola

package main

import (
	"fmt"
)

type Phrase struct {
	mots	[]*Mot
}

var (
	phrase *Phrase
)

func (p *Phrase) append(m *Mot) {
	p.mots = append(p.mots, m)
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

func (p *Phrase) rang(m *Mot) int {
	for i, mot := range p.mots {
		if mot == m {
			return i
		}
	}
	return -1
}
