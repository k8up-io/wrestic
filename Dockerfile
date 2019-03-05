FROM golang:1.11-alpine

ENV RESTIC_VERSION=0.9.3 \
    SHASUM=3c882962fc07f611a6147ada99c9909770d3e519210fd483cde9609c6bdd900c

WORKDIR /tmp
# TODO: re-enable this when APPU-1524 is done
# RUN set -x; apk add --no-cache wget bzip2 ca-certificates && \
#     wget -q -O restic.bz2 https://github.com/restic/restic/releases/download/v${RESTIC_VERSION}/restic_${RESTIC_VERSION}_linux_amd64.bz2 && \
#     echo "${SHASUM}  restic.bz2" | sha256sum -c - && \
#     bzip2 -d restic.bz2 && \
#     mv restic /usr/local/bin/restic && \
#     chmod +x /usr/local/bin/restic && \
#     mkdir /.cache && chmod -R 777 /.cache

RUN set -x; apk add --no-cache wget bzip2 ca-certificates git gcc && \
    git clone https://github.com/Kidswiss/restic && cd restic && \
    git checkout tar && go run -mod=vendor build.go -v && \
    mv restic /usr/local/bin/restic && chmod +x /usr/local/bin/restic && \
    mkdir /.cache && chmod -R 777 /.cache

WORKDIR /go/src/git.vshn.net/vshn/wrestic
COPY . .

RUN go install -v ./...

ENTRYPOINT [ "/go/bin/wrestic" ]
