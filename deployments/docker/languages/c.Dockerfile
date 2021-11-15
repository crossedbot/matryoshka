FROM matryoshka/runner

USER root
WORKDIR /tmp
RUN set -eux; \
    whoami && pwd; \
    apt-get update && \
    apt-get install -q -y --fix-missing \
        build-essential && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

USER nested
WORKDIR /home/nested
COPY ./deployments/docker/languages/dependecies/c/Makefile .
RUN gcc --version
