//     util.go - Publicola

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	lcas    = [...]string{"nominatif", "vocatif", "accusatif", "génitif", "datif", "ablatif", "locatif"}
	lgenre  = [...]string{"masculin", "féminin", "neutre"}
	lnombre = [...]string{"singulier", "pluriel"}
)

// renvoie le premier cas de la liste lcas contenu dans morpho
func cas(morpho string) string {
	for _, c := range lcas {
		if strings.Contains(morpho, c) {
			return c
		}
	}
	return ""
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

// renvoie le premier genre de la liste lnombre contenu dans morpho
func nombre(morpho string) string {
	for _, n := range lnombre {
		if strings.Contains(morpho, n) {
			return n
		}
	}
	return ""
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

// renvoie le premier élément du split(s, sep)
func PrimEl(s, sep string) string {
	eclats := strings.Split(s, sep)
	return eclats[0]
}
