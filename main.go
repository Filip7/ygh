package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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
			shouldSkip := handlePKGBUILDShowing(body, pkg)
			if shouldSkip {
				continue
			}
		}

		fmt.Println("Cloning the repo...")
		if _, err := os.Stat(installLoc + "/" + pkg); os.IsNotExist(err) {
			if err := runCommand("git", "clone", "--branch", pkg, "--single-branch", "https://github.com/archlinux/aur.git", installLoc+"/"+pkg); err != nil {
				fmt.Println("error happend cloning git repo")
				fmt.Println(err)
			}
		} else {
			fmt.Println("Repo exists...updating")
			if err := runCommand("git", "fetch", installLoc+"/"+pkg); err != nil {
				fmt.Println("error happend updating git repo")
				fmt.Println(err)
			}
		}

		fmt.Println("Building and installing the package")
		// instal dependencies, install, clean
		if err := runCommand("makepkg", "-f", "-s", "-i", "-c", "-D", installLoc+"/"+pkg); err != nil {
			fmt.Println("error happend making the packages")
			fmt.Println(err)
		}
	}
}

func handlePKGBUILDShowing(body []byte, pkg string) bool {
	pkgbuild := fmt.Sprintf("%s\n", body)
	fmt.Printf("%s", pkgbuild)
	fmt.Printf("Do you want to continue with the build for %s? [Y/n]: ", pkg)
	var buildInput string
	fmt.Scanln(&buildInput)
	switch buildInput {
	case "n", "N":
		fmt.Printf("Skipping instalation of %s...", pkg)
		return true
	}
	return false
}

func runCommand(name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	return command.Run()
}
