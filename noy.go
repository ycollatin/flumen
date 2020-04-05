// noy.go  --  gentes
package main

import "strings"

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
			n.idgr = pe
		} else {
			n.generique = true
			n.idgr = e
		}
		ln = append(ln, n)
	}
	return ln
}
