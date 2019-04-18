FROM arm32v7/golang:stretch

COPY qemu-arm-static /usr/bin/
WORKDIR /go/src/github.com/automatedhome/onewire
COPY . .
RUN go build -o onewire cmd/main.go

FROM arm32v7/busybox:1.30-glibc

COPY --from=0 /go/src/github.com/automatedhome/onewire/onewire /usr/bin/onewire

ENTRYPOINT [ "/usr/bin/onewire" ]
