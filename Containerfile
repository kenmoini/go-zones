FROM registry.access.redhat.com/ubi8/go-toolset:1.14.12 AS build

WORKDIR /opt/app-root/src
COPY . .
RUN go build

FROM scratch AS bin

COPY --from=build /opt/app-root/src/go-zones /usr/local/bin/
COPY container_root/ /

EXPOSE 8080

CMD [ "go-zones -mode server -config /etc/go-zones/config.yml" ]