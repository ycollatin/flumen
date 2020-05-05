//    element.go   -- gentes
package main

import (
	"strings"
)

// un el est un élément de Groupe
type El struct {
	groupe   *Groupe  // groupe propriétaire du el
	ids		 []string // identifiants des groupes possibles pour le noyau
	idsexcl  []string // ids exclus
	familles []string // le préfixe seulement des ces ids est requis
	famexcl  []string // préfixes exclus
	cles	 []string // clés des lemmes possibles
	clesexcl []string // clés exclues
	poss	 []string // pos des lemmes
	posexcl  []string // pos exclus
	lien     string   // étiquette du lien noyau -> el
	multi    bool     // armé : le lien peut être utilisé plusieurs fois
	morpho   []string // traits morphos requis
	accord   string   // accord el - noyau
	terminal bool     // le el est un mot
	lexsynt  []string // étiquettes lexicosyntaxiques
	lsexcl   []string // exclusions lexicosyntaxiques
}

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
// lien: identifiant du lien qui sera affiché dans le graphe
// morpho : morpho d'une lemmatisation (gocol.Sr.Morphos[i])
// accord : accord entre l'élément du groupe et son noyau : 'c' 'g' 'n' ou une combinaison des 3
// lexsynt : propriétés requises du lemme (lexsynt.la)
//
// identifiant, lemme, famille_groupe et groupe peuvent être préfixés en '!' pour 
// en faire des propriétés interdites.
//
func creeEl(v string, g *Groupe, t bool) *El {
	el := new(El)
	el.groupe = g
	vv := strings.Split(v, ";")
	for i, e := range vv {
		switch i {
		case 0: // noyaux
			// partage des éléments 
			ee := strings.Split(e, " ")
			for _, ecl := range ee {
				parts := strings.Split(ecl, ".")
				if len(parts) == 2 {
					part := parts[0]
					if part[0] == '!' {
						// idsexcl
						el.idsexcl = append(el.idsexcl, part[1:])
					} else {
						// ids
						el.ids = append(el.ids, part)
					}
				} else if strings.Contains(ecl, "\"") {
					// clés
					if ecl[0] == '!' {
						el.clesexcl = append(el.clesexcl, ecl[2:len(ecl)-1])
					} else {
						el.cles = append(el.clesexcl, ecl[1:len(ecl)-1])
					}
				} else if strings.Contains(ecl, "@") {
					if ecl[0] == '!' {
						el.posexcl = append(el.posexcl, ecl[2:len(ecl)-1])
					} else {
						el.poss = append(el.poss, ecl[1:])
					}
				} else {
					// familles
					if ecl[0] == '!' {
						el.famexcl = append(el.famexcl, ecl[1:])
					} else {
						el.familles = append(el.familles, ecl)
					}
				}
			}
		case 1: // id-lien
			if e > "" && e[0] == '+' {
				el.lien = e[1:]
				el.multi = true
			} else {
				el.lien = e
			}
		case 2: // morpho
			el.morpho = strings.Split(e, ",")
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
	el.terminal = t
	return el
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