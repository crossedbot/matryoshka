FROM matryoshka/runner

USER root
WORKDIR /tmp
RUN set -eux; \
    URL='https://dl.google.com/go/go1.17.8.linux-amd64.tar.gz'; \
    SHA256='980e65a863377e69fd9b67df9d8395fd8e93858e7a24c9f55803421e453f4f99'; \
    whoami && pwd; \
    wget -O go.tgz.asc "${URL}.asc" --progress=dot:giga; \
    wget -O go.tgz "${URL}" --progress=dot:giga; \
    echo "${SHA256} *go.tgz" | sha256sum --strict --check -; \
    tar -C /usr/local -xzf go.tgz

USER nested
WORKDIR /home/nested
ENV PATH /usr/local/go/bin:$PATH
RUN go version
