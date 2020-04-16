// nod.go  --  gentes

package main

import (
	"fmt"
	"strings"
)

type Nod struct {
	grp      *Groupe // groupe du noeud Nod
	mma, mmp []*Mot  // liste des mots avant et après le noyau
	nucl     *Mot    // noyau du Nod
	rang     int     //
}

// lignes graphviz du nœud
func (n *Nod) graf() []string {
	var ll []string
	for i, m := range n.mma {
		lien := n.grp.ante[i].lien
		// si le lien du sub est vide, c'est que c'est un élément étranger, appartenant à un autre groupe
		// (hyperbate). Il ne faut donc pas l'inclure dans le noeud.
		if lien > "" {
			ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, n.grp.ante[i].lien))
		}
	}
	diff := 0
	for i, m := range n.mmp {
		lien := n.grp.post[i+diff].lien
		if lien == "" {
			diff++
			lien = n.grp.post[i+diff].lien
		}
		ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, lien))
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
	mm = append(mm, " - "+n.grp.id)
	return strings.Join(mm, " ")
}
