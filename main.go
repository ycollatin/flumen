//   main.go --	Gentes
//	analyseur syntaxique du latin

package main

//     FIXME
//
// Il manque le moyen d'appliquer la règle suivante :
// dans une proposition infinitive,
//     le sujet précède l'objet
//     il est au premier rang de l'infinitive
//     s'il n'y a qu'un accusatif, c'est le sujet
//
// ter:v.capuam
// n:@v;;act;;mvm
// 1 a:@NP "domus";lieu;acc;;mvm
// 2 a:"domus";lieu;acc
// 3 a:@NP;lieu;acc;;mvm
//
// si la dernière ligne est 1, rien ne marche
// si 2 ou 3, OK
//
//  idem pour n.hnAdj
//
//   Tri :
//    - favoriser les arcs non croisés
//
// XXX
//
// - dans vargraph, commenté iui:i et i$:ii
// - CONSTRUCTIONS LEXICALES chez Cic.
//      * tempus agendi et cogitandi
//
// TODO GROUPES
// - omni officio : n.app apparaît avant n.det. Comment résoudre ?
// - Deux sujets : Sustulimus manus et ego et Balbus (m6)
// - pos multiples ?
// - traiter les praenomina M. L. etc.
// - Le vocatif confondu avec nomin ou acc : Vere loquar, iudices.
// - dans la périphrase inf. futur, l'aux. esse est souvent omis :
//   /responsurum hominem nemo arbitrabatur./
// - Dans les phrases brèves, "est" final souvent omis : quid autem absurdum ? quid enim molestius ?
// - sujet d'une inf. futur elliptique : omnia se facturum recepit.
// - épithète d'un sujet elliptique : nunquam uidelicet sitiens biberat.
// - Hiérarchie des règles :
//   /in Italia speramus fore/ La règle ter:hgprep permet d'avoir la bonne solution,
//   mais en 2ème position. Il faut une solution pour donner un indice
//   de priorité d'une règle par rapport à d'autres.
// - la coordination -que est difficile pour "linquamus naturam artisque uideamus."
// - quand un verbe a un sujet et un objet et qu'on ne peut dire lequel
//   est sujet et lequel l'objet, c'est le premier qui est sujet :
//   Cic. Div. I, 50 : Ita res somnium conprobauit.
// - quand un groupe peut être sujet et objet d'un v. transitif, le groupe
//   le premier est plutôt objet : lapidibus duo consules ceciderunt.
// - POS des romains : dies xxx nondum fuerant.
// - lexsynt.la : identifier et changer les initiales majuscules u > V
// - groupe pour les adj neutres + est : malum, opus, necesse est
// - parasitage de /sum/ par /edo/ : comment supprimer "excl" dans lexsynt
// - parasitage de /do/ par /dato/ :   "
//
// TODO COMMANDES
// - option "c" -  analyse : donner la fonction (lien)
// - a: p: essayer un préfixe ap: pour économiser le nombre de groupes.
// - pseudovariables pour groupes.la
// - une option pour charger une base de groupes différente ?
// - une commande pour atteindre une phrase dans les textes longs ?
// - un champ groupe.anrel - analyses du relatif ?
// - accord :de personne sujet-verbe et verbe-verbe coord; de mode v-v, etc.

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
f->enregistrer la sortie;
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
	scrb, ajout  string    // tampon et résultat d'analyse
	vert  func(...interface{}) string
)

// affiche les arcs syntaxique de la phrase
// expl : la source du graphe est affichée
// j : débogage de l'arbre
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
	scrb = ""
	ajout = ""
	if tronc.vendange == nil {
		scrb = "échec de l'analyse"
		fmt.Println(scrb)
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
		ajout = "\n---- source ----\n"
		scrb += ajout
		fmt.Println(ajout)
		for _, n := range sol.nods {
			scrb += alin(n.doc(false))
			fmt.Println(n.doc(true))
		}
	}
	if j {
		ajout = "--- journal ----\n"
		scrb += ajout
		fmt.Print(ajout)
		ajout = strings.Join(journal, "\n")
		scrb += alin(ajout)
		fmt.Println(ajout)
	}
	ajout = "----------------\n"
	fmt.Print(ajout)
	scrb += ajout
	// numérotation de la solution
	ajout = fmt.Sprintf("%d/%d\n", ibr+1, len(tronc.vendange))
	fmt.Print(ajout)
	scrb += ajout
	// graphe en arcs
	initArcs()
	ajout = strings.Join(graphe(src), "\n")
	fmt.Println(ajout)
	scrb += alin(ajout)
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

// TODO : factoriser le calcul de src
func dot() {

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
	scrb = ""
	ajout = ""
	if tronc.vendange == nil {
		scrb = "échec de l'analyse"
		fmt.Println(scrb)
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

	texte.affiche(aidePh)
	fmt.Println(strings.Join(srcDot(src), "\n"))
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
		case "a":
			modeA = false
			modeJ = false
			analyse(false, false)
		case "c":
			lemmatise()
		case "d":
			modeJ = true
			analyse(true, true)
		case "g":
			modeA = true
			analyse(true, false)
		case "f":
			log(scrb)
		case "h":
			motprec()
		case "i":
			saisie()
		case "j":
			texte.porro()
			ibr = 0
		case "k":
			texte.retro()
		case "l":
			motsuiv()
			ibr = 0
		case "p":
			ibr--
			analyse(modeA, modeJ)
		case "r":
			chxTexte()
		case "s":
			ibr++
			analyse(modeA, modeJ)
		case "t":
			dot()
		case "x":
			fmt.Println("\nVale")
			os.Exit(0)
		}
	}
}
