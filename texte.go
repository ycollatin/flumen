//  texte.go - Gentes

// Partage d'un texte en []Branche

package main

import (
	"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

type Texte struct {
	nom      string
	compteur int
	imot     int       // rang du mot courant
	phrases  []string
}

var tronc *Branche

func (t *Texte) append(p string) {
	t.phrases = append(t.phrases, p)
}

// Efface l'écran, affiche un entête, la Branche, et
// le texte du param aide.
func (t Texte) affiche(aide string) {
	ClearScreen()
	fmt.Printf("%s\n%s, phrase %d, mot %d\n", entete, t.nom, t.compteur+1, texte.imot)
	fmt.Println(tronc.enClair())
	fmt.Println(aide)
}

// crée un texte à partir du fichier nommé nf
func CreeTexte(nf string) *Texte {
	var tp string // texte de la Branche
	ll := gocol.Lignes(chCorpus + nf)
	t := new(Texte)
	t.nom = nf
	for _, l := range ll {
		for {
			ifp := strings.IndexAny(l, ".?;!:")
			if ifp < 0 {
				tp += l + " "
				break
			} else {
				tp += l[:ifp+1] //+ " "
				// supprimer l'espace initiale
				if tp > "" && tp[0] == ' ' {
					tp = tp[1:]
				}
				// ajouter la nouvelle phrase
				t.append(tp)
				tp = ""
				l = l[ifp+1:]
			}
		}
	}
	return t
}

// initialise la phrase sur laquelle
// pointe t.compteur
func (t *Texte) majPhrase() {
	initArcs()
	tronc = creeTronc(t.phrases[t.compteur])
}

// avance d'une phrase
func (t *Texte) porro() {
	if len(t.phrases)-t.compteur > 1 {
		t.compteur++
		t.majPhrase()
		t.affiche(aidePh)
	}
}

// recule d'une phrase
func (t *Texte) retro() {
	if t.compteur > 0 {
		t.compteur--
		t.majPhrase()
		t.affiche(aidePh)
	}
}
