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
	//"fmt"
	"strconv"
	"strings"
)

type Word struct {
	gr   string
	len  int
	rang int
	d, f int // pos de début et fin de mot
	Pg   int // pos de départ et d'arrivée gauche du mot
	Pd   int // et droite
}

type Arc struct {
	motA  *Word  // mot de départ (noyau)
	motB  *Word  // mot d'arrivée (sub)
	label string // étiquette de l'arc
	dist  int    // distance abs(motA.rang - motB.rang)
	ecrit bool
}

var (
	arcs    []*Arc
	gabarit []rune
	lignes  []string
	words   []*Word
)

const (
	hh rune = '─'
	vv rune = '│'
	dr rune = '┌'
	dl rune = '┐'
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
		arra = a.motA.Pd
		a.motA.decd()
		arrb = a.motB.Pg
		a.motB.incg()
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
		for j := arra + 1; j < arrb; j++ {
			if k < len(a.label) {
				lignes[i] = place(lignes[i], rune(a.label[k]), j)
				k++
			} else {
				lignes[i] = place(lignes[i], hh, j)
			}
		}
		lignes[i] = place(lignes[i], dl, arrb)
	} else { // part vers la gauche
		arra = a.motA.Pg
		a.motA.incg()
		arrb = a.motB.Pd
		a.motB.decd()
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
		for j := arrb + 1; j < arra; j++ {
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

// réinitialisation des arcs
func initArcs() {
	arcs = nil
	gabarit = nil
	lignes = nil
	words = nil
}

// vrai si aucun caractère autres que ' '
// n'est dans nl entre ma.d et mb.d
func libre(nl int, a int, b int) bool {
	for nl > len(lignes)-1 {
		lignes = append(lignes, string(gabarit))
	}
	var seg []rune
	runes := []rune(lignes[nl])
	if a < b {
		seg = runes[a+1 : b-1]
	} else {
		if a-b > 1 {
			seg = runes[b+1 : a-1]
		} else {
			return true
		}
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
		// calcul de la colonne de l'initiale du mot
		nm.d = report
		nm.f = nm.d + nm.len
		report += nm.len + 1
		// points de départ et d'arrivée des arcs
		switch nm.len {
		case 1, 2:
			nm.Pg = nm.d
			nm.Pd = nm.d
		case 3:
			nm.Pg = nm.d
			nm.Pd = nm.f
		case 4:
			nm.Pg = nm.d + 1
			nm.Pd = nm.d + 2
		default:
			nm.Pg = nm.d + nm.len/4
			nm.Pd = nm.Pg + nm.len/2
		}
		words = append(words, nm)
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
		na.motA = words[ia]
		na.motB = words[ib]
		dif := na.motA.rang - na.motB.rang
		if dif < 0 {
			dif = -dif
		}
		na.dist = dif
		arcs = append(arcs, na)
	}

	// gabarit des lignes où sont tracés les arcs
	lenll := len(ll[0])
	for i := 0; i < lenll; i++ {
		gabarit = append(gabarit, ' ')
	}
	// ajout éventuel de la phrase
	if len(lignes) == 0 {
		lignes = append(lignes, ll[0])
		// ajout de la première ligne (V ou vv)
		lignes = append(lignes, string(gabarit))
	}
	// calcul des arcs, remplissage des lignes
	// en commençant par les plus courts
	for i := 1; i <= len(words); i++ {
		for _, a := range arcs {
			if !a.ecrit && a.dist == i {
				arcus(a)
				a.ecrit = true
			}
		}
	}
	// génération des lignes en commençant par le haut
	var retour []string
	for i := len(lignes) - 1; i > -1; i-- {
		retour = append(retour, lignes[i])
	}
	return retour
}
