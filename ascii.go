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
	5 -> 4
	6 -> 3
	6 -> 0
	6 -> 5

	doit afficher :

    ┌──────────────────────────────────────────┐
    │                        ┌────────────────┐│
    │┌────────────────┐      │    ┌───┐       ││
    ││        ┌──────┐│      │    │   │┌─────┐││
    V│        V      │V      V    V   │V     │││
Prometheus Iapeti filius homines ex luto finxit.

*/

import (
	"strconv"
	"strings"
)

type Word struct {
	d		int		// positions de départ et d'arrivée du mot
	gr		string
	len		int
	nba		int		// nombre d'arcs partant du mot ou y aboutissant
	rang	int
}

type Arc struct {
	motA	*Word
	motB	*Word
	label	string	// étiquette de l'arc
	dist	int		// distance abs(motA.rang - motB.rang)
	ecrit	bool
}

var (
	arcs	[]*Arc
	gabarit	[]rune
	lignes  []string
	mots	[]*Word
)

const (
	hh rune = '─'
    vv rune = '│'
	/*
	dr rune = '┌'
	dl rune = '┐'
	*/
	dr rune = '.'
	dl rune = '.'
	V  rune = 'V'
)

// trace l'arc a
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
		var k int
		for j := a.motA.d+1; j < a.motB.d; j++ {
			if k < len(a.label) {
				lignes[i] = place(lignes[i], rune(a.label[k]), j)
				k++
			} else {
				lignes[i] = place(lignes[i], hh, j)
			}
		}
		lignes[i] = place(lignes[i], dl, a.motB.d)
		// calcul des prochains points de départ/arrivée
		a.motA.d--
		a.motB.d++
	} else {
		var k int
		lignes[i] = place(lignes[i], dl, a.motA.d)
		for j := a.motB.d+1; j < a.motA.d; j++ {
			if k < len(a.label) {
				lignes[i] = place(lignes[i], rune(a.label[k]), j)
				k++
			} else {
				lignes[i] = place(lignes[i], hh, j)
			}
		}
		lignes[i] = place(lignes[i], dr, a.motB.d)
		// départ/arrivée
		a.motA.d++
		a.motB.d--
	}
}

// renvoie le mot le plus à droite de l'arc
func (a *Arc) dernier() *Word {
	if a.motA.rang > a.motB.rang {
		return a.motA
	}
	return a.motB
}

// réinitialisation des arcs
func initArcs() {
	arcs = nil
	gabarit = nil
	lignes = nil
	mots = nil
}

// vrai si aucun caractère autres que ' '
// n'est dans nl entre ma.d et mb.d
func libre(nl int, a int, b int) bool {
	for len(lignes) < nl+1 {
		lignes = append(lignes, string(gabarit))
		return true
	}

	var seg string
	if a < b {
		seg = lignes[nl][a:b]
	} else {
		seg = lignes[nl][b:a]
	}
	for i := 0; i < len(seg); i++ {
		if seg[i] != ' ' { //&& !strings.Contains(seg, "┐") {
			return false
		}
	}
	return true
}

// place le caractère ch à la position ou dans l
func place(l string, ch rune, ou int) string {
	rr := []rune(l)
	for ou >= len(rr) {
		rr = append(rr, ' ')
	}
	lg := rr[:ou]
	lg = append(lg, ch)
	ld := rr[ou+1:]
	lg = append(lg, ld...)
	return string(lg)
}

// graphe en ascii du code dot ll
func graphe(ll []string) []string {
	lm := strings.Split(ll[0], " ")
	// création des mots
	var report int
	for i, ecl := range lm {
		nm := new(Word)
		nm.gr = ecl
		nm.rang = i
		nm.len = len(ecl)
		nm.d = report + nm.len/2
		// calcul de la colonne de l'initiale du mot
		if i == 0 {
			report = len(ecl)
		} else {
			report += len(ecl) + 1
		}
		//nm.d = report - len(ecl)/2
		mots = append(mots, nm)
	}
	// création des arcs
	for i, l := range ll {
		if i == 0 {
			continue
		}
		na := new(Arc)
		// séparation arc - étiquette
		ecl := strings.Split(l, " [")
		if len(ecl) > 1 {
			na.label = ecl[1][:len(ecl[1])-1]
		}
		l = ecl[0]
		ecl = strings.Split(l, " -> ")
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
	lenll := len(ll[0])
	for i := 0; i<lenll; i++ {
		gabarit = append(gabarit, ' ')
	}
	//gabarit = strings.Repeat(" ", len(ll[0]))
	// ajout éventuel de la phrase
	if len(lignes) == 0 {
		lignes = append(lignes, ll[0])
		// ajout de la première ligne (V ou vv)
		lignes = append(lignes, string(gabarit))
	}
	// calcul des arcs, remplissage des lignes
	for i := 1; i <= len(mots); i++ {
		for _, a := range arcs {
			if !a.ecrit && a.dist == i {
				arcus(a)
				a.ecrit = true
			}
		}
	}
	// génération des lignes en commençant par le haut
	var retour []string
	for i := len(lignes)-1; i > -1; i-- {
		retour = append(retour, lignes[i])
	}
	return retour
}
