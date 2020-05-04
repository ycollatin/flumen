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
	ids		 []string // identifiants des groupes possibles pour le noyau
	idsexcl  []string // ids exclus
	familles []string // le préfixe seulement des ces ids est requis
	cles	 []string // clés des lemmes possibles
	clesexcl []string // clés exclues
	poss	 []string // pos des lemmes
	posexcl  []string // pos exclus
	lien     string   // étiquette du lien noyau -> sub
	multi    bool     // armé : le lien peut être utilisé plusieurs fois
	morpho   []string // traits morphos requis
	accord   string   // accord sub - noyau
	terminal bool     // le sub est un mot
	lexsynt  []string // étiquettes lexicosyntaxiques
	exclls   []string // exclusions lexicosyntaxiques
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
			// partage des éléments 
			els := strings.Split(e, " ")
			for _, el := range els {
				parts := strings.Split(el, ".")
				if len(parts) == 2 {
					if el[0] == '!' {
						// idsexcl
						sub.idsexcl = append(sub.idsexcl, el[1:])
					} else {
						// ids
						sub.ids = append(sub.ids, el)
					}
				} else if strings.Contains(el, "\"") {
					// clés
					if el[0] == '!' {
						sub.clesexcl = append(sub.clesexcl, el[2:len(el)-1])
					} else {
						sub.cles = append(sub.clesexcl, el[1:len(el)-1])
					}
				} else if strings.Contains(el, "@") {
					if el[0] == '!' {
						sub.posexcl = append(sub.posexcl, el[2:len(el)-1])
					} else {
						sub.poss = append(sub.poss, el[1:])
					}
				} else {
				// familles
					sub.familles = append(sub.familles, el)
				}
			}
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
			els := strings.Split(e, " ")
			for _, el := range els {
				if el[0] == '!' {
					sub.exclls = append(sub.exclls, el[1:])
				} else {
					sub.lexsynt = append(sub.lexsynt, el)
				}
			}
		}
	}
	sub.terminal = t
	return sub
}

func (s *Sub) vaId(id string) bool {
	for _, ne := range s.noyexcl {
		if ne.id == id {
			return false
		}
	}
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
