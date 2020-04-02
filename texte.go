//  texte.go - Publicola 

// Partage d'un texte en []Phrase

package main

import (
	"fmt"
	"strings"
	"github.com/ycollatin/gocol"
)

type Texte  struct {
	nom			string
	compteur	int
	phrases		[]string
}

func (t *Texte) append(p string) {
  t.phrases = append(t.phrases, p)
}

// Efface l'écran, affiche un entête, la phrase, et
// le texte du param aide.
func (t Texte) affiche(aide string) {
	ClearScreen()
	fmt.Printf("%s, phrase %d, mot %d\n", t.nom, t.compteur, phrase.imot)
	fmt.Println(phrase.enClair())
	fmt.Println(aide)
}

// crée un texte à partir du fichier nommé nf
func CreeTexte (nf string) *Texte {
	var tp	string	// texte de la phrase
	ll := gocol.Lignes(chCorpus+nf)
	t := new(Texte)
	t.nom = nf
	for _, l := range ll {
		for {
			ifp := strings.IndexAny(l, ".?;!")
			if ifp < 0 {
				tp += l + " "
				break;
			} else {
				tp += l[:ifp+1] //+ " "
				// supprimer l'espace initiale
				if tp > "" && tp[0] == ' ' {
					tp = tp[1:]
				}
				// créer et ajouter la nouvelle phrase
				t.phrases = append(t.phrases, tp)
				tp = ""
				l =	l[ifp+1:]
			}
		}
	}
	return t
}

func (t *Texte) majPhrase() {
	initArcs()
	phrase = creePhrase(t.phrases[t.compteur])
	t.affiche(aidePh)
}

func (t *Texte) porro() {
	if t.compteur < len(t.phrases) {
		t.compteur++
		t.majPhrase()
	}
}

func (t *Texte) retro() {
	if t.compteur > 0 {
		t.compteur--
		t.majPhrase()
	}
}
