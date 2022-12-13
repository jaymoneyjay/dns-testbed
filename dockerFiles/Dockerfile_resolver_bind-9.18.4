FROM ubuntu:kinetic

RUN apt-get update \
  && apt-get install -y \
  vim \
    iproute2\
    locate  \
    gnupg \
    tar \
    wget

#RUN add-apt-repository universe
RUN apt-get install -y \
    build-essential \
    libtool  \
    autoconf  \
    automake  \
    libssl-dev



RUN apt-get install -y \
    libuv1-dev \
    libnghttp2-dev \
    libcap-dev \
    libprotobuf-c-dev libprotobuf-c1 \
    protobuf-c-compiler \
    fstrm-bin libfstrm0 libfstrm-dev \
    libprotobuf-c-dev libprotobuf-c1


RUN wget https://ftp.isc.org/isc/bind9/9.18.4/bind-9.18.4.tar.xz

RUN tar xf bind-9.18.4.tar.xz
WORKDIR bind-9.18.4
RUN ./configure --enable-dnstap --sysconfdir=/etc/bind --localstatedir=/  \
    --enable-threads --enable-largefile --with-libtool --enable-shared --with-gnu-ld

RUN make
RUN make install
RUN ldconfig

RUN apt-get install -y bind9
RUN apt remove -y bind9

RUN mkdir /var/cache/bind
RUN chgrp bind /var/cache/bind
RUN chmod g+w /var/cache/bind

