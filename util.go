//     util.go - Publicola

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func afflin(ll []string, h int) {
	fmt.Println("afflin, len",len(ll), "h", h)
	var d int
	if h > len(ll) {
		h = len(ll)
	}
	for {
		for i:= d; i < len(ll) && i < d+h; i++ {
			fmt.Println(ll[i])
		}
		fmt.Println("d page suiv., u page préc., q quitter")
		k := GetKey()
		switch (k) {
		case "d":
			ClearScreen()
			d += h
			if d > len(ll) {
				d = d - h
				if d < 0 {
					d = 0
				}
			}
		case "u":
			ClearScreen()
			d -= h
			if d < 0 {
				d = 0
			}
		case "q":
			return
		}
	}
}

var (
	lcas = [...]string {"nominatif","vocatif","accusatif","génitif","datif","ablatif","locatif"}
	lgenre = [...]string {"masculin","féminin","neutre"}
	lnombre = [...]string {"singulier", "pluriel"}
)

func cas(morpho string) string {
	for _, c := range lcas {
		if strings.Contains(morpho, c) {
			return c
		}
	}
	return ""
}

func genre(morpho string) string {
	for _, g := range lgenre {
		if strings.Contains(morpho, g) {
			return g
		}
	}
	return ""
}

func nombre(morpho string) string {
	for _, n := range lnombre {
		if strings.Contains(morpho, n) {
			return n
		}
	}
	return ""
}

func chxMultiple(q string, ll []string) string {
	fmt.Println(q)
	for i := 0; i < len(ll); i++ {
		fmt.Println(i, ll[i])
	}
	c := InputInt("N° choix ?")
	return ll[c]
}

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

func sauveLignes(nf string, ll []string) {
	f, err := os.Create(nf)
	if err != nil {
		fmt.Println("erreur d'écriture")
		f.Close()
		return
	}
	for _, l := range ll {
		fmt.Fprintln(f, l)
	}
	f.Close()
}

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

func InputInt(q string) int {
	var i int
	fmt.Printf("%s ", q)
	_, err := fmt.Scanf("%d", &i)
	if err != nil {
		return InputInt(q)
	}
	return i
}

func InputString(q string) string {
	in := bufio.NewReader(os.Stdin)
	fmt.Print(q," ")
	s, err := in.ReadString('\n')
	if err == nil {
		return strings.TrimSpace(s)
	}
	return "err"
}

// renvoie le premier élément du split(s, sep)
func PrimEl(s, sep string) string {
	eclats := strings.Split(s, sep)
	return eclats[0]
}
