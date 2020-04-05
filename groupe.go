// groupe.go - Gentes

package main

// signets
// grvapos

import (
	//"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

type Groupe struct {
	id			string
	noyaux		[]*Noy		// pos du noyau
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
		case "pos":
			g.noyaux = creeNoy(v)
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

// la morpho morf est-elle compatible avec le noyau du groupe g ?
func (g *Groupe) vaMorph(morf string) bool {
	//debog := g.id=="n.genabl"
	//if debog {fmt.Println("vamorph",g.id,"morf",morf)}
	va := true
	for _, tr := range g.morph {
		va = va && strings.Contains(morf, tr)
	}
	return va
}

/*
func (g *Groupe) vaPos(p string) bool {
	// signet grvapos
	//debog := g.id == "v.prepAbl"
	//if debog {fmtp.Println("Groupe.vaPos, p", p)}
	pgen := strings.Index(p, ".") < 0
	for _, pos := range g.pos {
		primel := PrimEl(pos, ".")
		if primel == pos {
			// le pos est générique
			if pgen && pos == p {
				return true
			}
		} else {
			// le pos est suffixé
			if pos == p {
				return true
			}
		}
	}
	return false
}
*/
