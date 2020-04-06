// noy.go  --  gentes
package main

import (
	//"fmt"
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
		n.generique = pe == e
		if n.generique {
			n.idgr = pe
		}
		ln = append(ln, n)
	}
	return ln
}

func (n *Noy) vaPos(p string) bool {
	//debog := n.id=="n" && p=="n.gen"
	pel := PrimEl(p, ".")
	//if debog {fmt.Println("Noy.vaPos p="+p,"n.id="+n.id,"n.idgr="+n.idgr,"pel="+pel,"n.generique",n.generique)}
	if n.generique {
		return n.id == pel
	}
	if pel == p {
		return n.idgr == p
	}
	return p == n.id
}
