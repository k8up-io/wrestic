FROM docker.io/golang:1.13-alpine as build

RUN set -x; apk add --no-cache wget bzip2 ca-certificates git gcc && \
    git clone https://github.com/vshn/restic && cd restic && \
    git checkout 2319-dump-dir-tar-rebase && go run -mod=vendor build.go -v && \
    mv restic /usr/local/bin/restic && chmod +x /usr/local/bin/restic && \
    mkdir /.cache && chmod -R 777 /.cache

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
