// noy.go  --  gentes
package main

import (
	"github.com/ycollatin/gocol"
	"strings"
)

// Une définition de groupe peut donner un choix de noyaux
type Noy struct {
	id, idgr  string // identifiant
	canon     string // canon du lemme, entre " dans groupes.la
	generique bool   // vrai si l'id est suffixé
}

// créateur du noyau
func creeNoy(s string) (ln, lnExcl []*Noy) {
	ecl := strings.Split(s, " ")
	for _, e := range ecl {
		ex := e[0] == '!'
		if ex {
			e = e[1:]
		}
		n := new(Noy)
		if e[0] == '"' {
			n.canon = e[1:len(e)-1]
			n.generique = true
		} else {
			n.id = e
			pe := PrimEl(e, ".")
			n.generique = pe == e
			if n.generique {
				n.idgr = pe
			}
		}
		if ex {
			lnExcl = append(lnExcl, n)
		} else {
			ln = append(ln, n)
		}
	}
	return
}

// vérifie que le *Mot m de lemmatisation Sr peut être un noyau du groupe
func (n *Noy) vaSr(sr gocol.Sr) bool {
	return sr.Lem.Gr[0] == n.canon
}

// vérifie que p peut être un noyau du groupe
func (n *Noy) vaPos(p string) bool {
	pel := PrimEl(p, ".")
	pgen := pel == p
	if pel == p || pgen {
		return n.idgr == p
	}
	if n.generique {
		return n.id == pel
	}
	return p == n.id
}
