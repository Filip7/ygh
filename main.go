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
	S      bool
	editor string
}

var args Args

func main() {
	flag.BoolVar(&args.S, "S", false, "Do the sync operation")
	flag.StringVar(&args.editor, "editor", "", "Define editor to use when editing the PKGBUILD")
	flag.Parse()
	packageNames := flag.Args()

	installLoc := setupInstallDir()

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
		fmt.Printf("Do you want to read the PKGBUILD for %s before downloading? [y/N]: ", pkg)
		fmt.Scanln(&input)
		switch input {
		case "y", "Y":
			shouldSkip := handlePKGBUILDShowing(body, pkg)
			if shouldSkip {
				continue
			}
		}

		fmt.Println("Cloning the repo...")
		pkgInstall := installLoc + "/" + pkg
		if _, err := os.Stat(pkgInstall); os.IsNotExist(err) {
			if err := runCommand("git", "clone", "--branch", pkg, "--single-branch", "https://github.com/archlinux/aur.git", pkgInstall); err != nil {
				fmt.Println("error happend cloning git repo")
				fmt.Println(err)
			}
		} else {
			fmt.Println("Repo exists...updating")
			if err := runCommand("git", "fetch", pkgInstall); err != nil {
				fmt.Println("error happend updating git repo")
				fmt.Println(err)
			}
		}

		fmt.Printf("Do you want to edit the PKGBUILD for %s? [y/N]: ", pkg)
		fmt.Scanln(&input)
		switch input {
		case "y", "Y":
			editor := getEditor()
			if err := runCommand(editor, pkgInstall+"/PKGBUILD"); err != nil {
				fmt.Println("error happend trying to open the editor")
				fmt.Println(err)
			}
		}

		fmt.Println("Building and installing the package")
		// instal dependencies, install, clean
		if err := runCommand("makepkg", "-f", "-s", "-i", "-c", "-D", pkgInstall); err != nil {
			fmt.Println("error happend making the packages")
			fmt.Println(err)
		}
	}
}

func setupInstallDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Error getting user home directory")
	}

	installLoc := homeDir + "/.cache/ygh"
	if _, err := os.Stat(installLoc); os.IsNotExist(err) {
		err := os.Mkdir(installLoc, 0755)
		if err != nil {
			log.Fatalln("error creating" + installLoc + " directory")
		}
	}
	return installLoc
}

func getEditor() string {
	if args.editor != "" {
		return args.editor
	}

	// Fall back to $EDITOR environment variable
	editor := os.Getenv("EDITOR")
	if editor != "" {
		return editor
	}

	// Final fallback
	return "vi"
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
