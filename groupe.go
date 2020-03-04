// groupe.go - Gentes

package main

import (
	"github.com/ycollatin/gocol"
	"strings"
)

type Sub struct {
	pos		string
	lien	string
	morpho	[]string
	accord	string
}

func (s *Sub) idGr() string {
	ee := strings.Split(s.pos, ".")
	return ee[0]
}

type Groupe struct {
	id,idGr	string
	pos		[]string	// pos du noyau
	morph	[]string	// traits morpho du noyau
	lexSynt	[]string	// étiquettes lexicosyntaxiques du noyyau
	ante	[]*Sub
	post	[]*Sub
}

var grpTerm, grp []*Groupe

func creeSub(v string) *Sub {
	sub := new(Sub)
	vv := strings.Split(v, ";")
	for i, e := range(vv) {
		switch i {
			case 0:	// pos
			sub.pos = e
			case 1:	// id-lien
			sub.lien = e
			case 2: // morpho
			sub.morpho = strings.Split(e, " ")
			case 3: // accord
			sub.accord = e
		}
	}
	return sub
}

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
		case "id":
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
			g.ante = append(g.ante, creeSub(v))
		case "p":
			g.post = append(g.post, creeSub(v))
		}
	}
	return g
}

func (g *Groupe) estNoyau(m *Mot) bool {
	//pos		[]string	// pos du noyau
	//morph	[]string	// traits morpho du noyau
	//lexSynt	[]string	// étiquettes lexicosyntaxiques du noyyau
	for _, an := range m.ans {
		// pos
		if !contient(g.pos, an.Lem.Pos) {
			return false
		}
		// morpho
		var va bool
		for _, morf := range an.Morphos {
			va = true
			for _, gmorf := range g.morph {
				va = va && strings.Contains(morf, gmorf)
			}
		}
		if !va {
			return false
		}
		for _, ls := range(g.lexSynt) {
			va = va && contient(m.lexsynt, ls)
		}
		if !va {
			return false
		}
	}
	return true
}

func lisGroupes(nf string) {
	llin := gocol.Lignes(nf)
	var ll []string
	for _, l := range llin {
		deb := l[:4]
		switch deb {
		case "grp:":
			g := creeGroupe(ll)
			grp = append(grp, g)
			ll = nil
		case "ter:":
			g := creeGroupe(ll)
			grpTerm = append(grpTerm, g)
			ll = nil
		default:
			ll = append(ll, l)
		}
	}
	grp = append(grp, creeGroupe(ll))
}
