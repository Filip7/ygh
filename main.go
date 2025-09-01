package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type Args struct {
	S bool
}

var args Args

func main() {
	flag.BoolVar(&args.S, "S", false, "Do the sync operation")
	flag.Parse()
	packageNames := flag.Args()

	installLoc := "/tmp/ygh"
	if _, err := os.Stat(installLoc); os.IsNotExist(err) {
		err := os.Mkdir(installLoc, 0755)
		if err != nil {
			log.Fatalln("error creating '/tmp/ygh' directory")
		}
	}

	for _, pkg := range packageNames {
		url := fmt.Sprintf("https://raw.githubusercontent.com/archlinux/aur/refs/heads/%s/PKGBUILD", pkg)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("error happened fetching PKGBUILD for " + pkg)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error happend ", err)
		}

		var input string
		fmt.Printf("Do you want to read the PKGBUILD for %s? [y/N]: ", pkg)
		fmt.Scanln(&input)
		switch input {
		case "y", "Y":
			pkgbuild := fmt.Sprintf("%s\n", body)
			fmt.Printf("%s", pkgbuild)
		case "N", "n", "":
		default:
			fmt.Println("Continue")
		}

		fmt.Println("Cloning the repo...")
		if err := runCommand("git", "clone", "--branch", pkg, "--single-branch", "https://github.com/archlinux/aur.git", installLoc+"/"+pkg); err != nil {
			fmt.Println("error happend cloning git repo")
			fmt.Println(err)
		}

		fmt.Println("Building the package")
		if err := runCommand("makepkg", "-D", installLoc+"/"+pkg); err != nil {
			fmt.Println("error happend making the packages")
			fmt.Println(err)
		}

		fmt.Println("Installing the package")
		builtPkg := findPkgName(pkg, installLoc)
		if err := runCommand("sudo", "pacman", "-U", installLoc+"/"+pkg+"/"+builtPkg); err != nil {
			fmt.Println("error happend installing the packages")
			fmt.Println(err)
		}
	}
}

func findPkgName(pkgName string, location string) string {
	root := os.DirFS(location + "/" + pkgName)
	pkgs, err := fs.Glob(root, "*.pkg.tar.zst")
	if err != nil {
		log.Fatalln("error happend locating built package", err)
	}

	for i, f := range pkgs {
		if strings.Contains(f, "debug") {
			pkgs = slices.Delete(pkgs, i, i)
		}
	}

	return pkgs[0]
}

func runCommand(name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	return command.Run()
}
