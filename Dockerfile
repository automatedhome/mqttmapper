FROM arm32v7/golang:stretch

COPY qemu-arm-static /usr/bin/
WORKDIR /go/src/github.com/automatedhome/mqttmapper
COPY . .
RUN make build

FROM arm32v7/busybox:1.30-glibc

COPY --from=0 /go/src/github.com/automatedhome/mqttmapper/mqttmapper /usr/bin/mqttmapper
COPY config.yaml /config.yaml

ENTRYPOINT [ "/usr/bin/mqttmapper" ]
