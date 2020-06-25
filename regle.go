// regle.go - Flumen

package main

// signets
// vaMorph

import (
	"strings"

	"github.com/ycollatin/gocol"
)

type Regle struct {
	id             string
	nucl, nuclexcl *El
	morph          []string // traits morpho du noyau
	lexsynt        []string // étiquettes lexicosyntaxiques du noyau
	ante           []*El    // éléments précédant le noyau
	post           []*El    // éléments suivant le noyau
	multi          bool     // le groupe possède au moins un sub multi
	nbsubs         int      // nombre de subs
}

var grpTerm, grp []*Regle

// créateur de Groupe
func creeRegle(ll []string) *Regle {
	if len(ll) == 0 {
		return nil
	}
	g := new(Regle)
	for _, l := range ll {
		kv := strings.Split(l, ":")
		k := kv[0]
		v := kv[1]
		switch k {
		case "ter", "grp":
			g.id = v
		case "n":
			// test sur v.objsuj et n.sr
			g.nucl = creeEl(k, v, g)
		case "morph":
			g.morph = strings.Split(v, ",")
		case "lexsynt":
			lecl := strings.Split(v, " ")
			for _, ecl := range lecl {
				g.lexsynt = append(g.lexsynt, ecl)
			}
		case "a":
			g.ante = append(g.ante, creeEl(k, v, g))
		case "ag":
			g.ante = append(g.ante, creeEl(k, v, g))
		case "p":
			g.post = append(g.post, creeEl(k, v, g))
		case "pg":
			g.post = append(g.post, creeEl(k, v, g))
		}
	}
	// multi
	for _, sa := range g.ante {
		if sa.multi {
			g.multi = true
		}
	}
	for _, sp := range g.post {
		if sp.multi {
			g.multi = true
		}
	}
	g.nbsubs = len(g.ante) + len(g.post)
	return g
}

// Renvoie vrai si, dans le même ordre, tous les
// éléments de ga portent le même lien que
// les éléments correspondants de gb.
func (ga *Regle) equiv(gb *Regle) bool {
	if ga.id == gb.id {
		return true
	}
	if len(ga.ante) != len(gb.ante) {
		return false
	}
	if len(ga.post) != len(gb.post) {
		return false
	}
	for i, el := range ga.ante {
		if gb.ante[i].lien != el.lien {
			return false
		}
	}
	for i, el := range ga.post {
		if gb.post[i].lien != el.lien {
			return false
		}
	}
	return true
}

// lecture des groupes dans le fichier nf
func lisRegles(nf string) {
	llin := gocol.Lignes(nf)
	var ll []string
	for _, l := range llin {
		vk := strings.Split(l, ":")
		deb := vk[0]
		if deb == "ter" || deb == "grp" {
			g := creeRegle(ll)
			if g != nil {
				if PrimEl(ll[0], ":") == "grp" {
					grp = append(grp, g)
				} else {
					grpTerm = append(grpTerm, g)
				}
				ll = nil
			}
		}
		ll = append(ll, l)
	}
	grp = append(grp, creeRegle(ll))
}

// la morpho morf est-elle compatible avec le noyau du groupe g ?
func (g *Regle) vaMorph(morf string) bool {
	for _, gmorf := range g.morph {
		va := true
		ecl := strings.Split(gmorf, " ")
		for _, trait := range ecl {
			va = va && strings.Contains(morf, trait)
		}
		if va {
			return true
		}
	}
	return false
}
