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

type Phrase struct {
	gr		string
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
			n := p.noeud(m, g)
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
			n := p.noeud(m, g)
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

/*
// id du Nod dont m est le noyau
func (p *Phrase) estNuclDe(m *Mot) string {
	var ret string
	for _, nod := range p.nods {
		if nod.nucl == m {
			ret = nod.grp.id
		}
	}
	return ret
}
*/

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

// si m peut être noyau d'un gourpe g, un Nod est renvoyé, sinon nil.
func (p *Phrase) noeud(m *Mot, g *Groupe) *Nod {
	//debog := g.id=="n.fam" && m.gr == "filius"
	//if debog {fmt.Println("noeud", m.gr, g.id)}
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
	//if debog {fmt.Println("   noeud oka, estNoyau",m.gr,g.id,m.estNoyau(g))}
	// m peut-il être noyau du groupe g ?
	if !m.estNoyau(g) {
		return nil
	}

	// création du noeud de retour
	nod := new(Nod)
	nod.grp = g
	nod.nucl = m
	nod.rang = m.rang
	// vérif des subs
	// ante
	r := rang - 1
	//if debog {fmt.Println("   noeud okb")}
	// reгcherche rétrograde des subs ante
	for ia := lante-1; ia > -1; ia-- {
		sub := g.ante[ia]
		ma := p.mots[r]
		// passer les mots
		for ma.dejaSub() && r > 0 {
			r--
			ma = p.mots[r]
		}
		//if debog {fmt.Println("  ma",ma.gr,"nl/nm",ma.nl,ma.nm,"estSub",m.gr,"id grup",sub.groupe.id,ma.estSub(sub, m))}
		if !ma.estSub(sub, m) {
			// réinitialiser lemme et morpho de ma
			return nil
		}
		ma.sub = sub
		nod.mma = append(nod.mma, ma)
		r--
		//if debog {fmt.Println("    vu",ma.gr)}
	}
	//if debog {fmt.Println("   okd",len(g.post),"g.post")}
	// post
	for ip, sub := range g.post {
		r := rang + ip + 1
		mp := p.mots[r]
		//if debog {fmt.Println("post, mp",mp.gr)}
		for mp.dejaSub() && r < len(p.mots) - 1 {
			r++
			mp = p.mots[r]
		}
		//if debog {fmt.Println("     mp", mp.gr,"estSub",m.gr,sub.groupe.id,mp.estSub(sub, m))}
		if !mp.estSub(sub, m) {
			// réinitialiser lemme et morpho de mp
			return nil
		}
		mp.sub = sub
		nod.mmp = append(nod.mmp, mp)
	}
	if len(nod.mma) + len(nod.mmp) > 0 {
		m.pos = g.id
		// fixer lemme et morpho de m, ma et mp
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
