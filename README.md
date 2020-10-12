# Fuse

Use fuse to patch/create local files (text based) into remote repositories hosted on different providers and create pull requests.

## Build

The build process depends on golangci-lint for enforce sensible standards. Thus, to build fuse you'll need it.
Download instructions can be found (here)[https://golangci-lint.run/usage/install/#local-installation].

Once installer, run:

    make
    
To release fuse run:

    make <version>

Binaries for linux, MacOs and Windows 64bit are generated and located at releases/<version> .

## Usage

Basic usage example for azure dev ops:

    fuse azdevops --orgUrl <organization url> --project <my-azdevops-project> --pat <personal-auth-token>  --repoName <target repo name> --contentDir <directory-with-files-to-patch>

The _<directory-with-files-to-patch>_ must be an *absolute path* and it's structure should be the same as the target repo, otherwise, expected patches will be interpreted as new files.

To see all available commands and flags, run:

    fuse --help
    fuse azdevops --help
    fuse github --help
    
## Supported providers

- [Azure Devops](https://dev.azure.com/)
- [Github](https://github.com/)

## TODO's
    * Tests