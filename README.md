# The Go Programming Language

Go is an open source programming language that makes it easy to build simple,
reliable, and efficient software.

The canonical Git repository is located at https://go.googlesource.com/go.
There is a mirror of the repository at https://github.com/golang/go.

Unless otherwise noted, the Go source files are distributed under the
BSD-style license found in the LICENSE file.

##  Fork alert!

You are viewing a fork right now that adds support for older Windows Versions.
Currently this is Windows XP (NT 5.1) and Windows 7 (NT 6.1).
The code is in the desired branches you can check out.

## Current Branches

### Windows XP
- `release-branch.go1.19-nt51` Go 1.19 for Windows XP (32/64) - overwhelmingly functional
- `release-branch.go1.20-nt51` Go 1.20 for Windows XP (32/64) - overwhelmingly functional
- `release-branch.go1.21-nt51` Go 1.21 for Windows XP (32/64) - it keeps getting harder, but it still works

### Windows 7
- `release-branch.go1.21-nt61` Go 1.21 for Windows 7 (the first one officially unsupported)

## Updates

I will update a branch on request, please open an issue about it on Github if you need an updated toolchain.


#### Bootstrapping from another OS Example

Adapt the parameters as you need them, for example `amd64`->`386`.

	git clone --branch release-branch.go1.19-nt51 --single-branch https://github.com/BieHDC/go-backports
	cd go-backports/src/
	GOOS=windows GOARCH=amd64 ./bootstrap.bash
	cd ..
	cd ..
	-> copy go-windows-amd64-bootstrap.tbz to the destination os
	-> zip up the "go-backports" folder and copy it also to the destination os
	-> unzip them to C:\
	-> run these commands (and adapt the paths if you changed them)
	set GOROOT_BOOTSTRAP=C:\go-windows-amd64-bootstrap
	set GOROOT=C:\go
	set CGO_ENABLED=0
	cd C:\go\src
	.\all.bat (or make.bat if you dont want tests)
	-> when all goes right, you just need to follow "Install From Source" below and its done.
	NOTES: Yes you need both and yes they need to be different folders, you cant bootstrap into bootstrap!

#### Install From Source

Run `set "PATH=%PATH%;c:\go\bin` in cmd.exe exactly like that and you have `go` available after bootstrapping.
You can make that permanent with `setx` or use the gui.

### Contributing

See CONTRIBUTING.md

