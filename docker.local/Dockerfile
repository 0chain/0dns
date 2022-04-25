FROM golang:1.18.0-alpine3.15 as zdns_build


RUN apk add --update --no-cache build-base git

ENV SRC_DIR=/0dns
ENV GO111MODULE=on
#ENV GOPROXY=http://10.10.10.100:3080

# cached go mod download
ADD ./code/go/0dns.io/go.mod $SRC_DIR/go/0dns.io/
ADD ./code/go/0dns.io/go.sum $SRC_DIR/go/0dns.io/

WORKDIR $SRC_DIR/go/0dns.io
RUN go mod download -x

# Add the source code:
ADD ./code/go/0dns.io $SRC_DIR/go/0dns.io

# Download the dependencies:
WORKDIR $SRC_DIR/go/0dns.io/zdnscore/zdns




RUN go build -v -tags "bn256 development" -ldflags "-X 0chain.net/core/build.BuildTag=$GIT_COMMIT"

# Copy the build artifact into a minimal runtime image:
FROM alpine:3.15
ENV APP_DIR=/0dns
WORKDIR $APP_DIR

COPY --from=zdns_build $APP_DIR/go/0dns.io/zdnscore/zdns/zdns $APP_DIR/bin/zdns

