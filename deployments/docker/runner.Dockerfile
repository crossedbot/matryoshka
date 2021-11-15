FROM debian:buster-slim

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
COPY --from=matryoshka/builder /go/bin/runner ./bin/runner

CMD [ "runner" ]
