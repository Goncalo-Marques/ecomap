FROM postgis/postgis

RUN apt-get update
RUN apt-get -y install build-essential wget git libboost-all-dev cmake postgresql-server-dev-16

RUN wget -O pgrouting-3.6.1.tar.gz https://github.com/pgRouting/pgrouting/archive/v3.6.1.tar.gz
RUN tar xfz pgrouting-3.6.1.tar.gz \
    && cd pgrouting-3.6.1 \
    && mkdir build \
    && cd build \
    && cmake .. \
    && make \
    && make install 
