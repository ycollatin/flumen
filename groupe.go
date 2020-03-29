// groupe.go - Gentes

package main

import (
	"github.com/ycollatin/gocol"
	"strings"
)

type Noy struct {
	id, idgr	string
	generique	bool
}

func creeNoy(s string) []*Noy {
	var ln []*Noy
	ecl := strings.Split(s, " ")
	for _, e := range ecl {
		n := new(Noy)
		n.id = e
		eg := strings.Split(e, ".")
		n.generique = len(eg) == 1
		if n.generique {
			n.idgr = eg[0]
		} else {
			n.idgr = n.id
		}
		ln = append(ln, n)
	}
	return ln
}

// un Sub est un élément de Groupe
type Sub struct {
	groupe		*Groupe		// groupe propriétaire du sub
	noyaux		[]*Noy		// Noyaux possibles du sub
	lien		string		// étiquette du lien noyau -> sub
	morpho		[]string	// traits morphos requis
	accord		string		// accord sub - noyau
	generique	bool		// le pos n'a pas de sousgroupe (séparé par '.')
	terminal	bool		// le sub est un mot
	lexsynt		[]string	// étiquettes lexicosyntaxiques 
}

func creeSub(v string, g *Groupe, t bool) *Sub {
	sub := new(Sub)
	sub.groupe = g
	vv := strings.Split(v, ";")
	for i, e := range(vv) {
		switch i {
			case 0:	// pos
			sub.noyaux = creeNoy(e)
			case 1:	// id-lien
			sub.lien = e
			case 2: // morpho
			sub.morpho = strings.Split(e, ",")
			case 3: // accord
			sub.accord = e
			case 4: //lexsynt
			sub.lexsynt = strings.Split(e, " ")
		}
	}
	sub.terminal = t
	return sub
}

func (s *Sub) vaId(id string) bool {
	for _, n := range s.noyaux {
		if n.idgr == id {
			return true
		}
	}
	return false
}

func (s *Sub) vaPos(sr gocol.Sr) bool {
	for _, n := range s.noyaux {
		if n.generique && n.id == sr.Lem.Pos {
			return true
		} else if n.idgr == sr.Lem.Pos {
			return true
		}
	}
	return false
}

type Groupe struct {
	id,idGr		string
	pos			[]string	// pos du noyau
	morph		[]string	// traits morpho du noyau
	lexSynt		[]string	// étiquettes lexicosyntaxiques du noyau
	ante		[]*Sub
	post		[]*Sub
}

var grpTerm, grp []*Groupe

func creeGroupe(ll []string) *Groupe {
	if len(ll) == 0 {
		return nil
	}
	g := new(Groupe)
	for _, l := range ll {
		kv := strings.Split(l, ":")
		k := kv[0]
		v := kv[1]
		switch k {
		case "ter", "grp":
			g.id = v
			ee := strings.Split(v, ".")
			g.idGr = ee[0]
		case "pos":
			g.pos = strings.Split(v, " ")
		case "morph":
			g.morph = strings.Split(v, ";")
		case "lexSynt":
			g.lexSynt = strings.Split(v, " ")
		case "a":
			g.ante = append(g.ante, creeSub(v, g, true))
		case "ag":
			g.ante = append(g.ante, creeSub(v, g, false))
		case "p":
			g.post = append(g.post, creeSub(v, g, true))
		case "pg":
			g.post = append(g.post, creeSub(v, g, false))
		}
	}
	return g
}

func lisGroupes(nf string) {
	llin := gocol.Lignes(nf)
	var ll []string
	for _, l := range llin {
		deb := l[:4]
		if deb == "ter:" || deb == "grp:" {
			g := creeGroupe(ll)
			if g != nil {
				if ll[0][:4] == "grp:" {
					grp = append(grp, g)
				} else {
					grpTerm = append(grpTerm, g)
				}
				ll = nil
			}
		}
		ll = append(ll, l)
	}
	grp = append(grp, creeGroupe(ll))
}

func (g *Groupe) nbSubs() int {
	return len(g.ante) + len(g.post)
}
