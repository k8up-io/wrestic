FROM golang:1.10-alpine

ENV RESTIC_VERSION=0.9.1 \
    SHASUM=f7f76812fa26ca390029216d1378e5504f18ba5dde790878dfaa84afef29bda7

WORKDIR /tmp
RUN set -x; apk add --no-cache wget bzip2 ca-certificates && \
    wget -q -O restic.bz2 https://github.com/restic/restic/releases/download/v${RESTIC_VERSION}/restic_${RESTIC_VERSION}_linux_amd64.bz2 && \
    echo "${SHASUM}  restic.bz2" | sha256sum -c - && \
    bzip2 -d restic.bz2 && \
    mv restic /usr/local/bin/restic && \
    chmod +x /usr/local/bin/restic && \
    mkdir /.cache && chmod -R 777 /.cache

WORKDIR /go/src/git.vshn.net/vshn/wrestic
COPY . .

RUN go install -v ./...

ENTRYPOINT [ "wrestic" ]
