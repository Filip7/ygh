package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Args struct {
	S      bool
	R      bool
	Syu    bool
	editor string
}

var args Args

func main() {
	flagProcessing()
	packageNames := flag.Args()

	installLoc := setupInstallDir()

	if args.S {
		doInstall(packageNames, installLoc)
	}

	if args.Syu {
		doUpdate(installLoc)
	}

	if args.R {
		doRemove(packageNames, installLoc)
	}

	flag.Usage()
}

func flagProcessing() {
	flag.BoolVar(&args.S, "S", false, "Install packages")
	flag.BoolVar(&args.Syu, "Syu", false, "Update packages")
	flag.BoolVar(&args.R, "R", false, "Remove packages")
	flag.StringVar(&args.editor, "editor", "", "Define editor to use when editing the PKGBUILD")
	flag.Parse()
}

func doInstall(packageNames []string, installLoc string) {
	for _, pkg := range packageNames {
		url := fmt.Sprintf("https://raw.githubusercontent.com/archlinux/aur/refs/heads/%s/PKGBUILD", pkg)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("error happened fetching PKGBUILD for " + pkg)
			os.Exit(1)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error happend ", err)
			os.Exit(1)
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
				os.Exit(1)
			}
		} else {
			fmt.Println("Repo exists...updating " + pkg)
			if err := runCommand("git", "pull", pkgInstall); err != nil {
				fmt.Println("error happend updating git repo")
				fmt.Println(err)
				os.Exit(1)
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
				os.Exit(1)
			}
		}

		fmt.Println("Building and installing the package " + pkg)
		// instal dependencies, install, clean
		if err := runCommand("makepkg", "-f", "-s", "-i", "-c", "-D", pkgInstall); err != nil {
			fmt.Println("error happend making the packages")
			fmt.Println(err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

func doUpdate(installLoc string) {
	entries, err := os.ReadDir(installLoc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, dir := range entries {
		pkgInstall := installLoc + "/" + dir.Name()
		pkg := dir.Name()

		fmt.Println("Updating " + dir.Name())
		if err := runCommand("git", "pull", pkgInstall); err != nil {
			fmt.Println("error happend updating git repo")
			fmt.Println(err)
			os.Exit(1)
		}
		var input string
		fmt.Printf("Do you want to edit the PKGBUILD for %s? [y/N]: ", pkg)
		fmt.Scanln(&input)
		switch input {
		case "y", "Y":
			editor := getEditor()
			if err := runCommand(editor, pkgInstall+"/PKGBUILD"); err != nil {
				fmt.Println("error happend trying to open the editor")
				fmt.Println(err)
				os.Exit(1)
			}
		}

		fmt.Println("Building and installing the package " + pkg)
		// instal dependencies, install, clean
		if err := runCommand("makepkg", "-f", "-s", "-i", "-c", "-D", pkgInstall); err != nil {
			fmt.Println("error happend making the packages")
			fmt.Println(err)
			os.Exit(1)
		}

	}
	os.Exit(0)
}

func doRemove(packageNames []string, installLoc string) {
	var input string
	fmt.Printf("Are you sure you want to delete %s? [Y/n]:", strings.Join(packageNames, " "))
	fmt.Scanln(&input)
	if input == "n" || input == "N" {
		fmt.Println("Stoping...")
		os.Exit(1)
	}

	for _, pkg := range packageNames {
		fmt.Println("Uninstalling " + pkg)
		if err := runCommand("sudo", "pacman", "-R", pkg); err != nil {
			fmt.Println("error happend deleting the package")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Deleting local git repo of " + pkg)
		if err := runCommand("rm", "-rf", installLoc+"/"+pkg); err != nil {
			fmt.Println("error happend deleting the package")
			fmt.Println(err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

func setupInstallDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory")
		os.Exit(1)
	}

	installLoc := homeDir + "/.cache/ygh"
	if _, err := os.Stat(installLoc); os.IsNotExist(err) {
		err := os.Mkdir(installLoc, 0755)
		if err != nil {
			fmt.Println("error creating" + installLoc + " directory")
			os.Exit(1)
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
