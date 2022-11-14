FROM ubuntu:bionic

RUN apt update
RUN apt install -y wget \
    dpkg \
    socat \
    systemd \

RUN wget https://secure.nic.cz/files/knot-resolver/knot-resolver-release.deb
RUN dpkg -i knot-resolver-release.deb
RUN apt install -y knot-resolver

RUN useradd -ms /bin/bash knot-resolver

USER knot-resolver
WORKDIR /home/knot-resolver