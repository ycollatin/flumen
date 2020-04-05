// groupe.go - Gentes

package main

// signets
// typegroupe
// subvapos
// grvapos

import (
	//"fmt"
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
		pe := PrimEl(e, ".")
		if pe != e {
			n.generique = true
			n.idgr = pe
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
	//debog := id=="n.appFam"
	//if debog {fmt.Println("   debog vaId",id)}
	for _, n := range s.noyaux {
		//if debog{fmt.Println("   generique",n.generique,"idgr",n.idgr)}
		if n.generique {
			ecl := strings.Split(id, ".")
			if n.idgr == ecl[0] {
				return true
			}
		} else if n.id == id {
			return true
		}
	}
	return false
}

func (s *Sub) aPos(sr gocol.Sr) bool {
	for _, n := range s.noyaux {
		if n.generique && n.id == sr.Lem.Pos {
			return true
		} else if n.idgr == sr.Lem.Pos {
			return true
		}
	}
	return false
}

func (s *Sub) vaMorpho(m string) bool {
	for _, sm := range s.morpho {
		if !strings.Contains(m, sm) {
			return false
		}
	}
	return true
}

func (s *Sub) vaPos(p string) bool {
	// signet subvapos
	//debog := s.groupe.id=="v.prepobj" && p=="n.prepAbl"
	//if debog {fmt.Println("Sub.vaPos g",s.groupe.id,"p",p)}
	pgen := strings.Index(p, ".") > -1
	for _, n := range s.noyaux {
		//if debog {fmt.Println("  .vaPos, n.idgr",n.idgr,"n.id",n.id)}
		if pgen {
			if n.id == p {
				return true
			}
		} else {
			if n.idgr == p {
				return true
			}
		}
	}
	return false
}

// typegroupe

type Groupe struct {
	id			string
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

func (g *Groupe) vaPos(p string) bool {
	// signet grvapos
	//debog := g.id == "v.prepobj"
	//if debog {fmtp.Println("Groupe.vaPos, p", p)}
	for _, pos := range g.pos {
		prel := PrimEl(pos, ".")
		if prel == pos && pos == PrimEl(p, ".") {
			return true
		} else if pos == p {
			return true
		}
	}
	return false
}
