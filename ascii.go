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
	"strconv"
	"strings"
)

type Forme struct {
	d		int		// positions de départ et d'arrivée du mot
	//gr		string
	len		int
	nba		int		// nombre d'arcs partant du mot ou y aboutissant
	rang	int
}

type Arc struct {
	motA	*Forme
	motB	*Forme
	dist	int  // distance abs(motA.rang - motB.rang)
	ecrit	bool
}

var (
	arcs	[]*Arc
	gabarit string
	lignes  []string
	mots	[]*Forme
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
	lignes[1] = place(lignes[1], a.motA.d, vv)
	lignes[1] = place(lignes[1], a.motB.d, V)
	// placer les verticales si nécessaire
	i := 2
	for !libre(i, a.motA.d, a.motB.d) {
		lignes[i] = place(lignes[i], a.motA.d, vv)
		lignes[i] = place(lignes[i], a.motB.d, vv)
		i++
	}
	if a.motA.d < a.motB.d {
		lignes[i] = place(lignes[i], a.motA.d, dr)
		for j := a.motA.d+1; j < a.motB.d; j++ {
			lignes[i] = place(lignes[i], j, hh)
		}
		lignes[i] = place(lignes[i], a.motB.d, dl)
		// calcul des prochains points de départ/arrivée
		a.motA.d--
		a.motB.d++
	} else {
		lignes[i] = place(lignes[i], a.motA.d, dl)
		for j := a.motB.d+1; j < a.motA.d; j++ {
			lignes[i] = place(lignes[i], j, hh)
		}
		lignes[i] = place(lignes[i], a.motB.d, dr)
		// départ/arrivée
		a.motA.d++
		a.motB.d--
	}
}

func (a *Arc) dernier() *Forme {
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
func place(l string, ou int, ch rune) string {
	rr := []rune(l)
	rr[ou] = ch
	return string(rr)
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
		nm := new(Forme)
		//nm.gr = ecl
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
