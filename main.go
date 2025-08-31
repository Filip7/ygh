package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

type Args struct {
	s bool
	q bool
}

var args Args

func main() {
	flag.BoolVar(&args.s, "S", false, "Do the sync operation")
	flag.Parse()
	packageNames := flag.Args()

	for _, pkg := range packageNames {
		resp, err := http.Get("https://raw.githubusercontent.com/archlinux/aur/refs/heads/" + pkg + "/PKGBUILD")
		if err != nil {
			fmt.Println("error happened fetching PKGBUILD for " + pkg)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error happend ", err)
		}

		var input string
		fmt.Print("Do you want to read the PKGBUILD for " + pkg + "? [y/N]:")
		fmt.Scanln(&input)
		switch input {
		case "y", "Y":
			fmt.Printf("%s \n", body)
		case "N", "n", "":
			continue
		default:
			fmt.Println("Option not recognized, defaulting to NO")
		}
	}
}
