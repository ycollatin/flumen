package main

/*
	Représentation ascii des arcs syntaxiques d'une phrase
	Données en entrée :
	- une liste ordonnée de mots
	- une liste d'arcs partant d'un mot A et aboutissant à un mot B

	exemple :

	     0        1     2       3    4   5     6
	Prometheus Iapeti filius homines ex luto finxit
	0 -> 2
	2 -> 1
	4 -> 5
	6 -> 0
	6 -> 3
	6 -> 4

	doit afficher :

    ┌───────────────────────────────────────┐						   
    │                       ┌──────────────┐│ 
    │┌───────────────┐      │    ┌────────┐││
    ││       ┌──────┐│	    │    │┌──┐    │││
	V│       V      │V      V    V│  V    │││
Prometheus Iapeti filius homines ex luto finxit
*/

import (
	//"fmt"
	"strconv"
	"strings"
)

type Word struct {
	d		int		// positions de départ et d'arrivée du mot
	len		int
	nba		int		// nombre d'arcs partant du mot ou y aboutissant
	rang	int
}

type Arc struct {
	motA	*Word
	motB	*Word
	dist	int  // distance abs(motA.rang - motB.rang)
	ecrit	bool
}

var (
	arcs	[]*Arc
	gabarit string
	lignes  []string
	mots	[]*Word
)

const (
	hh rune = '─'
    vv rune = '│'
	dr rune = '┌'
	dl rune = '┐'
	V  rune = 'V'
)

func arcus(a *Arc) {
	// première ligne : départ vv et arrivée V
	lignes[1] = place(lignes[1], vv, a.motA.d)
	lignes[1] = place(lignes[1], V, a.motB.d)
	// placer les verticales si nécessaire
	i := 2
	for !libre(i, a.motA.d, a.motB.d) {
		lignes[i] = place(lignes[i], vv, a.motA.d)
		lignes[i] = place(lignes[i], vv, a.motB.d)
		i++
	}
	if a.motA.d < a.motB.d {
		lignes[i] = place(lignes[i], dr, a.motA.d)
		for j := a.motA.d+1; j < a.motB.d; j++ {
			lignes[i] = place(lignes[i], hh, j)
		}
		lignes[i] = place(lignes[i], dl, a.motB.d)
		// calcul des prochains points de départ/arrivée
		a.motA.d--
		a.motB.d++
	} else {
		lignes[i] = place(lignes[i], dl, a.motA.d)
		for j := a.motB.d+1; j < a.motA.d; j++ {
			lignes[i] = place(lignes[i], hh, j)
		}
		lignes[i] = place(lignes[i], dr, a.motB.d)
		// départ/arrivée
		a.motA.d++
		a.motB.d--
	}
}

func (a *Arc) dernier() *Word {
	if a.motA.rang > a.motB.rang {
		return a.motA
	}
	return a.motB
}

// vrai si aucun caractère autres que ' '
// n'est dans nl entre ma.d et mb.d
func libre(nl int, a int, b int) bool {
	for len(lignes) < nl+1 {
		lignes = append(lignes, gabarit)
		return true
	}
	var seg string
	if a < b {
		seg = lignes[nl][a:b]
	} else {
		seg = lignes[nl][b:a]
	}
	for i := 0; i < len(seg); i++ {
		if seg[i] != ' ' {
			return false
		}
	}
	return true
}

// place le caractère ch à la position ou dans l
func place(l string, ch rune, ou int) string {
	/*
	lg := l[:ou]
	ld := l[ou+1:]
	return fmt.Sprintf("%s%s%s", lg, "@", ld)
	*/
	lr := []rune(l)
	lg := lr[:ou]
	ld := lr[ou+1:]
	lg = append(lg, ch)
	lg = append(lg, ld...)
	return string(lg)
}

func graphe(ll []string) []string {
	/*
	ll := []string {
		"Prometheus Iapeti filius homines ex luto finxit",
		"0 -> 2",
		"2 -> 1",
		"4 -> 5",
		"6 -> 0",
		"6 -> 3",
		"6 -> 4",
	}
	*/
	lm := strings.Split(ll[0], " ")
	// création des mots
	var report int
	for i, ecl := range lm {
		nm := new(Word)
		nm.rang = i
		nm.len = len(ecl)
		// calcul de la colonne de l'initiale du mot
		if i == 0 {
			report = len(ecl)
		} else {
			report += len(ecl) + 1
		}
		nm.d = report - len(ecl)/2
		mots = append(mots, nm)
	}
	// création des arcs
	for i := 1; i < len(ll); i++ {
		l := ll[i]
		// suppression provisoire de l'étiquette
		pcr := strings.Index(l, " [")
		if pcr > 4 {
			l = l[:pcr]
		}
		ecl := strings.Split(l, " -> ")
		na := new(Arc)
		ia, _ := strconv.Atoi(ecl[0])
		ib, _ := strconv.Atoi(ecl[1])
		na.motA = mots[ia]
		na.motB = mots[ib]
		dif := na.motA.rang - na.motB.rang
		if dif < 0 {
			dif = -dif
		}
		na.dist = dif
		arcs = append(arcs, na)
	}

	// gabarit des lignes où sont tracés les arcs
	gabarit = strings.Repeat(" ", len(ll[0]))
	// ajout de la phrase
	lignes = append(lignes, ll[0])
	// ajout de la première ligne (V ou vv)
	lignes = append(lignes, gabarit)
	// calcul des arcs, remplissage des lignes
	for i := 1; i <= len(mots); i++ {
		for _, a := range arcs {
			if !a.ecrit && a.dist == i {
				arcus(a)
				a.ecrit = true
			}
		}
	}
	return lignes
}
