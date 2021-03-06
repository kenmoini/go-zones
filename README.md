# GoZones

[![Tests](https://github.com/kenmoini/go-zones/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kenmoini/go-zones/actions/workflows/test.yml) [![release](https://github.com/kenmoini/go-zones/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/kenmoini/go-zones/actions/workflows/release.yml)

GoZones is an application that will take DNS Zones as defined in YAML and generate BIND-compatable DNS Zone files and the configuration required to load the zone file.

GoZones can operate in single-file input/output batches, or via an HTTP server.

## Example Commands & Parameters

```bash
# File Mode - input source, output target (default mode)
$ ./go-zones -source=./zones.yml -dir=./generated
# Server Mode
$ ./go-zones -mode server -config=./config.yml
```

## Deployment - As a Container

GoZones comes with a set of `Containerfile`s that can be built with Docker or Podman with the following commands:

### Server Mode Container

```bash
# Build the container
podman build -f Containerfile -t go-zones .

# Create a config directory locally with a server configuration YAML file
mkdir config && cp config.yml.example config/config.yml

# Mount that directory and run a container
podman run -p 8080:8080 -v config:/etc/go-zones go-zones
```

### File Mode Fronting to BIND

There is an extra `Containerfile.file-to-bind` container definition file that will set up a container image that starts GoZones to generate zone files and BIND configuration, then starts a BIND DNS Server.

```bash
# Build the container
podman build -f Containerfile.file-to-BIND -t go-zones:file-to-bind .

# Create a config directory locally with a Zones configuration YAML file for file mode operation
mkdir config && cp zones.yml.example config/zones.yml

# Mount that directory and run a container
podman run -p 8053:8053 -v config:/etc/go-zones go-zones:file-to-bind
```

### Adding extra files to the container image

If you need additional assets along side the Golang binary in the built container you can simply place them in the `container_root` directory - directories/files in this `container_root` directory will be copied to the root of the container file system.  You can find an example of using a touchfile to create the `/etc/go-zones/` directory in the built container.

## Deployment - Building From Source

Since this is just a Golang application, as long as you have Golang v1.15+ then the following commands will do the job:

```bash
go build

./go-zones
```