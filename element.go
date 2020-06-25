//    element.go   -- Gentes
package main

import (
	"fmt"
	"os"
	"strings"
)

// un el est la définition d'un élément de Groupe
type El struct {
	regle     *Regle   // regle propriétaire du el
	ids       []string // identifiants des groupes possibles pour le noyau
	idsexcl   []string // ids exclus
	familles  []string // le préfixe seulement des ces ids est requis
	famexcl   []string // familles exclues (une famille regroupe les groupes de même préfixe)
	cles      []string // clés des lemmes possibles
	clesexcl  []string // clés exclues
	poss      []string // pos des lemmes
	posexcl   []string // pos exclus
	lien      string   // étiquette du lien noyau -> el
	lienN     string   // si l'élément est noyau, il doit être sub par un lien d'id lienN
	lienNexcl string   // si l'élément est noyau, il ne doit pas être sub par un lein d'id lienN
	multi     bool     // armé : le lien peut être utilisé plusieurs fois
	morpho    []string // traits morpho requis
	morphexcl []string // traits morpho exclus
	accord    string   // accord el - noyau
	terminal  bool     // le el est un mot
	lexsynt   []string // étiquettes lexicosyntaxiques
	lsexcl    []string // exclusions lexicosyntaxiques
}

// func creeEl(v string, g *Groupe, t bool) *El
// crée un el du groupe g à partir de la ligne v, terminal si v armé
// type_groupe;identifiant;lien;morpho;accord;lexsynt
// type_groupe: n|a|p|ag|pg
// 		n noyau
// 		a mot antéposé
// 		p mot postposé
// 		ag groupe antéposé
// 		pg groupe postposé
// identifiant: @pos|"lemme"|famille_groupe|groupe
// 		plusieurs identifiants possibles séparés par un espace
// 		@pos : pos du lemme du mot ou du mot-noyau
// 		"lemme" : clé du lemme d'une lemmatisation (gocol.Sr) du mot
// 		famille_groupe : la partie précédant le point '.' dans l'identifiant du groupe
// 		groupe : l'identifiant complet du groupe
// lien: si l'élément n'est pas noyau, identifiant du lien qui sera affiché dans le graphe
//       si l'élément est noyau, identifiant du lien entre l'élément et son propre noyau.
// morpho : morpho d'une lemmatisation (gocol.Sr.Morphos[i])
// accord : accord entre l'élément du groupe et son noyau : 'c' 'g' 'n' ou une combinaison des 3
// lexsynt : propriétés requises du lemme (lexsynt.la)
//
// identifiant, lemme, famille_groupe et groupe peuvent être préfixés en '!' pour
// en faire des propriétés interdites.
//
func creeEl(k, v string, g *Regle) *El {
	el := new(El)
	el.regle = g
	var noyau bool
	switch k {
	case "ter":
		el.terminal = true
	case "grp":
		el.terminal = false
	case "n":
		noyau = true
	}
	vv := strings.Split(v, ";")
	for i, e := range vv {
		if e == "" {
			continue
		}
		if strings.HasSuffix(e, " ") {
			fmt.Printf("groupe %s, ligne %s, espace(s) de fin à supprimer.\n", g.id, v)
			os.Exit(1)
		}
		switch i {
		case 0: // noyaux
			// partage des éléments
			ee := strings.Split(e, " ")
			for _, ecl := range ee {
				parts := strings.Split(ecl, ".")
				switch {
				case len(parts) == 2:
					// l'élément est un groupe ex. p.prepAcc
					part := parts[0]
					if part[0] == '!' {
						// idsexcl
						el.idsexcl = append(el.idsexcl, ecl[1:])
					} else {
						// ids
						el.ids = append(el.ids, ecl)
					}
				case strings.HasPrefix(ecl, "!\""):
					// clé de lemme exclu
					el.clesexcl = append(el.clesexcl, strings.Trim(ecl, "!\""))
				case strings.HasPrefix(ecl, "\""):
					// clé de lemme
					el.cles = append(el.cles, strings.Trim(ecl, "\""))
				case strings.HasPrefix(ecl, "!@"):
					// pos de lemme exclu
					el.posexcl = append(el.posexcl, ecl[2:])
				case strings.HasPrefix(ecl, "@"):
					el.poss = append(el.poss, ecl[1:])
				default:
					// familles (i.e. préfixes de groupes)
					if strings.HasPrefix(ecl, "!") {
						el.famexcl = append(el.famexcl, ecl[1:])
					} else {
						el.familles = append(el.familles, ecl)
					}
				}
			}
		case 1: // id-lien
			if noyau {
				if e[0] == '!' {
					el.lienNexcl = e[1:]
				} else {
					el.lienN = e
				}
			} else {
				if e[0] == '+' {
					el.lien = e[1:]
					el.multi = true
				} else {
					el.lien = e
				}
			}
		case 2: // morpho
			elsm := strings.Split(e, ",")
			for _, elm := range elsm {
				if elm[0] == '!' {
					el.morphexcl = append(el.morphexcl, elm[1:])
				} else {
					el.morpho = append(el.morpho, elm)
				}
			}
			if len(el.morpho) == 1 && el.morpho[0] == "" {
				el.morpho = nil
			}
		case 3: // accord
			el.accord = e
		case 4: //lexsynt
			els := strings.Split(e, " ")
			for _, ecl := range els {
				if ecl[0] == '!' {
					el.lsexcl = append(el.lsexcl, ecl[1:])
				} else {
					el.lexsynt = append(el.lexsynt, ecl)
				}
			}
		}
	}
	return el
}

func (el *El) generique() bool {
	return len(el.poss)+len(el.cles) == 0
}

// vrai si la morpho est acceptée par l'une des morphos du el
func (s *El) vaMorpho(m string) bool {
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
