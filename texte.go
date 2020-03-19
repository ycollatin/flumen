//  texte.go - Publicola 

// Partage d'un texte en []Phrase

package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"
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
func CreeTexte(nf string) *Texte {
	t, _ := ioutil.ReadFile("./corpus/" + nf)
	contenu := string(t)
	var (
		mot string
		p	*Phrase
		tp	string	// texte de la phrase
	)
	texte := new(Texte)
	for i :=0; i < len(contenu); i++ {
		r := contenu[i]
		// sauter les lignes "!.*$"
		if  r == '!' {
			for r != '\n' {
				i++
				r = contenu[i]
			}
		}
		s := string(r)
		if s != "\n" {
			tp += s
		}
		if unicode.IsLetter(rune(r)) {
			mot += s
		} else if mot > "" {
			if p == nil {
				p = new(Phrase)
			}
			p.append(creeMot(mot))
			mot = ""
			if strings.ContainsAny(".;?!", s) {
				p.gr = tp
				texte.append(p)
				p = nil
				tp = ""
			}
		}
	}
	texte.append(p)
	texte.nom = nf
	return texte
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
