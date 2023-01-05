## Useful git tools:
-	`git difftool -d go1.10.8 go1.19.4`
		-> compare 2 git tags as folders in meld
			useful for quick overlook of changes
-	`git log -p -- runtime/netpoll_windows.go`
		-> see all commits and changes to a file
			use this after you found differences with difftool
-	`git show 50f4896b72d16b6538178c8ca851b20655075b7f:src/runtime/netpoll_windows.go`
		-> view the file as a whole when it shows diffs from log -p
			sometimes commits can look messy, this makes it much better to read

From experience this is the best way to go about it, since you run into trouble once the fallback code is removed and the main code changes.

Search for a string anywhere
	`git log -S <search string> --source --all`


## How to strap:
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


## Cross compile with the new go to the dest system
	cd into go-backports/src
	run ./make.bash
	this gives you a new compiler that has the changes you made to go applied
	so you then can set the env vars and use that compiler to generate the correct output
	#!/bin/bash
	export GOMODCACHE=/path/to/different/cache
	export GOPATH=/path/to/different/path
	export PATH="/path/to/your/go-backports/bin:$PATH"
	$SHELL
	-> then you have a terminal set to the new go install where you can `GOOS=windows GOARCH=amd64 go build ...` from (dont use this in bootstrap!)