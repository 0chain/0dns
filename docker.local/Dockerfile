FROM golang:1.11.4-alpine3.8 as zdns_build


RUN apk add --update --no-cache build-base git

ENV SRC_DIR=/0dns
ENV GO111MODULE=on


# Add the source code:
ADD ./code/go/0dns.io $SRC_DIR/go/0dns.io

# Download the dependencies:
WORKDIR $SRC_DIR/go/0dns.io/zdnscore/zdns


RUN go mod download

RUN go build -v -tags "bn256 development" -ldflags "-X 0chain.net/core/build.BuildTag=$GIT_COMMIT"

# Copy the build artifact into a minimal runtime image:
FROM golang:1.11.4-alpine3.8
ENV APP_DIR=/0dns
WORKDIR $APP_DIR

COPY --from=zdns_build $APP_DIR/go/0dns.io/zdnscore/zdns/zdns $APP_DIR/bin/zdns

