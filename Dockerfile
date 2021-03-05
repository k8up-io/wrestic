FROM docker.io/golang:1.15 as build

ENV RESTIC_VERSION=0.12.0

RUN set -x; \
    apt-get update \
 && apt-get install -y \
      bzip2 \
      ca-certificates \
      gcc \
      git \
      wget \
 && wget "https://github.com/restic/restic/releases/download/v${RESTIC_VERSION}/restic_${RESTIC_VERSION}_linux_amd64.bz2" \
 && bunzip2 "restic_${RESTIC_VERSION}_linux_amd64.bz2" \
 && mkdir /build \
 && mv "restic_${RESTIC_VERSION}_linux_amd64" /build/restic \
 && chmod +x /build/restic

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test -v ./...
ENV CGO_ENABLED=0
RUN go install -v ./...

# runtime image
FROM docker.io/alpine:3
WORKDIR /app

RUN mkdir /.cache && chmod -R g=u /.cache
RUN apk --no-cache add ca-certificates

COPY --from=build /build/restic /usr/local/bin/restic
COPY --from=build /go/bin/wrestic /app/

ENTRYPOINT [ "./wrestic" ]
