# Changelog

## [1.2.0](https://github.com/Filip7/ygh/compare/v1.1.0...v1.2.0) (2025-09-03)


### Features

* add option to edit the PKBUILD before installing ([2758f89](https://github.com/Filip7/ygh/commit/2758f89004ecc32578647e60f9d37f4028327f85))

## [1.1.0](https://github.com/Filip7/ygh/compare/v1.0.0...v1.1.0) (2025-09-01)


### Features

* use $HOME/.cache/ygh as install dir for pkgs ([#3](https://github.com/Filip7/ygh/issues/3)) ([1bb0807](https://github.com/Filip7/ygh/commit/1bb08078a452498ae8988f1ab91c63594016cf04))
* use makepkg for installing instead of calling pacman ([4831448](https://github.com/Filip7/ygh/commit/48314486e2bf62a62e8ab0c304ebe44eb02bd96b))

## 1.0.0 (2025-09-01)


### Features

* ask user if he wants to continue after checking the PKGBUILD ([c30a41b](https://github.com/Filip7/ygh/commit/c30a41b8c6eaa381b8e1fb6132b800ed565290ad))
* handle repo already existing - fetch instead ([c30a41b](https://github.com/Filip7/ygh/commit/c30a41b8c6eaa381b8e1fb6132b800ed565290ad))
* implement fetching of PKGBUILD and displaying it to the user ([82ac394](https://github.com/Filip7/ygh/commit/82ac394ee2b32cee8270731b7d382b4580887d5d))
* initial workflow, does the entier install ([9a34f9c](https://github.com/Filip7/ygh/commit/9a34f9c551902e713555e32fe048bf7557444f7d))


### Bug Fixes

* skip instead of exit in case we have multiple packages that we need to install ([0d0b9e9](https://github.com/Filip7/ygh/commit/0d0b9e965980f0fc8c16bce5225f811179c4cdcf))
