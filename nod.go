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

/*
func (na *Nod) egale(nb *Nod) bool {
	if na.nucl != nb.nucl {
		return false
	}
	va := true
	for _, ma := range na.mma {
		va = false
		for _, mb := range nb.mma {
			va = va || ma == mb
		}
		if !va {
			return false
		}
	}
	for _, ma := range na.mmp {
		va = false
		for _, mb := range nb.mmp {
			va = va || ma == mb
		}
		if !va {
			return false
		}
	}
	return true
}
*/

// lignes graphviz du nœud
func (n *Nod) graf() []string {
	var ll []string
	for i, m := range n.mma {
		lien := n.grp.ante[i].lien
		// si le lien du sub est vide, c'est que c'est un élément étranger, appartenant à un autre groupe
		// (hyperbate). Il ne faut donc pas l'inclure dans le noeud.
		j := len(n.mma) - i - 1
		if lien > "" {
			ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, n.grp.ante[j].lien))
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

func (n *Nod) inclut(m *Mot) bool {
	for _, el := range n.mma {
		if el == m {
			return true
		}
	}
	if m == n.nucl {
		return true
	}
	for _, el := range n.mmp {
		if el == m {
			return true
		}
	}
	return false
}
