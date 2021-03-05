# Golang Boilerplate - GoZones

This is a boilerplate application to quickly create new Golang applications.  It can support File mode or Server mode transactions and includes additional helper functions for logging, file interactions, and other common operations - in addition to other developmental building blocks such as a Container defintion, GoReleaser configuration and matching GitHub Actions that are ready to use with no or little set up.

*This document is a work in progress and this boilerplate is likely to evolve over time*

[![Tests](https://github.com/kenmoini/golang-boilerplate/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kenmoini/golang-boilerplate/actions/workflows/test.yml) [![release](https://github.com/kenmoini/golang-boilerplate/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/kenmoini/golang-boilerplate/actions/workflows/release.yml)

***In this documentation, "GoZones/go-zones" can be replaced with your application name***

## Example Commands & Parameters

```bash
# File Mode - input source, output target (default mode)
$ ./go-zones -source=./zones.yml -dir=./generated
# Server Mode
$ ./go-zones -mode server -config=./config.yml
```

## Deployment - As a Container

This boilerplate comes with a `Containerfile` that can be built with Docker or Podman with the following command:

```bash
# Build the container
podman build -f Containerfile -t go-zones .

# Create a config directory locally with a server configuration YAML file
mkdir config && cp config.yml.example config/config.yml

# Mount that directory and run a container
podman run -p 8080:8080 -v config/:/etc/go-zones go-zones
```

### Adding extra files to the container image

If you need additional assets along side the Golang binary in the built container you can simply place them in the `container_root` directory - directories/files in this `container_root` directory will be copied to the root of the container file system.  You can find an example of using a touchfile to create the `/etc/go-zones/` directory in the built container.

## Deployment - Building From Source

Since this is just a Golang application, as long as you have Golang v1.15+ then the following commands will do the job:

```bash
go build

./go-zones
```

Of course, once you change the name of the application the executable name will change as well.

## Starting Development

### Initial Changes

You'll likely want to change a few things, namely the name.

- You can find the `appName` defined in the `variables.go` file.  This is what the application references internally for logs and so on.
- You will also need to change the name/path in the `go.mod` file to match your repository path/name.  This is where the executable package gets its name.
- Rename the template touch folder in `container_root/etc/go-zones/`.
- Change the references in this `README.md` file to match your application name.
- Adjust the `.gitignore` file to match what will be the name of your executable package.
- Modify the `go-zones.service` file to match if you're utilizing this in Server Mode with SystemD as a Service

## Lifecycle

As the versioning of your application progresses, make sure to keep that semantic version up to date in the `appVersion` variable defined in the `variables.go` file.

Once you are ready to release a new version of your application, you can utilize GoReleaser and GitHub Actions to create packaged Go binaries of your application across a matrix of operating systems and architectures.

## Creating a Release

### Generating GPG Keys for Signing

### Creating the Repository Secrets in GitHub

### Creating a New Release in GitHub

