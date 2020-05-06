// nod.go  --  gentes

package main

import (
	"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

type Nod struct {
	groupe      *Groupe           // groupe du noeud Nod
	mma, mmp []*Mot            // liste des mots avant et après le noyau
	rra, rrp map[int]gocol.Res // liste des lemmatisations de chaque mot
	nucl     *Mot              // noyau du Nod
	rnucl    gocol.Res         // lemmatisations du noyau
	rang     int
	valide	bool
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
	mm = append(mm, " - "+n.groupe.id)
	return strings.Join(mm, " ")
}

func (na Nod) egale(nb Nod) bool {
	if na.nucl != nb.nucl && egaleRes(na.rnucl, nb.rnucl) {
		return false
	}
	if na.groupe.id != nb.groupe.id {
		return false
	}
	va := true
	for _, ma := range na.mma {
		va = false
		for _, mb := range nb.mma {
			va = va || (ma == mb && egaleRes(na.rra[ma.rang], nb.rra[mb.rang]))
		}
		if !va {
			return false
		}
	}
	for _, ma := range na.mmp {
		va = false
		for _, mb := range nb.mmp {
			va = va || (ma == mb && egaleRes(na.rrp[ma.rang], nb.rrp[mb.rang]))
		}
		if !va {
			return false
		}
	}
	return true
}

func egaleRes(resa, resb gocol.Res) bool {
	va := true
	for _, sra := range resa {
		for _, srb := range resb {
			va = va && sra.Lem.Cle == srb.Lem.Cle
			if !va {
				return false
			}
			for _, morfa := range sra.Morphos {
				for _, morfb := range srb.Morphos {
					va = va && morfa == morfb
					if !va {
						return false
					}
				}
			}
		}
	}
	return true
}

// lignes graphviz du nœud
func (n *Nod) graf() []string {
	var ll []string
	for i, m := range n.mma {
		lien := n.groupe.ante[i].lien
		// si le lien du sub est vide, c'est que c'est un élément étranger, appartenant à un autre groupe
		// (hyperbate). Il ne faut donc pas l'inclure dans le noeud.
		j := len(n.mma) - i - 1
		if lien > "" {
			ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, n.groupe.ante[j].lien))
		}
	}
	diff := 0
	for i, m := range n.mmp {
		lien := n.groupe.post[i+diff].lien
		if lien == "" {
			diff++
			lien = n.groupe.post[i+diff].lien
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

func (n *Nod) nbEl()  int {
	return len(n.mma) + len(n.mmp) + 1
}
