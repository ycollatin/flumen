//  texte.go - Publicola

// Partage d'un texte en []Phrase

package main

import (
	"fmt"
	"github.com/ycollatin/gocol"
	"strings"
)

type Texte struct {
	nom      string
	compteur int
	phrases  []string
	phrase   *Phrase
}

func (t *Texte) append(p string) {
	t.phrases = append(t.phrases, p)
}

// Efface l'écran, affiche un entête, la phrase, et
// le texte du param aide.
func (t Texte) affiche(aide string) {
	ClearScreen()
	fmt.Printf("%s, phrase %d, mot %d\n", t.nom, t.compteur, texte.phrase.imot)
	fmt.Println(t.phrase.enClair())
	fmt.Println(aide)
}

// crée un texte à partir du fichier nommé nf
func CreeTexte(nf string) *Texte {
	var tp string // texte de la phrase
	ll := gocol.Lignes(chCorpus + nf)
	t := new(Texte)
	t.nom = nf
	for _, l := range ll {
		for {
			ifp := strings.IndexAny(l, ".?;!")
			if ifp < 0 {
				tp += l + " "
				break
			} else {
				tp += l[:ifp+1] //+ " "
				// supprimer l'espace initiale
				if tp > "" && tp[0] == ' ' {
					tp = tp[1:]
				}
				// créer et ajouter la nouvelle phrase
				t.phrases = append(t.phrases, tp)
				tp = ""
				l = l[ifp+1:]
			}
		}
	}
	return t
}

func (t *Texte) majPhrase() {
	initArcs()
	t.phrase = creePhrase(t.phrases[t.compteur])
}

func (t *Texte) porro() {
	if len(t.phrases)-t.compteur > 1 {
		t.compteur++
		t.majPhrase()
		t.affiche(aidePh)
	}
}

func (t *Texte) retro() {
	if t.compteur > 0 {
		t.compteur--
		t.majPhrase()
		t.affiche(aidePh)
	}
}
