// groupe.go - Gentes

package main

import (
	"github.com/ycollatin/gocol"
	"strings"
)

type Sub struct {
	pos			string		// pos du sub
	lien		string		// étiquette du lien noyau -> sub
	morpho		[]string	// traits morphos requis
	accord		string		// accord sub - noyau
	terminal	bool	// le sub est un mot
	// lexsynt ?
}

func creeSub(v string, t bool) *Sub {
	sub := new(Sub)
	vv := strings.Split(v, ";")
	for i, e := range(vv) {
		switch i {
			case 0:	// pos
			sub.pos = e
			case 1:	// id-lien
			sub.lien = e
			case 2: // morpho
			sub.morpho = strings.Split(e, ",")
			case 3: // accord
			sub.accord = e
		}
	}
	sub.terminal = t
	return sub
}

func (s *Sub) idGr() string {
	ee := strings.Split(s.pos, ".")
	return ee[0]
}

type Groupe struct {
	id,idGr	string
	pos		[]string	// pos du noyau
	morph	[]string	// traits morpho du noyau
	lexSynt	[]string	// étiquettes lexicosyntaxiques du noyau
	ante	[]*Sub
	post	[]*Sub
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
			g.ante = append(g.ante, creeSub(v, true))
		case "ag":
			g.ante = append(g.ante, creeSub(v, false))
		case "p":
			g.post = append(g.post, creeSub(v, true))
		case "pg":
			g.post = append(g.post, creeSub(v, false))
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
