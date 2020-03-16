//   main.go --	Gentes
/*
	analyseur syntaxique du latin

	Phrase de test :
	Prometheus Iapeti filius homines ex luto finxit.
		0		1		2		3	  4	  5		6
*/

// main.go - Publicola

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/ycollatin/gocol"
)

const (
	version = "Alpha"
	aidePh =
	`l->mot suivant ; h->mot précédent ;
	j->phrase suivante ; k->phrase précédente ;
	a->arbre de la phrase ; x->quitter`
	//s-> définir une suite morphosyntaxique ; x->Exit`
	//aideS =
	//`i-> id de la suite ; n-> n° du noyau ; 
	// l-> liens (n°départ.fonction.n°sub,n°etc.);
	// v-> valider ; r-> retour`
)

var (
	ch, chData	string	// chemins du binaire et des données
	dat			bool	// drapeau de chargement des données
	imot		int		// n° de mot
	ip			int		// n° de phrase
	module		string
	modules		[]string
	rouge		func(...interface{}) string
	texte		*Texte
)

func analyse() {
	//fmt.Println(phrase.arbre())
	lingraphe := graphe(phrase.arbre())
	for i := len(lingraphe)-1; i > -1; i-- {
		fmt.Println(lingraphe[i])
	}
}

// choix du texte latin
func chxTexte() {
	files, err := ioutil.ReadDir(ch + "/corpus/")
	if err != nil {
		fmt.Println("Répertoire", ch+"/corpus/", "introuvable")
		return
	}
	textes := []string{}
	for _, fileInfo := range files {
		textes = append(textes, fileInfo.Name())
	}
	for i:= 0; i < len(files); i++ {
		fmt.Println(i+1, textes[i])
	}
	chx := InputInt("n° du texte")
	ftexte := textes[chx-1]
	texte = CreeTexte(ftexte)
	texte.affiche(aidePh)
	majPhrase()
}

func motprec() {
	if texte == nil {
		txtNil()
		return
	}
	if imot > 0 {
		imot--
		texte.affiche(aidePh)
	}
}

func motsuiv() {
	if texte == nil {
		txtNil()
		return
	}
	if imot < len(phrase.mots)-1 {
		imot++
		texte.affiche(aidePh)
	}
}

func porro() {
	if texte == nil {
		txtNil()
		return
	}
	if ip < len(texte.phrases) {
		ip++
		majPhrase()
		imot = 0
		texte.affiche(aidePh)
	}
}

func retro() {
	if texte == nil {
		txtNil()
		return
	}
	if ip > 0 {
		ip--
		majPhrase()
		imot = 0
		texte.affiche(aidePh)
	}
}

func txtNil() {
	fmt.Println("Il faut d'abord charger un texte : commande txt")
}

func main() {
	ClearScreen()
    fmt.Println("Suites, grammaire latine")
    fmt.Println("Yves Ouvrard, GPL3")

	// couleur
	rouge = color.New(color.FgRed, color.Bold).SprintFunc()

	// lecture des données Collatinus
	dir, _ := os.Executable()
	ch = path.Dir(dir)
	chData = path.Dir(dir) + "/data/"
	go gocol.Data(chData)
	// lecture des données syntaxiques
	lisGroupes(chData+"groupes.la")
	// choix du texte
	chxTexte()
	texte.affiche(aidePh)
	for {
		k := GetKey()
		switch k {
		case "l":
			motsuiv()
		case "h":
			motprec()
		case "j":
			porro()
		case "k":
			retro()
		case "a":
			analyse()
		case "x":
			fmt.Println("\nVale")
			os.Exit(0)
		}
	}
	return
}
