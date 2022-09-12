FROM matryoshka/runner

USER root
WORKDIR /tmp
RUN set -eux; \
    URL='https://dl.google.com/go/go1.19.1.linux-amd64.tar.gz'; \
    SHA256='acc512fbab4f716a8f97a8b3fbaa9ddd39606a28be6c2515ef7c6c6311acffde'; \
    whoami && pwd; \
    wget -O go.tgz.asc "${URL}.asc" --progress=dot:giga; \
    wget -O go.tgz "${URL}" --progress=dot:giga; \
    echo "${SHA256} *go.tgz" | sha256sum --strict --check -; \
    tar -C /usr/local -xzf go.tgz

USER nested
WORKDIR /home/nested
ENV PATH /usr/local/go/bin:$PATH
RUN go version
