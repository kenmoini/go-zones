# GoZones

[![Docker Repository on Quay](https://quay.io/repository/kenmoini/go-zones/status "Docker Repository on Quay")](https://quay.io/repository/kenmoini/go-zones) [![Tests](https://github.com/kenmoini/go-zones/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kenmoini/go-zones/actions/workflows/test.yml) [![release](https://github.com/kenmoini/go-zones/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/kenmoini/go-zones/actions/workflows/release.yml)

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

Pre-built container images can be found at https://quay.io/repository/kenmoini/go-zones

GoZones comes with a set of `Containerfile`s that can be built with Docker or Podman with the following commands:

### Server Mode Container

```bash
# Build the container
podman build -f Containerfile -t go-zones .

# Create a config directory locally with a server configuration YAML file
mkdir -p config && cp config.yml.example config/config.yml

# Mount that directory and run a container
podman run -p 8080:8080 -v "$(pwd)"/config:/etc/go-zones/ go-zones
```

### File Mode Fronting to BIND

There is an extra `Containerfile.file-to-bind` container definition file that will set up a container image that starts GoZones to generate zone files and BIND configuration, then starts a BIND DNS Server.

```bash
# Build the container
podman build -f Containerfile.file-to-BIND -t go-zones:file-to-bind .

# Create a config directory locally with a Zones configuration YAML file for file mode operation
mkdir -p config && cp zones.yml.example config/zones.yml

# Mount that directory and run a container
podman run -d -p 8053:8053 -v "$(pwd)"/config:/etc/go-zones/ go-zones:file-to-bind
```

#### Container-as-a-Service

If you're interested in running a GoZones File-to-BIND container as a service with a static IP, you can do so with the following (tested on RHEL 8.3 with Podman 2.2.1):

##### 1. Create Bridge Device

First you must create a Bridge Network Device on your system - this creates a virtual device, the bridge, that allows containers and VMs on your system to connect through to the network that system is connected to.

Creating a Bridge device is outside of the scope of this document, find the different ways to create one here: https://www.tecmint.com/create-network-bridge-in-rhel-centos-8/

##### 2. Create a new Podman Bridge Network

By default, containers will have access to a bridge device that connects the Pods to a NAT'd network.  This is not ideal for running static services for your network - instead, use Podman to create a new network that uses a macvlan-style container network to connect to a bridge device.  Run the following commands:

```bash
sudo podman create network lanBridge
sudo nano /etc/cni/net.d/lanBridge.conflist
```

Make the file look something like this, substituting for ***bridge*** (your bridge device), and your bridged network subnet and ***range*** Podman can utilize (it can overlap with your full subnet, DHCP would be passed off to the gateway through the bridge):

```json
{
   "cniVersion": "0.4.0",
   "name": "lanBridge",
   "plugins": [
      {
         "type": "bridge",
         "bridge": "LANbr0",
         "ipam": {
            "type": "host-local",
            "ranges": [
                [
                    {
                        "subnet": "192.168.42.0/24",
                        "rangeStart": "192.168.42.2",
                        "rangeEnd": "192.168.42.254",
                        "gateway": "192.168.42.1"
                    }
                ]
            ],
            "routes": [
                {"dst": "0.0.0.0/0"}
            ]
         }
      },
      {
         "type": "portmap",
         "capabilities": {
            "portMappings": true
         }
      },
      {
         "type": "firewall",
         "backend": ""
      },
      {
         "type": "tuning",
         "capabilities": {
            "mac": true
         }
      }
   ]
}
```

##### 3. Test the Container and Network

Before firing up a service wraped deployment, test the plumbing so far:

```bash
# Create needed directories
mkdir -p /opt/service-containers/dns-core-1/volumes/etc-conf

# Download an example zones file
curl https://raw.githubusercontent.com/kenmoini/go-zones/main/zones.yml.example -o /opt/service-containers/dns-core-1/volumes/etc-conf/zones.yml

# Test the container, assign it an IP in your bridged subnet range
podman run -d --name dns-core-1 --network lanBridge --ip 192.168.42.10 -p 53 -v /opt/service-containers/dns-core-1/volumes/etc-conf:/etc/go-zones/ quay.io/kenmoini/go-zones:file-to-bind
```

Note that the `-d` option will launch the container into the background - test the BIND DNS Server running in the container, the query for example.net should look like the following from an internal network:

```bash
dig @192.168.42.10 example.net

; <<>> DiG 9.11.20-RedHat-9.11.20-5.el8_3.1 <<>> @192.168.42.10 example.net
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 7346
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 2, ADDITIONAL: 3

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: 710f246c6fde139b45d1957e60453e6b526aded7b99e739f (good)
;; QUESTION SECTION:
;example.net.                   IN      A

;; ANSWER SECTION:
example.net.            3600    IN      A       192.168.42.100

;; AUTHORITY SECTION:
example.net.            86400   IN      NS      dns-core-1.example.labs.
example.net.            86400   IN      NS      dns-core-2.example.labs.

;; ADDITIONAL SECTION:
dns-core-1.example.labs. 86400  IN      A       192.168.42.2
dns-core-2.example.labs. 86400  IN      A       192.168.42.3

;; Query time: 0 msec
;; SERVER: 192.168.42.10#53(192.168.42.10)
;; WHEN: Sun Mar 07 15:58:19 EST 2021
;; MSG SIZE  rcvd: 178
```

Note that you may need to handle some SELinux contexts *(like disabling it lol jk kinda not jk)* - also don't forget to clean up the running container with `podman ps` and `podman kill dns-core-1 && podman rm dns-core-1`

##### 4. Creating a Service

Now that the service is tested to work properly you can create a service file and launch the container at system boot.

Reference the file `go-zones-file-to-bind-podman.service` in this repository for defining your service.  The file is ready to work with these example steps as a service called `dns-core-1` and the following steps will produce that resulting service:

```bash
# Download the service file
curl https://raw.githubusercontent.com/kenmoini/go-zones/main/go-zones-file-to-bind-podman.service -o /etc/systemd/system/dns-core-1.service

# Reload systemd services
systemctl daemon-reload

# Start service
systemctl start dns-core-1

# Check service status and running container
systemctl status dns-core-1
podman ps
```

### Adding extra files to the container image

If you need additional assets along side the Golang binary in the built container you can simply place them in the `container_root` directory - directories/files in this `container_root` directory will be copied to the root of the container file system.  You can find an example of using a touchfile to create the `/etc/go-zones/` directory in the built container.

## Deployment - Building From Source

Since this is just a Golang application, as long as you have Golang v1.15+ then the following commands will do the job:

```bash
go build

./go-zones
```