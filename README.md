# ygh

**y**et another aur helper fetching from the **g**it**h**ub mirror

`y` is doing a lot of heavy lifting here :D

At the time of creating this project AUR is very often under DDoS attack. Since other AUR helpers do not support GitHub mirror this project was created.

> [!IMPORTANT]
> Do not use this if you are not familiar with all the implications of building and using packages from AUR.
> Consult this wiki page <https://wiki.archlinux.org/title/Arch_User_Repository>

## Usage

Currently application only supports `-S` command for installing packages.  
Decided to continue with the `pacman` style flags.

If you supply multiple package names, it will try and install all of them.

Use it like this

```sh
./ygh -S yay pacaur
```

If the app detects that the app is already downloaded, it will try to update it instead.

`ygh` supports editing the PKGBUILD before installing the package. It will either use `$EDITOR` variable or you can specify it like this:

```sh
./ygh --editor nvim -S yay
```

## Future plan

The plan is to support more of the `pacman` commands. Look into the ALPM project and maybe use it instead of calling makepkg via command execute.  
Investigate how to implement updates.  
Etc.

## Should you use this as your main AUR helper?

No, absolutely not, this is not fully featured nor battle tested.  
If you are looking for recommendations, look at [yay](https://github.com/Jguer/yay)

## Contributing

Contributions are welcome
