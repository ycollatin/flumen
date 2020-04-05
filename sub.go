//    sub.go   -- gentes
package main

import (
	"strings"
	"github.com/ycollatin/gocol"
)

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

// Vérifie la conformité du pos du sub
// le sub peut avoir un pos générique (n, v, NP, Adv...) ou
// suffixé (n.fam, s.obj)
// Le paramètre p est la pos du candidat. il peut être lui
// aussi suffixé ou non
func (s *Sub) vaPos(p string) bool {
	// signet subvapos
	//debog := s.groupe.id=="v.suj" && p=="n.appFam"
	//if debog {fmt.Println("Sub.vaPos g",s.groupe.id,"p",p)}
	pgen := strings.Index(p, ".") < 0
	for _, n := range s.noyaux {
		//if debog {fmt.Println("  .vaPos, pgen",pgen,"n.idgr",n.idgr,"n.id",n.id)}
		if pgen {
			if n.generique {
				if n.idgr == p {
					return true
				}
			}
		} else {
			if n.generique {
				if n.idgr == PrimEl(p, ".") {
					return true
				}
			} else {
				if p == n.id {
					return true
				}
			}
		}
	}
	return false
}
