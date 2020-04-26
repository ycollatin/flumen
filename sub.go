//    sub.go   -- gentes
package main

import (
	"strings"
)

// un Sub est un élément de Groupe
type Sub struct {
	groupe   *Groupe  // groupe propriétaire du sub
	noyaux   []*Noy   // Noyaux possibles du sub
	noyexcl  []*Noy   // Noyaux exclus
	lien     string   // étiquette du lien noyau -> sub
	multi    bool     // armé : le lien peut être utilisé plusieurs fois
	morpho   []string // traits morphos requis
	accord   string   // accord sub - noyau
	terminal bool     // le sub est un mot
	lexsynt  []string // étiquettes lexicosyntaxiques
}

// crée un sub du groupe g à partir de la ligne v, terminal si v armé
func creeSub(v string, g *Groupe, t bool) *Sub {
	sub := new(Sub)
	sub.groupe = g
	vv := strings.Split(v, ";")
	for i, e := range vv {
		switch i {
		case 0: // noyaux
			sub.noyaux, sub.noyexcl = creeNoy(e)
		case 1: // id-lien
			if e > "" && e[0] == '+' {
				sub.lien = e[1:]
				sub.multi = true
			} else {
				sub.lien = e
			}
		case 2: // morpho
			sub.morpho = strings.Split(e, ",")
			if len(sub.morpho) == 1 && sub.morpho[0] == "" {
				sub.morpho = nil
			}
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

// vrai si la morpho est acceptée par l'une des morphos du sub
func (s *Sub) vaMorpho(m string) bool {
	if len(s.morpho) == 0 {
		return true
	}
	for _, sm := range s.morpho {
		lt := strings.Split(sm, " ")
		va := true
		for _, trait := range lt {
			va = va && strings.Contains(m, trait)
		}
		if va {
			return true
		}
	}
	return false
}
