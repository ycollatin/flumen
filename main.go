//   main.go --	Gentes
//	analyseur syntaxique du latin

package main

// FIXME
// [Coll] uereare (uerearis)
// Tri :
//    - Albianum te confecisse arbitror : mauvais tri des solutions.
//    - favoriser les arcs non croisés
// Pyrrha dicitur esse creata. boucle infinie, due au lexsynt dico2:...,attrp
// Pour d'autres phrases aussi.
//
// XXX
//
// - CONSTRUCTIONS LEXICALES chez Cic.
// 		* pietate erga te
//      * tempus agendi et cogitandi
// TODO
// - in Italia speramus fore : la règle ter:hgprep permet d'avoir la bonne solution,
//   mais en 2ème position. Il faudrait une ligne supplémentaire pour donner un indice
//   de priorité d'une règle par rapport à d'autres.
// - la coordination -que est difficile pour "linquamus naturam artisque uideamus."
// - quand un verbe a un sujet et un objet et qu'on ne peut dire lequel
//   est sujet et lequel l'objet, c'est le premier qui est sujet :
//   Cic. Div. I, 50 : Ita res somnium conprobauit.
// - quand un groupe peut être sujet et objet d'un v. transitif, le groupe
//   le prmier est plutôt objet : lapidibus duo consules ceciderunt.
// - POS des romains : dies xxx nondum fuerant.
// - lexsynt.la : identifier et changer les initiales majuscules u > V
// - pseudovariables pour groupes.la
// - groupe pour les adj neutres + est : malum, opus, necesse est
// - a: p: essayer un préfixe ap: pour économiser le nombre de groupes.
// - une option pour charger une base de groupes différente ?
// - un signet pour les textes longs ?
// - dicitur prima mortalis. avec lexsynt dico1:attrp, boucle infinie.
// - verbes se construisant avec locatif. Caesarem Sinuessae mansurum nuntiabant.
// - Pour éviter la pléthore :
//   - durcir les conditions des règles
//   - élaguer les branches avant la fin
//   - élaguer les branches après récolte
//   - utiliser les goroutines : go bf.explore()
//   - enrichir la syntaxe de groupes.la
// - Si pléthore, trouver un moyen de navigation en fixant des arcs.
// - un champ groupe.anrel - analyses du relatif ?
// - accord :de personne sujet-verbe et verbe-verbe coord; de mode v-v, etc.
// - fonction de sortie au format GraphViz.
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
	ibr   int // rang de l'analyse (branche) courante
	rouge func(...interface{}) string
	vert  func(...interface{}) string
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
	mc := texte.motCourant()
	fmt.Println("lemmatisation", rouge(mc.gr))
	var res gocol.Res
	if tronc.vendange != nil {
		sol := tronc.vendange[ibr]
		for _, nod := range sol.nods {
			res = append(res, nod.toRes(mc)...)
		}
		fmt.Println(resToString(res, mc.ans))
	} else {
		fmt.Println(gocol.Restostring(mc.ans))
	}
}

// surlignage du mot précédent
func motprec() {
	if texte.imot > 0 {
		texte.imot--
		texte.affiche(aidePh)
	}
}

// surlignage du mot suivant
func motsuiv() {
	if texte.imot < len(mots)-1 {
		texte.imot++
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
	vert = color.New(color.FgGreen, color.Bold).SprintFunc()
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
	var modeA bool // source du graphe
	var modeJ bool // débogage de l'arbre
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
