//   main.go --	Gentes
//	analyseur syntaxique du latin

package main

// FIXME
// - esse inf a un objet !
// - tibiis facere : dat au lieu d'abl. quod Marsya tibii facere non potuit.
// - in certamen provocavit : certamen objet au lieu de gprep
// - consul populi Romani : Romani gén au lieu d'épith
// - tantum nocte crescebat : tantum nocte n.epith, nocte crescebat, v.obj.
// - Marcellinum *tibi esse* iratum scis : tibi devrait fixer esse à sum.
// - la distinction lemmatisation avant/après analyse synt. n'est pas complète.
// - subiciunt veribus prunas et viscera torrent :
//   AmbiguÏté entre la coord prunas et viscera    (faux)
//				  et la coord subiciunt et torrent (juste)
//
// XXX
// - beaucoup de confusions entre n+app et n+gén.
// - Trouver une solution pour la construction personnelle de l'infinitive
//   Homerus dicitur caecus fuisse. (ou caecum)
// - CONSTRUCTIONS LEXICALES
// 		* pietate erga te

// TODO
// - traiter la coordination par -que := et -
// - traiter de la même manière le noyau et les subs, aussi bien dans le code
//   que dans les données ?
// - tenir compte de la morpho unique (voluptatem. acc. sing.)
// - une version qui explore toutes les possibilités de groupes.
//   en privilégiant les groupes qui ont le plus de mots
// - un champ groupe.anrel - analyses du relatif ?
// - donner une POS distincte aux verbes intransitifs. v. gocol.indMorph
// - accord de personne sujet-verbe ?
// - saisie d'une phrase ?
// - fonction de sortie au format GraphViz
// - parasitage de /sum/ par /edo/ : comment supprimer "excl" dans lexsynt
// - parasitage de /do/ par /dato/ :   "

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	//"strings"

	"github.com/fatih/color"
	"github.com/ycollatin/gocol"
)

const (
	version = "Alpha"
	aidePh  = `    l->mot suivant ; h->mot précédent ;
    j->phrase suivante ; k->phrase précédente ;
    c->lemmatisation du mot courant ;
    a->arbre de la phrase ; r->retour; x->quitter`
	//s-> définir une suite morphosyntaxique ; x->Exit`
	//aideS =
	//`i-> id de la suite ; n-> n° du noyau ;
	// l-> liens (n°départ.fonction.n°sub,n°etc.);
	// v-> valider ; r-> retour`
)

var (
	ch, chData string // chemins du binaire et des données
	chCorpus   string // chemin du corpus
	dat        bool   // drapeau de chargement des données
	//module		string
	//modules		[]string
	rouge func(...interface{}) string
	texte *Texte
)

// affiche les arcs syntaxique de la phrase
func analyse(expl bool) {
	texte.affiche(aidePh)
	texte.tronc.explore()
	/*
		if expl {
			for _, n := range texte.tronc.nods {
				fmt.Println(n.doc())
			}
			fmt.Println("\n----- source ---\n",
				strings.Join(ar, "\n"),
				"\n----------------")
		}
		fmt.Println(strings.Join(gr, "\n"))
	*/
}

// choix du texte latin
func chxTexte() {
	ClearScreen()
	fmt.Println("Suites, grammaire latine")
	fmt.Println(" © Yves Ouvrard 2020, licence GPL3")
	texte = nil
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
		for i := 0; i < len(files); i++ {
			fmt.Println(i+1, textes[i])
		}
		chx = InputInt("n° du texte")
	}
	if chx < 0 {
		main()
	}
	if chx > len(textes) {
		chx = len(textes)
	}
	ftexte := textes[chx-1]
	texte = CreeTexte(ftexte)
	texte.majPhrase()
	texte.affiche(aidePh)
}

func lemmatise() {
	texte.affiche(aidePh)
	im := texte.tronc.imot
	texte.tronc.imot = im
	mc := texte.tronc.motCourant()
	fmt.Println("lemmatisation", rouge(mc.gr))
	/*
		if len(mc.ans2) > 0 {
			ll2 := gocol.Restostring(mc.ans2)
			fmt.Println(rouge(ll2))
			ll3 := strings.Split(ll2, "\n")
			ll := strings.Split(gocol.Restostring(mc.ans), "\n")
			for _, l := range ll {
				if !contient(ll3, l) {
					fmt.Println(l)
				}
			}
		} else {*/
	fmt.Println(gocol.Restostring(mc.ans))
	//}
}

func motprec() {
	if texte.tronc.imot > 0 {
		texte.tronc.imot--
		texte.affiche(aidePh)
	}
}

func motsuiv() {
	if texte.tronc.imot < len(texte.tronc.mots)-1 {
		texte.tronc.imot++
		texte.affiche(aidePh)
	}
}

func main() {
	// couleur
	rouge = color.New(color.FgRed, color.Bold).SprintFunc()
	// chemins
	dir, _ := os.Executable()
	ch = path.Dir(dir)
	chData = ch + "/data/"
	chCorpus = ch + "/corpus/"
	// lecture des données Collatinus
	gocol.Data(chData)
	// lecture des données syntaxiques
	lisGroupes(chData + "groupes.la")
	lisLexsynt()
	// choix du texte
	chxTexte()
	// capture des touches
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
		case "r":
			chxTexte()
		case "x":
			fmt.Println("\nVale")
			os.Exit(0)
		}
	}
	return
}
