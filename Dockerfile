FROM docker.io/golang:1.13-alpine as build

ENV RESTIC_VERSION=0.9.5 \
    SHASUM=08cd75e56a67161e9b16885816f04b2bf1fb5b03bc0677b0ccf3812781c1a2ec

WORKDIR /tmp
RUN set -x; apk add --no-cache wget bzip2 ca-certificates && \
    wget -q -O restic.bz2 https://github.com/restic/restic/releases/download/v${RESTIC_VERSION}/restic_${RESTIC_VERSION}_linux_amd64.bz2 && \
    echo "${SHASUM}  restic.bz2" | sha256sum -c - && \
    bzip2 -d restic.bz2 && \
    mv restic /usr/local/bin/restic && \
    chmod +x /usr/local/bin/restic && \
    mkdir /.cache && chmod -R 777 /.cache

RUN set -x; apk add --no-cache wget bzip2 ca-certificates git gcc && \
    git clone https://github.com/vshn/restic && cd restic && \
    git checkout 2319-dump-dir-tar && go run -mod=vendor build.go -v && \
    mv restic /usr/local/bin/restic && chmod +x /usr/local/bin/restic

WORKDIR /go/src/git.vshn.net/vshn/wrestic
COPY . .

RUN go install -v ./...

# runtime image
FROM docker.io/alpine:3
WORKDIR /app

RUN mkdir /.cache && chmod -R g=u /.cache
RUN apk --no-cache add ca-certificates

COPY --from=build /usr/local/bin/restic /usr/local/bin/restic
COPY --from=build /go/bin/wrestic /app/

USER 1001

ENTRYPOINT [ "./wrestic" ]
