// nod.go  --  gentes

package main

import (
	"fmt"
	"strings"
)

type Nod struct {
	grp		*Groupe		// groupe du noeud Nod
	mma,mmp	[]*Mot		// liste des mots avant et après le noyau
	nucl	*Mot		// noyau du Nod
	rang	int			// 
}

// lignes graphviz du nœud
func (n *Nod) graf() ([]string) {
	//debog := n.grp.id == "v.coordobjv"
	var ll []string
	for i, m := range n.mma {
		//if debog {fmt.Println("Nod.graf",n.grp.id,"m",m.gr,"lien",n.grp.ante[i].lien)}
		ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, n.grp.ante[i].lien))
	}
	for i, m := range n.mmp {
		ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, n.grp.post[i].lien))
	}
	return ll
}

// liste des éléments du noeud, noyau en rouge
func (n *Nod) doc() string {
	var mm []string
	for _, m := range n.mma {
		mm = append(mm, m.gr)
	}
	mm = append(mm, rouge(n.nucl.gr))
	for _, m := range n.mmp {
		mm = append(mm, m.gr)
	}
	mm = append (mm, " - " + n.grp.id)
	return strings.Join(mm, " ")
}
