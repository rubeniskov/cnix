ARG GO_IMAGE=golang
ARG GO_VERSION

FROM ${GO_IMAGE}:${GO_VERSION}-alpine AS builder

ARG VERSION
ARG VC_REF

WORKDIR /tmp/gobuild
COPY ./ .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v \
      -ldflags="-w -s \
       -X 'main.BuildTime=$(date +%s)' \
	 -X 'main.BuildVcRef=${VC_REF}' \
	 -X 'main.BuildVersion=${VERSION}' \
       -X 'main.BuildSource=docker' \
      " .

FROM scratch
ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
ENTRYPOINT [ "/bin/cnix" ]
COPY --from=builder /tmp/gobuild/cnix /bin/cnix
