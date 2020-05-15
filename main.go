//   main.go --	Gentes
//	analyseur syntaxique du latin

package main

// FIXME
//
// Atticae plurimam salutem : Atticae vu comme un génitif
//
// Pyrrha dicitur esse creata. boucle infinie, due au lexsynt dico2:...,attrp
// Pour d'autres phrases aussi.
// XXX
//
// Attention ! l'étiquette 'n' ne désigne pas un pos, mais une famille de groupes.
// Si un nom nominatif est requis, il faut écrire @n
//
// - CONSTRUCTIONS LEXICALES chez Cic.
// 		* pietate erga te
//      * tempus agendi et cogitandi
// TODO
// - une option pour charger une base de groupes différente ?
// - dicitur prima mortalis. avec lexsynt dico1:attrp, boucle infinie.
// - verbes se construisant avec locatif. Caesarem Sinuessae mansurum nuntiabant.
// - groupes isolés : Bene hercle faciunt.
// - Pour éviter la pléthore :
//   - durcir les conditions des règles
//   - élaguer les branches avant la fin
//   - élaguer les branches après récolte
//   - utiliser les goroutines : go bf.explore()
//   - enrichir la syntaxe de groupes.la
// - Si pléthore, trouver un moyen de navigation en fixant des arcs.
// - factoriѕer la négation ? neg(ch string) string {}
// - factoriser les " : guil(ch string) string {}
// - surlignage des lemmatisations : la récolte doit aussi rapporter les nods
//   des branches terminales
// - traiter la coordination par -que := et -
// - un champ groupe.anrel - analyses du relatif ?
// - accord :de personne sujet-verbe et verbe-verbe coord; de mode v-v, etc.
// - fonction de sortie au format GraphViz.
// - sortie de lemmatisation exacte.
// - parasitage de /sum/ par /edo/ : comment supprimer "excl" dans lexsynt
// - parasitage de /do/ par /dato/ :   "

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/ycollatin/gocol"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

const (
	version = "Alpha"
	entete  = "Gentes - grammaire latine"
	licence = " licence GPL3 © Yves Ouvrard 2020"
	aidePh  = `l->mot suivant ; h->mot précédent ;
j->phrase suivante ; k->phrase précédente ;
c->lemmatisation du mot courant ;
a->arbre de la phrase ; g->arbre ٍ& sa source ;
s->solution suivante ; p->solution précédente ;
r->retour; x->quitter.`
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
	texte.affiche(aidePh)
	if tronc.vendange == nil {
		texte.majPhrase()
		tronc.explore()
		tronc.recolte()
		tronc.elague()
		// tri
		sort.SliceStable(tronc.vendange, func(i, j int) bool {
			return tronc.vendange[i].nbarcs < tronc.vendange[j].nbarcs
		})
	}
	if tronc.vendange == nil {
		fmt.Println("échec de l'analyse")
		return
	}
	if ibr < 0 {
		ibr = 0
	}
	if ibr >= len(tronc.vendange) {
		ibr = len(tronc.vendange) - 1
	}
	sol := tronc.vendange[ibr]
	// graphe
	var src []string
	src = append(src, tronc.gr)
	for _, n := range sol.nods {
		src = append(src, n.graf()...)
	}
	if expl {
		fmt.Println("\n---- source ----\n")
		for _, n := range sol.nods {
			fmt.Println(n.doc())
		}
	}
	if j {
		fmt.Println("--- journal ----")
		fmt.Println(strings.Join(journal, "\n"))
	}
	fmt.Println("----------------")
	// numérotation de la solution
	fmt.Printf("%d/%d\n", ibr+1, len(tronc.vendange))
	// graphe en arcs
	initArcs()
	fmt.Println(strings.Join(graphe(src), "\n"))
}

// choix du texte latin
func chxTexte() {
	ClearScreen()
	fmt.Println(entete)
	fmt.Println(licence)
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

// lemmatisation du mot courant
func lemmatise() {
	texte.affiche(aidePh)
	im := tronc.imot
	tronc.imot = im
	mc := tronc.motCourant()
	fmt.Println("lemmatisation", rouge(mc.gr))
	fmt.Println(restostring(mc.ans))
}

// surlignage du mot précédent
func motprec() {
	if tronc.imot > 0 {
		tronc.imot--
		texte.affiche(aidePh)
	}
}

// surlignage du mot suivant
func motsuiv() {
	if tronc.imot < len(mots)-1 {
		tronc.imot++
		texte.affiche(aidePh)
	}
}

// saisie d'une phrase, et ajout au début de la liste
func saisie() {
	ClearScreen()
	fmt.Println(entete)
	fmt.Println("Phrase à analyser :")
	reader := bufio.NewReader(os.Stdin)
	p, _ := reader.ReadString('\n')
	texte.phrases = append([]string{p}, texte.phrases...)
	texte.compteur = 0
	texte.majPhrase()
	texte.affiche(aidePh)
	initArcs()
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
	var modeA bool  // source du graphe
	var modeJ bool  // débogage de l'arbre
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
			modeA = false
			modeJ = false
			analyse(false, false)
		case "g":
			modeA = true
			analyse(true, false)
		case "d":
			modeJ = true
			analyse(true, true)
		case "p":
			ibr--
			analyse(modeA, modeJ)
		case "s":
			ibr++
			analyse(modeA, modeJ)
		case "r":
			chxTexte()
		case "i":
			saisie()
		case "x":
			fmt.Println("\nVale")
			os.Exit(0)
		}
	}
	return
}
