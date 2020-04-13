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
    ▽│        ▽      │▽      ▽    ▽   │▽     │││
Prometheus Iapeti filius homines ex luto finxit.

*/

import (
	"strconv"
	"strings"
)

type Word struct {
	gr		string
	len		int
	nba		int		// nombre d'arcs partant du mot ou y aboutissant
	rang	int
	d, f	int		// pos de début et fin de mot
	Pg		int		// pos de départ et d'arrivée gauche du mot
	Pd		int		// et droite
	dirder	int		// le dernier arc arrive de ou part vers la g.(-1) ou la dr(1)
}

type Arc struct {
	motA	*Word	// mot de départ (noyau)
	motB	*Word	// mot d'arrivée (sub)
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
	dr rune = '┌'
	dl rune = '┐'
	//dr rune = '.'
	//dl rune = '.'
	V  rune = '▽'
)

// décrémente, incrémente les points de départ/arrivée
func (w *Word) decg() {
	if w.Pg > w.d {
		w.Pg--
	}
}

func (w *Word) incg() {
	if w.Pg < w.f {
		w.Pg++
	}
}

func (w *Word) decd() {
	if w.Pd > w.d {
		w.Pd--
	}
}

func (w *Word) incd() {
	if w.Pd < w.f {
		w.Pd++
	}
}

// trace l'arc a
func arcus(a *Arc) {
	var arrb, arra int
	// si le dernier arc partait vers la gauche
	//   et que le nouveau part vers la gauche
	if a.motA.rang < a.motB.rang { // part vers la droite
		// pointѕ de départ et d'arrivée
		if a.motA.dirder > 0 {
			a.motA.decg()
			arra = a.motA.Pg
		} else {
			arra = a.motA.Pd
		}
		if a.motB.dirder <= 0 {
			a.motB.incd()
			arrb = a.motB.Pd
		} else {
			arrb = a.motB.Pg
		}
		// calcul des prochains points de départ/arrivée
		a.motA.dirder = 1
		a.motB.dirder = -1
		// première ligne : départ vv et arrivée V
		lignes[1] = place(lignes[1], vv, arra)
		lignes[1] = place(lignes[1], V, arrb)
		// placer les verticales si nécessaire
		i := 2
		for !libre(i, arra, arrb) {
			lignes[i] = place(lignes[i], vv, arra)
			lignes[i] = place(lignes[i], vv, arrb)
			i++
		}
		lignes[i] = place(lignes[i], dr, arra)
		var k int
		for j := arra+1; j < arrb; j++ {
			if k < len(a.label) {
				lignes[i] = place(lignes[i], rune(a.label[k]), j)
				k++
			} else {
				lignes[i] = place(lignes[i], hh, j)
			}
		}
		lignes[i] = place(lignes[i], dl, arrb)
	} else { // part vers la gauch
		if a.motA.dirder < 0 {
			a.motA.Pd++
			arra = a.motA.Pd
		} else {
			arra = a.motA.Pg
		}
		if a.motB.dirder >= 0 {
			a.motB.Pg--
			arrb = a.motB.Pg
		} else {
			//a.motB.Pd++
			arrb = a.motB.Pd
		}
		a.motA.dirder = -1
		a.motB.dirder = 1
		lignes[1] = place(lignes[1], V, arrb)
		lignes[1] = place(lignes[1], vv, arra)
		i := 2
		for !libre(i, arra, arrb) {
			lignes[i] = place(lignes[i], vv, arra)
			lignes[i] = place(lignes[i], vv, arrb)
			i++
		}
		var k int
		lignes[i] = place(lignes[i], dl, arra)
		for j := arrb+1; j < arra; j++ {
			if k < len(a.label) {
				lignes[i] = place(lignes[i], rune(a.label[k]), j)
				k++
			} else {
				lignes[i] = place(lignes[i], hh, j)
			}
		}
		lignes[i] = place(lignes[i], dr, arrb)
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
	if ou < 0 {
		ou = 0
	}
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
		nm.Pg = report + nm.len/4
		nm.Pd = nm.Pg + nm.len/2
		// calcul de la colonne de l'initiale du mot
		nm.d = report
		if i == 0 {
			report = nm.len
		} else {
			report += nm.len + 1
		}
		nm.f = nm.d + nm.len
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
	// en commençant par les plus courts
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
