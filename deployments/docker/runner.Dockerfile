ARG OS=debian:bullseye-slim
ARG GOLANG_VERSION=1.19-bullseye
ARG CGO=0
ARG GOOS=linux
ARG GOARCH=amd64

#-------------------------------------------------------------------------------
FROM golang:${GOLANG_VERSION} AS gobuilder

ARG CGO
ARG GOOS
ARG GOARCH

RUN go version
WORKDIR /go/src/
COPY . .
RUN cd cmd/runner && \
    CGO_ENABLED='${CGO}' GOOS='${GOOS}' GOARCH='${GOARCH}' \
    make -f /go/src/Makefile build

#-------------------------------------------------------------------------------
FROM ${OS}

ENV NESTED_HOME /usr/local/nested
ENV PATH ${NESTED_HOME}/bin:$PATH
ENV LANG C.UTF-8
WORKDIR ${NESTED_HOME}

# Install packages
RUN apt update && \
    apt-get install -y \
        bash \
        coreutils \
        wget && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create user(s)
RUN groupadd nested
RUN useradd -m -d /home/nested -g nested -s /bin/bash nested

# Set up environment
USER nested

RUN mkdir -vp ${NESTED_HOME}
COPY --from=gobuilder /go/bin/runner ./bin/runner

CMD [ "runner" ]
