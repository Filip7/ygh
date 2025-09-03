# ygh

**y**et another aur helper fetching from the **g**it**h**ub mirror

`y` is doing a lot of heavy lifting here :D

At the time of creating this project AUR is very often under DDoS attack. Since other AUR helpers do not support GitHub mirror this project was created.

> [!IMPORTANT]
> Do not use this if you are not familiar with all the implications of building and using packages from AUR.
> Consult this wiki page <https://wiki.archlinux.org/title/Arch_User_Repository>

## Install

Download the package from GitHub release page <https://github.com/Filip7/ygh/releases>  
Unpack the `ygh_Linux_x86_64.tar.gz` and move the binary to `/usr/bin`

You can also download the checksum file to verify that the package is not tempered with and is generated with release-please gh action.

For every update you have to repeat those steps.

Here is an example:

```sh
# Download the package
curl -LO https://github.com/Filip7/ygh/releases/download/v1.2.0/ygh_Linux_x86_64.tar.gz

# Download the checksums file
curl -LO https://github.com/Filip7/ygh/releases/download/v1.2.0/ygh_1.2.0_checksums.txt

# Generate the SHA256 checksum of the downloaded file
sha256sum ygh_Linux_x86_64.tar.gz > downloaded_checksum.txt

# Extract the expected checksum for ygh_Linux_x86_64.tar.gz from the checksums file
grep ygh_Linux_x86_64.tar.gz ygh_1.2.0_checksums.txt > expected_checksum.txt

# Compare the generated checksum with the expected checksum
diff downloaded_checksum.txt expected_checksum.txt

# If the files are identical, there will be no output.
# If they differ, diff will show the differences.
```

### Install from source

Clone the repo using git.  
Run `go build` inside the downloaded directory then start the program like this `./ygh`

```sh
go build
./ygh -S yay
```

## Usage

Currently application only supports `-S` command for installing packages.  
Decided to continue with the `pacman` style flags.

If you supply multiple package names, it will try and install all of them.

Use it like this

### Installation of packages

```sh
ygh -S yay pacaur
```

If the app detects that the app is already downloaded, it will try to update it instead.

`ygh` supports editing the PKGBUILD before installing the package. It will either use `$EDITOR` variable or you can specify it like this:

```sh
ygh --editor nvim -S yay
```

### Update packages

Use `-Syu` option to update packages installed using `ygh`

```sh
ygh -Syu
```

This will walk all the installed packages, git fetch and install them.

### Deletion of packages

```sh
ygh -R pacaur
```

This will use `pacman` in the background to delete the package and will also delete the `git` repo stored in the `~/.cache/ygh` directory.

## Future plan

- If git pull has nothing to pull, skip update
- Make code a bit better, centralize logging of err msgs, move code to it's own packages or files
- For update, fetch updates from multiple goroutines instead of doing it in a for loop
- Look into the ALPM and see if that could be used directly
- Support any other `pacman` style commands?

Etc, etc.

Basically learn the Arch ecosystem

## Should you use this as your main AUR helper?

No, **absolutely not**, this is not fully featured nor battle tested.  
Also, it's recommended and advised to use only the official AUR repository as a source, even though the GitHub mirror is official, it's marked as experimental!  
If you are looking for recommendations, look at [yay](https://github.com/Jguer/yay)

## Contributing

Contributions are welcome
