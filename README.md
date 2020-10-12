# Fuse

Use fuse to patch/create local files (text based) into remote repositories hosted on different providers and create pull requests.

## Build

From the root of the project run:

    go build
    go test
    
    ./fuse azdevops ...

If you make any contribution, run this before submitting your PR:

    golangci-lint run --fast --exclude-use-default=false
    
Download instructions are (here)[https://golangci-lint.run/usage/install/#local-installation].

## Usage

Basic usage example for azure dev ops:

    fuse azdevops --orgUrl <organization url> --project <my-azdevops-project> --pat <personal-auth-token>  --repoName <target repo name> --contentDir <directory-with-files-to-patch>

The content directory structure should be the same as the target repo, otherwise, expected patches will be interpreted as new files.

To see all available commands and flags, run:

    fuse --help
    fuse azdevops --help
    
## Release

To release fuse, at the moment just run the script located at `scripts/release.sh` with the proper version:

    bash scripts/release.sh <version>

At the moment, linux, MacOs and Windows 64bit binaries are generated.

## TODO's

    * Improve README (include dependencies and differente usage patterns)
    * Create a MAKEFILE
    * Implemement tests
    * Support more providers, github, gitlab, bitbucket ...