//   main.go --	Gentes
//	analyseur syntaxique du latin

package main

// FIXME pas de lemmatisation pour la première phrase
// TODO Omission du pos d'un sub dans groupes.la
// TODO La lemmatisation ne peut pas toujours être fixée.
// 		par l'adoption d'un noeud. Cette adoption
//		peut soit éliminer quelques lemmatisations,
//		soit n'en éliminer aucune.
// TODO filtre lexsynt pour les subs

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
	c->lemmatisation du mot courant ;
	a->arbre de la phrase ; x->quitter`
	//s-> définir une suite morphosyntaxique ; x->Exit`
	//aideS =
	//`i-> id de la suite ; n-> n° du noyau ; 
	// l-> liens (n°départ.fonction.n°sub,n°etc.);
	// v-> valider ; r-> retour`
)

var (
	ch, chData	string	// chemins du binaire et des données
	chCorpus	string	// chemin du corpus
	dat			bool	// drapeau de chargement des données
	module		string
	modules		[]string
	rouge		func(...interface{}) string
	texte		*Texte
)

// affiche les arcs syntaxique de la phrase
func analyse(expl bool) {
	texte.majPhrase()
	ar, _ := phrase.arbre()
	gr := graphe(ar)
	if expl {
		for _, n := range phrase.nods {
			fmt.Println(n.doc())
		}
	}
	for _, lin := range gr {
		fmt.Println(lin)
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
	nbf := len(files)
	chx := 1
	if nbf > 1 {
		for i:= 0; i < len(files); i++ {
			fmt.Println(i+1, textes[i])
		}
		chx = InputInt("n° du texte")
	}
	ftexte := textes[chx-1]
	texte = CreeTexte(ftexte)
	texte.majPhrase()
}

func lemmatise() {
	fmt.Println("lemmatisation", rouge(phrase.mots[phrase.imot].gr))
	fmt.Println(gocol.Restostring(phrase.mots[phrase.imot].ans))
}

func motprec() {
	if texte == nil {
		txtNil()
		return
	}
	if phrase.imot > 0 {
		phrase.imot--
		texte.affiche(aidePh)
	}
}

func motsuiv() {
	if texte == nil {
		txtNil()
		return
	}
	if phrase.imot < len(phrase.mots)-1 {
		phrase.imot++
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
	//chData = path.Dir(dir) + "/data/"
	chData = ch + "/data/"
	chCorpus = ch + "/corpus/"
	go gocol.Data(chData)
	// lecture des données syntaxiques
	lisGroupes(chData+"groupes.la")
	lisLexsynt()
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
			texte.porro()
		case "k":
			texte.retro()
		case "c":
			lemmatise()
		case "a":
			analyse(false)
		case "g":
			analyse(true)
		case "x":
			fmt.Println("\nVale")
			os.Exit(0)
		}
	}
	return
}
