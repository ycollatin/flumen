// nod.go  --  gentes

package main

import (
	"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

type Nod struct {
	groupe   *Groupe // groupe du noeud Nod
	mma, mmp []*Mot  // liste des mots avant et après le noyau
	nbsubs   int
	rra, rrp map[int]gocol.Res // liste des lemmatisations de chaque mot
	nucl     *Mot              // noyau du Nod
	rnucl    gocol.Res         // lemmatisations du noyau
	rang     int
	valide   bool
}

func (n *Nod) copie() Nod {
	var nn Nod
	nn.groupe = n.groupe
	for _, m := range n.mma {
		nn.mma = append(nn.mma, m)
	}
	for _, m := range n.mmp {
		nn.mmp = append(nn.mmp, m)
	}
	nn.nucl = n.nucl
	for _, r := range n.rnucl {
		var nr gocol.Sr
		nr.Lem = r.Lem
		for _, morf := range r.Morphos {
			nr.Morphos = append(nr.Morphos, morf)
		}
		nn.rnucl = append(nn.rnucl, nr)
	}
	nn.rang = n.rang
	nn.nbsubs = nn.groupe.nbsubs
	return nn
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

// Compare les noeuds na et nb et renvoie vrai
// s'ils sont égaux
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

// Compare les lemmatisations resa et resb et renvoie
// vrai si elles sont égales
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
		if lien != "" {
			ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, lien))
		}
	}
	return ll
}
