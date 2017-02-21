FROM registry.clickyab.ae/clickyab/ubuntu
MAINTAINER f0ruD <fzero@rubi.gd>

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y \
    wget \
    && wget https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64.deb \
    && dpkg -i dumb-init_*.deb \
    && mkdir /app \
    && rm dumb-init_*.deb \
    && rm -rf /var/lib/apt/lists/*

ADD ./init-script.sh /sbin/entry-point.sh
RUN chmod a+x /sbin/entry-point.sh
VOLUME /app

EXPOSE 80
ENTRYPOINT ["/sbin/entry-point.sh"]

CMD ["/app/bin/server"]