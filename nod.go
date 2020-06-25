// nod.go  --  gentes

package main

import (
	"fmt"
	"strings"

	"github.com/ycollatin/gocol"
)

type Nod struct {
	regle   *Regle // règle du noeud Nod
	mma, mmp []*Mot  // liste des mots avant et après le noyau
	nbsubs   int
	nucl     *Mot              // noyau du Nod
	rra, rrp map[int]gocol.Res // liste des lemmatisations de chaque mot
	lla, llp map[int]string    // liste des liens entre le noyau et chaque mot
	rnucl    gocol.Res         // lemmatisations du noyau
	rang     int
	valide   bool
}

func (n *Nod) copie() *Nod {
	nn := new(Nod)
	nn.regle = n.regle
	for _, m := range n.mma {
		nn.mma = append(nn.mma, m)
	}
	for _, m := range n.mmp {
		nn.mmp = append(nn.mmp, m)
	}
	nn.nucl = n.nucl
	// ajout des lemmatisations réduites
	nn.rra = make(map[int]gocol.Res)
	for k, v := range n.rra {
		nn.rra[k] = v
	}
	nn.lla = make(map[int]string)
	for k, v := range n.lla {
		nn.lla[k] = v
	}
	nn.rrp = make(map[int]gocol.Res)
	for k, v := range n.rrp {
		nn.rrp[k] = v
	}
	nn.llp = make(map[int]string)
	for k, v := range n.llp {
		nn.llp[k] = v
	}
	// lemmatisation réduite du noyau
	for _, r := range n.rnucl {
		var nr gocol.Sr
		nr.Lem = r.Lem
		for _, morf := range r.Morphos {
			nr.Morphos = append(nr.Morphos, morf)
		}
		nn.rnucl = append(nn.rnucl, nr)
	}
	nn.rang = n.rang
	nn.nbsubs = nn.regle.nbsubs
	return nn
}

// liste des éléments du noeud, noyau en rouge
func (n *Nod) doc(color bool) string {
	var mm []string
	for _, m := range n.mma {
		mm = append(mm, m.gr)
	}
	if color {
		mm = append(mm, rouge(n.nucl.gr))
	} else {
		mm = append(mm, fmt.Sprintf("*%s*", n.nucl.gr))
	}
	for _, m := range n.mmp {
		mm = append(mm, m.gr)
	}
	mm = append(mm, fmt.Sprintf("- %s", n.regle.id))
	for _, v := range n.lla {
		mm = append(mm, fmt.Sprintf("- %s", v))
	}
	for _, v := range n.llp {
		mm = append(mm, fmt.Sprintf("- %s", v))
	}
	return strings.Join(mm, " ")
}

// Compare les noeuds na et nb et renvoie vrai
// s'ils sont égaux
func (na *Nod) egale(nb *Nod) bool {
	//if na.nucl != nb.nucl && egaleRes(na.rnucl, nb.rnucl) {
	if na.nucl != nb.nucl {
		return false
	}
	//if na.regle.id != nb.regle.id {
	if !na.regle.equiv(nb.regle) {
		return false
	}
	va := true
	for _, ma := range na.mma {
		va = false
		for _, mb := range nb.mma {
			//va = va || (ma == mb && egaleRes(na.rra[ma.rang], nb.rra[mb.rang]))
			va = va || ma == mb
		}
		if !va {
			return false
		}
	}
	for _, ma := range na.mmp {
		va = false
		for _, mb := range nb.mmp {
			//va = va || (ma == mb && egaleRes(na.rrp[ma.rang], nb.rrp[mb.rang]))
			va = va || ma == mb
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
	lenmma := len(n.mma) - 1
	for i, m := range n.mma {
		lien := n.regle.ante[lenmma-i].lien
		// si le lien du sub est vide, c'est que c'est un élément étranger, appartenant à une autre règle
		// (hyperbate). Il ne faut donc pas l'inclure dans le noeud.
		if lien > "" {
			ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, lien))
		}
	}
	for i, m := range n.mmp {
		lien := n.regle.post[i].lien
		if lien != "" {
			ll = append(ll, fmt.Sprintf("%d -> %d [%s]", n.nucl.rang, m.rang, lien))
		}
	}
	return ll
}

// renvoie les lemmatisations réduites de tous les mots du Nod
func (n *Nod) toRes(m *Mot) gocol.Res {
	if n.nucl == m {
		return n.rnucl
	}
	for i, ra := range n.rra {
		if i == m.rang {
			return ra
		}
	}
	for i, rp := range n.rrp {
		if i == m.rang {
			return rp
		}
	}
	return nil
}
