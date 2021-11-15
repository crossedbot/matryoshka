FROM matryoshka/runner

USER root
WORKDIR /tmp
RUN set -eux; \
    URL='https://dl.google.com/go/go1.16.5.linux-amd64.tar.gz'; \
    SHA256='b12c23023b68de22f74c0524f10b753e7b08b1504cb7e417eccebdd3fae49061'; \
    whoami && pwd; \
    wget -O go.tgz.asc "${URL}.asc" --progress=dot:giga; \
    wget -O go.tgz "${URL}" --progress=dot:giga; \
    echo "${SHA256} *go.tgz" | sha256sum --strict --check -; \
    tar -C /usr/local -xzf go.tgz

USER nested
WORKDIR /home/nested
ENV PATH /usr/local/go/bin:$PATH
RUN go version
