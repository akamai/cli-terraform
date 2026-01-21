#!/bin/bash
# Creates unsigned binaries for macOS (64bit), and Linux/Windows (32 and 64bit) 

if [[ -z "$1" ]]
then
	echo "Version not supplied."
	echo "Usage: build.sh <version>"
	exit 1
fi

mkdir -p build/"$1"

GOOS=darwin GOARCH=arm64 go build -o build/"$1"/akamai-terraform-"$1"-macarm64 -ldflags="-X 'github.com/akamai/cli-terraform/v2/cli.Version=$1'" .
shasum -a 256 build/"$1"/akamai-terraform-"$1"-macarm64 | awk '{print $1}' > build/"$1"/akamai-terraform-"$1"-macarm64.sig
GOOS=darwin GOARCH=amd64 go build -o build/"$1"/akamai-terraform-"$1"-macamd64 -ldflags="-X 'github.com/akamai/cli-terraform/v2/cli.Version=$1'" .
shasum -a 256 build/"$1"/akamai-terraform-"$1"-macamd64 | awk '{print $1}' > build/"$1"/akamai-terraform-"$1"-macamd64.sig
GOOS=linux GOARCH=arm64 go build -o build/"$1"/akamai-terraform-"$1"-linuxarm64 -ldflags="-X 'github.com/akamai/cli-terraform/v2/cli.Version=$1'" .
shasum -a 256 build/"$1"/akamai-terraform-"$1"-linuxarm64 | awk '{print $1}' > build/"$1"/akamai-terraform-"$1"-linuxarm64.sig
GOOS=linux GOARCH=amd64 go build -o build/"$1"/akamai-terraform-"$1"-linuxamd64 -ldflags="-X 'github.com/akamai/cli-terraform/v2/cli.Version=$1'" .
shasum -a 256 build/"$1"/akamai-terraform-"$1"-linuxamd64 | awk '{print $1}' > build/"$1"/akamai-terraform-"$1"-linuxamd64.sig
GOOS=linux GOARCH=386 go build -o build/"$1"/akamai-terraform-"$1"-linux386 -ldflags="-X 'github.com/akamai/cli-terraform/v2/cli.Version=$1'" .
shasum -a 256 build/"$1"/akamai-terraform-"$1"-linux386 | awk '{print $1}' > build/"$1"/akamai-terraform-"$1"-linux386.sig
GOOS=windows GOARCH=386 go build -o build/"$1"/akamai-terraform-"$1"-windows386.exe -ldflags="-X 'github.com/akamai/cli-terraform/v2/cli.Version=$1'" .
shasum -a 256 build/"$1"/akamai-terraform-"$1"-windows386.exe | awk '{print $1}' > build/"$1"/akamai-terraform-"$1"-windows386.exe.sig
GOOS=windows GOARCH=amd64 go build -o build/"$1"/akamai-terraform-"$1"-windowsamd64.exe -ldflags="-X 'github.com/akamai/cli-terraform/v2/cli.Version=$1'" .
shasum -a 256 build/"$1"/akamai-terraform-"$1"-windowsamd64.exe | awk '{print $1}' > build/"$1"/akamai-terraform-"$1"-windowsamd64.exe.sig
