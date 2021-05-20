FROM docker.io/golang:1.16 as build

ARG RESTIC_VERSION=0.12.0

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
RUN BUILD_VERSION=$(git describe --tags --always --dirty --match=v* || (echo "command failed $$?"; exit 1)) \
 && go install -v -ldflags "-X main.Version=$BUILD_VERSION -X 'main.BuildDate=$(date)'" ./...

# nonroot image
FROM docker.io/alpine:3 as nonroot
WORKDIR /app

RUN mkdir /.cache && chmod -R g=u /.cache
RUN apk --no-cache add ca-certificates

COPY --from=build /build/restic /usr/local/bin/restic
COPY --from=build /go/bin/wrestic /app/

USER 1001

ENTRYPOINT [ "./wrestic" ]

# root image
FROM nonroot

USER 0
