//   main.go --	Gentes
//	analyseur syntaxique du latin

package main

// FIXME
// - l'exclusion des subs ne marche pas
// - se in mare praecipitauit : conflit sujet-prepAcc
// - eo Romam : une étiquette lexsynt ? pour les 2 membres ?
// - Immo locuples erat, muet.
// - Lentulus consul fuit : Lentulus sujet absent.
// - Ceres cum vellet : alumnum 2 fois sub.
// - dicitur prima mortalis : l'attr dicitur mortalis
//   + épith prima mortalis non reconnue
//   mais reconnue si dicitur a un sujet
// - neque in perpetuum servare sciebant :
//	 servare infobj de sciebant est reconnu, mais
//   sans ses subs.
// - Solutions redondantes : sans doute règles redondantes.
// - subiciunt veribus prunas et viscera torrent :
//   AmbiguÏté entre la coord prunas et viscera    (faux)
//				  et la coord subiciunt et torrent (juste)
//	 prunas 2 fois objet
//
// XXX
// - Trouver une solution pour la construction personnelle de l'infinitive
//   Homerus dicitur caecus fuisse. (ou caecum)
// - CONSTRUCTIONS LEXICALES
// 		* pietate erga te

// TODO
// - un journal d'analyse : niveau, nœud avec mots
// - surlignage des lemmatisations : la récolte doit aussi rapporter les nods 
//   des branches terminales
// - éviter une réanalyse ?
// - traiter la coordination par -que := et -
// - un champ groupe.anrel - analyses du relatif ?
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
	"strings"

	"github.com/fatih/color"
	"github.com/ycollatin/gocol"
)

const (
	version = "Alpha"
	aidePh  = `    l->mot suivant ; h->mot précédent ;
    j->phrase suivante ; k->phrase précédente ;
    c->lemmatisation du mot courant ;
    a->arbre de la phrase ; g->arbre ٍ& sa source ;
	s->solution suivante ; p->solution précédente ;
	r->retour; x->quitter.`
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
	ibr   int
	rouge func(...interface{}) string
	texte *Texte
)

// affiche les arcs syntaxique de la phrase
func analyse(expl bool, j bool) {
	texte.majPhrase()
	texte.affiche(aidePh)
	tronc.explore()
	recolte := tronc.recolte()
	if recolte == nil {
		fmt.Println("échec de l'analyse")
		return
	}
	if ibr < 0 {
		ibr = 0
	}
	if ibr >= len(recolte) {
		ibr = len(recolte) - 1
	}
	nods := recolte[ibr]
	// graphe
	var src []string
	src = append(src, tronc.gr)
	for _, n := range nods {
		src = append(src, n.graf()...)
	}
	if expl {
		for _, n := range nods {
			fmt.Println(n.doc())
		}
		fmt.Println("\n---- source ----\n")
		fmt.Println(strings.Join(src, "\n"))
	}
	if j {
		fmt.Println("--- journal ----")
		fmt.Println(strings.Join(journal, "\n"))
	}
	fmt.Println("----------------")
	// numérotation de la solution
	fmt.Printf("%d/%d\n", ibr+1, len(recolte))
	// graphe en arcs
	initArcs()
	fmt.Println(strings.Join(graphe(src), "\n"))
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
	im := tronc.imot
	tronc.imot = im
	mc := tronc.motCourant()
	fmt.Println("lemmatisation", rouge(mc.gr))
	fmt.Println(restostring(mc.ans))
}

func motprec() {
	if tronc.imot > 0 {
		tronc.imot--
		texte.affiche(aidePh)
	}
}

func motsuiv() {
	if tronc.imot < len(mots)-1 {
		tronc.imot++
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
	var modeA bool
	var modeJ bool
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
			ibr = 0
		case "k":
			texte.retro()
			ibr = 0
		case "c":
			lemmatise()
		case "a":
			analyse(false, false)
			modeA = false
			modeJ = false
		case "g":
			analyse(true, false)
			modeA = true
			modeJ = false
		case "d":
			analyse(true, true)
			modeJ = true
		case "p":
			ibr--
			analyse(modeA, modeJ)
		case "s":
			ibr++
			analyse(modeA, modeJ)
		case "r":
			chxTexte()
		case "x":
			fmt.Println("\nVale")
			os.Exit(0)
		}
	}
	return
}
