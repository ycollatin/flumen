//     util.go - flumen

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ycollatin/gocol"
)

var (
	lcas    = [...]string{"nominatif", "vocatif", "accusatif", "génitif", "datif", "ablatif", "locatif"}
	lgenre  = [...]string{"masculin", "féminin", "neutre"}
	lnombre = [...]string{"singulier", "pluriel"}
)

func alin(s string) string {
	return fmt.Sprintf("%s\n", s)
}

// renvoie le premier cas de la liste lcas contenu dans morpho
func cas(morpho string) string {
	for _, c := range lcas {
		if strings.Contains(morpho, c) {
			return c
		}
	}
	return ""
}

func appendRes(resa, resb gocol.Res) gocol.Res {
	for _, srb := range resb {
		var ai bool
		for i, sra := range resa {
			if sra.Lem == srb.Lem {
				ai = true
				for _, mb := range srb.Morphos {
					if !contient(sra.Morphos, mb) {
						resa[i].Morphos = append(resa[i].Morphos, mb)
					}
				}
			}
		}
		if !ai {
			resa = append(resa, resb...)
		}
	}
	return resa
}

// efface l'écran
func ClearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

// vrai si s est un élément de ls
func contient(ls []string, s string) bool {
	for _, e := range ls {
		if e == s {
			return true
		}
	}
	return false
}

func log(s string) {
	nf := "log-flumen.txt"
	f, err := os.OpenFile(nf, os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		f, _ = os.OpenFile(nf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	}
	defer f.Close() // on ferme automatiquement à la fin de notre programme
	f.WriteString(s)
}

// renvoie le premier genre de la liste lgenre contenu dans morpho
func genre(morpho string) string {
	for _, g := range lgenre {
		if strings.Contains(morpho, g) {
			return g
		}
	}
	return ""
}

// capture de la dernière touche enfoncée
func GetKey() string {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	// réactiver le retour arrière que l'une des dernières commandes a désactivé
	defer exec.Command("stty", "-F", "/dev/tty", "icanon").Run()
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	return string(b)
}

// saisie d'un entier
func InputInt(q string) int {
	var i int
	fmt.Printf("%s ", q)
	_, err := fmt.Scanf("%d", &i)
	if err != nil {
		return -1
	}
	return i
}

// renvoie le premier genre de la liste lnombre contenu dans morpho
func nombre(morpho string) string {
	for _, n := range lnombre {
		if strings.Contains(morpho, n) {
			return n
		}
	}
	return ""
}

// renvoie le premier élément du split(s, sep)
func PrimEl(s, sep string) string {
	eclats := strings.Split(s, sep)
	return eclats[0]
}

// Compare les lemmatisations Flumen et les lemmatisations Collatinus
// Les premières sont un sous-ensemble des secondes, et seront
// colorées en vert.
func resToString(resFlumen, resCol gocol.Res) (res string) {
	mapg := make(map[string][]string) // clé : Clé de lemme
	mapc := make(map[string][]string) // valeur : morphos
	var clesg, clesc []string         // clés des deux maps
	var trsc []string
	// map et clés Flumen
	for _, srg := range resFlumen {
		k, ll := srToString(srg)
		clesg = append(clesg, k)
		mapg[k] = ll
	}
	// map et clés Collatinus
	for _, src := range resCol {
		k, ll := srToString(src)
		clesc = append(clesc, k)
		trsc = append(trsc, srToTr(src))
		mapc[k] = ll
	}
	var lres []string
	// pour chaque cle de clesc
	for i, clec := range clesc {
		if contient(clesg, clec) {
			// si elle est dans clg : vert
			lres = append(lres, vert(trsc[i]))
			//  afficher toutes les morphos de cleg
			for _, morf := range mapg[clec] {
				lres = append(lres, vert(morf))
			}
			// afficher les clés de mapc absentes de mapg
			for _, morf := range mapc[clec] {
				if !contient(mapg[clec], morf) {
					lres = append(lres, morf)
				}
			}
		} else {
			// sinon, tout en normal
			lres = append(lres, trsc[i])
			morfc := mapc[clec]
			for _, mc := range morfc {
				lres = append(lres, mc)
			}
		}
	}
	return strings.Join(lres, "\n")
}

func srcDot(src []string) (dot []string) {
	dot = append(dot, "digraph D {")
	for i, m := range mots {
		dot = append(dot, fmt.Sprintf("%d [label=\"%s\"];", i, m.gr))
	}
	for i, ls := range src {
		if i == 0 {
			dot = append(dot, fmt.Sprintf("label=\"%s\";", ls))
		} else {
			ls = strings.Replace(ls, "[", "[label=\"", 1)
			ls = strings.Replace(ls, "]", "\"]", 1)
			dot = append(dot, ls+";")
		}
	}
	dot = append(dot, "}")
	return dot
}

// renvoie la lemmatisation sr sous forme de chaîne
func srToString(sr gocol.Sr) (k string, ll []string) {
	for _, l := range sr.Morphos {
		ll = append(ll, fmt.Sprintf("   %s", l))
	}
	return sr.Lem.Cle, ll
}

// renvoie la traduction du lemme de sr
func srToTr(sr gocol.Sr) (tr string) {
	return fmt.Sprintf("%s [%s] : %s", sr.Lem.Gr, sr.Lem.Pos, sr.Lem.Traduction)
}
