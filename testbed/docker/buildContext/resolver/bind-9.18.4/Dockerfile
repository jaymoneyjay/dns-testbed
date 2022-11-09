FROM ubuntu:bionic

RUN apt-get update \
  && apt-get install -y \
  bind9 \
  bind9utils \
  bind9-doc \
  vim \
  wget \
  tar \
  iproute2 \
  libexpat1-dev \
    libssl-dev \
    build-essential \
    pkg-config \
    libuv1.dev \
    libnghttp2-dev \
    flex \
    bison \
    libcap-dev \
    openssl \
    libxml2

RUN wget https://downloads.isc.org/isc/bind9/9.18.4/bind-9.18.4.tar.xz
RUN tar -xf bind-9.18.4.tar.xz
RUN ./bind-9.18.4/configure
RUN make
RUN make install

# Enable IPv4
RUN sed -i 's/OPTIONS=.*/OPTIONS="-4 -u bind"/' /etc/default/bind9
