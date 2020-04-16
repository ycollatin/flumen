// groupe.go - Gentes

package main

// signets
// vaMorph

import (
	//"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

type Groupe struct {
	id              string
	noyaux, noyexcl []*Noy   // pos du noyau
	morph           []string // traits morpho du noyau
	lexsynt         []string // étiquettes lexicosyntaxiques du noyau
	exclls          []string // propriétés lexicosyntaxiques exclues
	ante            []*Sub   // éléments précédant le noyau
	post            []*Sub   // éléments suivant le noyau
}

var grpTerm, grp []*Groupe

// créateur de Groupe
func creeGroupe(ll []string) *Groupe {
	if len(ll) == 0 {
		return nil
	}
	g := new(Groupe)
	for _, l := range ll {
		kv := strings.Split(l, ":")
		k := kv[0]
		v := kv[1]
		switch k {
		case "ter", "grp":
			g.id = v
		case "pos":
			g.noyaux, g.noyexcl = creeNoy(v)
		case "morph":
			g.morph = strings.Split(v, ",")
		case "lexsynt":
			lecl := strings.Split(v, " ")
			for _, ecl := range lecl {
				if ecl[0] != '!' {
					g.lexsynt = append(g.lexsynt, ecl)
				} else {
					g.exclls = append(g.exclls, ecl[1:])
				}
			}
		case "a":
			g.ante = append(g.ante, creeSub(v, g, true))
		case "ag":
			g.ante = append(g.ante, creeSub(v, g, false))
		case "p":
			g.post = append(g.post, creeSub(v, g, true))
		case "pg":
			g.post = append(g.post, creeSub(v, g, false))
		}
	}
	return g
}

// lecture des groupes dans le fichier nf
func lisGroupes(nf string) {
	llin := gocol.Lignes(nf)
	var ll []string
	for _, l := range llin {
		deb := l[:4]
		if deb == "ter:" || deb == "grp:" {
			g := creeGroupe(ll)
			if g != nil {
				if ll[0][:4] == "grp:" {
					grp = append(grp, g)
				} else {
					grpTerm = append(grpTerm, g)
				}
				ll = nil
			}
		}
		ll = append(ll, l)
	}
	grp = append(grp, creeGroupe(ll))
}

func (g *Groupe) nbSubs() int {
	return len(g.ante) + len(g.post)
}

// la morpho morf est-elle compatible avec le noyau du groupe g ?
func (g *Groupe) vaMorph(morf string) bool {
	//debog := g.id=="v.sumpp" //&& morf=="indic, subj"
	//if debog {fmt.Println(" -vamorph",g.id,"morf",morf,"g.morph",g.morph)}
	for _, gmorf := range g.morph {
		//if debog {fmt.Println("   .vamorph, gmorf",gmorf)}
		va := true
		ecl := strings.Split(gmorf, " ")
		for _, trait := range ecl {
			va = va && strings.Contains(morf, trait)
			//if debog {fmt.Println("   .vamorph, morf",morf,"trait",trait,"va",va)}
		}
		if va {
			return true
		}
	}
	//if debog {fmt.Println("   .vamorph, false")}
	return false
}
