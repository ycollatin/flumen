//    sub.go   -- gentes
package main

import (
	"strings"
)

// un Sub est un élément de Groupe
type Sub struct {
	groupe		*Groupe		// groupe propriétaire du sub
	noyaux		[]*Noy		// Noyaux possibles du sub
	noyexcl		[]*Noy		// Noyaux exclus
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
			case 0:	// noyaux
				sub.noyaux, sub.noyexcl = creeNoy(e)
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

func (s *Sub) vaMorpho(m string) bool {
	for _, sm := range s.morpho {
		if !strings.Contains(m, sm) {
			return false
		}
	}
	return true
}
