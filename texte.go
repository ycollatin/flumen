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
	phrases		[]*Phrase
}

func (t *Texte) append(p *Phrase) {
  t.phrases = append(t.phrases, p)
}

func (t *Texte) phrase() *Phrase {
	return t.phrases[t.compteur]
}

func (t Texte) affiche(aide string) {
	ClearScreen()
	fmt.Printf("%s, phrase %d, mot %d\n", t.nom, t.compteur, imot)
	p := t.phrase()
	fmt.Println(p.enClair())
	fmt.Println(aide)
}

// crée un texte à partir du fichier nommé nf
func CreeTexte (nf string) *Texte {
	var tp	string	// texte de la phrase
	fmt.Println("creeTexte", nf,"-", chCorpus+nf)
	ll := gocol.Lignes(chCorpus+nf)
	t := new(Texte)
	t.nom = nf
	for _, l := range ll {
		for {
			ifp := strings.IndexAny(l, ".?;!")
			if ifp < 0 {
				tp += l
				fmt.Println("tp",tp)
				break;
			} else {
				tp += l[:ifp] + " "
				fmt.Println("else tp", tp)
				// créer et ajouter la nouvelle phrase
				p := creePhrase(tp)
				t.phrases = append(t.phrases, p)
				tp = ""
				l =	l[ifp+1:]
			}
		}
	}
	return t
}

func (t *Texte) majPhrase() {
	initArcs()
	phrase = t.phrases[t.compteur]
	t.affiche(aidePh)
}

func (t *Texte) porro() {
	if t.compteur < len(t.phrases) {
		t.compteur++
		t.majPhrase()
		imot = 0
		t.affiche(aidePh)
	}
}

func (t *Texte) retro() {
	if t.compteur > 0 {
		t.compteur--
		t.majPhrase()
		imot = 0
		t.affiche(aidePh)
	}
}
