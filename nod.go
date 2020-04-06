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
	var ll []string
	for _, m := range n.mma {
		ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, m.sub.lien))
	}
	for _, m := range n.mmp {
		ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, m.sub.lien))
	}
	return ll
}

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

func (n *Nod) fixeRes() {
	for _, m := range n.mma {
		fmt.Print(m.gr, " ");
	}
	fmt.Print("\n")
	for _, m := range n.mmp {
		fmt.Print(m.gr, " ");
	}
	fmt.Print("\n")
}
