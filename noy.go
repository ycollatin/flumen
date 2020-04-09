// noy.go  --  gentes
package main

import (
	"strings"
)

// Une définition de groupe peut donner un choix de noyaux
type Noy struct {
	id, idgr	string	// identifiant
	generique	bool	// vrai si l'id est suffixé
}

// créateur du noyau
func creeNoy(s string) []*Noy {
	var ln []*Noy
	ecl := strings.Split(s, " ")
	for _, e := range ecl {
		n := new(Noy)
		n.id = e
		pe := PrimEl(e, ".")
		n.generique = pe == e
		if n.generique {
			n.idgr = pe
		}
		ln = append(ln, n)
	}
	return ln
}

// vérifie que p peut être un noyau du groupe
func (n *Noy) vaPos(p string) bool {
	pel := PrimEl(p, ".")
	if n.generique {
		return n.id == pel
	}
	if pel == p {
		return n.idgr == p
	}
	return p == n.id
}
