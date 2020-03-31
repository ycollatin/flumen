//     phrase.go - Publicola

package main

import (
	"fmt"
	"strings"
)

type Nod struct {
	grp		*Groupe		// groupe du noeud Nod
	mma,mmp	[]*Mot		// liste des mots avant et après le noyau
	nucl	*Mot		// noyau du Nod
	rang	int			// 
}

// lignes graphviz du nœud
func (n *Nod) graf() ([]string) {
	var ll []string
	for _, m := range n.mma {
		ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, m.sub.lien))
	}
	for _, m := range n.mmp {
		ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, m.sub.lien))
	}
	return ll
}

func (n *Nod) valide() {
	if !n.nucl.elucide() {
		n.nucl.valideTmp()
	}
	for _, m := range n.mma {
		if !m.elucide() {
			m.valideTmp()
		}
	}
	for _, m := range n.mmp {
		if !m.elucide() {
			m.valideTmp()
		}
	}
}

type Phrase struct {
	gr		string
	imot	int
	mots	[]*Mot
	nods	[]*Nod
}

var phrase *Phrase

func creePhrase(t string) *Phrase {
	p := new(Phrase)
	p.gr = t
	mm := strings.Split(t, " ")
	for _, nm := range(mm) {
		p.append(creeMot(nm))
	}
	return p
}

func (p *Phrase) append(m *Mot) {
	m.rang = len(p.mots)
	p.mots = append(p.mots, m)
}

func (p *Phrase) arbre() ([]string, []string) {
	var lexpl []string
	var ll []string
	// réinitialisation des noeuds
	p.nods = nil
	// recherche des noyaux
	// groupes terminaux
	for _, m := range p.mots {
		if m.dejaNoy() {
			continue
		}
		for _, g := range grpTerm {
			n := m.noeud(g)
			if n != nil {
				p.nods = append(p.nods, n)
				lexpl = append(lexpl, n.grp.id)
			}
		}
		// résolution des conflits (à écrire)
	}

	// groupes non terminaux
	for _, m := range p.mots {
		// pour chaque déf. de groupe non terminal
		for _, g := range grp {
			// m noyau ?
			n := m.noeud(g)
			if n != nil {
				p.nods = append(p.nods, n)
				lexpl = append(lexpl, n.grp.id)
			}
		}
		// résolution des conflits (à écrire)
	}

	// graphe
	ll = append(ll, p.gr)
	for _, n := range p.nods {
		ll = append(ll, n.graf()...)
	}
	return ll, lexpl
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

func (p *Phrase) enClair() string {
	var lm []string
	for i:=0; i<len(p.mots); i++ {
		m := p.mots[i].gr
		if i == p.imot {
			m = rouge(m)
		}
		//ec = fmt.Sprintf("%s %s", ec, m)
		lm = append(lm, m)
	}
	return strings.Join(lm, " ")+"."
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

func (p *Phrase) rang(m *Mot) int {
	for i, mot := range p.mots {
		if mot == m {
			return i
		}
	}
	return -1
}
