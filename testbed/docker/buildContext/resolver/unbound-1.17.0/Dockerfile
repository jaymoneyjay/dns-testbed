FROM ubuntu:bionic

RUN apt-get update \
    && apt-get install -y \
    vim \
    unbound \
    libssl-dev

# Run eternal loop
CMD ["/bin/bash", "-c", "while :; do sleep 10; done"]