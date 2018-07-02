FROM golang:1.10 as builder

COPY . $GOPATH/src/git.vshn.net/vshn/wrestic
RUN cd $GOPATH/src/git.vshn.net/vshn/wrestic/cmd/wrestic && go build -o /build/wrestic

FROM alpine:3.7

ENV RESTIC_VERSION=0.9.1

RUN apk update && apk add wget bzip2 && cd / && \
    wget https://github.com/restic/restic/releases/download/v${RESTIC_VERSION}/restic_${RESTIC_VERSION}_linux_amd64.bz2 && \
    bzip2 -d restic_${RESTIC_VERSION}_linux_amd64.bz2 && \
    rm restic_${RESTIC_VERSION}_linux_amd64.bz2 && \
    mkdir -p /usr/local/bin/ && \
    mv /restic_${RESTIC_VERSION}_linux_amd64 /usr/local/bin/restic && \
    chmod +x /usr/local/bin/restic

COPY --from=builder /build/wrestic /usr/local/bin/wrestic

ENTRYPOINT [ "/usr/local/bin/wrestic" ]
