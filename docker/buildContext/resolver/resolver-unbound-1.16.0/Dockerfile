FROM ubuntu:bionic

RUN apt-get update \
    && apt-get install -y \
    vim \
    wget \
    tar \
    flex \
    bison \
    libexpat1-dev \
    libssl-dev \
    build-essential

RUN wget https://www.nlnetlabs.nl/downloads/unbound/unbound-1.16.0.tar.gz
RUN tar xzf unbound-1.16.0.tar.gz
RUN ./unbound-1.16.0/configure
RUN make
RUN make install

# Run eternal loop
CMD ["/bin/bash", "-c", "while :; do sleep 10; done"]