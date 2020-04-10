//       mot.go - Publicola

package main

import (
	//"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

// signets :
//
// motnoeud
// motestnoyau
// motestSub
// motestSubde

// rappel de la lemmatisation dans gocol :
// type Sr struct {
//	Lem     *Lemme
//	Nmorph  []int
//	Morphos []string
// }
//
// type Res []Sr

type Lm struct {
	l	*gocol.Lemme
	m	string
}

type Mot struct {
	gr			string		// graphie du mot
	rang		int			// rang du mot dans la phrase à partir de 0
	ans, ans2	gocol.Res	// ensemble des lemmatisations, ans provisoire
	dejasub		bool		// le mot est déjà l'élément d'n nœud
	llm			[]Lm		// liste des lemmes ٍ+ morpho possibles
	tmpl, tmpm	int			// n°s provisoires de Sr et morpho
	pos			string		// id du groupe dont le mot est noyau
							// ou à défaut pos du mot, si elle est décidée
	lexsynt		[]string	// propriétés lexicosyntaxiques
}

func creeMot(m string) *Mot {
	mot := new(Mot)
	mot.gr = m
	var echec bool
	mot.ans, echec = gocol.Lemmatise(m)
	if echec {
		mot.ans, echec = gocol.Lemmatise(gocol.Majminmaj(m))
	}
	// ajout du genre pour les noms
	if !echec {
		for i, a := range mot.ans {
			mot.ans[i] = genus(a)
		}
	}
	return mot
}

func accord(lma, lmb, cgn string) bool {
va := true
for i:=0; i<len(cgn); i++ {
	switch cgn[i] {
	case 'c':
		k := cas(lma)
		va = va && strings.Contains(lmb, k)
	case 'g':
		g := genre(lma)
		va = va && strings.Contains(lmb, g)
	case 'n':
		n := nombre(lma)
			va = va && strings.Contains(lmb, n)
		}
	}
	return va
}

func (m *Mot) dejaNoy() bool {
	for _, n := range texte.phrase.nods {
		if n.nucl == m {
			return true
		}
	}
	return false
}

func domine(ma, mb *Mot) bool {
	return false
}

// teste si m peut être le noyau du groupe groupe g
func (m *Mot) estNoyau(g *Groupe) gocol.Res {
	//signet motestnoyau
	//debog := m.gr=="dedit" && g.id=="v.sov"
	//debog := m.gr=="dedit" && g.id=="v.sov"
	//if debog {fmt.Println(" -estNoyau",m.gr,g.id,"ans",len(m.ans),"pos",m.pos)}

	var ans3 gocol.Res
	// vérif du pos
	if m.pos != "" {
		// 1. La pos définitif n'est pas encore fixé
		va := false
		for _, noy := range g.noyaux {
			va = va || noy.vaPos(m.pos)
		}
		if !va {
			return ans3
		}
		ans3 = m.ans2
	} else {
		for _, a := range m.ans {
			for _, noy := range g.noyaux {
				if noy.canon > "" {
					if noy.vaSr(a) {
						ans3 = append(ans3, a)
						break
					}
				} else {
					if noy.vaPos(a.Lem.Pos) {
						ans3 = append(ans3, a)
						break
					}
				}
			}
		}
	}
	//if debog {fmt.Println("  .estNoyau, oka, len ans3", len(ans3))}

	// vérif lexicosyntaxique
	var ans4 gocol.Res
	for _, a := range ans3 {
		va := true
		for _, ls := range g.lexSynt {
			va = va && lexsynt(a.Lem.Gr[0], ls)
		}
		if va {
			ans4 = append(ans4, a)
		}
	}
	//if debog {fmt.Println("  .estNoyau, okb, len ans4",len(ans4))}

	// vérif morpho. Si aucune n'est requise, renvoyer true
	if len(g.morph) == 0 {
		return ans4
	}

	var ans5 gocol.Res
	for _, sr := range ans4 {
		var morfos []string  // morphos de sr acceptées par g
		for _, morf := range sr.Morphos {
			//if debog {fmt.Println("  .estNoyau, morf",morf,"g.morph",g.morph)}
			if g.vaMorph(morf) {
				morfos = append(morfos, morf)
			}
		}
		//if debog {fmt.Println("  .estNoyau, morfos",len(morfos))}
		if len(morfos) > 0 {
			sr.Morphos = morfos
			ans5 = append(ans5, sr)
		}
	}
	//if debog {fmt.Println("  .estNoyau, len ans5",len(ans5))}
	return ans5
}

// id des Nod dont m est déjà le noyau
func (m *Mot) estNuclDe() []string {
	var ret []string
	for _, nod := range texte.phrase.nods {
		if nod.nucl == m {
			ret = append(ret, nod.grp.id)
		}
	}
	return ret
}

// vrai si m est compatible avec Sub et le noyau mn
func (m *Mot) estSub(sub *Sub, mn *Mot) gocol.Res {
	//debog := sub.groupe.id=="v.conj" && m.gr == "imposuit" && mn.gr=="et"
	//if debog {fmt.Println(" -estSub m",m.gr,"pos",m.pos,"sub",sub.groupe.id,"mn",mn.gr)}
	// signet motestSub
	var ans2 gocol.Res
	// vérification des pos
	if m.pos != "" {
		//if debog {fmt.Println("   .estSub, !=")}
		// 1. La pos du mot est définitive
		va := false
		for _, noy := range sub.noyaux {
			//if debog {fmt.Println("   .estSub, noy",noy.id,"pos",m.pos)}
			va = va || noy.vaPos(m.pos)
		}
		if !va {
			return ans2
		}
		ans2 = m.ans2
		//if debog {fmt.Println("   .estSub, noy, len ans2",len(ans2))} //,ans2[0].Lem.Gr)}
	} else {
		//if debog {fmt.Println("   .estSub, else")}
		// 2. La pos définitif n'est pas encore fixée
		for _, an := range m.ans {
			//if debog {fmt.Println("   .estSub,lemme",an.Lem.Gr,"pos",an.Lem.Pos)}
			for _, noy := range sub.noyaux {
				if noy.canon > "" {
					//if debog {fmt.Println("   .estSub, noy.canon",noy.canon)}
					if noy.vaSr(an) {
						ans2 = append(ans2, an)
						break
					}
				} else {
					if noy.vaPos(an.Lem.Pos) {
						ans2 = append(ans2, an)
						break
					}
				}
			}
		}
	}
	//if debog {fmt.Println("   .estSub, len(ans2)",len(ans2))} //,"lemme",ans2[0].Lem.Gr)}
	if len(ans2) == 0 {
		return ans2
	}

	//morphologie
	var ans3 gocol.Res
	for _, an := range ans2 {
		//if debog {fmt.Println("   .estSub2",an.Lem.Gr,"morphos",len(an.Morphos))}
		var lmorf []string
		for _, morfs := range an.Morphos {
			//if debog {fmt.Println("   .estSub3 morfs",morfs,"sub",sub.morpho,sub.lien)}
			// pour toutes les morphos valides de m
			if sub.vaMorpho(morfs) {
				lmorf = append(lmorf, morfs)
			}
		}
		if len(lmorf) > 0 {
			an.Morphos = lmorf
			ans3 = append(ans3, an)
		}
	}
	//if debog {fmt.Println("   .estSub1, oka, len ans3",len(ans3))}

	// accord
	var ans4 gocol.Res
	// pour toutes les morphos valides de mn
	for _, ann := range mn.ans2 {
		for _, morfn := range ann.Morphos {
			// pour toutes les morphos valides de m
			for _, an := range ans3 {
				var lmorf []string
				for _, morfs := range an.Morphos {
					if accord(morfn, morfs, sub.accord) {
						lmorf = append(lmorf, morfs)
					}
				}
				if len(lmorf) > 0 {
					an.Morphos = lmorf
					ans4 = append(ans4, an)
				}
			}
		}
	}
	//if debog {fmt.Println("   .estSub3, sortie ans3",len(ans3))}
	return ans4
}

func (ma *Mot) estSubDe(mb *Mot) bool {
	// signet motestSubDe
	for _, n := range texte.phrase.nods {
		if mb == n.nucl {
			for _, m := range n.mma {
				if m == mb {
					return true
				}
			}
			for _, m := range n.mmp {
				if m == mb {
					return true
				}
			}
		}
	}
	return false
}

// ajoute le genre à la morpho d'un nom
func genus(sr gocol.Sr) gocol.Sr {
	if sr.Lem.Pos != "n" && sr.Lem.Pos != "NP" {
		return sr
	}
	inc := 12
	switch sr.Lem.Genre {
	case "féminin":
		inc += 12
	case "neutre":
		inc += 24
	}
	for i, _ := range sr.Nmorph {
		sr.Nmorph[i] += inc
	}
	return sr
}

/*
!FIXME
grp:v.conj
pos:v
a:"et";conj

!FIXME
grp:v.coord
pos:v
pg:v;coord;;n
*/

// si m peut être noyau d'un gourpe g, un Nod est renvoyé, sinon nil.
func (m *Mot) noeud(g *Groupe) *Nod {
	// signet motnoeud
	//debog := g.id=="n.coord" && m.gr == "Argiopes"
	//if debog {fmt.Println("-noeud",g.id,m.gr,"pos",m.pos)}
	rang := m.rang
	lante := len(g.ante)
	//if debog {fmt.Println("  .noeud, lante",len(g.ante))}
	// mot de rang trop faible
	if rang < lante {
		return nil
	}
	// ou trop élevé
	if texte.phrase.nbmots - rang < len(g.post) {
		return nil
	}
	// m peut-il être noyau du groupe g ?
	res2 := m.estNoyau(g)
	if len(res2) == 0 {
		return nil
	}
	//if debog {fmt.Println(" .noeud oka, res2", len(res2),res2[0].Lem.Gr)}
	m.ans2 = res2
	// création du noeud de retour
	nod := new(Nod)
	nod.grp = g
	nod.nucl = m
	nod.rang = rang
	// vérif des subs
	// ante
	r := rang - 1
	//if debog {fmt.Println("  .noeud okb",lante,"lante, r",r,nod.doc())}
	// reгcherche rétrograde des subs ante
	for ia := lante-1; ia > -1; ia-- {
		if r < 0 {
			// le rang du mot est < 0 : impossible
			return nil
		}
		sub := g.ante[ia]
		ma := texte.phrase.mots[r]
		//if debog {fmt.Println("  .noeud, avant dejasub",ma.gr,"lien",sub.lien,"ia",ia,"r",r)}
		// passer les mots déjà subordonnés
		for ma.dejasub {
			r--
			if r < 0 {
				return nil
			}
			ma = texte.phrase.mots[r]
		}
		//if debog {fmt.Println(" .noeud ma",ma.gr,"estSub",m.gr,"grup",sub.groupe.id)}
		// vérification de réciprocité, puis du lien lui-même
		if m.estSubDe(ma) {
			return nil
		}
		res3 := ma.estSub(sub, m)
		//if debog {fmt.Println(" .noeud estSub, res3",len(res3))} //,res3[0].Lem.Gr)}
		if len(res3) == 0 {
			//if debog {fmt.Println("  .noeud ma",ma.gr,"n'est pas sub",sub.lien,"de",m.gr)}
			return nil
		}
		nod.mma = append(nod.mma, ma)
		r--
		//if debog {fmt.Println("    vu",ma.gr)}
	}
	//if debog {fmt.Println("  .noeud okd",len(g.post),"g.post, rang",rang,"nbmots",texte.phrase.nbmots)}
	// post
	for ip, sub := range g.post {
		r := rang + ip + 1
		if r >= texte.phrase.nbmots {
			break
		}
		mp := texte.phrase.mots[r]
		//if debog {fmt.Println("  .noeud avant dejasub, post, mp",mp.gr,"dejasub",mp.dejasub)}
		for mp.dejasub && r < len(texte.phrase.mots) -1 {
			r++
			mp = texte.phrase.mots[r]
		}
		//if debog {fmt.Println("  .noeud apres dejasub mp", mp.gr)}
		if m.estSubDe(mp) {
			return nil
		}
		res4 := mp.estSub(sub, m)
		if len(res4) == 0 {
			return nil
		}
		//mp.sub2 = sub
		nod.mmp = append(nod.mmp, mp)
		r++
	}
	// fixer les pos et sub des mots du noeud
	if len(nod.mma) + len(nod.mmp) > 0 {
		m.pos = g.id
		//fmt.Println("   .noeud, noy",m.gr,"pos",m.pos)
		for _, m := range nod.mma {
			m.dejasub = true
		}
		for _, m := range nod.mmp {
			m.dejasub = true
		}
		//if debog {fmt.Println("   .noeud m.ans2", len(m.ans2),m.ans2[0].Lem.Gr)}
		return nod
	}
	return nil
}

/*
// vrai si ma est sub de mb
func (ma *Mot) SubDe(mb *Mot) bool {
	// chercher le groupe dont mb est noyau
	for _, n := range texte.phrase.nods {
		if mb == n.nucl {
			//return ma.elDe(n)
			for _, m := range n.mma {
				if m == ma {
					return true
				}
			}
			for _, m := range n.mmp {
				if m == ma {
					return true
				}
			}
		}
	}
	return false
}
*/
