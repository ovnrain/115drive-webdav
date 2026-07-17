FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

# GoReleaser dockers_v2 places binaries under $TARGETPLATFORM/
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/115drive-webdav /usr/bin/115drive-webdav

ENTRYPOINT ["/usr/bin/115drive-webdav"]
