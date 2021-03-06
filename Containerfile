FROM registry.access.redhat.com/ubi8/go-toolset:1.14.12 AS build

WORKDIR /opt/app-root/src
COPY . .
RUN go build

FROM registry.fedoraproject.org/fedora-minimal:latest

COPY --from=build /opt/app-root/src/go-zones /usr/local/bin/
COPY container_root/ /

RUN microdnf update -y \
 && rm -rf /var/cache/yum \
 && chmod +x /opt/app-root/start-server.sh

EXPOSE 8080

CMD [ "/opt/app-root/start-server.sh" ]